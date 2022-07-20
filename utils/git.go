package utils

import (
	"os/exec"
	"strings"

	"github.com/tuingking/tong/pkg/logger"
	"go.uber.org/zap"
)

func GetLatestGitTag() string {
	var (
		stdout = new(strings.Builder)
		query  = "describe --tags --abbrev=0"
	)

	// fmt.Println("command: git", query)

	cmd := exec.Command("git", strings.Split(query, " ")...)
	cmd.Stdout = stdout
	if err := cmd.Run(); err != nil {
		if err.Error() == "exit status 128" {
			return "tag not found\n"
		}
		logger.Logger.Error("failed get latest git tag", zap.Error(err))
		panic(err)
	}

	return stdout.String()
}
