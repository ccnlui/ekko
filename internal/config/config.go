package config

import (
	"time"
)

var (
	QuoteMsgRate time.Duration = time.Second
	TradeMsgRate time.Duration = 2 * time.Second

	// Media driver
	AeronDir           string
	MediaDriverTimeout time.Duration = 10 * time.Second
	StreamID           int           = 42
	Channel            string        = "aeron:udp?endpoint=localhost:40123"
)
