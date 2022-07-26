package transceiver

import "log"

func reportProgress(startTimeNs int64, nowNs int64, sentMsg uint64) uint32 {
	elapsedSec := uint32((nowNs - startTimeNs)) / NANOS_PER_SECOND
	var sendRate uint32
	switch elapsedSec {
	case 0:
		sendRate = uint32(sentMsg)
	default:
		sendRate = uint32(sentMsg) / elapsedSec
	}
	log.Printf("[info] send rate %d msg/sec\n", sendRate)
	return elapsedSec
}
