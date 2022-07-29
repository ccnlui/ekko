package transceiver

import (
	"bytes"
	"context"
	"ekko/internal/config"
	"ekko/internal/util"
	"log"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
)

type AeronTransceiver struct {
	aeron     *aeron.Aeron
	pub       *aeron.Publication
	sub       *aeron.Subscription
	assembler *aeron.FragmentAssembler
	rcvdMsg   uint64
	histogram *hdrhistogram.Histogram
}

func NewAeronTransceiver(histogram *hdrhistogram.Histogram) *AeronTransceiver {
	log.Println("[info] new aeron transceiver!")

	tcv := &AeronTransceiver{
		histogram: histogram,
	}
	tcv.init()
	return tcv
}

func (tcv *AeronTransceiver) init() {
	aeronCtx := aeron.NewContext().AeronDir(config.AeronDir).MediaDriverTimeout(config.MediaDriverTimeout)
	a, err := aeron.Connect(aeronCtx)
	if err != nil {
		log.Fatalln("[fatal] failed to connect to media driver: ", config.AeronDir, err.Error())
	}
	tcv.aeron = a

	tcv.pub = <-a.AddPublication(config.ServerChannel, int32(config.ServerStreamID))
	for !tcv.pub.IsConnected() {
		time.Sleep(time.Millisecond)
	}
	log.Println("[info] publication connected to media driver:", tcv.pub)

	tcv.sub = <-a.AddSubscription(config.ClientChannel, int32(config.ClientStreamID))
	for !tcv.sub.IsConnected() {
		time.Sleep(time.Millisecond)
	}
	log.Println("[info] subscription connected to media driver:", tcv.sub)

	inBuf := bytes.NewBuffer(make([]byte, 0, config.MaxMessageSize))
	onMessage := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		// timestamp
		timestamp := buffer.GetInt64(offset)
		tcv.histogram.RecordValue(time.Now().UnixNano() - timestamp)

		// message
		inBuf.Reset()
		buffer.WriteBytes(inBuf, offset, length)
		tcv.rcvdMsg += 1
		// log.Printf("[debug] %8.d Got a fragment timestamp: %d offset: %d length: %d payload: %s\n",
		// 	tcv.rcvdMsg, timestamp, offset, length, string(inBuf.Next(int(length))),
		// )
	}
	tcv.assembler = aeron.NewFragmentAssembler(onMessage, 512)
}

func (tcv *AeronTransceiver) Close() {
	if tcv.pub != nil {
		tcv.pub.Close()
	}
	if tcv.sub != nil {
		tcv.sub.Close()
	}
	if tcv.aeron != nil {
		tcv.aeron.Close()
	}
}

func (tcv *AeronTransceiver) Reset() {
	tcv.rcvdMsg = 0
	tcv.histogram.Reset()
}

func (tcv *AeronTransceiver) SendAndReceive(ctx context.Context, msg []byte, iteration uint64, numMsg uint64) uint64 {
	var (
		startTimeNs      int64 = time.Now().UnixNano()
		endTimeNs        int64 = startTimeNs + int64(iteration)*NANOS_PER_SECOND
		sendIntervalNs   int64 = NANOS_PER_SECOND / int64(numMsg)
		nextSendTimeNs   int64 = startTimeNs
		nextReportTimeNs int64 = startTimeNs + NANOS_PER_SECOND
		nowNs            int64 = startTimeNs

		totalNumMsg uint64 = iteration * numMsg
		sentMsg     uint64 = 0
		batchSize   uint64 = config.BatchSize
	)

	outBuf := make([]byte, config.MaxMessageSize)
	idleStrategy := idlestrategy.Busy{}
	for {

		sent := tcv.SendWithNoRetry(atomic.MakeBuffer(outBuf), msg, batchSize)
		sentMsg += sent
		if totalNumMsg == sentMsg {
			reportProgress(startTimeNs, nowNs, sentMsg)
			break
		}

		nowNs = time.Now().UnixNano()
		if sent == batchSize {

			// next batch
			batchSize = util.Min(totalNumMsg-sentMsg, config.BatchSize)
			nextSendTimeNs += sendIntervalNs

			// receive batch size
			for nowNs < nextSendTimeNs && nowNs < endTimeNs {
				switch {
				case tcv.rcvdMsg < sentMsg:
					f := tcv.Receive()
					idleStrategy.Idle(f)

				default:
					// received batch already
					idleStrategy.Idle(0)
				}

				nowNs = time.Now().UnixNano()
			}

		} else {
			// next batch
			batchSize -= sent
			tcv.Receive()
		}

		if ctx.Err() != nil || nowNs >= endTimeNs {
			break
		}
		if nowNs >= nextReportTimeNs {
			elapsedSec := reportProgress(startTimeNs, nowNs, sentMsg)
			nextReportTimeNs = startTimeNs + int64((elapsedSec+1)*NANOS_PER_SECOND)
		}
	}

	// receive before exit
	deadline := time.Now().UnixNano() + RECEIVE_DEADLINE_NS
	for tcv.rcvdMsg < sentMsg {
		// log.Printf("[debug] try receive before exit: received: %d sent: %d\n", tcv.rcvdMsg, sentMsg)
		f := tcv.Receive()
		idleStrategy.Idle(f)
		if time.Now().UnixNano() >= deadline {
			log.Printf("[warn] not all messages were received after %ds deadline", RECEIVE_DEADLINE_NS/NANOS_PER_SECOND)
			break
		}
	}
	log.Printf("[info] messages sent: %v", sentMsg)
	return sentMsg
}

func (tcv *AeronTransceiver) SendWithNoRetry(outBuf *atomic.Buffer, msg []byte, num uint64) uint64 {
	sent := uint64(0)

	// timestamp
	outBuf.PutInt64(0, time.Now().UnixNano())

	// message
	msgLength := int32(len(msg))
	if msgLength > 0 {
		outBuf.PutBytesArray(8, &msg, 0, msgLength)
	}

	for i := uint64(0); i < num; i++ {
		res := tcv.pub.Offer(outBuf, 0, 8+msgLength, nil)
		if err := util.CheckPublicationResult(res); err != nil {
			log.Println("[error]", err)
			break
		}
		if res > 0 {
			sent += 1
		}
		// log.Printf("[debug] sent: %v size: %v", sent, len(msg))
	}
	return sent
}

func (tcv *AeronTransceiver) Receive() int {
	return tcv.sub.Poll(tcv.assembler.OnFragment, 10)
}
