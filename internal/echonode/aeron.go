package echonode

import (
	"bytes"
	"context"
	"ekko/internal/config"
	"ekko/internal/util"
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
	piped := 0
	dropped := 0
	onMessage := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		// Don't create new bytes everytime. This is only an example
		// bytes := buffer.GetBytesArray(offset, length)

		inBuf.Reset()
		buffer.WriteBytes(inBuf, offset, length)
		// log.Printf("[info] %8.d: Got a fragment offset: %d length: %d payload: %s\n",
		// 	piped, offset, length,
		// 	string(inBuf.Next(int(length))),
		// )

		var res int64
		for {
			if res = node.pub.Offer(buffer, offset, length, nil); res > 0 {
				piped += 1
				// log.Printf("[debug] piped: %v size: %v", piped, length)
				break
			}
			if !util.RetryPublicationResult(res) {
				dropped += 1
				log.Println("[debug] dropped:", util.PublicationErrorString(res), string(inBuf.Next(int(length))))
				break
			}
		}
	}
	assembler := aeron.NewFragmentAssembler(onMessage, 512)

	idleStrategy := idlestrategy.Busy{}
	for {
		if ctx.Err() != nil {
			fmt.Println("bye!")
			return
		}

		fragmentsRead := node.sub.Poll(assembler.OnFragment, 10)
		idleStrategy.Idle(fragmentsRead)
	}
}
