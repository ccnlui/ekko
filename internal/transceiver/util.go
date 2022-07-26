package transceiver

import "log"

func reportProgress(startTimeNs int64, nowNs int64, sentMsg uint64) uint64 {
	elapsedSec := uint64((nowNs - startTimeNs)) / NANOS_PER_SECOND
	var sendRate uint64
	switch elapsedSec {
	case 0:
		sendRate = sentMsg
	default:
		sendRate = sentMsg / elapsedSec
	}
	log.Printf("[info] send rate %d msg/sec\n", sendRate)
	return elapsedSec
}
