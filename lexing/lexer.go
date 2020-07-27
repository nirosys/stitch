package lexing

import (
	"errors"
	"fmt"
	"io"
)

var ErrUnexepctedChar = errors.New("unexpected character")

type TokenType uint8

const (
	TOKEN_NONE  TokenType = iota
	K_LET                 /* let - statement start */
	K_IMPORT              /* import - statement start */
	K_NODE                /* node - node definition */
	K_FUNCTION            /* fn - function definition */
	K_MODIFIER            /* mod - node modifier */
	K_INTERNAL            /* internal - alias to hosted value/node/etc. */
	K_IF                  /* if - start of conditional */
	K_ELSE                /* else - */
	K_TRUE                /* true - for.. boolean true */
	K_FALSE               /* false - for.. boolean false */
	K_AND                 /* and - for logic */
	K_OR                  /* or - for logic */
	K_FOREACH             /* foreach - for looping */
	K_IN                  /* in - for foreach loops */
	L_INTEGER             //
	L_FLOAT               //
	L_STRING              //
	D_RPARENTH            //
	D_LPARENTH            //
	D_RBRACE              //
	D_LBRACE              //
	D_LBRACKET            //
	D_RBRACKET            //
	D_SEMICOLON           //
	O_MINUS               //
	O_PLUS                //
	O_STAR                //
	O_SLASH               //
	O_MODULUS             //
	O_COMMA               //
	O_COLON               //
	O_DOT                 //
	O_LESSTHAN            //
	O_GT                  // >
	O_LT                  // <
	O_GTEQ                // >=
	O_LTEQ                // <=
	O_EQ                  // ==
	O_NEQ                 // !=
	O_TAGMARKER           //
	O_BANG                // !
	O_ARROW               //
	O_ASSIGN              //
	IDENT                 //
	COMMENT               //
	EOF                   //
)

var TokenStrings = map[TokenType]string{
	K_LET:       "keyword 'let'",
	K_IMPORT:    "keyword 'import'",
	K_NODE:      "keyword 'node'",
	K_FUNCTION:  "keyword 'fn'",
	K_MODIFIER:  "keyword 'mod'",
	K_INTERNAL:  "keyword 'internal'",
	K_IF:        "keyword 'if'",
	K_ELSE:      "keyword 'else'",
	K_TRUE:      "keyword 'true'",
	K_FALSE:     "keyword 'false'",
	K_AND:       "keyword 'and'",
	K_OR:        "keyword 'or'",
	K_FOREACH:   "keyword 'foreach'",
	K_IN:        "keyword 'in'",
	L_INTEGER:   "INTEGER literal",
	L_FLOAT:     "FLOAT literal",
	L_STRING:    "STRING literal",
	D_RPARENTH:  "')'",
	D_LPARENTH:  "'('",
	D_RBRACE:    "'}'",
	D_LBRACE:    "'{'",
	D_LBRACKET:  "'['",
	D_RBRACKET:  "']'",
	D_SEMICOLON: "';'",
	O_MINUS:     "'-'",
	O_PLUS:      "'+'",
	O_STAR:      "'*'",
	O_SLASH:     "'/'",
	O_MODULUS:   "'%'",
	O_COMMA:     "','",
	O_COLON:     "':'",
	O_DOT:       "'.'",
	O_LESSTHAN:  "'<'",
	O_GT:        "'>'",
	O_GTEQ:      "'>='",
	O_LT:        "'<'",
	O_LTEQ:      "'<='",
	O_EQ:        "'=='",
	O_NEQ:       "'!='",
	O_TAGMARKER: "'@'",
	O_BANG:      "'!'",
	O_ARROW:     "'->'",
	O_ASSIGN:    "'='",
	IDENT:       "IDENTIFIER",
	COMMENT:     "COMMENT",
	EOF:         "EOF",
}

type Position struct {
	Line   int
	Column int
}

func (p *Position) AdvanceLine() {
	p.Column = 0
	p.Line += 1
}

func (p *Position) AdvanceChar() {
	p.Column += 1
}

type Token struct {
	Text     string
	Position Position
	Type     TokenType
}

type Lexer struct {
	input io.Reader

	buffer   []byte
	readPos  int
	position Position
	ch       byte
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		input:    r,
		buffer:   make([]byte, 0, 256),
		readPos:  0,
		position: Position{Line: 0, Column: 0},
		ch:       0,
	}
}

func (l *Lexer) readMore() error {
	l.buffer = l.buffer[:cap(l.buffer)] // Reset our buffer to max size, jic

	n, err := l.input.Read(l.buffer)
	if err != nil && err != io.EOF { // Error reading, tell parser and exit
		return err
	} else if err == io.EOF {
		l.buffer = l.buffer[:0]
		return io.EOF
	}
	l.buffer = l.buffer[:n] // Resize to what we just read
	l.readPos = 0
	return nil
}

func (l *Lexer) NextToken() (Token, error) {
	var token Token
	for token.Type == TOKEN_NONE {
		char, err := l.peekChar()
		pos := l.position

		if len(l.buffer) == 0 {
			return Token{Text: "", Position: pos, Type: EOF}, nil
		} else if err != nil {
			return Token{}, err
		}

		if char <= '9' && char >= '0' {
			b, dec, err := l.slurpNumeric()
			if err != nil {
				return Token{}, err
			}
			numType := L_INTEGER
			if dec {
				numType = L_FLOAT
			}
			return Token{Text: string(b), Position: pos, Type: numType}, nil
		} else if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_' {
			b, err := l.slurpIdentifier()
			if err != nil {
				return Token{}, err
			}
			t := Token{Text: string(b), Position: pos}
			switch t.Text {
			case "let":
				t.Type = K_LET
			case "import":
				t.Type = K_IMPORT
			case "node":
				t.Type = K_NODE
			case "fn":
				t.Type = K_FUNCTION
			case "internal":
				t.Type = K_INTERNAL
			case "if":
				t.Type = K_IF
			case "else":
				t.Type = K_ELSE
			case "true":
				t.Type = K_TRUE
			case "false":
				t.Type = K_FALSE
			case "and":
				t.Type = K_AND
			case "or":
				t.Type = K_OR
			case "foreach":
				t.Type = K_FOREACH
			case "in":
				t.Type = K_IN
			case "mod":
				t.Type = K_MODIFIER
			default:
				t.Type = IDENT
			}
			return t, nil
		} else {
			switch char {
			case '+':
				_, _ = l.takeChar()
				return Token{Text: "+", Position: pos, Type: O_PLUS}, nil
			case '-':
				_, _ = l.takeChar()
				cur, err := l.peekChar()
				if err != nil || cur != '>' {
					return Token{Text: "-", Position: pos, Type: O_MINUS}, nil
				} else if cur == '>' {
					_, _ = l.takeChar()
					return Token{Text: "->", Position: pos, Type: O_ARROW}, nil
				}
			case '*':
				_, _ = l.takeChar()
				return Token{Text: "*", Position: pos, Type: O_STAR}, nil
			case '/':
				_, _ = l.takeChar()
				return Token{Text: "/", Position: pos, Type: O_SLASH}, nil
			case '%':
				_, _ = l.takeChar()
				return Token{Text: "%", Position: pos, Type: O_MODULUS}, nil
			case '(':
				_, _ = l.takeChar()
				return Token{Text: "(", Position: pos, Type: D_LPARENTH}, nil
			case ')':
				_, _ = l.takeChar()
				return Token{Text: ")", Position: pos, Type: D_RPARENTH}, nil
			case '{':
				_, _ = l.takeChar()
				return Token{Text: "{", Position: pos, Type: D_LBRACE}, nil
			case '}':
				_, _ = l.takeChar()
				return Token{Text: "}", Position: pos, Type: D_RBRACE}, nil
			case '[':
				_, _ = l.takeChar()
				return Token{Text: "[", Position: pos, Type: D_LBRACKET}, nil
			case ']':
				_, _ = l.takeChar()
				return Token{Text: "]", Position: pos, Type: D_RBRACKET}, nil
			case ';':
				_, _ = l.takeChar()
				return Token{Text: ";", Position: pos, Type: D_SEMICOLON}, nil
			case ',':
				_, _ = l.takeChar()
				return Token{Text: ",", Position: pos, Type: O_COMMA}, nil
			case ':':
				_, _ = l.takeChar()
				return Token{Text: ":", Position: pos, Type: O_COLON}, nil
			case '.':
				_, _ = l.takeChar()
				return Token{Text: ".", Position: pos, Type: O_DOT}, nil
			case '<':
				_, _ = l.takeChar()

				next, err := l.peekChar()
				if err != nil {
					return Token{}, err
				} else if next != '=' {
					return Token{Text: "<", Position: pos, Type: O_LT}, nil
				} else {
					_, _ = l.takeChar()
					return Token{Text: "<=", Position: pos, Type: O_LTEQ}, nil
				}
			case '>':
				_, _ = l.takeChar()

				next, err := l.peekChar()
				if err != nil {
					return Token{}, err
				} else if next != '=' {
					return Token{Text: ">", Position: pos, Type: O_GT}, nil
				} else {
					_, _ = l.takeChar()
					return Token{Text: ">=", Position: pos, Type: O_GTEQ}, nil
				}
			case '!':
				_, _ = l.takeChar()

				next, err := l.peekChar()
				if err != nil {
					return Token{}, err
				} else if next != '=' {
					return Token{Text: "!", Position: pos, Type: O_BANG}, nil
				} else {
					_, _ = l.takeChar()
					return Token{Text: "!=", Position: pos, Type: O_NEQ}, nil
				}
			case '=':
				_, _ = l.takeChar()

				next, err := l.peekChar()

				if err != nil {
					return Token{}, err
				} else if next != '=' {
					return Token{Text: "=", Position: pos, Type: O_ASSIGN}, nil
				} else if next == '=' {
					_, _ = l.takeChar()
					return Token{Text: "==", Position: pos, Type: O_EQ}, nil
				}
			case '"':
				b, err := l.slurpString()
				if err != nil {
					return Token{}, err
				}
				return Token{Text: string(b), Position: pos, Type: L_STRING}, nil
			case '#':
				b, err := l.slurpComment()
				if err != nil {
					return Token{}, err
				}
				return Token{Text: string(b), Position: pos, Type: COMMENT}, nil
			case '@':
				_, _ = l.takeChar()
				return Token{Text: "@", Position: pos, Type: O_TAGMARKER}, nil
			case ' ', '\t':
				_, _ = l.takeChar()
			case '\n':
				_, _ = l.takeChar()
				l.position.AdvanceLine()
			default:
				ch, _ := l.takeChar()
				return Token{}, fmt.Errorf("%w: %c", ErrUnexepctedChar, ch)
			}
		}
	}
	return Token{}, nil
}

func (l *Lexer) peekChar() (byte, error) {
	if l.readPos >= len(l.buffer) {
		err := l.readMore()
		if err != nil || len(l.buffer) == 0 {
			return 0, err
		}
	}
	return l.buffer[l.readPos], nil
}

func (l *Lexer) takeChar() (byte, error) {
	if l.readPos >= len(l.buffer) {
		err := l.readMore()
		if err != nil || len(l.buffer) == 0 {
			return 0, err
		}
	}
	b := l.buffer[l.readPos]
	l.readPos++
	l.position.AdvanceChar()
	return b, nil
}

func (l *Lexer) slurpNumeric() ([]byte, bool, error) {
	var haveDecimal bool
	var bytes []byte

	cur, err := l.peekChar()
	for ((cur <= '9' && cur >= '0') || (cur == '.' && !haveDecimal)) && err == nil {
		cur, err = l.takeChar()
		if !haveDecimal {
			haveDecimal = cur == '.'
		}
		bytes = append(bytes, cur)
		cur, err = l.peekChar()
	}
	if err == io.EOF {
		err = nil
	}
	return bytes, haveDecimal, err
}

func (l *Lexer) slurpIdentifier() ([]byte, error) {
	var bytes []byte
	cur, err := l.peekChar()
	for ((cur <= 'z' && cur >= 'a') || (cur <= 'Z' && cur >= 'A') || (cur <= '9' && cur >= '0') || cur == '_') && err == nil {
		cur, err = l.takeChar()
		bytes = append(bytes, cur)
		cur, err = l.peekChar()
	}
	if err == io.EOF {
		err = nil
	}
	return bytes, err
}

func (l *Lexer) slurpString() ([]byte, error) {
	var escaped bool = false
	var done bool = false
	var bytes []byte
	_, err := l.takeChar() // take leading double-quote
	if err != nil {
		return []byte{}, err
	}

	cur, err := l.peekChar()
	for !done && err == nil {
		if cur == '"' && !escaped {
			done = true
			_, _ = l.takeChar()
		} else if cur == '"' && escaped {
			escaped = false
			bytes = append(bytes, '"')
			_, _ = l.takeChar()
		} else if cur == '\\' && !escaped {
			escaped = true
			_, _ = l.takeChar()
		} else if escaped {
			_, _ = l.takeChar()
			escaped = false
			bytes = append(bytes, '\\')
			bytes = append(bytes, cur)
		} else {
			cur, _ = l.takeChar()
			bytes = append(bytes, cur)
		}
		cur, err = l.peekChar()
	}
	if err == io.EOF && done {
		err = nil
	}
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func (l *Lexer) slurpComment() ([]byte, error) {
	var bytes []byte
	_, _ = l.takeChar() // Take leading #

	cur, err := l.peekChar()
	for cur != '\n' && err == nil {
		cur, _ = l.takeChar()
		bytes = append(bytes, cur)
		cur, err = l.peekChar()
	}
	if err == io.EOF {
		err = nil
	}
	return bytes, err
}
