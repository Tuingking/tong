package ddl

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tuingking/tong/pkg/gogen"
	"github.com/tuingking/tong/pkg/logger"
	"go.uber.org/zap"
)

const howToUse = `tong ddl -f [PATH_TO_GO_FILE]`

var (
	// flag
	file string
)

var Cmd = &cobra.Command{
	Use:   "ddl",
	Short: "Generate MySQL DDL query from go struct.",
	Long: fmt.Sprintf(`Generate MySQL DDL query from go struct.
	example: %s`, howToUse),
	RunE: runE,
}

func init() {
	Cmd.PersistentFlags().StringVarP(&file, "file", "f", "", "path to go file")
}

func runE(cmd *cobra.Command, args []string) error {
	if err := valid(); err != nil {
		return err
	}

	result, err := gogen.GenerateDDL(file, ".")
	if err != nil {
		logger.Logger.Error("failed to generate DDL", zap.Error(err))
		return err
	}
	fmt.Println(color.GreenString("Output:"))
	fmt.Println(result)

	return nil
}

func valid() error {
	if file == "" {
		return fmt.Errorf("please specify the go file path.\nUsage example:\n\t%s", howToUse)
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("directory `%s` not found", file)
	}

	return nil
}
