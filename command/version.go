package command

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var Version = &cobra.Command{
	Use:   "version",
	Short: "show cli version",
	Long:  "show cli version",
	Run:   run,
}

func init() {

}

func run(cmd *cobra.Command, args []string) {
	getLatestGitTag()
}

func getLatestGitTag() {
	q := `git describe --tags --abbrev=0`

	cmd := exec.Command("git", q)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	fmt.Printf("🐶🐶[DEBUG]🐶🐶 cmd: %+v\n", cmd)
}
