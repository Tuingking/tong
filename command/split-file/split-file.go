package splitfile

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Tuingking/tong/config"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var (
	// flag
	file   string
	maxRow int
)

var cmd = &cobra.Command{
	Use:   "split-file",
	Short: "Split huge csv file into multiple files.",
	Long: `Split huge csv file into multiple files.
	example: tong split-file -f [FILE]`,
	RunE: runE,
}

func NewCmd(cfg config.Config) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&file, "file", "f", "", "path to file")
	cmd.PersistentFlags().IntVarP(&maxRow, "max", "m", 10_000, "max row per file")
	cmd.MarkPersistentFlagRequired("file")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	return splitCsvFile()
}

func splitCsvFile() error {
	base := getWorkingDir()

	s := spinner.New(spinner.CharSets[26], 300*time.Millisecond)
	s.Color("white", "bold")
	s.Prefix = "preparing"
	s.Start()

	f, err := os.Open(file)
	if err != nil {
		return err
	}

	stat, _ := f.Stat()
	filename := stat.Name()

	var ext string

	reader := csv.NewReader(f)
	if strings.HasSuffix(filename, ".tsv") {
		reader.Comma = '\t'
		ext = ".tsv"
	} else if strings.HasSuffix(filename, ".csv") {
		reader.Comma = ','
		ext = ".csv"
	} else {
		return errors.New("invalid file format")
	}

	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	header := rows[0]
	data := rows[1:]

	batchSize := maxRow
	totalBatch := len(data) / batchSize
	if len(data)%batchSize > 0 {
		totalBatch++
	}

	s.Prefix = fmt.Sprintf("processing %d rows", len(data))
	time.Sleep(2 * time.Second)

	for i := 0; i < totalBatch; i++ {
		batchNo := i
		start := batchNo * batchSize
		end := start + batchSize
		if end >= len(data)-1 {
			end = len(data)
		}

		s.Prefix = fmt.Sprintf("creating file #%d", batchNo)

		records := data[start:end]
		filename := strings.ReplaceAll(filename, ext, fmt.Sprintf("_%d%s", batchNo, ext))
		dest := path.Join(base, filename)
		if err := createFile(header, records, dest); err != nil {
			return err
		}

		time.Sleep(300 * time.Millisecond)
	}

	s.Stop()
	fmt.Printf("total rows: %d\n", len(data))
	fmt.Printf("total file: %d\n", totalBatch)
	fmt.Printf("location: %s\n", base)

	return nil
}

func createFile(header []string, records [][]string, filename string) error {
	f, _ := os.Create(filename)
	writter := csv.NewWriter(f)
	defer writter.Flush()

	writter.Write(header)
	return writter.WriteAll(records)
}

func getWorkingDir() string {
	wd, _ := os.Getwd()
	return wd
}
