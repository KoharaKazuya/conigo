package conigo

import (
	"fmt"
	"reflect"
)

func (c *Container) resolveArgs(resolver interface{}) ([]reflect.Value, error) {
	resolverType := reflect.TypeOf(resolver)

	var argValues []reflect.Value
	for i := 0; i < resolverType.NumIn(); i++ {
		argType := resolverType.In(i)
		resolved, err := c.resolveByType(argType)
		if err != nil {
			return nil, fmt.Errorf("resolve failed (for %v): %s", argType, err.Error())
		}
		argValues = append(argValues, resolved)
	}

	return argValues, nil
}

func (c *Container) resolveByType(t reflect.Type) (reflect.Value, error) {
	for k, v := range c.registry {
		if k == t {
			if !v.resolved {
				if v.constructing {
					return reflect.Value{}, fmt.Errorf("detected cyclic dependency for %v", t)
				}
				v.constructing = true
				value, err := c.construct(v.constructor)
				if err != nil {
					return reflect.Value{}, err
				}
				v.value = value
				v.constructing = false
			}
			return v.value, nil
		}
	}

	return reflect.Value{}, fmt.Errorf("cannot find provider for %v", t)
}
