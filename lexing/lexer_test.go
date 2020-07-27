package lexing

import (
	"fmt"
	"strings"
	"testing"
)

func Test_Lexer(t *testing.T) {
	tests := []struct {
		in     string
		tokens []Token
	}{
		{"100", []Token{
			Token{Text: "100", Position: Position{Line: 0, Column: 0}, Type: L_INTEGER},
		}},
		{"\n\n100", []Token{
			Token{Text: "100", Position: Position{Line: 2, Column: 0}, Type: L_INTEGER},
		}},
		{"10.5", []Token{
			Token{Text: "10.5", Position: Position{Line: 0, Column: 0}, Type: L_FLOAT},
		}},
		{"hello", []Token{
			Token{Text: "hello", Position: Position{Line: 0, Column: 0}, Type: IDENT},
		}},
		{"fn foo", []Token{
			Token{Text: "fn", Position: Position{Line: 0, Column: 0}, Type: K_FUNCTION},
			Token{Text: "foo", Position: Position{Line: 0, Column: 3}, Type: IDENT},
		}},
		{"fn foo(): string { }", []Token{
			Token{Text: "fn", Position: Position{Line: 0, Column: 0}, Type: K_FUNCTION},
			Token{Text: "foo", Position: Position{Line: 0, Column: 3}, Type: IDENT},
			Token{Text: "(", Position: Position{Line: 0, Column: 6}, Type: D_LPARENTH},
			Token{Text: "(", Position: Position{Line: 0, Column: 7}, Type: D_RPARENTH},
			Token{Text: ":", Position: Position{Line: 0, Column: 8}, Type: O_COLON},
			Token{Text: "string", Position: Position{Line: 0, Column: 10}, Type: IDENT},
			Token{Text: "{", Position: Position{Line: 0, Column: 17}, Type: D_LBRACE},
			Token{Text: "}", Position: Position{Line: 0, Column: 19}, Type: D_RBRACE},
		}},
		{"node", []Token{
			Token{Text: "node", Position: Position{Line: 0, Column: 0}, Type: K_NODE},
		}},
		{"let foo = bar", []Token{
			Token{Text: "let", Position: Position{Line: 0, Column: 0}, Type: K_LET},
			Token{Text: "foo", Position: Position{Line: 0, Column: 4}, Type: IDENT},
			Token{Text: "=", Position: Position{Line: 0, Column: 8}, Type: O_ASSIGN},
			Token{Text: "bar", Position: Position{Line: 0, Column: 10}, Type: IDENT},
		}},
		{"foo == bar", []Token{
			Token{Text: "foo", Position: Position{Line: 0, Column: 0}, Type: IDENT},
			Token{Text: "==", Position: Position{Line: 0, Column: 4}, Type: O_EQ},
			Token{Text: "bar", Position: Position{Line: 0, Column: 7}, Type: IDENT},
		}},

		{"foo:snmp.get(\"foo\")", []Token{
			Token{Text: "foo", Position: Position{Line: 0, Column: 0}, Type: IDENT},
			Token{Text: ":", Position: Position{Line: 0, Column: 3}, Type: O_COLON},
			Token{Text: "snmp", Position: Position{Line: 0, Column: 4}, Type: IDENT},
			Token{Text: ".", Position: Position{Line: 0, Column: 8}, Type: O_DOT},
			Token{Text: "get", Position: Position{Line: 0, Column: 9}, Type: IDENT},
			Token{Text: "(", Position: Position{Line: 0, Column: 12}, Type: D_LPARENTH},
			Token{Text: "foo", Position: Position{Line: 0, Column: 13}, Type: L_STRING},
			Token{Text: ")", Position: Position{Line: 0, Column: 18}, Type: D_RPARENTH},
		}},
		{"# This is a comment", []Token{
			Token{Text: " This is a comment", Position: Position{Line: 0, Column: 0}, Type: COMMENT},
		}},
		{"foo() -> bar()", []Token{
			Token{Text: "foo", Position: Position{Line: 0, Column: 0}, Type: IDENT},
			Token{Text: "(", Position: Position{Line: 0, Column: 3}, Type: D_LPARENTH},
			Token{Text: ")", Position: Position{Line: 0, Column: 4}, Type: D_RPARENTH},
			Token{Text: "->", Position: Position{Line: 0, Column: 6}, Type: O_ARROW},
			Token{Text: "bar", Position: Position{Line: 0, Column: 9}, Type: IDENT},
			Token{Text: "(", Position: Position{Line: 0, Column: 12}, Type: D_LPARENTH},
			Token{Text: ")", Position: Position{Line: 0, Column: 13}, Type: D_RPARENTH},
		}},
		{"foo() -> { bar() }", []Token{
			Token{Text: "foo", Position: Position{Line: 0, Column: 0}, Type: IDENT},
			Token{Text: "(", Position: Position{Line: 0, Column: 3}, Type: D_LPARENTH},
			Token{Text: ")", Position: Position{Line: 0, Column: 4}, Type: D_RPARENTH},
			Token{Text: "->", Position: Position{Line: 0, Column: 6}, Type: O_ARROW},
			Token{Text: "{", Position: Position{Line: 0, Column: 9}, Type: D_LBRACE},
			Token{Text: "bar", Position: Position{Line: 0, Column: 11}, Type: IDENT},
			Token{Text: "(", Position: Position{Line: 0, Column: 14}, Type: D_LPARENTH},
			Token{Text: ")", Position: Position{Line: 0, Column: 15}, Type: D_RPARENTH},
			Token{Text: "}", Position: Position{Line: 0, Column: 17}, Type: D_RBRACE},
		}},
		{"foo() -> {\n bar()\n}", []Token{
			Token{Text: "foo", Position: Position{Line: 0, Column: 0}, Type: IDENT},
			Token{Text: "(", Position: Position{Line: 0, Column: 3}, Type: D_LPARENTH},
			Token{Text: ")", Position: Position{Line: 0, Column: 4}, Type: D_RPARENTH},
			Token{Text: "->", Position: Position{Line: 0, Column: 6}, Type: O_ARROW},
			Token{Text: "{", Position: Position{Line: 0, Column: 9}, Type: D_LBRACE},
			Token{Text: "bar", Position: Position{Line: 1, Column: 1}, Type: IDENT},
			Token{Text: "(", Position: Position{Line: 1, Column: 4}, Type: D_LPARENTH},
			Token{Text: ")", Position: Position{Line: 1, Column: 5}, Type: D_RPARENTH},
			Token{Text: "}", Position: Position{Line: 2, Column: 0}, Type: D_RBRACE},
		}},
		{"true", []Token{
			Token{Text: "true", Position: Position{Line: 0, Column: 0}, Type: K_TRUE},
		}},
		{"false", []Token{
			Token{Text: "false", Position: Position{Line: 0, Column: 0}, Type: K_FALSE},
		}},
		{"if true {}", []Token{
			Token{Text: "if", Position: Position{Line: 0, Column: 0}, Type: K_IF},
			Token{Text: "true", Position: Position{Line: 0, Column: 3}, Type: K_TRUE},
			Token{Text: "{", Position: Position{Line: 0, Column: 8}, Type: D_LBRACE},
			Token{Text: "}", Position: Position{Line: 0, Column: 9}, Type: D_RBRACE},
		}},
		{"if true {} else {}", []Token{
			Token{Text: "if", Position: Position{Line: 0, Column: 0}, Type: K_IF},
			Token{Text: "true", Position: Position{Line: 0, Column: 3}, Type: K_TRUE},
			Token{Text: "{", Position: Position{Line: 0, Column: 8}, Type: D_LBRACE},
			Token{Text: "}", Position: Position{Line: 0, Column: 9}, Type: D_RBRACE},
			Token{Text: "else", Position: Position{Line: 0, Column: 11}, Type: K_ELSE},
			Token{Text: "{", Position: Position{Line: 0, Column: 16}, Type: D_LBRACE},
			Token{Text: "}", Position: Position{Line: 0, Column: 17}, Type: D_RBRACE},
		}},
	}

	for i := range tests {
		lex := NewLexer(strings.NewReader(tests[i].in))
		j := 0
		for cur, err := lex.NextToken(); cur.Type != EOF; j++ {
			if err != nil {
				t.Error(err.Error())
				return
			}
			if cur.Type != tests[i].tokens[j].Type {
				fmt.Printf("%+v\n", cur)
				t.Errorf("[%d] Unexpected token: %s != %s", i, TokenStrings[tests[i].tokens[j].Type], TokenStrings[cur.Type])
			}
			if cur.Position.Line != tests[i].tokens[j].Position.Line {
				fmt.Printf("%+v\n", cur)
				t.Errorf("[%d] Unexpected line #: %d != %d", i, tests[i].tokens[j].Position.Line, cur.Position.Line)
			}
			if cur.Position.Column != tests[i].tokens[j].Position.Column {
				fmt.Printf("%+v\n", cur)
				t.Errorf("[%d] Unexpected column #: %d != %d", i, tests[i].tokens[j].Position.Column, cur.Position.Column)
			}
			cur, err = lex.NextToken()
		}
	}
}
