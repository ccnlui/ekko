package transceiver

import "log"

func NewTransceiver(transport string) {
	switch transport {
	case "aeron":
		log.Println("[info] new aeron transport!")
	default:
		log.Println("[error] unknown transport", transport)
	}
}
