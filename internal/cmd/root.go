// Package cmd provides command line interface for the server application.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "server",
		Short: "Go template project server",
		Long:  `A template project demonstrating Go best practices with gRPC/Connect, Viper, and structured logging.`,
	}
)

// Execute 执行根命令.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/config.yaml", "config file path")
}
