package conigo_test

import (
	"fmt"

	"github.com/KoharaKazuya/conigo"
)

type HogeService interface {
	Say() string
}

type HogeImpl struct{}

func (s *HogeImpl) Say() string {
	return "Hoge"
}

func Example_interface() {
	container := conigo.New()

	// provide concrete implementation as HogeService
	if err := container.Provide(func() HogeService {
		return &HogeImpl{}
	}); err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	// resolve HogeService
	if err := container.Resolve(func(s HogeService) {
		fmt.Printf("Say: %s", s.Say())
	}); err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	// Output:
	// Say: Hoge
}
