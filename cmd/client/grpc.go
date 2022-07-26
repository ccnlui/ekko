package client

import (
	"ekko/internal/loadtestrig"

	"github.com/spf13/cobra"
)

var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "grpc transport",
	Run: func(cmd *cobra.Command, args []string) {
		loadtestrig.Run(cmd.Context(), "grpc")
	},
}
