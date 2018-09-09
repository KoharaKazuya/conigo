package conigo_test

import (
	"fmt"

	"github.com/KoharaKazuya/conigo"
)

type DB interface {
	Name() string
}

type ReadonlyDB DB

type ReadWriteDB DB

type DBImpl struct {
	name string
}

func (d *DBImpl) Name() string {
	return d.name
}

func Example_namedProvider() {
	container := conigo.New()

	// provide ReadonlyDB
	if err := container.Provide(func() ReadonlyDB {
		return &DBImpl{name: "ro"}
	}); err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	// provide ReadWriteDB
	if err := container.Provide(func() ReadWriteDB {
		return &DBImpl{name: "rw"}
	}); err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	// resolver can distinguish each other
	if err := container.Resolve(func(ro ReadonlyDB, rw ReadWriteDB) {
		fmt.Printf("%s, %s", ro.Name(), rw.Name())
	}); err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	// Output:
	// ro, rw
}
