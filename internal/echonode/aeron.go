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
	pub   *aeron.Publication
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

	node.sub = <-a.AddSubscription(config.ServerChannel, int32(config.ServerStreamID))
	for !node.sub.IsConnected() {
		time.Sleep(time.Millisecond)
	}
	log.Println("[info] subscription connected to media driver:", node.sub)

	node.pub = <-a.AddPublication(config.ClientChannel, int32(config.ClientStreamID))
	for !node.pub.IsConnected() {
		time.Sleep(time.Millisecond)
	}
	log.Println("[info] publication connected to media driver:", node.pub)
}

func (node *AeronEchoNode) Close() {
	if node.sub != nil {
		node.sub.Close()
	}
	if node.pub != nil {
		node.pub.Close()
	}
	if node.aeron != nil {
		node.aeron.Close()
	}
}

func (node *AeronEchoNode) Run(ctx context.Context) {
	log.Println("[info] running aeron echo node!")

	inBuf := &bytes.Buffer{}
	count := 1
	onMessage := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		// Don't create new bytes everytime. This is only an example
		// bytes := buffer.GetBytesArray(offset, length)

		inBuf.Reset()
		buffer.WriteBytes(inBuf, offset, length)
		log.Printf("[info] %8.d: Got a fragment offset: %d length: %d payload: %s\n",
			count, offset, length,
			string(inBuf.Next(int(length))),
		)

		count += 1
	}
	assembler := aeron.NewFragmentAssembler(onMessage, 512)

	idleStrategy := idlestrategy.Sleeping{SleepFor: time.Millisecond}
	for {
		if ctx.Err() != nil {
			fmt.Println("bye!")
			return
		}

		fragmentsRead := node.sub.Poll(assembler.OnFragment, 10)
		idleStrategy.Idle(fragmentsRead)
	}
}
