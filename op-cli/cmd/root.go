package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "op",
	Short: "op-cli is a CLI tool to interact with the superchain",
}

func Execute() {
	rootCmd.PersistentFlags().String("fmt", "table", "define whether the output should be json or table")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
