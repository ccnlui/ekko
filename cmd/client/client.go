package client

import "github.com/spf13/cobra"

// Cmd is the echamber client command
var Cmd = &cobra.Command{
	Use:   "client",
	Short: "echo client",
}
