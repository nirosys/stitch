package stitch

import (
	"bytes"
	"fmt"
	//	"encoding/json"
	"io"
	//	"os"

	"github.com/nirosys/stitch/analysis"
	"github.com/nirosys/stitch/ast"
	"github.com/nirosys/stitch/parsing"
)

const stitchVersion = "v0.0.1"

func Version() string {
	return stitchVersion
}

// Program ////////////////////////////////////////////////////////////////////
type Program struct {
	Tree    *ast.ASTree
	Symbols *analysis.SymbolTable
}

func NewProgram(r io.Reader) *Program {
	parser := parsing.NewParser(r)
	tree := parser.Parse()

	prog := &Program{Tree: tree}
	if symbols, err := analysis.Analyze(tree); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	} else {
		prog.Symbols = symbols
	}

	//enc := json.NewEncoder(os.Stdout)
	//enc.SetIndent("", "  ")
	//enc.Encode(tree)

	return prog
}

func ExtendProgram(prog *Program, r io.Reader) (*Program, error) {
	parser := parsing.NewParser(r)
	tree := parser.Parse()

	if symbols, err := analysis.AnalyzeWithSymbols(tree, prog.Symbols); err == nil {
		if prog.Tree != nil {
			tree.Statements = append(prog.Tree.Statements, tree.Statements...)
		}
		newProg := &Program{
			Tree:    tree,
			Symbols: symbols,
		}
		return newProg, nil
	} else {
		return nil, err
	}
}

func (p *Program) String() string {
	var buffer bytes.Buffer
	stmts := p.Tree.Statements
	for _, stmt := range stmts {
		buffer.WriteString(stmt.String())
		buffer.WriteString(";\n")
	}
	return buffer.String()
}

func (p *Program) Errors() []string {
	return []string{}
}
