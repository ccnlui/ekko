package server

import "github.com/spf13/cobra"

// Cmd is the ekko client command
var Cmd = &cobra.Command{
	Use:   "server",
	Short: "server that echoes back every message",
}
