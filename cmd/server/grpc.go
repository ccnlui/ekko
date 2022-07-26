package server

import (
	"ekko/internal/echonode"

	"github.com/spf13/cobra"
)

var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "grpc transport",
	Run: func(cmd *cobra.Command, args []string) {
		n := echonode.NewEchoNode("grpc")
		defer n.Close()
		n.Run(cmd.Context())
	},
}
