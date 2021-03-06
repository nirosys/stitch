package analysis

import (
	"errors"
	"fmt"

	"github.com/nirosys/stitch/ast"
)

type StitchType uint

const (
	TypeUnknown  StitchType = 0
	TypeInteger  StitchType = 1
	TypeFloat    StitchType = 2
	TypeString   StitchType = 3
	TypeNode     StitchType = 4
	TypeNodeType StitchType = 5
	TypeList     StitchType = 6
	TypeMap      StitchType = 7
	TypeNodeSlot StitchType = 8
	TypeFunction StitchType = 9
	TypeBoolean  StitchType = 10
)

var typeStrings = map[StitchType]string{
	TypeUnknown:  "UNKNOWN",
	TypeInteger:  "INTEGER",
	TypeFloat:    "FLOAT",
	TypeString:   "STRING",
	TypeBoolean:  "BOOL",
	TypeNode:     "NODE",
	TypeNodeType: "NODE TYPE",
	TypeList:     "LIST",
	TypeMap:      "MAP",
	TypeNodeSlot: "SLOT",
	TypeFunction: "FUNCTION",
}

type Symbol struct {
	Name       *ast.Identifier
	Type       StitchType
	ParamTypes []StitchType // For Functions
	ReturnType StitchType   // For Functions.
	// TODO: Add origin.. File/etc Line & Column
}

var ErrSymbolExists = errors.New("symbol already exists")
var ErrTypeMismatch = errors.New("type mismatch")

type SymbolTable struct {
	symbols map[string]*Symbol
}

func (s *SymbolTable) Add(name string, sym *Symbol) error {
	if _, have := s.symbols[name]; have {
		return fmt.Errorf("%w: %s", ErrSymbolExists, name)
	} else {
		s.symbols[name] = sym
	}
	return nil
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols: make(map[string]*Symbol),
	}
}

func Analyze(tree *ast.ASTree) (*SymbolTable, error) {
	table := &SymbolTable{symbols: map[string]*Symbol{}}
	return AnalyzeWithSymbols(tree, table)
}

func AnalyzeWithSymbols(tree *ast.ASTree, table *SymbolTable) (*SymbolTable, error) {
	if table == nil {
		table = NewSymbolTable()
	}
	for _, stmt := range tree.Statements {
		if tpe, err := analyzeStatement(stmt, table); err != nil {
			return nil, err
		} else {
			fmt.Printf("ANALYSIS: %s => %s\n", stmt, typeStrings[tpe])
		}
	}
	return table, nil
}

func analyzeStatement(stmt ast.Statement, symTable *SymbolTable) (StitchType, error) {
	switch t := stmt.(type) {
	case *ast.LetStatement:
		if tpe, err := analyzeExpression(t.Value, symTable); err != nil {
			return TypeUnknown, err
		} else {
			err := symTable.Add(t.Name.String(), &Symbol{Name: t.Name, Type: tpe})
			return tpe, err
		}
	case *ast.FunctionLiteral:
		return analyzeFunctionLiteral(t, symTable)
	case ast.Expression:
		return analyzeExpression(t, symTable)
	case *ast.CommentStatement:
		return TypeUnknown, nil
	default:
		fmt.Printf("ANALYSIS: Unknown statement type\n")
	}
	return TypeUnknown, nil
}

func analyzeFunctionLiteral(fn *ast.FunctionLiteral, symTable *SymbolTable) (StitchType, error) {
	return TypeUnknown, nil
}

func analyzeExpression(exp ast.Expression, symTable *SymbolTable) (StitchType, error) {
	switch t := exp.(type) {
	case *ast.IntegerLiteral:
		return TypeInteger, nil
	case *ast.StringLiteral:
		return TypeString, nil
	case *ast.BoolLiteral:
		return TypeBoolean, nil
	case *ast.InfixExpression:
		return analyzeInfixExpression(t, symTable)
	case *ast.Identifier:
		if sym, ok := symTable.symbols[t.Identifier]; ok {
			return sym.Type, nil
		}
		return TypeUnknown, nil
	default:
		return TypeUnknown, nil
	}
}

func analyzeInfixExpression(infix *ast.InfixExpression, symTable *SymbolTable) (StitchType, error) {
	lType, err := analyzeExpression(infix.Left, symTable)
	if err != nil {
		return TypeUnknown, err
	}
	rType, err := analyzeExpression(infix.Right, symTable)
	if err != nil {
		return TypeUnknown, err
	}

	switch infix.Operator {
	case "+", "-", "/", "*":
		if lType != rType {
			return TypeUnknown, fmt.Errorf("%w: operator '%s' not defined for %s and %s", ErrTypeMismatch, infix.Operator, typeStrings[lType], typeStrings[rType])
		}
		return lType, nil
	}
	return TypeUnknown, nil
}
