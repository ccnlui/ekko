package client

import "github.com/spf13/cobra"

// Cmd is the ekko client command
var Cmd = &cobra.Command{
	Use:   "client",
	Short: "client that sends messages to server, and measure RTT latency in microseconds",
}
