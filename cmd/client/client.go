package client

import "github.com/spf13/cobra"

// Cmd is the ekko client command
var Cmd = &cobra.Command{
	Use:   "client",
	Short: "Send messages to server, and measure RTT latency in microseconds",
}

func init() {
	Cmd.AddCommand(AeronCmd)
}
