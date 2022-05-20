package transceiver

import (
	"context"
	"log"
)

type Transceiver interface {
	SendAndReceive(ctx context.Context, msg []byte, num int)
	Close()
}

func NewTransceiver(transport string) Transceiver {
	switch transport {
	case "aeron":
		return NewAeronTransceiver()

	default:
		log.Fatal("[fatal] unknown transport", transport)
	}
	return nil
}
