package gotong

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"text/template"
)

func getPackageName(appName string) string {
	return "github.com/tuingking/" + appName
}

func cdToProjectDir() {
	err := os.Chdir(appName)
	if err != nil {
		panic(err)
	}

	log.Printf("cd to %s\n", appName)
}

func initGoMod() {
	log.Println("go mod init")

	q := "mod init " + getPackageName(appName)

	var stderr = new(strings.Builder)

	cmd := exec.Command("go", strings.Split(q, " ")...)
	cmd.Stderr = stderr
	cmd.Run()

	log.Printf("stderr: %s\n", stderr)
}

func goModTidy() {
	log.Println("go mod tidy")

	q := "mod tidy"

	var stderr = new(strings.Builder)

	cmd := exec.Command("go", strings.Split(q, " ")...)
	cmd.Stderr = stderr
	cmd.Run()

	log.Printf("stderr: %s\n", stderr)
}

// root of this repository
func getBasePath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), "..", "..")
}

func copyTemplate(dest string, src string, data map[string]string) error {
	dst, err := os.Create(fmt.Sprintf("%s/%s", appName, dest))
	if err != nil {
		return err
	}
	defer dst.Close()

	tpl := template.Must(template.ParseFiles(path.Join(getBasePath(), "/template/gotong/"+src)))
	return tpl.Execute(dst, data)
}
