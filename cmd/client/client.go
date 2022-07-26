package client

import "github.com/spf13/cobra"

// Cmd is the ekko client command
var Cmd = &cobra.Command{
	Use:   "client",
	Short: "Send messages to echo server, measure RTT latencies in microseconds",
}

func init() {
	Cmd.AddCommand(AeronCmd)
	Cmd.AddCommand(GrpcCmd)
}
