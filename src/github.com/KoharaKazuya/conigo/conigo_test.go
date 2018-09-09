package conigo_test

import (
	"errors"
	"testing"

	"github.com/KoharaKazuya/conigo"
)

type a string
type b string

func TestProvideValidation(t *testing.T) {
	cases := []struct {
		input   interface{}
		success bool
	}{
		{func() int { return 0 }, true},
		{func() (int, error) { return 0, nil }, true},
		{func(int) int { return 0 }, true},
		{func(int, string) int { return 0 }, true},
		{1, false},
		{func() {}, false},
		{func() (int, string, error) { return 0, "", nil }, false},
		{func() (int, string) { return 0, "" }, false},
	}

	for i, c := range cases {
		container := conigo.New()

		err := container.Provide(c.input)
		if (err == nil) != c.success {
			t.Errorf("[case%03d] got: %v\ncase: %#v", i, err, c)
		}
	}
}

func TestResolveValidation(t *testing.T) {
	cases := []struct {
		input   interface{}
		success bool
	}{
		{func() {}, true},
		{func() error { return nil }, true},
		{func(int, string) {}, true},
		{func(int, string) error { return nil }, true},
		{1, false},
		{func() int { return 0 }, false},
		{func() (int, error) { return 0, nil }, false},
		{func() (error, int) { return nil, 0 }, false},
	}

	for i, c := range cases {
		container := conigo.New()
		container.Provide(func() int { return 0 })
		container.Provide(func() string { return "" })

		err := container.Resolve(c.input)
		if (err == nil) != c.success {
			t.Errorf("[case%03d] got: %v\ncase: %#v", i, err, c)
		}
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

func TestDoubleRegistrationError(t *testing.T) {
	container := conigo.New()

	container.Provide(func() int { return 0 })
	err := container.Provide(func() int { return 1 })

	if err == nil {
		t.Error("double registration error pass through")
	}
}

func TestCyclicDependencyDetection(t *testing.T) {
	container := conigo.New()

	container.Provide(func(_ b) a { return "a" })
	container.Provide(func(_ a) b { return "b" })
	err := container.Resolve(func(_ b) {})

	if err == nil {
		t.Error("cyclic dependency detection pass through")
	}
}
