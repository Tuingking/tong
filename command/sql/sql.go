package sql

import (
	"github.com/Tuingking/tong/command/sql/ddl"
	findfield "github.com/Tuingking/tong/command/sql/find-field"
	"github.com/Tuingking/tong/command/sql/migrate"
	"github.com/Tuingking/tong/config"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "sql",
	Short: "MySQL related utility command.",
	Long: `MySQL related utility command.
	example: tong sql`,
}

func NewCmd(cfg config.Config) *cobra.Command {
	cmd.AddCommand(migrate.NewCmd(cfg))
	cmd.AddCommand(ddl.Cmd)
	cmd.AddCommand(findfield.Cmd)

	return cmd
}
