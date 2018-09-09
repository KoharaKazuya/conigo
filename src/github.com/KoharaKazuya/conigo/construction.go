package conigo

import (
	"reflect"
)

func (c *Container) construct(constructor interface{}) (reflect.Value, error) {
	argValues, err := c.resolveArgs(constructor)
	if err != nil {
		return reflect.Value{}, err
	}

	constructorReturns := reflect.ValueOf(constructor).Call(argValues)
	if len(constructorReturns) > 1 {
		if err, ok := constructorReturns[1].Interface().(error); !ok || err != nil {
			return reflect.Value{}, err
		}
	}

	return constructorReturns[0], nil
}
