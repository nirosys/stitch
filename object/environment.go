package object

import (
	"fmt"

	"github.com/rs/xid"
)

type Environment struct {
	packages     map[string]*Package
	store        map[string]Object
	unboundNodes map[string]Object
	parent       *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store:        map[string]Object{},
		unboundNodes: map[string]Object{},
		packages:     map[string]*Package{},
		parent:       nil,
	}
}

func (e *Environment) IsGlobal() bool {
	return e.parent == nil
}

func (e *Environment) Clone() *Environment {
	env := NewEnvironment()
	env.parent = e
	return env
}

func (e *Environment) GetPackage(name string) (*Package, error) {
	if e.parent != nil {
		return e.parent.GetPackage(name)
	} else {
		if pkg, ok := e.packages[name]; ok {
			return pkg, nil
		} else {
			return nil, fmt.Errorf("unknown identifier '%s'", name)
		}
	}
}

func (e *Environment) PutPackage(pkg *Package) error {
	if e.parent != nil {
		return e.parent.PutPackage(pkg)
	} else {
		if _, ok := e.packages[pkg.Name]; !ok {
			e.packages[pkg.Name] = pkg
		}
		return nil
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	// Check this scope
	if obj, ok := e.store[name]; !ok {
		if obj, ok := e.unboundNodes[name]; ok {
			return obj, ok
		}
	} else {
		return obj, ok
	}
	if e.parent == nil {
		obj, ok := e.packages[name]
		return obj, ok
	} else {
		return e.parent.Get(name)
	}
}

func (e *Environment) Put(name string, val Object) Object {
	var retObj Object

	if v, ok := e.store[name]; ok { // Check local scope..
		if v.Type() == NodeObjectType { // Track unbound nodes..
			found := false
			for _, o := range e.store {
				found = found || (o == val)
			}
			if !found {
				e.PutUnboundNode(v)
			}
		}
		e.store[name] = val
		retObj = val
	} else if _, has := e.parent_get(name); has {
		retObj = e.parent.Put(name, val)
	} else {
		e.store[name] = val
		retObj = val
	}

	for unboundIdent, unbound := range e.unboundNodes {
		if unbound == val {
			delete(e.unboundNodes, unboundIdent)
		}
	}
	return retObj
}

func (e *Environment) parent_get(name string) (Object, bool) {
	if e.parent == nil {
		return nil, false
	} else {
		return e.parent.Get(name)
	}
}

// This function will track flow nodes, if our env is the global environment,
// in order to allow us to wire up top-level nodes with no Input referenced.
func (e *Environment) PutUnboundNode(val Object) Object {
	if val.Type() == NodeObjectType && e.IsGlobal() {
		name := "_unbound" + xid.New().String()
		e.unboundNodes[name] = val
		return val
	}
	return nil
}

func (e *Environment) GetNames() []string {
	names := []string{}
	for k, _ := range e.store {
		names = append(names, k)
	}
	return names
}

func (e *Environment) GetUnboundNodes() []string {
	names := []string{}
	for k, _ := range e.unboundNodes {
		names = append(names, k)
	}
	return names
}
