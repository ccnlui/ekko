package server

import (
	"github.com/spf13/cobra"
)

// Cmd is the ekko client command
var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Echo back every message to client",
}

func init() {
	Cmd.AddCommand(AeronCmd)
}
