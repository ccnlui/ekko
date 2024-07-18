package config

import (
	"time"
)

var (
	// Benchmark params
	WarmUpIterations  uint64 = 10
	WarmUpMessageRate uint64 = 20_000
	Iterations        uint64 = 10
	MessageRate       uint64 = 500_000
	MessageLength     uint64 = 128
	BatchSize         uint64 = 1
	MaxMessageSize    uint64 = 1028

	// Media driver
	AeronDir           string
	MediaDriverTimeout time.Duration = 10 * time.Second
	ServerStreamID     int           = 8000
	ServerChannel      string        = "aeron:ipc"
	ClientStreamID     int           = 9000
	ClientChannel      string        = "aeron:ipc"
	// ServerChannel      string        = "aeron:udp?endpoint=localhost:40123"
	// ClientChannel      string        = "aeron:udp?endpoint=localhost:40321"

	// Grpc
	ServerAddr string = "localhost:9090"
)
