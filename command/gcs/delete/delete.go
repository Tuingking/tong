package delete

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"github.com/spf13/cobra"
	"github.com/tuingking/tong/config"
	"github.com/tuingking/tong/pkg/gcs"
)

var (
	client *storage.Client

	// flag
	bucket   string // GCS bucket name
	filename string
)

var cmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete file in GCS bucket.",
	Long: `Delete file in GCS bucket.
	example: tong gcs delete -b [BUCKET_NAME] -f [FILENAME]`,
	RunE: runE,
}

func NewCmd(cfg config.Config) *cobra.Command {
	var err error
	client, err = gcs.NewClient(context.Background(), cfg.GCS)
	if err != nil {
		log.Fatal(err)
	}

	cmd.PersistentFlags().StringVarP(&bucket, "bucket", "b", "", "GCS bucket name")
	cmd.PersistentFlags().StringVarP(&filename, "filename", "f", "", "file name (including the dir path)")
	cmd.MarkPersistentFlagRequired("bucket")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	return client.Bucket(bucket).Object(filename).Delete(cmd.Context())
}
