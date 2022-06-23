package gcs

import (
	"github.com/Tuingking/tong/command/gcs/list"
	"github.com/Tuingking/tong/command/gcs/upload"
	"github.com/Tuingking/tong/config"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "gcs",
	Short: "Google Cloud Storage related utility command.",
	Long: `Google Cloud Storage related utility command.
	example: tong gcs`,
}

func NewCmd(cfg config.Config) *cobra.Command {
	cmd.AddCommand(upload.NewCmd(cfg))
	cmd.AddCommand(list.NewCmd(cfg))

	return cmd
}
