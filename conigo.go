/*
Package conigo provides a simple DI (Dependency Injection) container.

Conigo helps you to inject the dependencies into constructors.
Conigo automatically resolves the dependency graph on demand,
by the argument types and the return types.
All you have to do is writing constructors with the dependencies as arguments.

Conigo is simple. It has only one constructor and two methods and depends on no third-party package.

If you need rich DI support library, I recommend to use [dig](https://github.com/uber-go/dig).
Conigo is the almost subset of dig.
*/
package conigo

import (
	"fmt"
	"reflect"
)

// Container is a DI (Dependency Injection) container
type Container struct {
	registry map[reflect.Type]*entry
}

// New generates Container
func New() *Container {
	return &Container{
		registry: make(map[reflect.Type]*entry),
	}
}

// Provide registers a constructor on Container.
//
// constructor must be function.
// constructor can have any arguments. The arguments are resolved and injected by type.
// constructor must have a return type. The return type is used as dependency key.
// constructor can have an error. If error is not nil, Container fails Resolve().
//
// constructor examples
//
//     func() *StructA { return &StructA{} }                              // OK.
//     func(a *StructA) *StructB { return &StructB{a} }                   // OK.
//     func() (*StructA, error) { return &StructA{}, nil }                // OK.
//     func(a *StructA, b InterfaceB) InterfaceC { return &StructC{a,b} } // OK.
//     func() {}                                                          // Invalid.
//     func() error { return nil }                                        // Invalid.
//     func() (int, string) { return 0, "" }                              // Invalid.
//     func() (int, string, error) { return 0, "", nil }                  // Invalid.
func (c *Container) Provide(constructor interface{}) error {
	if err := c.validateConstructorFunction(constructor); err != nil {
		return err
	}

	if err := c.register(constructor); err != nil {
		return err
	}

	return nil
}

// Resolve injects dependencies by providers on Container.
//
// resolver must be function.
// resolver can have any arguments. The arguments are resolved and injected by type.
// resolver can have an error. If error is not nil, Container fails Resolve().
//
// resolver examples
//
//     func() {}                                           // OK.
//     func(a *StructA) {}                                 // OK.
//     func() error { return nil }                         // OK.
//     func(a *StructA, b InterfaceB) error { return nil } // OK.
//     func() int { return 0 }                             // Invalid.
//     func() (int, error) { return 0, nil }               // Invalid.
func (c *Container) Resolve(resolver interface{}) error {
	if err := c.validateResolverFunction(resolver); err != nil {
		return err
	}

	argValues, err := c.resolveArgs(resolver)
	if err != nil {
		return err
	}

	resolverReturns := reflect.ValueOf(resolver).Call(argValues)
	if len(resolverReturns) > 0 && !resolverReturns[0].IsNil() {
		return resolverReturns[0].Interface().(error)
	}

	return nil
}

type entry struct {
	constructor  interface{}
	constructing bool
	resolved     bool
	value        reflect.Value
}

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

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
