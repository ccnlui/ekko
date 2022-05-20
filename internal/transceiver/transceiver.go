package transceiver

import (
	"log"
)

type Transceiver interface {
	SendAndReceive(msg []byte, num int)
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
