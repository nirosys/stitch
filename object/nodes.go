package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nirosys/stitch/ast"
)

// NodeType ///////////////////////////////////////////////////////////////////
type NodeType struct {
	Name        string                   // Name of the type (eg. "snmp:get", "std:template", "usr:foo")
	NodeArgs    []*ast.FunctionParameter // Argument names
	InputSlots  []string                 // Input Slot Names
	OutputSlots []string                 // Output Slot Names

	// For user supplied node types
	Body *ast.BlockExpression
	Env  *Environment
}

func (n *NodeType) Type() ObjectType { return NodeTypeObjectType }
func (n *NodeType) Inspect() string {
	var buffer bytes.Buffer
	buffer.WriteString("node ")
	buffer.WriteString(n.Name)
	buffer.WriteString(" {Inputs:[")
	buffer.WriteString(strings.Join(n.InputSlots, ","))
	buffer.WriteString("],Outputs:[")
	buffer.WriteString(strings.Join(n.OutputSlots, ","))
	buffer.WriteString("],Arguments:[")
	argnames := make([]string, 0, len(n.NodeArgs))
	for _, a := range n.NodeArgs {
		argnames = append(argnames, a.Identifier.String())
	}
	buffer.WriteString(strings.Join(argnames, ","))
	buffer.WriteString("]}")
	return buffer.String()
}

func (n *NodeType) Identifier(name string) (Object, error) {
	return nil, fmt.Errorf("'%s' not defined for node types", name)
}

func (n *NodeType) Arguments() []*ast.FunctionParameter {
	return n.NodeArgs
}

func (n *NodeType) Construct(args []Object) (Object, error) {
	node := NewNode()
	node.NodeType = n
	node.Arguments = args
	for _, input := range n.InputSlots {
		node.InputSlots[input] = struct{}{}
	}
	for _, output := range n.OutputSlots {
		node.OutputSlots[output] = struct{}{}
	}
	return node, nil
}

// Node ///////////////////////////////////////////////////////////////////
type Node struct {
	NodeType    *NodeType
	Arguments   []Object
	InputSlots  map[string]struct{}
	OutputSlots map[string]struct{}
	TagName     *string
	FieldName   *string

	connections map[string][]*NodeSlot
}

func NewNode() *Node {
	return &Node{
		Arguments:   []Object{},
		InputSlots:  map[string]struct{}{},
		OutputSlots: map[string]struct{}{},
		connections: map[string][]*NodeSlot{},
	}
}

func (f *Node) GetConnections() []*Connection {
	conns := []*Connection{}
	for n, slots := range f.connections {
		for _, slot := range slots {
			conns = append(conns, NewConnection(f.GetSlot(n), slot))
		}
	}
	return conns
}

func (f *Node) GetSlot(name string) *NodeSlot {
	if _, ok := f.InputSlots[name]; ok {
		return &NodeSlot{Name: name, IsInput: true, Node: f}
	} else if _, ok := f.OutputSlots[name]; ok {
		return &NodeSlot{Name: name, IsInput: false, Node: f}
	} else {
		return nil
	}
}

func (f *Node) Type() ObjectType {
	return NodeObjectType
}

func (f *Node) Inspect() string {
	var buffer bytes.Buffer
	buffer.WriteString("Node {Type=")
	buffer.WriteString(f.NodeType.Name)
	buffer.WriteString(",Args=[")
	args := make([]string, 0, len(f.Arguments))
	for _, arg := range f.Arguments {
		args = append(args, arg.Inspect())
	}
	buffer.WriteString(strings.Join(args, ","))
	buffer.WriteString("],InputSlots=[")
	inputs := make([]string, 0, len(f.InputSlots))
	for n, _ := range f.InputSlots {
		inputs = append(inputs, n)
	}
	buffer.WriteString(strings.Join(inputs, ","))

	buffer.WriteString("],OutputSlots=[")
	outputs := make([]string, 0, len(f.OutputSlots))
	for n, _ := range f.OutputSlots {
		outputs = append(outputs, n)
	}
	buffer.WriteString(strings.Join(outputs, ","))
	buffer.WriteByte(']')

	if f.FieldName != nil {
		buffer.WriteString(",Field=\"")
		buffer.WriteString(*f.FieldName)
		buffer.WriteByte('"')
	} else if f.TagName != nil {
		buffer.WriteString(",Tag=\"")
		buffer.WriteString(*f.TagName)
		buffer.WriteByte('"')
	}

	/*
		buffer.WriteString("],Connections = [")
		for slotname, clist := range f.connections {
			for _, c := range clist {
				buffer.WriteString("    ")
				buffer.WriteString(slotname)
				buffer.WriteString(" -> ")
				buffer.WriteString(c.Inspect())
				buffer.WriteByte('\n')
			}
		}
		buffer.WriteString("  ]\n}")
	*/
	buffer.WriteByte('}')
	return buffer.String()
}

func (f *Node) ConnectSlots(mine string, otherSlot *NodeSlot) (Object, error) {
	f.connections[mine] = append(f.connections[mine], otherSlot)
	return otherSlot.Node, nil
}

func (f *Node) Connect(other Connectable) (Object, error) {
	var otherSlot *NodeSlot
	switch t := other.(type) {
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
				f.Connect(c)
			}
		}
		return other, nil
	default:
		return nil, fmt.Errorf("cannot connect to type '%s'", other.Type())
	}

	if s := f.GetSlot("Output"); s == nil {
		return nil, fmt.Errorf("no channel 'Output' on node")
	} else {
		return f.ConnectSlots("Output", otherSlot)
	}
}

func (f *Node) Identifier(name string) (Object, error) {
	if _, ok := f.InputSlots[name]; ok {
		return &NodeSlot{Name: name, IsInput: true, Node: f}, nil
	} else if _, ok := f.OutputSlots[name]; ok {
		return &NodeSlot{Name: name, IsInput: false, Node: f}, nil
	}
	return nil, fmt.Errorf("'%s' not defined for node", name)
}
