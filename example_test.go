package conigo_test

import (
	"fmt"

	"github.com/KoharaKazuya/conigo"
)

func Example() {
	// use any type as dependency key
	type Key string

	// generate DI container
	container := conigo.New()

	// provide "Value" string as Key
	err := container.Provide(func() Key { return "Value" })
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	// resolve arguments by type automatically
	err = container.Resolve(func(dep Key) {
		fmt.Printf("dep: %s", dep)
	})
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	// Output:
	// dep: Value
}
