package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nirosys/stitch/ast"
)

type ObjectType string

const (
	UnknownObjectType  = "WUT"
	IntegerObjectType  = "INTEGER"
	StringObjectType   = "STRING"
	NodeObjectType     = "NODE"
	NodeTypeObjectType = "NODE TYPE"
	NodeSlotType       = "NODESLOT"
	ConnectionType     = "CONNECTION"
	PackageObjectType  = "PACKAGE"
	FunctionObjectType = "FUNCTION"
	InternalObjectType = "INTERNAL FUNCTION"
	ListObjectType     = "LIST"
	BoolObjectType     = "BOOL"
	ModifierObjectType = "MODIFIER"
	MapObjectType      = "MAP"
)

func (t ObjectType) IsPrimitive() bool {
	switch t {
	case IntegerObjectType:
		return true
	case StringObjectType:
		return true
	default:
		return false
	}
}

// Package ////////////////////////////////////////////////////////////////////
type Package struct {
	Name        string
	Environment *Environment
}

func NewPackage(name string, env *Environment) (*Package, error) {
	return &Package{Name: name, Environment: env}, nil
}

func (p *Package) Type() ObjectType {
	return PackageObjectType
}

func (p *Package) Inspect() string {
	return fmt.Sprintf("import \"%s.stitch\"", p.Name)
}

func (p *Package) Identifier(name string) (Object, error) {
	if obj, ok := p.Environment.Get(name); !ok {
		return nil, fmt.Errorf("unknown identifier '%s' in package '%s'", name, p.Name)
	} else {
		return obj, nil
	}
}

// Connection /////////////////////////////////////////////////////////////////
type Connection struct {
	Start *NodeSlot
	End   *NodeSlot
}

func NewConnection(start, end *NodeSlot) *Connection {
	return &Connection{Start: start, End: end}
}

func (c *Connection) Type() ObjectType {
	return ConnectionType
}

func (c *Connection) Inspect() string {
	var buffer bytes.Buffer
	buffer.WriteString(c.Start.Inspect())
	buffer.WriteString(" -> ")
	buffer.WriteString(c.End.Inspect())
	return buffer.String()
}

// Function //////////////////////////////////////////////////////////////////
type Function struct {
	Parameters []*ast.FunctionParameter
	Body       *ast.BlockExpression
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FunctionObjectType }
func (f *Function) Inspect() string {
	var buffer bytes.Buffer
	buffer.WriteString("fn ")
	buffer.WriteString(f.Body.String())
	return buffer.String()
}

func (f *Function) Identifier(name string) (Object, error) {
	return nil, fmt.Errorf("'%s' not defined for function", name)
}

func (f *Function) FuncBody() ast.Expression {
	return f.Body
}

func (f *Function) FuncParameters() []*ast.FunctionParameter {
	return f.Parameters
}

func (f *Function) Scope() *Environment {
	return f.Env
}

// InternalObject /////////////////////////////////////////////////////////////

type NativeFunction struct {
	Fn     func(env *Environment, args []Object) (Object, error)
	Params []*ast.FunctionParameter
	//Params []object.ObjectType
}

type InternalFunction struct {
	Fn  *NativeFunction
	Env *Environment
}

func (i *InternalFunction) Type() ObjectType { return InternalObjectType }
func (i *InternalFunction) Inspect() string {
	fmt.Printf("Internal: %#v\n", i)
	var buffer bytes.Buffer
	buffer.WriteString("fn (")
	params := make([]string, 0, len(i.Fn.Params))
	for _, p := range i.Fn.Params {
		params = append(params, p.Identifier.Identifier)
	}
	buffer.WriteString(strings.Join(params, ", "))
	buffer.WriteByte(')')
	return buffer.String()
}

func (i *InternalFunction) Identifier(name string) (Object, error) {
	return nil, fmt.Errorf("'%s' not defined for function", name)
}

func (i *InternalFunction) FuncBody() ast.Expression {
	return nil
}

func (i *InternalFunction) FuncParameters() []*ast.FunctionParameter {
	return i.Fn.Params
}

func (i *InternalFunction) Scope() *Environment {
	return i.Env
}

// ListObject /////////////////////////////////////////////////////////////////
type List struct {
	Contents  []Object
	InnerType ObjectType
}

func (l *List) Type() ObjectType { return ListObjectType }
func (l *List) Inspect() string {
	var buffer bytes.Buffer
	buffer.WriteByte('[')
	objs := make([]string, 0, len(l.Contents))
	for _, obj := range l.Contents {
		if obj != nil {
			objs = append(objs, obj.Inspect())
		}
	}
	buffer.WriteString(strings.Join(objs, ", "))
	buffer.WriteByte(']')
	return buffer.String()
}

func (l *List) Identifier(name string) (Object, error) {
	return nil, fmt.Errorf("'%s' not defined for list", name)
}

func (l *List) Connect(other Connectable) (Object, error) {
	if l.InnerType != NodeObjectType && l.InnerType != NodeSlotType {
		return nil, fmt.Errorf("connections can not be made with type %s", l.InnerType)
	}

	for _, obj := range l.Contents {
		if c, ok := obj.(Connectable); !ok {
			// shouldn't happen, but jic
			return nil, fmt.Errorf("connections can not be made with type %s", obj.Type())
		} else {
			c.Connect(other)
		}
	}

	return other, nil
}

func (l *List) Add(other Object) (Object, error) {
	if t, ok := other.(*List); ok {
		ret := &List{InnerType: l.InnerType}

		if (len(l.Contents) > 0 && len(t.Contents) > 0) && t.InnerType != l.InnerType {
			return nil, fmt.Errorf("cannot concat list of %s to list of %s", t.InnerType, l.InnerType)
		}

		if len(t.Contents) == 0 && len(l.Contents) == 0 {
			return ret, nil
		}

		ret.Contents = make([]Object, 0, len(l.Contents)+len(t.Contents))
		ret.Contents = append(ret.Contents, l.Contents...)

		if len(t.Contents) == 0 {
			return ret, nil
		}

		ret.InnerType = t.InnerType // To handle the case that we're empty.
		ret.Contents = append(ret.Contents, t.Contents...)
		return ret, nil
	} else {
		return nil, fmt.Errorf("can only concatenate list to list (not %s)", other.Type())
	}
}

func (l *List) Subtract(other Object) (Object, error) {
	return nil, nil
}

func (l *List) Multiply(other Object) (Object, error) {
	return nil, nil
}

func (l *List) Divide(other Object) (Object, error) {
	return nil, nil
}

func (l *List) Modulus(other Object) (Object, error) {
	return nil, nil
}
