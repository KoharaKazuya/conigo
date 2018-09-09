package conigo

import (
	"fmt"
	"reflect"
)

var (
	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

func (c *Container) validateConstructorFunction(constructor interface{}) error {
	constructorType := reflect.TypeOf(constructor)

	if constructorType.Kind() != reflect.Func {
		return fmt.Errorf("must provide function: got %v (type %v)", constructor, constructorType)
	}

	if constructorType.NumOut() < 1 {
		return fmt.Errorf("must provide function has return type: got %v", constructor)
	}

	if constructorType.NumOut() > 2 {
		return fmt.Errorf("must provide function has less return types than 3: got %v", constructor)
	}

	if constructorType.NumOut() > 1 {
		returnErrorType := constructorType.Out(1)
		if !returnErrorType.Implements(errorInterface) {
			return fmt.Errorf("must provide function that second return type is error: got %v (sencod return type %v)", constructor, returnErrorType)
		}
	}

	return nil
}

func (c *Container) validateResolverFunction(resolver interface{}) error {
	resolverType := reflect.TypeOf(resolver)

	if resolverType.Kind() != reflect.Func {
		return fmt.Errorf("must provide function: got %v (type %v)", resolver, resolverType)
	}

	if resolverType.NumOut() > 1 || resolverType.NumOut() > 0 && !resolverType.Out(0).Implements(errorInterface) {
		return fmt.Errorf("must provide function returns only error: got %v", resolver)
	}

	return nil
}
