package loadtestrig

import (
	"bytes"
	"context"
	"ekko/internal/transceiver"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/HdrHistogram/hdrhistogram-go"
)

func Run(ctx context.Context, transport string) {
	histogram := hdrhistogram.New(1, 60_000_000_000, 3)
	t := transceiver.NewTransceiver(transport, histogram)
	defer t.Close()

	msg := generateMsg(64)
	t.SendAndReceive(ctx, msg, 5, 500_000)
	log.Println("[info] Histogram of RTT latencies in microseconds.")
	histogram.PercentilesPrint(os.Stdout, 5, 1000.0)
	fmt.Println("Bye!")
}

func generateMsg(n int) []byte {
	buf := bytes.Buffer{}
	for i := 0; i < n; i++ {
		buf.WriteString(strconv.FormatInt(int64(i), 10))
	}
	return buf.Bytes()
}
