package cmd

import (
	"context"
	"ekko/cmd/client"
	"ekko/cmd/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	envPrefix = "ekko"
)

// global flags
var (
	configPath string
)

func Execute() {
	cmd := newRootCmd()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-signals
		cancel()
		time.Sleep(time.Second)
		log.Println("[warn] not all commands finished completely")
		os.Exit(1)
	}()

	if err := cmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ekko",
		Short: "Marketdata message transport RTT benchmark tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := initConfig(cmd); err != nil {
				return err
			}
			configFileUsed := viper.ConfigFileUsed()
			if configFileUsed != "" {
				log.Println("[info] using config file: " + configFileUsed)
			}
			return nil
		},
	}

	rootCmd.AddCommand(client.Cmd)
	rootCmd.AddCommand(server.Cmd)
	return rootCmd
}

func initConfig(cmd *cobra.Command) error {
	// config from file
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
		viper.SetConfigName(".ekko")
		viper.SetConfigType("yaml")

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				log.Fatalln("[fatal]", err.Error())
			}
		}
	}

	// config from env
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	return bindFlags(cmd)
}

// bindFlags sets the value of all unset cobra flags from viper.
// This way the user can set the flags in this priority order:
// command line flag > environment variable > config > default
func bindFlags(cmd *cobra.Command) error {
	var returnErr error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// environment variables cannot contain dashes
		name := strings.ReplaceAll(f.Name, "-", "_")
		if err := viper.BindPFlag(name, f); err != nil {
			returnErr = err
			return
		}
		if !f.Changed && viper.IsSet(name) {
			if err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", viper.Get(name))); err != nil {
				returnErr = err
			}
		}
	})
	return returnErr
}
