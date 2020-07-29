package eval

import (
	"fmt"

	"github.com/nirosys/stitch"
	"github.com/nirosys/stitch/ast"
	"github.com/nirosys/stitch/object"

	"github.com/nirosys/gaufre/graph"
)

// Evaluator //////////////////////////////////////////////////////////////////
type Evaluator struct {
	Resolver ObjectResolver
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		Resolver: nil,
	}
}

func (e *Evaluator) evalNamedNode(i *ast.NamedNodeExpression, env *object.Environment) (object.Object, error) {
	// Valid uses:
	//   left is identifier, right is Node
	//   left is tag identifier, right is Node
	rightObj, err := e.eval(i.Expression, env)
	if err != nil {
		return nil, err
	}

	if node, ok := rightObj.(*object.Node); !ok {
		return nil, fmt.Errorf("cannot assign tag or field to type %s", rightObj.Type())
	} else if i.FieldName != nil {
		name := i.FieldName.String()
		node.FieldName = &name
		return rightObj, nil
	} else if i.TagName != nil {
		name := i.TagName.String()
		node.TagName = &name
		return rightObj, nil
	}
	return nil, fmt.Errorf("invalid named expression")
}

func (e *Evaluator) evalConditional(c *ast.ConditionalExpression, env *object.Environment) (object.Object, error) {
	cond, err := e.eval(c.Condition, env)
	if err != nil {
		return nil, err
	}
	if b, ok := cond.(*object.BoolObject); !ok {
		return nil, fmt.Errorf("expected boolean expression, found %s", cond.Type())
	} else if bool(*b) {
		return e.eval(c.Block, env)
	} else if c.Else != nil {
		return e.eval(c.Else, env)
	}
	return nil, nil
}

func (e *Evaluator) evalInternalFunc(l *ast.InternalExpression, env *object.Environment) (object.Object, error) {
	if obj, err := e.Resolver.Resolve(l.Name.Value); err != nil {
		return nil, err
	} else {
		switch t := obj.(type) {
		case *object.InternalFunction:
			t.Env = env
		}
		return obj, nil
	}
}

func (e *Evaluator) evalList(l *ast.ListLiteral, env *object.Environment) (object.Object, error) {
	list := &object.List{}
	var tpe object.ObjectType = object.UnknownObjectType
	contents := make([]object.Object, 0, len(l.Contents))

	for _, exp := range l.Contents {
		if obj, err := e.eval(exp, env); err != nil {
			return nil, err
		} else {
			if tpe == object.UnknownObjectType {
				tpe = obj.Type()
			} else if tpe != obj.Type() {
				return nil, fmt.Errorf("mixed types for list")
			}
			contents = append(contents, obj)
		}
	}
	list.InnerType = tpe
	list.Contents = contents
	return list, nil
}

func (e *Evaluator) evalMap(m *ast.MapLiteral, env *object.Environment) (object.Object, error) {
	mapObj := &object.MapObject{}
	mapObj.Fields = make(map[string]object.Object)
	for _, assign := range m.Assignments {
		if obj, err := e.eval(assign.Value, env); err != nil {
			return nil, err
		} else {
			mapObj.Fields[assign.Identifier.String()] = obj
		}
	}
	return mapObj, nil
}

func (e *Evaluator) evalBlockExpression(b *ast.BlockExpression, env *object.Environment) (object.Object, error) {
	var last object.Object
	for _, stmt := range b.Statements {
		if obj, err := e.eval(stmt, env); err != nil {
			return nil, err
		} else {
			last = obj
		}
	}

	return last, nil
}

func (e *Evaluator) evalFunctionDefinition(fun *ast.FunctionLiteral, env *object.Environment) (object.Object, error) {
	fn := &object.Function{
		Parameters: fun.Parameters,
		Body:       fun.Body,
		Env:        env, // Parent scope, we capture whatever is around us.. allows for nested funcs...
	}
	if fun.Identifier != nil {
		env.Put(fun.Identifier.String(), fn)
		return nil, nil
	} else {
		return fn, nil
	}
}

func (e *Evaluator) evalNotExpression(in *ast.NotExpression, env *object.Environment) (object.Object, error) {
	exp, err := e.eval(in.Expression, env)
	if err != nil {
		return nil, err
	}

	if boolExp, ok := exp.(*object.BoolObject); !ok {
		return nil, fmt.Errorf("'!' operator not defined for '%s'", exp.Type())
	} else {
		return boolExp.Not(), nil
	}
}

func (e *Evaluator) evalInfixComparison(in *ast.InfixExpression, env *object.Environment) (object.Object, error) {
	leftObj, err := e.eval(in.Left, env)
	if err != nil {
		return nil, err
	}

	rightObj, err := e.eval(in.Right, env)
	if err != nil {
		return nil, err
	}

	left, lok := leftObj.(object.Comparable)
	if !lok {
		return nil, fmt.Errorf("%s is not comparable", leftObj.Type())
	}
	right, rok := rightObj.(object.Comparable)
	if !rok {
		return nil, fmt.Errorf("%s is not comparable", rightObj.Type())
	}

	if !left.IsComparable(right) {
		return nil, fmt.Errorf("%s is not comparable to %s", left.Type(), right.Type())
	}

	switch in.Operator {
	case "==":
		if v, err := left.Equals(right); err != nil {
			return nil, err
		} else {
			return object.NewBoolObject(v), nil
		}
	case ">=": // Same as !(l < )
		if v, err := left.LessThan(right); err != nil {
			return nil, err
		} else {
			return object.NewBoolObject(!v), nil
		}
	case "<=": // Same as !(l > )
		if v, err := left.GreaterThan(right); err != nil {
			return nil, err
		} else {
			return object.NewBoolObject(!v), nil
		}
	case ">":
		if v, err := left.GreaterThan(right); err != nil {
			return nil, err
		} else {
			return object.NewBoolObject(v), nil
		}
	case "<":
		if v, err := left.LessThan(right); err != nil {
			return nil, err
		} else {
			return object.NewBoolObject(v), nil
		}
	case "!=":
		if v, err := left.Equals(right); err != nil {
			return nil, err
		} else {
			return object.NewBoolObject(!v), nil
		}
	case "and":
		lbool, lok := left.(*object.BoolObject)
		if !lok {
			return nil, fmt.Errorf("expected bool expression, found: %s", left.Type())
		}
		rbool, rok := right.(*object.BoolObject)
		if !rok {
			return nil, fmt.Errorf("expected bool expression, found: %s", right.Type())
		}
		return object.NewBoolObject(bool(*lbool) && bool(*rbool)), nil
	case "or":
		lbool, lok := left.(*object.BoolObject)
		if !lok {
			return nil, fmt.Errorf("expected bool expression, found: %s", left.Type())
		}
		rbool, rok := right.(*object.BoolObject)
		if !rok {
			return nil, fmt.Errorf("expected bool expression, found: %s", right.Type())
		}
		return object.NewBoolObject(bool(*lbool) || bool(*rbool)), nil
	}

	return nil, nil
}

func (e *Evaluator) evalInfixComputation(in *ast.InfixExpression, env *object.Environment) (object.Object, error) {
	leftObj, err := e.eval(in.Left, env)
	if err != nil {
		return nil, err
	}

	if in.Operator == "." {
		if i, ok := in.Right.(*ast.Identifier); !ok {
			return nil, fmt.Errorf("expected identifier but found '%s'", in.Right.String())
		} else {
			return leftObj.Identifier(i.Identifier)
		}
	}

	rightObj, err := e.eval(in.Right, env)
	if err != nil {
		return nil, err
	}

	if lStr, ok := leftObj.(*object.String); ok && rightObj.Type().IsPrimitive() {
		switch r := rightObj.(type) {
		case *object.String:
			return &object.String{Value: lStr.Value + r.Value}, nil
		default:
			return &object.String{Value: lStr.Value + r.Inspect()}, nil
		}
	}

	var left object.Computable
	if l, ok := leftObj.(object.Computable); !ok {
		return nil, fmt.Errorf("operator '%s' not defined for type %s", in.Operator, leftObj.Type())
	} else {
		left = l
	}

	var right object.Computable
	if r, ok := rightObj.(object.Computable); !ok {
		return nil, fmt.Errorf("operator '%s' not defined for type %s", in.Operator, rightObj.Type())
	} else {
		right = r
	}

	switch in.Operator {
	case "+":
		return left.Add(right)
	case "-":
		return left.Subtract(right)
	case "*":
		return left.Multiply(right)
	case "/":
		return left.Divide(right)
	case "%":
		return left.Modulus(right)
	default:
		return nil, fmt.Errorf("'%s' operator not implemented yet", in.Operator)
	}
}

func (e *Evaluator) evalExpressions(expr []ast.Expression, env *object.Environment) ([]object.Object, error) {
	list := make([]object.Object, len(expr), len(expr))

	for i, exp := range expr {
		if ret, err := e.eval(exp, env); err != nil {
			return nil, err
		} else {
			list[i] = ret
		}
	}

	return list, nil
}

func (e *Evaluator) evalNodeStatement(n *ast.NodeStatement, env *object.Environment) (object.Object, error) {
	nodeType := &object.NodeType{}
	nodeType.Name = n.Identifier.String()

	literal := n.Literal
	inputs := make([]string, 0, len(literal.InputSlots))
	for _, slot := range literal.InputSlots {
		inputs = append(inputs, slot.Identifier.String())
	}
	nodeType.InputSlots = inputs

	outputs := make([]string, 0, len(literal.OutputSlots))
	for _, slot := range literal.OutputSlots {
		outputs = append(outputs, slot.Identifier.String())
	}
	nodeType.OutputSlots = outputs

	nodeType.NodeArgs = literal.Arguments

	nodeType.Body = literal.Block
	nodeType.Env = env

	env.Put(nodeType.Name, nodeType)

	return nodeType, nil
}

func (e *Evaluator) evalForeachStatement(f *ast.ForeachStatement, env *object.Environment) (object.Object, error) {
	if obj, err := e.eval(f.List, env); err != nil {
		return nil, err
	} else if obj.Type() != object.ListObjectType {
		return nil, fmt.Errorf("invalid type, expected a LIST, got a %s", obj.Type())
	} else if list, ok := obj.(*object.List); ok {
		scope := env.Clone()
		for _, o := range list.Contents {
			scope.Put(f.LoopVar.Identifier, o)
			e.evalBlockExpression(f.Block, scope)
		}
	}
	return nil, nil
}

func (e *Evaluator) evalNodeLiteral(n *ast.NodeLiteral, env *object.Environment) (object.Object, error) {
	// TODO: Fix this.
	obj := object.NewNode()
	//for _, assign := range n.Properties {
	//	ident := assign.Identifier.String()
	//	switch ident {
	//	case "type":
	//		if tpe, ok := assign.Value.(*ast.StringLiteral); ok {
	//			obj.NodeType = tpe.Value
	//		} else {
	//			return nil, fmt.Errorf("invalid value for node type")
	//		}
	//	case "arg":
	//		if arg, ok := assign.Value.(*ast.StringLiteral); ok {
	//			obj.Argument = arg.Value
	//		} else {
	//			return nil, fmt.Errorf("invalid value for node arg")
	//		}

	//	default:
	//		return nil, fmt.Errorf("unknown node property: '%s'", ident)
	//	}
	//}
	env.PutUnboundNode(obj) // We're a literal, so we start out unbound..

	return obj, nil
}

func (e *Evaluator) evalConnectExpression(c *ast.ArrowExpression, env *object.Environment) (object.Object, error) {
	var right object.Connectable
	if r, err := e.eval(c.Right, env); err != nil {
		return nil, err
	} else if t, ok := r.(object.Connectable); !ok {
		return nil, fmt.Errorf("connections can not be with type %s", r.Type())
	} else {
		right = t
	}

	var left object.Connectable
	if l, err := e.eval(c.Left, env); err != nil {
		return nil, err
	} else if t, ok := l.(object.Connectable); !ok {
		return nil, fmt.Errorf("connections can not be with type %s", l.Type())
	} else {
		left = t
	}

	_, err := left.Connect(right)

	return left, err
}

func (e *Evaluator) eval(n ast.Node, env *object.Environment) (object.Object, error) {
	//fmt.Printf("stmt: %#v\n", n)
	switch t := n.(type) {
	case *ast.CommentStatement:
		return nil, nil // Do nothing..
	case *ast.ConditionalExpression:
		return e.evalConditional(t, env)
		/*
			case *ast.ImportStatement:
				// TODO: Track directory, where import should base out of
				if f, err := os.Open(t.Path); err != nil {
					return nil, err
				} else {
					defer f.Close()

					p := path.Base(t.Path)
					splitIdx := strings.IndexFunc(p, func(r rune) bool {
						return !(unicode.IsLetter(r) || unicode.IsDigit(r))
					})
					name := p[:splitIdx]
					if len(name) < 1 {
						return nil, fmt.Errorf("invalid package name: '%s'", name)
					}
					parser := parsing.NewParser(f)
					prog := parser.ParseProgram()
					if prog == nil {
						//fmt.Printf("Parser error(s):\n")
						for _, e := range parser.Errors() {
							fmt.Printf("%s: %s\n", p, e)
						}
						return nil, fmt.Errorf("parser error(s), check stdout for errors")
					} else {
						pkgEnv := object.NewEnvironment()
						if _, err := e.EvalProgram(prog, pkgEnv); err != nil {
							return nil, err
						} else {
							pkg, err := object.NewPackage(name, pkgEnv)
							if err != nil {
								return nil, err
							}
							return nil, env.PutPackage(pkg)
						}
					}
				}
				return nil, nil // TODO: Implement me.
		*/
	case *ast.LetStatement:
		if obj, err := e.eval(t.Value, env); err == nil {
			env.Put(t.Name.String(), obj)
			return nil, nil
		} else {
			return nil, err
		}
	case *ast.Identifier:
		if obj, ok := env.Get(t.String()); ok {
			return obj, nil
		}
		return nil, fmt.Errorf("%d: unknown identifier '%s'", t.Token.Position.Line, t.String())
	case *ast.CallExpression:
		obj, err := e.eval(t.Function, env) // TODO: rename function 'identifier'
		if err != nil {
			return nil, err
		}

		switch tpe := obj.(type) {
		case object.Constructable:
			params := tpe.Arguments()
			if len(t.Arguments) != len(params) {
				return nil, fmt.Errorf("expected %d arguments but found %d", len(params), len(t.Arguments))
			}
			if args, err := e.evalExpressions(t.Arguments, env); err != nil {
				return nil, err
			} else {
				if obj, err := tpe.Construct(args); err != nil {
					return nil, err
				} else {
					if obj.Type() == object.NodeObjectType {
						env.PutUnboundNode(obj)
					}
					return obj, nil
				}
			}
		case object.Callable:
			params := tpe.FuncParameters()
			if len(t.Arguments) != len(params) {
				return nil, fmt.Errorf("expected %d arguments but found %d", len(params), len(t.Arguments))
			}
			if args, err := e.evalExpressions(t.Arguments, env); err != nil {
				return nil, err
			} else {
				return e.applyFunction(env, tpe, args)
			}
		default:
			return nil, fmt.Errorf("'%s' is not a function", obj.Type())
		}
	case *ast.AssignmentExpression:
		target := t.Identifier.Identifier
		if obj, err := e.eval(t.Value, env); err != nil {
			return nil, err
		} else if _, have := env.Get(target); !have {
			return nil, fmt.Errorf("unknown identifier '%s'", target)
		} else {
			env.Put(target, obj)
			return obj, nil
		}
	case *ast.NodeStatement:
		return e.evalNodeStatement(t, env)
	case *ast.ForeachStatement:
		return e.evalForeachStatement(t, env)
	case *ast.NodeLiteral:
		return e.evalNodeLiteral(t, env)
	case *ast.StringLiteral:
		return &object.String{Value: t.Value}, nil
	case *ast.BoolLiteral:
		obj := object.BoolObject(t.Value)
		return &obj, nil
	case *ast.IntegerLiteral:
		return &object.Integer{Value: t.Value}, nil
	case *ast.InfixExpression:
		switch t.Operator {
		case "==", "<", "<=", ">", ">=", "!=", "and", "or":
			return e.evalInfixComparison(t, env)
		default:
			return e.evalInfixComputation(t, env)
		}
	case *ast.NamedNodeExpression:
		return e.evalNamedNode(t, env)
	case *ast.ArrowExpression:
		return e.evalConnectExpression(t, env)
	case *ast.FunctionLiteral:
		return e.evalFunctionDefinition(t, env)
	case *ast.BlockExpression:
		return e.evalBlockExpression(t, env)
	case *ast.ListLiteral:
		return e.evalList(t, env)
	case *ast.MapLiteral:
		return e.evalMap(t, env)
	case *ast.InternalExpression:
		return e.evalInternalFunc(t, env)
	case *ast.NotExpression:
		return e.evalNotExpression(t, env)
	default:
		return nil, fmt.Errorf("unknown node type: %T", n)
	}
}

func (e *Evaluator) applyFunction(scope *object.Environment, callable object.Callable, args []object.Object) (object.Object, error) {
	var retObj object.Object
	var err error

	switch fn := callable.(type) { // Allow for other callables
	case *object.Function:
		env := extendFunctionEnv(fn, args)
		retObj, err = e.eval(fn.Body, env)
	case *object.InternalFunction:
		retObj, err = fn.Fn.Fn(fn.Env, args)
	default:
		return nil, fmt.Errorf("not a function")
	}

	if err != nil {
		return nil, err
	} else if retObj != nil {
		// If we're returning a node, we need to start tracking it.
		// If we're not in the global environment, we won't track this obj until
		// it bubbles up to it.
		if retObj.Type() == object.NodeObjectType {
			scope.PutUnboundNode(retObj)
		}
		return retObj, nil
	}
	return nil, nil
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := fn.Env.Clone()
	params := fn.FuncParameters()
	for i, p := range args {
		env.Put(params[i].Identifier.String(), p)
	}
	return env
}

func (e *Evaluator) EvalProgram(prog *stitch.Program, env *object.Environment) (object.Object, error) {
	var obj object.Object
	for _, stmt := range prog.Tree.Statements {
		if o, err := e.eval(stmt, env); err != nil {
			return nil, err
		} else {
			obj = o
		}
	}
	return obj, nil
}

/*
func (e *Evaluator) Compile(prog *ast.Program) (*graph.Graph, error) {
	env := object.NewEnvironment()
	g := graph.NewGraph("test")

	if obj, err := e.EvalProgram(prog, env); err != nil {
		return nil, err
	} else if root, ok := obj.(*object.Node); !ok {
		return nil, fmt.Errorf("must return root node.")
	} else {
		// Our ENV contains everything that we need to build the gaufre graph.
		visited := make(map[*object.Node]int)
		if _, err := visitNode(root, visited, g); err != nil {
			return nil, err
		}
	}

	return g, nil
}
*/

func (e *Evaluator) CompileObject(node *object.Node) (*graph.Graph, error) {
	visited := make(map[*object.Node]int)
	g := graph.NewGraph("test")

	if _, err := visitNode(node, visited, g); err != nil {
		return nil, err
	}
	return g, nil
}

func visitNode(n *object.Node, visited map[*object.Node]int, g *graph.Graph) (int, error) {
	var startNodeId int
	var startNode graph.Node

	if nodeId, done := visited[n]; done {
		return nodeId, nil
	} else {
		startNodeId = len(visited)
		visited[n] = startNodeId

		startNode = graph.Node{
			ID:   uint(startNodeId),
			Name: "",
			Type: n.NodeType.Name,
		}
		startNode.Inputs = graph.Input{
			ID:   0,
			Name: "Input", // TODO: Gaufre needs to support multiple inputs.
		}
		config := map[string]interface{}{}
		args := map[string]interface{}{}
		for i, arg := range n.NodeType.NodeArgs {
			ident := arg.Identifier.Identifier
			obj := n.Arguments[i]
			switch t := obj.(type) {
			case *object.Integer:
				args[ident] = t.Value
			case *object.String:
				args[ident] = t.Value
			case *object.BoolObject:
				args[ident] = bool(*t)
			default:
				return 0, fmt.Errorf("%s not supported as node argument", obj.Type())
			}
		}

		if n.FieldName != nil {
			config["tag"] = *n.FieldName
		}

		config["args"] = args
		startNode.Configuration = graph.NewNodeConfig(config)

		for i, n := range n.NodeType.OutputSlots {
			startNode.Outputs = append(startNode.Outputs, graph.Output{
				ID:   uint(i),
				Name: n,
			})
		}

		g.AddNode(startNode)
	}
	conns := n.GetConnections()
	for _, conn := range conns {
		if endNodeId, err := visitNode(conn.End.Node, visited, g); err != nil {
			return 0, err
		} else {
			if endNode, err := g.NodeById(uint(endNodeId)); err != nil {
				return 0, err
			} else {
				endRef, _ := endNode.OutputRefByName(conn.End.Name)

				slotId, _ := endRef.SocketId()
				g.Connect(&startNode, 0, endNode, slotId)
			}
		}
	}

	return startNodeId, nil
}
