package gcs

import (
	"context"
	"os"
	"path"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type Option struct {
	CredentialFile string // base path is $HOME
}

func NewClient(ctx context.Context, opt Option) (*storage.Client, error) {
	homedir, _ := os.UserHomeDir()
	credFile := path.Join(homedir, opt.CredentialFile)
	opts := []option.ClientOption{option.WithCredentialsFile(credFile)}
	return storage.NewClient(ctx, opts...)
}
