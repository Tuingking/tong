package gcs

import (
	"github.com/spf13/cobra"
	"github.com/tuingking/tong/command/gcs/delete"
	"github.com/tuingking/tong/command/gcs/list"
	"github.com/tuingking/tong/command/gcs/upload"
	"github.com/tuingking/tong/config"
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
	cmd.AddCommand(delete.NewCmd(cfg))

	return cmd
}
