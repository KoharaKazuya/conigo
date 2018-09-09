package conigo

import (
	"fmt"
	"reflect"
)

func (c *Container) register(constructor interface{}) error {
	constructorType := reflect.TypeOf(constructor)

	key := constructorType.Out(0)
	if _, ok := c.registry[key]; ok {
		return fmt.Errorf("cannot provide the same return type: got %v (return type %v)", constructor, key)
	}

	c.registry[key] = &entry{
		constructor: constructor,
	}

	return nil
}
