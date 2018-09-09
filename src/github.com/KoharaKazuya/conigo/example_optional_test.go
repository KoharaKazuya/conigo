package conigo_test

import (
	"fmt"

	"github.com/KoharaKazuya/conigo"
)

type Spy struct{}

func (s *Spy) LookInto() {}

func Example_optional() {
	// when Spy is provided
	containerA := conigo.New()
	containerA.Provide(func() *Spy {
		return &Spy{}
	})
	containerA.Resolve(func(spy *Spy) {
		if spy != nil {
			spy.LookInto()
			fmt.Print("A: Spy!")
		}
	})

	// when Spy is not provided
	containerB := conigo.New()
	containerB.Provide(func() *Spy {
		return nil
	})
	containerB.Resolve(func(spy *Spy) {
		if spy != nil {
			spy.LookInto()
			fmt.Print("B: Spy!")
		}
	})

	// Output:
	// A: Spy!
}
