package gotong

import (
	"fmt"
	"os"
)

type ConfigBuilder struct {
	name   string
	errors []BuilderError
}

func NewConfigBuilder(name string) *ConfigBuilder {
	return &ConfigBuilder{
		name:   name,
		errors: []BuilderError{},
	}
}

// root > config
func (b *ConfigBuilder) MakeBootstrap() error {
	return b.MakeDir().
		MakeFileConfigYaml().
		MakeFileConfigGo().
		Build()
}

func (b *ConfigBuilder) Build() error {
	if len(b.errors) == 0 {
		return nil
	}

	return b.errors[0].Error()
}

func (b *ConfigBuilder) MakeDir() *ConfigBuilder {
	dirs := []string{"config"}

	for _, dir := range dirs {
		if err := os.MkdirAll(fmt.Sprintf("%s/%s", appName, dir), os.ModePerm); err != nil {
			b.setError("ConfigBuilder.MakeDir", err)
			return b
		}
	}

	return b
}

func (b *ConfigBuilder) MakeFileConfigYaml() *ConfigBuilder {
	return b
}

func (b *ConfigBuilder) MakeFileConfigGo() *ConfigBuilder {
	if err := copyTemplate("config/config.go", "config/config.gotong", map[string]string{"packageName": getPackageName(appName)}); err != nil {
		b.setError("ApiBuilder.MakeFileApiRouter", err)
	}

	return b
}

func (b *ConfigBuilder) setError(where string, why error) {
	b.errors = append(b.errors, NewBuilderError(where, why))
}
