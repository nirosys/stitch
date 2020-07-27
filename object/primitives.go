package object

import (
	"fmt"
	"strings"
)

// Integer ////////////////////////////////////////////////////////////////////

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return IntegerObjectType
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Identifier(name string) (Object, error) {
	return nil, fmt.Errorf("'%s' not defined for Integer", name)
}

func (i *Integer) Add(other Object) (Object, error) {
	if other.Type() != i.Type() {
		return nil, fmt.Errorf("type mis-match")
	}
	o := other.(*Integer)
	return &Integer{Value: i.Value + o.Value}, nil
}

func (i *Integer) Subtract(other Object) (Object, error) {
	if other.Type() != i.Type() {
		return nil, fmt.Errorf("type mis-match")
	}
	o := other.(*Integer)
	return &Integer{Value: i.Value - o.Value}, nil
}

func (i *Integer) Multiply(other Object) (Object, error) {
	if other.Type() != i.Type() {
		return nil, fmt.Errorf("type mis-match")
	}
	o := other.(*Integer)
	return &Integer{Value: i.Value * o.Value}, nil
}

func (i *Integer) Divide(other Object) (Object, error) {
	if other.Type() != i.Type() {
		return nil, fmt.Errorf("type mis-match")
	}
	o := other.(*Integer)
	return &Integer{Value: i.Value / o.Value}, nil
}

func (i *Integer) Modulus(other Object) (Object, error) {
	if other.Type() != i.Type() {
		return nil, fmt.Errorf("type mis-match")
	}
	o := other.(*Integer)
	return &Integer{Value: i.Value % o.Value}, nil
}

func (i *Integer) IsComparable(other Comparable) bool {
	_, ok := other.(*Integer)
	return ok
}

func (i *Integer) Equals(other Comparable) (bool, error) {
	if v, err := getIntValue(other); err != nil {
		return false, err
	} else {
		return i.Value == v, nil
	}
}
func (i *Integer) GreaterThan(other Comparable) (bool, error) {
	if v, err := getIntValue(other); err != nil {
		return false, err
	} else {
		return i.Value > v, nil
	}
}
func (i *Integer) LessThan(other Comparable) (bool, error) {
	if v, err := getIntValue(other); err != nil {
		return false, err
	} else {
		return i.Value < v, nil
	}
}

func getIntValue(c Comparable) (int64, error) {
	if that, ok := c.(*Integer); ok {
		return that.Value, nil
	} else {
		return 0, fmt.Errorf("cannot compare an integer against %s", c.Type())
	}
}

// String /////////////////////////////////////////////////////////////////////
type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return StringObjectType
}

func (s *String) Inspect() string {
	return fmt.Sprintf("\"%s\"", s.Value)
}

func (i *String) Identifier(name string) (Object, error) {
	return nil, fmt.Errorf("'%s' not defined for String", name)
}

// NodeSlot ///////////////////////////////////////////////////////////////////
type NodeSlot struct {
	Name    string
	IsInput bool
	Node    *Node
}

func (n *NodeSlot) Type() ObjectType {
	return NodeSlotType
}

func (n *NodeSlot) Inspect() string {
	return fmt.Sprintf("Slot{Name: \"%s\", Input: %t} On Node@%p", n.Name, n.IsInput, n.Node)
}

func (n *NodeSlot) Identifier(name string) (Object, error) {
	return nil, fmt.Errorf("'%s' not defined for slot type")
}

func (n *NodeSlot) Connect(obj Connectable) (Object, error) {
	var otherSlot *NodeSlot
	switch t := obj.(type) {
	case *Node:
		if s := t.GetSlot("Input"); s == nil {
			return nil, fmt.Errorf("no channel 'Input' on node")
		} else {
			otherSlot = s
		}
	case *NodeSlot:
		otherSlot = t
	case *List:
		for _, obj := range t.Contents {
			if c, ok := obj.(Connectable); !ok {
				return nil, fmt.Errorf("cannot connect to type '%s'", obj.Type())
			} else {
				n.Connect(c)
			}
		}
	default:
		return nil, fmt.Errorf("cannot connect to type '%s'", obj.Type())
	}

	return n.Node.ConnectSlots(n.Name, otherSlot)
}

// MapObject //////////////////////////////////////////////////////////////////
type MapObject struct {
	Fields map[string]Object
}

func (m *MapObject) Type() ObjectType {
	return MapObjectType
}

func (m *MapObject) Inspect() string {
	var buff strings.Builder
	buff.WriteByte('{')
	for k, v := range m.Fields {
		buff.WriteString(k)
		buff.WriteByte('=')
		buff.WriteString(v.Inspect())
		buff.WriteByte(';')
	}
	buff.WriteByte('}')
	return buff.String()
}

func (m *MapObject) Identifier(name string) (Object, error) {
	if obj, has := m.Fields[name]; !has {
		return nil, fmt.Errorf("field not found: %s", name)
	} else {
		return obj, nil
	}
}

// BoolObject /////////////////////////////////////////////////////////////////
type BoolObject bool

func (b *BoolObject) Type() ObjectType { return BoolObjectType }
func (b *BoolObject) Inspect() string {
	if bool(*b) {
		return "true"
	} else {
		return "false"
	}
}
func (b *BoolObject) Identifier(name string) (Object, error) {
	return nil, fmt.Errorf("'%s' not defined for list", name)
}

func (b *BoolObject) IsComparable(other Comparable) bool {
	_, ok := other.(*BoolObject)
	return ok
}

func (b *BoolObject) Equals(other Comparable) (bool, error) {
	this := bool(*b)
	if that, ok := other.(*BoolObject); ok {
		return this == bool(*that), nil
	} else {
		return false, fmt.Errorf("unable to compare bool against type: %T", other)
	}
}

func (b *BoolObject) GreaterThan(other Comparable) (bool, error) {
	return false, fmt.Errorf("cannot compare bool relatively")
}

func (b *BoolObject) LessThan(other Comparable) (bool, error) {
	return false, fmt.Errorf("cannot compare bool relatively")
}

func (b *BoolObject) Not() *BoolObject {
	this := bool(*b)
	return NewBoolObject(!this)
}

func NewBoolObject(v bool) *BoolObject {
	obj := BoolObject(v)
	return &obj
}
