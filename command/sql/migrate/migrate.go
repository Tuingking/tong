package migrate

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Tuingking/tong/config"
	"github.com/Tuingking/tong/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const howToUse = `tong sql migrate -f [PATH_TO_MIGRATION_FILE] --db [DATABASE_NAME]`

var (
	file     string
	username string
	password string
	hostport string
	db       string
)

var cmd = &cobra.Command{
	Use:   "migrate",
	Short: "Runs database migrations",
	Long:  fmt.Sprintf("Runs database migrations.\nUsage example:\n\t%s", howToUse),
	RunE:  runE,
}

func NewCmd(cfg config.Config) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&file, "file", "f", "", "migration file path (required)")
	cmd.PersistentFlags().StringVarP(&username, "username", "u", cfg.Mysql.Username, "MySQL username")
	cmd.PersistentFlags().StringVarP(&password, "password", "p", cfg.Mysql.Password, "MySQL password")
	cmd.PersistentFlags().StringVar(&hostport, "host", cfg.Mysql.HostPort, "MySQL host:port (default is localhost:3306)")
	cmd.PersistentFlags().StringVar(&db, "db", "", "MySQL database name (required)")

	cmd.MarkPersistentFlagRequired("file")
	cmd.MarkPersistentFlagRequired("db")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	if err := valid(); err != nil {
		return err
	}

	// debug()

	return migrate()
}

func migrate() error {
	q := `-source file://{file} -database mysql://{username}:{password}@tcp({hostport})/{database} -verbose up`

	replacer := strings.NewReplacer(
		"{file}", file,
		"{username}", username,
		"{password}", password,
		"{hostport}", hostport,
		"{database}", db,
	)

	// stdout := new(strings.Builder)
	stderr := new(strings.Builder)

	cmd := exec.Command("migrate", strings.Split(replacer.Replace(q), " ")...)
	// cmd.Stdout = stdout
	cmd.Stderr = stderr
	// output, _ := cmd.CombinedOutput()
	if err := cmd.Run(); err != nil {
		logger.Logger.Error("failed migrate", zap.Error(err))
		return err
	}

	log.Print(stderr)

	return nil
}

func valid() error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("👻👻👻 path `%s` not found", file)
	}

	return nil
}

func debug() {
	fmt.Println("file: ", file)
	fmt.Println("username: ", username)
	fmt.Println("password: ", password)
	fmt.Println("hostport: ", hostport)
	fmt.Println("database: ", db)
}
