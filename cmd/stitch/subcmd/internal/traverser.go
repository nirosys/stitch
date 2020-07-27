package internal

import (
	"fmt"
	"io"
	"strings"

	"github.com/nirosys/stitch/object"

	"github.com/emicklei/dot"
)

type Traverser struct {
	Environment *object.Environment
}

func NewTraverser(env *object.Environment) *Traverser {
	return &Traverser{Environment: env}
}

func (t *Traverser) visit(n *object.Node, visited map[*object.Node]int, g *dot.Graph) (int, error) {
	var startNode int
	if dotNode, done := visited[n]; done {
		return dotNode, nil
	} else {
		inputs := make([]string, 0, len(n.InputSlots))
		for k, _ := range n.InputSlots {
			inputs = append(inputs, "<"+k+"> "+k)
		}
		outputs := make([]string, 0, len(n.OutputSlots))
		for k, _ := range n.OutputSlots {
			outputs = append(outputs, "<"+k+"> "+k)
		}
		nodestr := n.NodeType.Name
		args := make([]string, 0, len(n.Arguments))
		for i, fp := range n.NodeType.NodeArgs {
			name := fp.Identifier.String()
			value := strings.ReplaceAll(strings.ReplaceAll(n.Arguments[i].Inspect(), "{{", `%%`), "}}", "%%")
			args = append(args, name+"="+value)
		}
		tag_field := ""
		if n.TagName != nil {
			tag_field = fmt.Sprintf("\nTag=%q", *n.TagName)
		} else if n.FieldName != nil {
			tag_field = fmt.Sprintf("\nField=%q", *n.FieldName)
		}

		label := fmt.Sprintf("{{%s}|%s\n%s%s|{%s}}",
			strings.Join(inputs, "|"),
			nodestr,
			tag_field,
			strings.Join(args, "\n"),
			strings.Join(outputs, "|"),
		)

		startNode = len(visited)
		fmt.Printf("   n%d[label=%q, shape=\"Mrecord\"];\n", startNode, label)
		//startNode = g.Node(label).Attr("shape", "record")
		visited[n] = startNode
	}

	conns := n.GetConnections()
	for _, conn := range conns {
		if endNode, err := t.visit(conn.End.Node, visited, g); err != nil {
			return 0, err
		} else {
			fmt.Printf("   n%d:%s -> n%d:%s;\n", startNode, conn.Start.Name, endNode, conn.End.Name)
			//g.Edge(startNode, endNode).Attr("taillabel", conn.Start.Name).Attr("headlabel", conn.End.Name)
		}
	}

	return startNode, nil
}

func (t *Traverser) RenderDotAll(w io.Writer) error {
	g := dot.NewGraph(dot.Directed)
	visited := map[*object.Node]int{}

	fmt.Printf("digraph {\n   rankdir=LR;\n")

	names := t.Environment.GetNames()
	names = append(names, t.Environment.GetUnboundNodes()...)

	for _, ident := range names {
		if obj, has := t.Environment.Get(ident); !has {
			return fmt.Errorf("identifier not found '%s'", ident)
		} else if node, ok := obj.(*object.Node); ok {
			t.visit(node, visited, g)
		}
	}
	fmt.Printf("\n}\n")
	//fmt.Fprintf(w, "%s\n", g.String())
	return nil
}

//func (t *Traverser) RenderDot(w io.Writer, ident string) error {
//	g := dot.NewGraph(dot.Directed)
//	visited := map[*object.Node]dot.Node{}
//
//	if obj, has := t.Environment.Get(ident); !has {
//		return fmt.Errorf("identifier not found '%s'", ident)
//	} else if node, ok := obj.(*object.Node); ok {
//		t.visit(node, visited, g)
//	} else {
//		return fmt.Errorf("'%s' is not a NODE", ident)
//	}
//
//	fmt.Fprintf(w, "%s\n", g.String())
//
//	return nil
//}
