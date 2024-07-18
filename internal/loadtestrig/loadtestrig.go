package loadtestrig

import (
	"bytes"
	"context"
	"ekko/internal/config"
	"ekko/internal/transceiver"
	"fmt"
	"log"
	"os"

	"github.com/HdrHistogram/hdrhistogram-go"
)

func Run(ctx context.Context, transport string) {
	histogram := hdrhistogram.New(1, 60_000_000_000, 3)
	tcv := transceiver.NewTransceiver(transport, histogram)
	defer tcv.Close()

	msg := generateMsg(config.MessageLength)

	log.Printf("[info] Running warmup for %d iterations of %d messages each, with %d bytes payload and a burst size of %d...\n",
		config.WarmUpIterations,
		config.WarmUpMessageRate,
		config.MessageLength,
		config.BatchSize,
	)
	tcv.SendAndReceive(ctx, msg, config.WarmUpIterations, config.WarmUpMessageRate)
	tcv.Reset()

	log.Printf("[info] Running measurement for %d iterations of %d messages each, with %d bytes payload and a burst size of %d...\n",
		config.Iterations,
		config.MessageRate,
		config.MessageLength,
		config.BatchSize,
	)
	sentMessages := tcv.SendAndReceive(ctx, msg, config.Iterations, config.MessageRate)

	log.Println("[info] Histogram of RTT latencies in microseconds.")
	histogram.PercentilesPrint(os.Stdout, 5, 1000.0)

	expectedMessages := config.Iterations * config.MessageRate
	if sentMessages < expectedMessages {
		msg := "[warn] Target message rate not achieved: expected to send %d messages in total but managed to send only %d messages\n"
		log.Printf(msg, expectedMessages, sentMessages)
	}

	fmt.Println("Bye!")
}

func generateMsg(msgLength uint64) []byte {
	buf := bytes.Buffer{}
	for i := uint64(0); i < msgLength; i++ {
		buf.WriteByte(0x42)
	}
	return buf.Bytes()
}
