// core/application/meta.go
package application

import (
	"fmt"
	"reflect"
)

// Snapshot represents a reflection summary of the Application composition.
// It lists the registered modules and adapters in runtime.
type Snapshot struct {
	Modules  []string
	Adapters []string
}

// snapshotOf builds a Snapshot representation of the Application state.
func snapshotOf(a *Application) *Snapshot {
	modNames := []string{}
	for _, m := range a.modules {
		modNames = append(modNames, m.Name)
	}

	adapters := []string{}
	for _, ad := range a.adapters {
		adapters = append(adapters, reflect.TypeOf(ad).String())
	}

	fmt.Printf("[snapshot] modules=%v adapters=%v\n", modNames, adapters)

	return &Snapshot{
		Modules:  modNames,
		Adapters: adapters,
	}
}
