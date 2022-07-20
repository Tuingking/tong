package gotong

import (
	"github.com/spf13/cobra"
	"github.com/tuingking/tong/config"
)

var (
	// flag
	appName string
)

var cmd = &cobra.Command{
	Use:   "app",
	Short: "Create new App.",
	Long: `Create new App.
	example: tong app --name [APP_NAME]`,
	RunE: runE,
}

func NewCmd(cfg config.Config) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&appName, "name", "n", "", "create gotong app")
	cmd.MarkPersistentFlagRequired("name")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	builder := NewGotongBuilder(appName)
	if err := builder.MakeBootstrap(); err != nil {
		return err
	}

	cdToProjectDir()
	initGoMod()
	goModTidy()

	return nil
}
