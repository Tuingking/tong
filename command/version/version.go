package version

import (
	"fmt"

	"github.com/Tuingking/tong/config"
	"github.com/spf13/cobra"
)

var ver string

var cmd = &cobra.Command{
	Use:   "version",
	Short: "show cli version",
	Long: `show cli version.
	example: tong --version`,
	Run: run,
}

func init() {}

func run(cmd *cobra.Command, args []string) {
	fmt.Printf("tong version: %s\n", ver)
}

func NewCmd(cfg config.Config) *cobra.Command {
	ver = cfg.Version
	return cmd
}
