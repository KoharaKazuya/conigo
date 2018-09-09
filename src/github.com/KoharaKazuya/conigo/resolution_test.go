package conigo_test

import (
	"errors"
	"testing"

	"github.com/KoharaKazuya/conigo"
)

type a string
type b string

func TestCyclicDependencyDetection(t *testing.T) {
	container := conigo.New()

	container.Provide(func(_ b) a { return "a" })
	container.Provide(func(_ a) b { return "b" })
	err := container.Resolve(func(_ b) {})

	if err == nil {
		t.Error("cyclic dependency detection pass through")
	}
}

func TestNoProviderError(t *testing.T) {
	container := conigo.New()

	container.Provide(func() a { return "a" })
	err := container.Resolve(func(_ b) {})

	if err == nil {
		t.Error("no provider error pass through")
	}
}

func TestResolverError(t *testing.T) {
	container := conigo.New()
	err := container.Resolve(func() error {
		return errors.New("Test Error")
	})

	if err == nil {
		t.Error("resolver error ignored")
	}
}
