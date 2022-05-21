package transceiver

import (
	"context"
	"log"

	"github.com/HdrHistogram/hdrhistogram-go"
)

type Transceiver interface {
	SendAndReceive(ctx context.Context, msg []byte, iteration int, numMsg int) int
	Close()
	Reset()
}

func NewTransceiver(transport string, histogram *hdrhistogram.Histogram) Transceiver {
	switch transport {
	case "aeron":
		return NewAeronTransceiver(histogram)

	default:
		log.Fatal("[fatal] unknown transport", transport)
	}
	return nil
}
