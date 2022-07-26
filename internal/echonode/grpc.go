package echonode

import (
	"context"
	"ekko/proto"
	"log"
)

type GrpcEchoNode struct {
	proto.UnimplementedEkkoServer
}

func NewGrpcEchoNode() *GrpcEchoNode {
	return &GrpcEchoNode{}
}

func (node *GrpcEchoNode) Run(ctx context.Context) {
	log.Println("[info] running grpc echo node!")
}

func (node *GrpcEchoNode) Close() {

}
