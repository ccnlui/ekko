package loadtestrig

import (
	"bytes"
	"context"
	"ekko/internal/transceiver"
	"fmt"
	"log"
	"strconv"
)

func Run(ctx context.Context, tcv transceiver.Transceiver) {
	// msg := []byte("hello")
	msg := generateMsg(11500)
	tcv.SendAndReceive(ctx, msg, 1, 10)
	log.Println("[info] Histogram of RTT latencies in microseconds.")
	fmt.Println("Bye!")
}

func generateMsg(n int) []byte {
	buf := bytes.Buffer{}
	for i := 0; i < n; i++ {
		buf.WriteString(strconv.FormatInt(int64(i), 10))
	}
	return buf.Bytes()
}
