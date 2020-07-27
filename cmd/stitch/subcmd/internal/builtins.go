package internal

import (
	"fmt"

	"github.com/nirosys/stitch/ast"
	"github.com/nirosys/stitch/object"
)

var HostedFuncs = map[string]*object.NativeFunction{
	"std:println": &object.NativeFunction{
		Fn: builtin_println,
		Params: []*ast.FunctionParameter{
			&ast.FunctionParameter{Identifier: &ast.Identifier{Identifier: "msg"}},
		},
		//Params: []object.ObjectType{object.StringObjectType},
	}}

var HostedNodeTypes = map[string]*object.NodeType{
	"snmp:get": &object.NodeType{
		Name: "snmp:get", NodeArgs: []*ast.FunctionParameter{
			&ast.FunctionParameter{
				Identifier: &ast.Identifier{Identifier: "oid"},
			},
		},
		InputSlots:  []string{"Input"},
		OutputSlots: []string{"Output", "Error", "Missing"},
	},
	"snmp:walk": &object.NodeType{
		Name: "snmp:walk",
		NodeArgs: []*ast.FunctionParameter{
			&ast.FunctionParameter{
				Identifier: &ast.Identifier{Identifier: "oid"},
			},
		},
		InputSlots:  []string{"Input"},
		OutputSlots: []string{"Output", "Error"},
	},
	"std:passthru": &object.NodeType{
		Name:        "std:passthru",
		NodeArgs:    []*ast.FunctionParameter{},
		InputSlots:  []string{"Input"},
		OutputSlots: []string{"Output"},
	},
}

type Resolver struct{}

func NewResolver() *Resolver {
	return &Resolver{}
}

func (r *Resolver) Resolve(name string) (object.Object, error) {
	if fn, ok := HostedFuncs[name]; ok {
		return &object.InternalFunction{Fn: fn, Env: nil}, nil
	} else if nodeType, ok := HostedNodeTypes[name]; ok {
		return nodeType, nil
	} else {
		return nil, fmt.Errorf("unknown internal \"%s\"", name)
	}
}

// TODO: Need to use reflection to handle these..
func builtin_println(env *object.Environment, args []object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments")
	}
	msg := args[0]
	if msg.Type() != object.StringObjectType {
		return nil, fmt.Errorf("type mismatch: expected STRING, received %s", msg.Type())
	}
	s, _ := msg.(*object.String)

	fmt.Printf("%s\n", s.Value)
	return nil, nil
}

func builtin_snmpget(env *object.Environment, args []object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments")
	}
	msg := args[0]
	if msg.Type() != object.StringObjectType {
		return nil, fmt.Errorf("type mismatch: expected STRING, received %s", msg.Type())
	}

	node := object.NewNode()
	node.Arguments = args
	node.InputSlots = map[string]struct{}{"Input": struct{}{}}
	node.OutputSlots = map[string]struct{}{"Output": struct{}{}, "Error": struct{}{}}

	return node, nil
}

func builtin_snmpwalk(env *object.Environment, args []object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments")
	}
	msg := args[0]
	if msg.Type() != object.StringObjectType {
		return nil, fmt.Errorf("type mismatch: expected STRING, received %s", msg.Type())
	}

	node := object.NewNode()
	node.Arguments = args
	node.InputSlots = map[string]struct{}{"Input": struct{}{}}
	node.OutputSlots = map[string]struct{}{"Output": struct{}{}, "Error": struct{}{}}

	return node, nil
}
