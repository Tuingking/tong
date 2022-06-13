package main

import (
	"log"

	"github.com/Tuingking/igo-cli/command"
	"github.com/spf13/cobra"
)

const version = "v0.0.0"

var rootCmd = &cobra.Command{
	Use:     "igo",
	Short:   "",
	Long:    "",
	Version: version,
}

func init() {
	// TODO: add your command here.
	rootCmd.AddCommand(command.Version)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
