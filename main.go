package main

import (
	"log"

	"github.com/Tuingking/tong/command/kafka"
	"github.com/Tuingking/tong/command/sql"
	"github.com/Tuingking/tong/command/version"
	"github.com/Tuingking/tong/config"
	"github.com/Tuingking/tong/pkg/logger"
	"github.com/spf13/cobra"
)

var ver = "v0.0.1"

var rootCmd = &cobra.Command{
	Use:     "tong",
	Short:   "command-line tool facilitating development of gotong-based application.",
	Long:    "command-line tool facilitating development of gotong-based application.",
	Version: ver,
}

func init() {
	cfg := config.Init()

	rootCmd.AddCommand(version.NewCmd(cfg))
	rootCmd.AddCommand(kafka.NewCmd(cfg))
	rootCmd.AddCommand(sql.NewCmd(cfg))
}

func main() {
	logger.Init(logger.DefaultOption())
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
