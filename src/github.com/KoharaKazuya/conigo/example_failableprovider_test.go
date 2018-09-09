package conigo_test

import (
	"errors"
	"fmt"
	"os"

	"github.com/KoharaKazuya/conigo"
)

type DebugMode bool

func Example_failableProvider() {
	container := conigo.New()

	// constructor may fail by case
	container.Provide(func() (DebugMode, error) {
		v, ok := os.LookupEnv("DEBUG")
		if !ok {
			return false, errors.New("No DEBUG environment variable")
		}
		return v == "1", nil
	})

	// fail when resolved
	if err := container.Resolve(func(mode DebugMode) {
		fmt.Printf("DebugMode: %t", mode)
	}); err != nil {
		fmt.Print("Resolve Error")
	}

	// Output:
	// Resolve Error
}
