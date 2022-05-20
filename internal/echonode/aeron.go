package echonode

import (
	"context"
	"ekko/internal/config"
	"log"

	"github.com/lirm/aeron-go/aeron"
)

type AeronEchoNode struct {
	aeron *aeron.Aeron
}

func NewAeronEchoNode() *AeronEchoNode {
	node := &AeronEchoNode{}
	node.init()
	return node
}

func (node *AeronEchoNode) init() {
	aeronCtx := aeron.NewContext().AeronDir(config.AeronDir)
	a, err := aeron.Connect(aeronCtx)
	if err != nil {
		log.Fatalln("[fatal] failed to connect to media driver: ", config.AeronDir, err.Error())
	}
	node.aeron = a
}

func (node *AeronEchoNode) Close() {
	if node.aeron != nil {
		node.aeron.Close()
	}
}

func (node *AeronEchoNode) Run(ctx context.Context) {

	log.Println("[info] running aeron echo node!")
}
