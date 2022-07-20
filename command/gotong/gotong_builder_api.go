package gotong

import (
	"fmt"
	"os"
)

type ApiBuilder struct {
	name   string
	errors []BuilderError
}

func NewApiBuilder(name string) *ApiBuilder {
	return &ApiBuilder{
		name:   name,
		errors: []BuilderError{},
	}
}

// root > app > api
func (b *ApiBuilder) MakeBootstrap() error {
	return b.MakeDir().
		MakeFileApiHandler().
		MakeFileApiResponse().
		MakeFileApiRouter().
		MakeFileApiServer().
		Build()
}

func (b *ApiBuilder) Build() error {
	if len(b.errors) == 0 {
		return nil
	}

	return b.errors[0].Error()
}

func (b *ApiBuilder) MakeDir() *ApiBuilder {
	dirs := []string{"app/api"}

	for _, dir := range dirs {
		if err := os.MkdirAll(fmt.Sprintf("%s/%s", appName, dir), os.ModePerm); err != nil {
			b.setError("ApiBuilder.MakeDir", err)
			return b
		}
	}

	return b
}

func (b *ApiBuilder) MakeFileApiHandler() *ApiBuilder {
	if err := copyTemplate("app/api/api_handler.go", "app/api/api_handler.gotong", map[string]string{"packageName": getPackageName(appName)}); err != nil {
		b.setError("ApiBuilder.MakeFileApiRouter", err)
	}

	return b
}

func (b *ApiBuilder) MakeFileApiResponse() *ApiBuilder {
	if err := copyTemplate("app/api/api_response.go", "app/api/api_response.gotong", map[string]string{"packageName": getPackageName(appName)}); err != nil {
		b.setError("ApiBuilder.MakeFileApiRouter", err)
	}

	return b
}

func (b *ApiBuilder) MakeFileApiRouter() *ApiBuilder {
	if err := copyTemplate("app/api/api_router.go", "app/api/api_router.gotong", map[string]string{"packageName": getPackageName(appName)}); err != nil {
		b.setError("ApiBuilder.MakeFileApiRouter", err)
	}

	return b
}

func (b *ApiBuilder) MakeFileApiServer() *ApiBuilder {
	if err := copyTemplate("app/api/api_server.go", "app/api/api_server.gotong", map[string]string{"packageName": getPackageName(appName)}); err != nil {
		b.setError("ApiBuilder.MakeFileApiRouter", err)
	}

	return b
}

func (b *ApiBuilder) setError(where string, why error) {
	b.errors = append(b.errors, NewBuilderError(where, why))
}
