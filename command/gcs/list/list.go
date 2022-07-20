package list

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"github.com/spf13/cobra"
	"github.com/tuingking/tong/config"
	"github.com/tuingking/tong/pkg/gcs"
	"google.golang.org/api/iterator"
)

var (
	client *storage.Client

	// flag
	bucket string // GCS bucket name
	dir    string
)

var cmd = &cobra.Command{
	Use:   "list",
	Short: "Get list of files in GCS bucket.",
	Long: `Get list of files in GCS bucket.
	example: tong gcs list -b [BUCKET_NAME]`,
	RunE: runE,
}

func NewCmd(cfg config.Config) *cobra.Command {
	var err error
	client, err = gcs.NewClient(context.Background(), cfg.GCS)
	if err != nil {
		log.Fatal(err)
	}

	cmd.PersistentFlags().StringVarP(&bucket, "bucket", "b", "", "GCS bucket name")
	cmd.PersistentFlags().StringVar(&dir, "dir", "", "file directory")
	cmd.MarkPersistentFlagRequired("bucket")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {

	result, err := GetObjectNames(context.Background(), bucket, dir)
	if err != nil {
		return err
	}

	fmt.Printf("result:\n")
	for i, v := range result {
		fmt.Printf("%d. %+v\n", i, v)
	}

	return nil
}

func GetObjectNames(ctx context.Context, bucket, directory string) ([]string, error) {
	result := []string{}

	objs := client.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: directory})
	for {
		obj, err := objs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return result, err
		}
		result = append(result, obj.Name)
	}
	return result, nil
}
