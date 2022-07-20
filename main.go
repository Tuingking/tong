package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/tuingking/tong/command/gcs"
	"github.com/tuingking/tong/command/gotong"
	"github.com/tuingking/tong/command/kafka"
	splitfile "github.com/tuingking/tong/command/split-file"
	"github.com/tuingking/tong/command/sql"
	"github.com/tuingking/tong/command/version"
	"github.com/tuingking/tong/config"
	"github.com/tuingking/tong/pkg/logger"
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
	setVersion(cfg.Version)

	rootCmd.AddCommand(version.NewCmd(cfg))
	rootCmd.AddCommand(gotong.NewCmd(cfg))
	rootCmd.AddCommand(kafka.NewCmd(cfg))
	rootCmd.AddCommand(sql.NewCmd(cfg))
	rootCmd.AddCommand(splitfile.NewCmd(cfg))
	rootCmd.AddCommand(gcs.NewCmd(cfg))
}

func main() {
	logger.Init(logger.DefaultOption())
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func setVersion(v string) {
	ver = v
}

// func printLogo() {
// 	fmt.Println(` o
// <|>
// < >
//  |        o__ __o    \o__ __o     o__ __o/
//  o__/_   /v     v\    |     |>   /v     |
//  |      />       <\  / \   / \  />     / \
//  |      \         /  \o/   \o/  \      \o/
//  o       o       o    |     |    o      |
//  <\__    <\__ __/>   / \   / \   <\__  < >
//                                         |
//                                 o__     o
//                                 <\__ __/>  `)
// }
