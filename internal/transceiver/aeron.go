package transceiver

import (
	"context"
	"ekko/internal/config"
	"fmt"
	"log"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
)

type AeronTransceiver struct {
	quoteMsgRate time.Duration
	tradeMsgRate time.Duration

	quoteDelayer delayer
	tradeDelayer delayer

	aeron *aeron.Aeron
	pub   *aeron.Publication
}

func NewAeronTransceiver() *AeronTransceiver {
	log.Println("[info] new aeron transport!")

	qmr := config.QuoteMsgRate
	tmr := config.TradeMsgRate
	qd := delayer{
		interval: qmr.Nanoseconds(),
	}
	td := delayer{
		interval: tmr.Nanoseconds(),
	}
	tcv := &AeronTransceiver{
		quoteMsgRate: qmr,
		tradeMsgRate: tmr,
		quoteDelayer: qd,
		tradeDelayer: td,
	}
	tcv.init()
	return tcv
}

func (tcv *AeronTransceiver) init() {
	aeronCtx := aeron.NewContext().AeronDir(config.AeronDir)
	a, err := aeron.Connect(aeronCtx)
	if err != nil {
		log.Fatalln("[fatal] failed to connect to media driver: ", config.AeronDir, err.Error())
	}
	tcv.aeron = a
	tcv.pub = <-a.AddPublication(config.Channel, int32(config.StreamID))
	for !tcv.pub.IsConnected() {
		time.Sleep(time.Millisecond)
	}
	log.Println("[info] publication connected to media driver:", tcv.pub)
}

func (tcv *AeronTransceiver) Close() {
	if tcv.aeron != nil {
		tcv.aeron.Close()
	}
}

func (tcv *AeronTransceiver) SendAndReceive(ctx context.Context, msg []byte, num int) {

	sent := 0
	dropped := 0
	for sent < num {

		if ctx.Err() != nil {
			fmt.Println("Bye!")
			return
		}

		if tcv.pub.IsConnected() {
			now := time.Now().UnixNano()
			if tcv.quoteDelayer.onScheduleSend(now) {

				// msg := time.Now().Local().String()
				outBuf := atomic.MakeBuffer([]byte(msg))

				var res int64
				for {
					if res = tcv.pub.Offer(outBuf, 0, int32(len(msg)), nil); res > 0 {
						sent += 1
						log.Printf("[debug] sent: %v size: %v", sent, len(msg))
						break
					}
					if !retryPublicationResult(res) {
						dropped += 1
						log.Println("[debug] dropped:", publicationErrorString(res), msg)
						break
					}
				}
			}
		}
	}
	log.Printf("[info] messages sent: %v dropped: %v", sent, dropped)
}

func retryPublicationResult(res int64) bool {
	switch res {
	case aeron.AdminAction, aeron.BackPressured:
		// log.Println("[debug] retry offer:", publicationErrorString(res))
		return true
	case aeron.NotConnected, aeron.MaxPositionExceeded, aeron.PublicationClosed:
		log.Println("[error] failed to offer", publicationErrorString(res))
		return false
	}
	return false
}

func publicationErrorString(res int64) string {
	switch res {
	case aeron.AdminAction:
		return "ADMIN_ACTION"
	case aeron.BackPressured:
		return "BACK_PRESSURED"
	case aeron.PublicationClosed:
		return "CLOSED"
	case aeron.NotConnected:
		return "NOT_CONNECTED"
	case aeron.MaxPositionExceeded:
		return "MAX_POSITION_EXCEEDED"
	default:
		return "UNKNOWN"
	}
}
