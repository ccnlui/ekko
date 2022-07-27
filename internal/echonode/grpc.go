package echonode

import (
	"context"
	"ekko/internal/config"
	"ekko/proto"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type GrpcEchoNode struct {
	proto.UnimplementedEkkoServer
}

func NewGrpcEchoNode() *GrpcEchoNode {
	node := &GrpcEchoNode{}
	return node
}

func (node *GrpcEchoNode) Run(_ context.Context) {
	lis, err := net.Listen("tcp", config.ServerAddr)
	if err != nil {
		log.Fatalln("[fatal] failed to listen on address: ", config.ServerAddr, err.Error())
	}

	s := grpc.NewServer()
	proto.RegisterEkkoServer(s, node)

	log.Println("[info] running grpc echo node!")
	err = s.Serve(lis)
	if err != nil {
		log.Fatalln("[fatal] failed to start GRPC server", err.Error())
	}
}

func (node *GrpcEchoNode) Close() {

}

func (node *GrpcEchoNode) UnaryEcho(ctx context.Context, req *proto.EchoRequest) (*proto.EchoResponse, error) {
	resp := &proto.EchoResponse{
		Timestamp: req.Timestamp,
		Payload:   req.Payload,
	}
	return resp, nil
}

func (node *GrpcEchoNode) ServerStreamingEcho(req *proto.EchoRequest, stream proto.Ekko_ServerStreamingEchoServer) error {
	for {
		resp := &proto.EchoResponse{
			Timestamp: uint64(time.Now().UnixNano()),
			Payload:   req.Payload,
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}

func (node *GrpcEchoNode) ClientStreamingEcho(stream proto.Ekko_ClientStreamingEchoServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.EchoResponse{
				Timestamp: uint64(time.Now().UnixNano()),
				Payload:   req.Payload,
			})
		}
		if err != nil {
			return err
		}
	}
}

func (node *GrpcEchoNode) BidirectionalStreamingEcho(stream proto.Ekko_BidirectionalStreamingEchoServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		resp := &proto.EchoResponse{
			Timestamp: req.Timestamp,
			Payload:   req.Payload,
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}
