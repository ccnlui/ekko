package echonode

import (
	"bytes"
	"context"
	"ekko/internal/config"
	"fmt"
	"log"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
)

type AeronEchoNode struct {
	aeron *aeron.Aeron
	sub   *aeron.Subscription
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
	node.sub = <-a.AddSubscription(config.Channel, int32(config.StreamID))
	log.Println("[info] subscription:", node.sub)
}

func (node *AeronEchoNode) Close() {
	if node.aeron != nil {
		node.aeron.Close()
	}
	if node.sub != nil {
		node.sub.Close()
	}
}

func (node *AeronEchoNode) Run(ctx context.Context) {
	log.Println("[info] running aeron echo node!")

	inBuf := &bytes.Buffer{}
	counter := 0
	handler := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		bytes := buffer.GetBytesArray(offset, length)
		inBuf.Reset()
		buffer.WriteBytes(inBuf, offset, length)
		fmt.Printf("%8.d: Got a fragment offset: %d length: %d payload: %s (buf:%s)\n",
			counter, offset, length,
			string(bytes),
			string(inBuf.Next(int(length))),
		)
		counter += 1
	}

	idleStrategy := idlestrategy.Sleeping{SleepFor: time.Millisecond}
	for {
		fragmentsRead := node.sub.Poll(handler, 10)
		idleStrategy.Idle(fragmentsRead)
	}
}
