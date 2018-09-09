package conigo_test

import (
	"testing"

	"github.com/KoharaKazuya/conigo"
)

func TestDoubleRegistrationError(t *testing.T) {
	container := conigo.New()

	container.Provide(func() int { return 0 })
	err := container.Provide(func() int { return 1 })

	if err == nil {
		t.Error("double registration error pass through")
	}
}
