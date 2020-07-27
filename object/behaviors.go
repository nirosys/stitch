package object

import (
	"github.com/nirosys/stitch/ast"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	Identifier(name string) (Object, error)
}

type Computable interface {
	Object
	Add(Object) (Object, error)
	Subtract(Object) (Object, error)
	Multiply(Object) (Object, error)
	Divide(Object) (Object, error)
	Modulus(Object) (Object, error)
}

type Connectable interface {
	Object
	Connect(Connectable) (Object, error)
}

type Callable interface {
	Object
	Scope() *Environment
	FuncBody() ast.Expression
	FuncParameters() []*ast.FunctionParameter
}

type Comparable interface {
	Object
	IsComparable(Comparable) bool
	Equals(Comparable) (bool, error)
	GreaterThan(Comparable) (bool, error)
	LessThan(Comparable) (bool, error)
}

type Constructable interface {
	Object
	Arguments() []*ast.FunctionParameter
	Construct([]Object) (Object, error)
}
