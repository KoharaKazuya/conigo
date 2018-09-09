package conigo_test

import (
	"testing"

	"github.com/KoharaKazuya/conigo"
)

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
