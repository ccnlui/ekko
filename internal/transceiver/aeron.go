package transceiver

import (
	"ekko/internal/config"
	"log"
	"time"

	"github.com/lirm/aeron-go/aeron"
)

type AeronTransceiver struct {
	quoteMsgRate time.Duration
	tradeMsgRate time.Duration

	quoteDelayer delayer
	tradeDelayer delayer

	aeron *aeron.Aeron
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
	return &AeronTransceiver{
		quoteMsgRate: qmr,
		tradeMsgRate: tmr,
		quoteDelayer: qd,
		tradeDelayer: td,
	}
}

func (tcv *AeronTransceiver) Init() {
	aeronCtx := aeron.NewContext().AeronDir(config.AeronDir)
	a, err := aeron.Connect(aeronCtx)
	if err != nil {
		log.Fatalln("[fatal] failed to connect to media driver: ", config.AeronDir, err.Error())
	}
	tcv.aeron = a
}

func (tcv *AeronTransceiver) Close() {
	if tcv.aeron != nil {
		tcv.aeron.Close()
	}
}

func (tcv *AeronTransceiver) SendAndReceive(msg []byte, num int) {

	count := 0
	for count < num {

		now := time.Now().UnixNano()
		if tcv.quoteDelayer.onScheduleSend(now) {
			log.Println("send:", string(msg))
			count += 1
		}

	}
}
