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
	"reflect"
)

// Container is a DI (Dependency Injection) container
type Container struct {
	registry map[reflect.Type]*entry
}

type entry struct {
	constructor  interface{}
	constructing bool
	resolved     bool
	value        reflect.Value
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
