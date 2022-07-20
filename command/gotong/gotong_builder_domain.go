package gotong

import (
	"fmt"
	"os"
)

type DomainBuilder struct {
	name   string
	errors []BuilderError
}

func NewDomainBuilder(name string) *DomainBuilder {
	return &DomainBuilder{
		name:   name,
		errors: []BuilderError{},
	}
}

// root > domain
func (b *DomainBuilder) MakeBootstrap() error {
	return b.MakeDir().Build()
}

func (b *DomainBuilder) MakeDir() *DomainBuilder {
	dirs := []string{"domain"}

	for _, dir := range dirs {
		if err := os.MkdirAll(fmt.Sprintf("%s/%s", appName, dir), os.ModePerm); err != nil {
			b.setError("DomainBuilder.MakeDir", err)
			return b
		}
	}

	return b
}

func (b *DomainBuilder) Build() error {
	if len(b.errors) == 0 {
		return nil
	}

	return b.errors[0].Error()
}

func (b *DomainBuilder) setError(where string, why error) {
	b.errors = append(b.errors, NewBuilderError(where, why))
}
