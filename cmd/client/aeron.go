package client

import (
	"ekko/internal/transceiver"

	"github.com/spf13/cobra"
)

var AeronCmd = &cobra.Command{
	Use:   "aeron",
	Short: "aeron transport",
	RunE: func(cmd *cobra.Command, args []string) error {
		transceiver.NewTransceiver("aeron")
		return nil
	},
}
