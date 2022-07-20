package gotong

import (
	"fmt"
	"os"
)

type PkgBuilder struct {
	name   string
	errors []BuilderError
}

func NewPkgBuilder(name string) *PkgBuilder {
	return &PkgBuilder{
		name:   name,
		errors: []BuilderError{},
	}
}

// root > config
func (b *PkgBuilder) MakeBootstrap() error {
	return b.MakeDir().
		MakeFileLoggerGo().
		Build()
}

func (b *PkgBuilder) Build() error {
	if len(b.errors) == 0 {
		return nil
	}

	return b.errors[0].Error()
}

func (b *PkgBuilder) MakeDir() *PkgBuilder {
	dirs := []string{"pkg/logger"}

	for _, dir := range dirs {
		if err := os.MkdirAll(fmt.Sprintf("%s/%s", appName, dir), os.ModePerm); err != nil {
			b.setError("PkgBuilder.MakeDir", err)
			return b
		}
	}

	return b
}

func (b *PkgBuilder) MakeFileConfigYaml() *PkgBuilder {
	return b
}

func (b *PkgBuilder) MakeFileLoggerGo() *PkgBuilder {
	if err := copyTemplate("pkg/logger/logger.go", "pkg/logger/logger.gotong", map[string]string{"packageName": getPackageName(appName)}); err != nil {
		b.setError("ApiBuilder.MakeFileApiRouter", err)
	}

	return b
}

func (b *PkgBuilder) setError(where string, why error) {
	b.errors = append(b.errors, NewBuilderError(where, why))
}
