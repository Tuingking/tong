package kafka

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/tuingking/tong/config"
)

var (
	// default config
	zookeeperConfigFile = "$HOME/kafka/config/zookeeper.properties"
	brokerConfigFile    = "$HOME/kafka/config/server.properties"
)

var (
	// flag
	clear bool
)

var cmd = &cobra.Command{
	Use:   "kafka",
	Short: "Start kafka zookeeper and broker.",
	Long: `Start kafka zookeeper and broker.
	example: tong kafka`,
	Run: run,
}

func NewCmd(cfg config.Config) *cobra.Command {
	overrideConfigFile(cfg)

	cmd.PersistentFlags().BoolVarP(&clear, "clear", "c", false, "clear tmp")

	return cmd
}

func run(cmd *cobra.Command, args []string) {
	clear, _ := cmd.Flags().GetBool("clear")
	if clear {
		clearTmp()
	}

	log.Printf("Config File:")
	log.Printf("%-9s: %s\n", "Zookeeper", zookeeperConfigFile)
	log.Printf("%-9s: %s\n", "Broker", brokerConfigFile)

	s := spinner.New(spinner.CharSets[26], 300*time.Millisecond)
	s.FinalMSG = "Kafka zookeeper and broker is now running\n"
	s.Color("white", "bold")
	s.Start()

	runZookeeper(s)
	runBrokerServer(s)

	s.Stop()
}

func runZookeeper(s *spinner.Spinner) {
	s.Prefix = "running zookeeper"

	script := fmt.Sprintf(`tell application "Terminal" to do script "zookeeper-server-start.sh %s"`, zookeeperConfigFile)

	bashCmd := exec.Command("osascript", "-s", "h", "-e", script)
	stderr, err := bashCmd.StderrPipe()
	log.SetOutput(os.Stderr)
	if err != nil {
		log.Fatal(err)
	}

	if err := bashCmd.Start(); err != nil {
		log.Fatal(err)
	}

	slurp, _ := ioutil.ReadAll(stderr)
	if string(slurp) != "" {
		log.Printf("slurp: %s\n", string(slurp))
	}

	if err := bashCmd.Wait(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(7 * time.Second)
}

func runBrokerServer(s *spinner.Spinner) {
	s.Prefix = "running broker"

	script := fmt.Sprintf(`tell application "Terminal" to do script "kafka-server-start.sh %s"`, brokerConfigFile)

	bashCmd := exec.Command("osascript", "-s", "h", "-e", script)
	stderr, err := bashCmd.StderrPipe()
	log.SetOutput(os.Stderr)
	if err != nil {
		log.Fatal(err)
	}

	if err := bashCmd.Start(); err != nil {
		log.Fatal(err)
	}

	slurp, _ := ioutil.ReadAll(stderr)
	if string(slurp) != "" {
		log.Printf("slurp: %s\n", string(slurp))
	}

	if err := bashCmd.Wait(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(4 * time.Second)
}

func clearTmp() {
	bashCmd := exec.Command("rm", "-rf", "/tmp/kafka-logs", "/tmp/zookeeper")
	err := bashCmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Local data deleted")
}

func overrideConfigFile(cfg config.Config) {
	if cfg.Kafka.ZookeeperConfigFile != "" {
		zookeeperConfigFile = cfg.Kafka.ZookeeperConfigFile
	}

	if cfg.Kafka.BrokerConfigFile != "" {
		brokerConfigFile = cfg.Kafka.BrokerConfigFile
	}
}
