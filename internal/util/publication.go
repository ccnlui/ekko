package util

import (
	"fmt"
	"log"

	"github.com/lirm/aeron-go/aeron"
)

func RetryPublicationResult(res int64) bool {
	switch res {
	case aeron.AdminAction, aeron.BackPressured:
		// log.Println("[debug] retry offer:", publicationErrorString(res))
		return true
	case aeron.NotConnected, aeron.MaxPositionExceeded, aeron.PublicationClosed:
		log.Println("[error] failed to offer", PublicationErrorString(res))
		return false
	}
	return false
}

func CheckPublicationResult(res int64) error {
	switch res {
	case aeron.PublicationClosed, aeron.NotConnected, aeron.MaxPositionExceeded:
		return fmt.Errorf("publication error: %v", PublicationErrorString(res))
	default:
		return nil
	}
}

func PublicationErrorString(res int64) string {
	switch res {
	case aeron.AdminAction:
		return "ADMIN_ACTION"
	case aeron.BackPressured:
		return "BACK_PRESSURED"
	case aeron.PublicationClosed:
		return "CLOSED"
	case aeron.NotConnected:
		return "NOT_CONNECTED"
	case aeron.MaxPositionExceeded:
		return "MAX_POSITION_EXCEEDED"
	default:
		return "UNKNOWN"
	}
}
