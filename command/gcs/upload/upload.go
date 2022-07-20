package upload

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/tuingking/tong/config"
	"github.com/tuingking/tong/pkg/gcs"
	"github.com/tuingking/tong/pkg/logger"
	"go.uber.org/zap"
)

var (
	s      *spinner.Spinner
	client *storage.Client

	extAllowed = []string{".csv"}

	// progress
	percentProgress string

	// flag
	file   string // path to file which going to be uploaded
	bucket string // GCS bucket name
	dir    string // GCS directory
	all    bool   // set true if file is folder
)

var cmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload file to GCS.",
	Long: `Upload file to GCS.
	example: tong gcs upload -f [FILE]`,
	RunE: runE,
}

func NewCmd(cfg config.Config) *cobra.Command {
	// init gcs
	var err error
	client, err = gcs.NewClient(context.Background(), cfg.GCS)
	if err != nil {
		log.Fatal(err)
	}

	// init spinner
	s = spinner.New(spinner.CharSets[26], 300*time.Millisecond)
	s.FinalMSG = "upload done!\n"

	cmd.PersistentFlags().StringVarP(&file, "file", "f", "", "file or file inside folder which going to be uploaded")
	cmd.PersistentFlags().StringVarP(&bucket, "bucket", "b", "", "GCS bucket name")
	cmd.PersistentFlags().StringVar(&dir, "dir", "", "file directory")
	cmd.PersistentFlags().BoolVar(&all, "all", false, "file directory")
	cmd.MarkPersistentFlagRequired("file")
	cmd.MarkPersistentFlagRequired("bucket")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	s.Prefix = "preparing"
	s.Start()

	cwd, _ := os.Getwd()

	f, err := os.Open(path.Join(cwd, file))
	if err != nil {
		return err
	}

	baseDir := ""
	filesToUpload := []string{}

	fstat, _ := f.Stat()
	if fstat.IsDir() {
		if !all {
			return errors.New("this file is directory. please add flag `--all` to upload all files inside directory")
		}

		baseDir = fstat.Name()

		files, err := ioutil.ReadDir(file)
		if err != nil {
			return err
		}

		// add all files
		for _, f := range files {
			if isExtAllowed(f.Name()) {
				filesToUpload = append(filesToUpload, f.Name())
			}
		}

	} else {
		filesToUpload = append(filesToUpload, file)
	}

	fmt.Printf("total file: %d\n", len(filesToUpload))

	for i, v := range filesToUpload {

		percentProgress = fmt.Sprintf("[%.0f%%]", float64((i+1))/float64(len(filesToUpload))*100)

		f, err := os.Open(path.Join(baseDir, v))
		if err != nil {
			logger.Logger.Error("failed open file", zap.Error(err))
			continue
		}

		metadata := map[string]string{
			"Content-Type":        "text/csv",
			"Content-Disposition": fmt.Sprintf("attachment; filename=%s", file),
		}

		if err := Upload(context.Background(), bucket, f, getFilename(v), metadata); err != nil {
			return err
		}
	}

	s.Stop()

	return nil
}

func Upload(ctx context.Context, bucket string, file io.Reader, filename string, metadata map[string]string) error {
	s.Prefix = percentProgress + "uploading " + filename
	time.Sleep(1 * time.Second)

	w := client.Bucket(bucket).Object(path.Join(dir, filename)).NewWriter(ctx)
	w.ContentType = metadata["Content-Type"]
	w.ContentDisposition = metadata["Content-Disposition"]
	if _, err := io.Copy(w, file); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func isExtAllowed(filename string) bool {
	for _, ext := range extAllowed {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}

	return false
}

func getFilename(filename string) string {
	token := strings.Split(filename, "/")
	return token[len(token)-1]
}
