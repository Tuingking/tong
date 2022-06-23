package list

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"cloud.google.com/go/storage"
	"github.com/Tuingking/tong/config"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
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

	for i, v := range result {
		fmt.Printf("🐶🐶[DEBUG]🐶🐶 file-%d: %+v\n", i, v)
	}

	return nil
}

func GetObjectNames(ctx context.Context, bucket, directory string) ([]string, error) {
	result := []string{}

	opts := []option.ClientOption{option.WithCredentialsFile(getConfigFile())}
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return result, err
	}

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

func getConfigFile() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), "..", "..", "..", "config", "xwork-dev-5e7544287095.json")
	// return path.Join(path.Dir(filename), "..", "..", "..", "config", "tk-dev-micro-71a89e645c42.json")
}
