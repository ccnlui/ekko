package transceiver

import (
	"context"
	"ekko/internal/config"
	"ekko/internal/util"
	"ekko/proto"
	"io"
	"log"
	"sync/atomic"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcTransceiver struct {
	conn      *grpc.ClientConn
	client    proto.EkkoClient
	rcvdMsg   uint64
	histogram *hdrhistogram.Histogram
}

func NewGrpcTransceiver(histogram *hdrhistogram.Histogram) *GrpcTransceiver {
	log.Println("[info] new grpc transceiver!")
	tcv := &GrpcTransceiver{
		histogram: histogram,
	}
	tcv.init()
	return tcv
}

func (tcv *GrpcTransceiver) init() {
	dialCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(
		dialCtx,
		config.ServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.FailOnNonTempDialError(true),
	)
	if err != nil {
		log.Fatalln("[fatal] failed to connect to grpc server", err.Error())
	}

	tcv.conn = conn
	tcv.client = proto.NewEkkoClient(conn)
}

func (tcv *GrpcTransceiver) SendAndReceive(ctx context.Context, msg []byte, iteration uint64, numMsg uint64) uint64 {

	stream, err := tcv.client.BidirectionalStreamingEcho(ctx)
	if err != nil {
		log.Fatalln("[fatal] bidirectional streaming echo error: ", err.Error())
	}

	done := make(chan struct{})
	go tcv.Receive(done, stream)

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

	for {
		sent := tcv.SendWithNoRetry(stream, msg, batchSize)
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

			// spin until next batch
			for nowNs < nextSendTimeNs && nowNs < endTimeNs {
				nowNs = time.Now().UnixNano()
			}
		} else {
			// next batch
			batchSize -= sent
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
	for atomic.LoadUint64(&tcv.rcvdMsg) < sentMsg {
		// log.Printf("[debug] try receive before exit: received: %d sent: %d\n", tcv.rcvdMsg, sentMsg)
		if time.Now().UnixNano() >= deadline {
			log.Printf("[warn] not all messages were received after %ds deadline", RECEIVE_DEADLINE_NS/NANOS_PER_SECOND)
			break
		}
	}
	log.Printf("[info] messages sent: %v", sentMsg)

	stream.CloseSend()
	<-done

	return sentMsg
}

func (tcv *GrpcTransceiver) Receive(notifyDone chan struct{}, stream proto.Ekko_BidirectionalStreamingEchoClient) {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			close(notifyDone)
			return
		}
		if err != nil {
			log.Fatalln("[fatal] receive error: ", err.Error())
		}

		atomic.AddUint64(&tcv.rcvdMsg, 1)
		tcv.histogram.RecordValue(time.Now().UnixNano() - int64(resp.Timestamp))
		// log.Println("got: ", resp)
	}
}

func (tcv *GrpcTransceiver) SendWithNoRetry(stream proto.Ekko_BidirectionalStreamingEchoClient, msg []byte, num uint64) uint64 {
	sent := uint64(0)

	req := &proto.EchoRequest{
		Timestamp: uint64(time.Now().UnixNano()),
		Payload:   msg,
	}

	for i := uint64(0); i < num; i++ {
		if err := stream.Send(req); err != nil {
			log.Println("[error] send error: ", err.Error())
		} else {
			sent += 1
		}
	}

	return sent
}

func (tcv *GrpcTransceiver) Close() {
	tcv.conn.Close()
}

func (tcv *GrpcTransceiver) Reset() {
	atomic.StoreUint64(&tcv.rcvdMsg, 0)
	tcv.histogram.Reset()
}
