package gotong

import (
	"fmt"
	"os"
)

type GotongBuilder struct {
	api    *ApiBuilder
	config *ConfigBuilder
	domain *DomainBuilder
	pkg    *PkgBuilder
}

func NewGotongBuilder(appName string) *GotongBuilder {
	return &GotongBuilder{
		api:    NewApiBuilder(appName),
		config: NewConfigBuilder(appName),
		domain: NewDomainBuilder(appName),
		pkg:    NewPkgBuilder(appName),
	}
}

func (b *GotongBuilder) MakeBootstrap() error {
	if err := b.makeBootstrap(); err != nil {
		b.Revert()
		return err
	}

	return nil
}

func (b *GotongBuilder) makeBootstrap() error {
	if err := b.api.MakeBootstrap(); err != nil {
		return err
	}
	if err := b.config.MakeBootstrap(); err != nil {
		return err
	}
	if err := b.domain.MakeBootstrap(); err != nil {
		return err
	}
	if err := b.pkg.MakeBootstrap(); err != nil {
		return err
	}

	return nil
}

func (b *GotongBuilder) Revert() error {
	fmt.Printf("🐶🐶[DEBUG]🐶🐶 REVERT CALLED\n")
	return os.RemoveAll(appName)
}
