package client

import (
	"ekko/internal/config"
	"ekko/internal/loadtestrig"

	"github.com/lirm/aeron-go/aeron"
	"github.com/spf13/cobra"
)

var AeronCmd = &cobra.Command{
	Use:   "aeron",
	Short: "aeron transport",
	Run: func(cmd *cobra.Command, args []string) {
		loadtestrig.Run(cmd.Context(), "aeron")
	},
}

func init() {
	AeronCmd.Flags().StringVar(&config.AeronDir, "aeron-dir", aeron.DefaultAeronDir,
		"directory name for aeron media driver",
	)
}
