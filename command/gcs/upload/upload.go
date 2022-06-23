package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"

	"cloud.google.com/go/storage"
	"github.com/Tuingking/tong/config"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var (
	// flag
	file   string // path to file which going to be uploaded
	bucket string // GCS bucket name
	dir    string // directory
)

var cmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload file to GCS.",
	Long: `Upload file to GCS.
	example: tong gcs upload -f [FILE]`,
	RunE: runE,
}

func NewCmd(cfg config.Config) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&file, "file", "f", "", "path to file which going to be uploaded")
	cmd.PersistentFlags().StringVarP(&bucket, "bucket", "b", "", "GCS bucket name")
	cmd.PersistentFlags().StringVar(&dir, "dir", "", "file directory")
	cmd.MarkPersistentFlagRequired("file")
	cmd.MarkPersistentFlagRequired("bucket")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	cwd, _ := os.Getwd()

	f, err := os.Open(path.Join(cwd, file))
	if err != nil {
		return err
	}

	metadata := map[string]string{
		"Content-Type":        "text/csv",
		"Content-Disposition": fmt.Sprintf("attachment; filename=%s", file),
	}

	return Upload(context.Background(), bucket, f, file, metadata)
}

func Upload(ctx context.Context, bucket string, file io.Reader, filename string, metadata map[string]string) error {
	s := spinner.New(spinner.CharSets[26], 300*time.Millisecond)
	s.FinalMSG = "upload done!\n"
	s.Color("white", "bold")
	s.Start()
	s.Prefix = "preparing"

	opts := []option.ClientOption{option.WithCredentialsFile(getConfigFile())}
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return err
	}

	s.Prefix = "uploading"

	w := client.Bucket(bucket).Object(path.Join(dir, filename)).NewWriter(ctx)
	w.ContentType = metadata["Content-Type"]
	w.ContentDisposition = metadata["Content-Disposition"]
	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	s.Stop()

	return nil
}

func getConfigFile() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), "..", "..", "..", "config", "xwork-dev-5e7544287095.json")
	// return path.Join(path.Dir(filename), "..", "..", "..", "config", "tk-dev-micro-71a89e645c42.json")
}
