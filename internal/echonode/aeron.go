package echonode

import (
	"context"
	"log"
)

type AeronEchoNode struct{}

func NewAeronEchoNode() *AeronEchoNode {
	return &AeronEchoNode{}
}

func (node *AeronEchoNode) Run(ctx context.Context) {
	log.Println("[info] running aeron echo node!")
}
