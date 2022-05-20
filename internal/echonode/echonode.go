package echonode

import (
	"context"
	"log"
)

type EchoNode interface {
	Run(ctx context.Context)
	Close()
}

func NewEchoNode(transport string) EchoNode {
	switch transport {
	case "aeron":
		return NewAeronEchoNode()
	default:
		log.Fatal("[fatal] unknown transport", transport)
	}
	return nil
}
