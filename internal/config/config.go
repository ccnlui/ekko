package config

import (
	"time"
)

var (
	QuoteMsgRate time.Duration = time.Second
	TradeMsgRate time.Duration = 2 * time.Second

	// Benchmark params
	WarmUpIterations  = 10
	WarmUpMessageRate = 20_000
	Iterations        = 10
	MessageRate       = 500_000
	MessageLength     = 128
	BatchSize         = 1
	MaxMessageSize    = 1028

	// Media driver
	AeronDir           string
	MediaDriverTimeout time.Duration = 10 * time.Second
	ServerStreamID     int           = 8000
	ServerChannel      string        = "aeron:udp?endpoint=localhost:40123"
	ClientStreamID     int           = 9000
	ClientChannel      string        = "aeron:udp?endpoint=localhost:40321"
)
