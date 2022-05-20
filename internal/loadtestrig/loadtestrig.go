package loadtestrig

import (
	"context"
	"ekko/internal/transceiver"
	"log"
)

func Run(ctx context.Context, tcv transceiver.Transceiver) {
	msg := []byte("hello")
	tcv.SendAndReceive(msg, 10)
	log.Println("[info] Histogram of RTT latencies in microseconds.")
}
