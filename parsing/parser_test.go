package parsing

import (
	//	"fmt"
	"strings"
	"testing"
)

func Test_ParseSingleStatement(t *testing.T) {
	tests := []struct {
		prog       string
		statements int
	}{
		// let statement..
		{prog: "let foo = \"bar\"", statements: 1},
		// identifier access
		{prog: "snmp.foo", statements: 1},
		// function calls
		{prog: "get(\"sysDescr\")", statements: 1},
		// scoped function call
		{prog: "snmp.get(\"sysDescr\")", statements: 1},
		// function call with multiple arguments..
		{prog: "snmp.foo(\"sysDescr\", \"sysDescr\")", statements: 1},
		// multiple function calls with a connection operator
		{prog: "snmp.get(\"sysDescr\") -> snmp.get(\"foo\")", statements: 1},
		// import statement
		{prog: "import \"foo.stitch\"\n", statements: 1},
		// basic arithmetic
		{prog: "1+2", statements: 1},
		// list literal
		{prog: "[1,2,3,4]", statements: 1},
		// map literal
		{prog: "{foo = \"bar\"; bar = 2}", statements: 1},
	}

	for i, test := range tests {
		p := NewParser(strings.NewReader(test.prog))
		tree := p.Parse()
		if len(tree.Statements) != test.statements {
			t.Errorf("[%d] unexpected number of statements: %d != %d", i, len(tree.Statements), test.statements)
		}
		//if len(tree.Statements) > 0 {
		//	s := tree.Statements[0].String()
		//	if s != test.prog {
		//		t.Errorf("[%d] statement mismatch: '%s' != '%s'", i, s, test.prog)
		//		for j, s := range tree.Statements {
		//			fmt.Printf("[%d/%d] %T\n", i, j, s)
		//		}
		//	}
		//}
	}
}

func Test_ParseMultipleStatements(t *testing.T) {
	tests := []struct {
		prog       string
		statements int
	}{
		// let statement
		{prog: "let foo = \"bar\";\nsnmp.get(\"sysDescr\");\n", statements: 2},
		{prog: "{\n  let foo = \"bar\";\n}\n;\n", statements: 1},
		{prog: "let foo = {\n  snmp.get(\"sysDescr\");\n}\n;\n", statements: 1},
	}
	for i, test := range tests {
		p := NewParser(strings.NewReader(test.prog))
		prog := p.Parse()
		if len(prog.Statements) != test.statements {
			t.Errorf("[%d] unexpected number of statements: %d != %d", i, len(prog.Statements), test.statements)
		}
		//progStr := prog.String()
		//if progStr != test.prog {
		//	t.Errorf("[%d] Program mismatch:\nParsed ===\n%s\n==== Test\n%s\n===", i, progStr, test.prog)
		//	for j, s := range prog.Statements {
		//		fmt.Printf("[%d/%d] %T\n", i, j, s)
		//	}
		//}
	}
}

func Test_ParseWithErrors(t *testing.T) {
	tests := []struct {
		prog       string
		statements int
		err        error
	}{
		// let statement
		{prog: "foreach i in [1, 2, 3, 4] {", statements: 0},
	}
	for i, test := range tests {
		p := NewParser(strings.NewReader(test.prog))
		prog := p.Parse()
		if len(prog.Statements) != test.statements {
			t.Errorf("[%d] unexpected number of statements: %d != %d", i, len(prog.Statements), test.statements)
		}
		//progStr := prog.String()
		//if progStr != test.prog {
		//	t.Errorf("[%d] Program mismatch:\nParsed ===\n%s\n==== Test\n%s\n===", i, progStr, test.prog)
		//	for j, s := range prog.Statements {
		//		fmt.Printf("[%d/%d] %T\n", i, j, s)
		//	}
		//}
	}
}
