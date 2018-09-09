package conigo_test

import (
	"fmt"

	"github.com/KoharaKazuya/conigo"
)

type LibName string

type StructExample struct {
	value LibName
}

func Example_struct() {
	container := conigo.New()

	// provide LibName
	if err := container.Provide(func() LibName {
		return "conigo"
	}); err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	// provide *StructExample (depends on LibName)
	if err := container.Provide(func(n LibName) *StructExample {
		return &StructExample{value: n}
	}); err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	// resolve *StructExample
	if err := container.Resolve(func(s *StructExample) {
		fmt.Printf("s: %v", s)
	}); err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	// Output:
	// s: &{conigo}
}
