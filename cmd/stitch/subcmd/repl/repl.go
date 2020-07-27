package repl

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nirosys/stitch"
	"github.com/nirosys/stitch/cmd/stitch/subcmd/internal"
	"github.com/nirosys/stitch/cmd/stitch/subcmd/internal/shellcmd"
	"github.com/nirosys/stitch/eval"
	"github.com/nirosys/stitch/lexing"
	"github.com/nirosys/stitch/object"

	"github.com/gookit/color"
	"github.com/peterh/liner"
)

var ReplVersion string = "<invalid build>"

const NORMAL_PROMPT = ">> "
const CONTINUE_PROMPT = "...   "

type MatchCount struct {
	Start int
	End   int
}

const (
	MatchScope   int = 0
	MatchBrace   int = 1
	MatchString  int = 2
	MatchParenth int = 3
)

type MatchCounts map[int]*MatchCount

func NewMatchCounts() MatchCounts {
	return map[int]*MatchCount{
		MatchScope:   &MatchCount{},
		MatchBrace:   &MatchCount{},
		MatchString:  &MatchCount{},
		MatchParenth: &MatchCount{},
	}
}

type Repl struct {
	evaluator *eval.Evaluator
	env       *object.Environment
	commander *shellcmd.Parser
	quiet     bool
	matches   MatchCounts
	inString  bool
	escaped   bool
	quit      bool
}

func NewRepl() *Repl {
	repl := &Repl{
		env:       object.NewEnvironment(),
		evaluator: eval.NewEvaluator(),
		commander: shellcmd.NewParser(),
		quit:      false,
	}
	repl.evaluator.Resolver = repl
	repl.commander.Prefix = "."
	repl.commander.AddCommand(repl.listCommand())
	repl.commander.AddCommand(repl.quitCommand())
	repl.commander.AddCommand(repl.dotCommand())
	repl.commander.AddCommand(repl.compileCommand())
	repl.commander.AddCommand(&shellcmd.Command{
		Use:   "quiet",
		Short: "Toggle quiet mode",
		RunE: func(cmd *shellcmd.Command, args []string) error {
			repl.quiet = !repl.quiet
			fmt.Printf("Quiet mode: %t\n", repl.quiet)
			return nil
		},
	})
	return repl
}

func (r *Repl) Run() error {
	fmt.Printf("Stitch\n%s\n\n", ReplVersion)
	fmt.Printf("Use '.quit' to quit, or '.help' for help\n")

	l := liner.NewLiner()
	defer l.Close()
	l.SetHighlighter(r)

	l.SetCtrlCAborts(true)
	l.SetMultiLineMode(true)

	r.matches = NewMatchCounts()

	prompt := NORMAL_PROMPT

	var tmpLines []string
REPLOOP:
	for {
		if line, err := l.Prompt(prompt); err == nil {
			if strings.HasPrefix(line, ".") {
				if err := r.commander.Run(line); err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
				} else if r.quit {
					break REPLOOP
				}
				continue
			}

			tmpline := strings.TrimSpace(line)
			if len(tmpline) == 0 {
				continue
			}

			if tmpline[len(tmpline)-1:] == "\\" || !r.countPairs(line) {
				tmpLines = append(tmpLines, strings.TrimRight(tmpline, "\\"))
				prompt = CONTINUE_PROMPT
				continue
			} else {
				tmpLines = append(tmpLines, line)
			}

			resultLine := strings.Join(tmpLines, "")
			l.AppendHistory(resultLine)
			tmpLines = nil

			reader := strings.NewReader(resultLine)
			r.ExecuteCode(reader)

			prompt = NORMAL_PROMPT
		}
	}

	return nil
}

// return true if we have all pairs.
func (r *Repl) countPairs(line string) bool {
	inString := r.inString
	escaped := r.escaped
	for _, c := range line {
		switch c {
		case '{':
			if !inString {
				r.matches[MatchScope].Start += 1
			}
		case '}':
			if !inString {
				r.matches[MatchScope].End += 1
			}
		case '"':
			if !inString && !escaped {
				inString = true
				r.matches[MatchString].Start += 1
			} else if !escaped {
				inString = false
				r.matches[MatchString].End += 1
			}
		}
	}
	r.inString = inString
	r.escaped = escaped
	done := true
	for _, v := range r.matches {
		done = done && (v.Start == v.End)
	}
	return done
}

func (r *Repl) handleCommands(line string, l *liner.State) (bool, bool) {
	toks := strings.Split(line, " ")
	cmd := toks[0]
	args := toks[1:]
	switch cmd {
	case ".quit":
		history := filepath.Join(os.TempDir(), ".stitch_history")
		if f, err := os.Create(history); err == nil {
			l.WriteHistory(f)
			f.Close()
		}
		return true, true
	case ".ls":
		r.ListEnv(args)
		return true, false
	case ".dot":
		trav := internal.NewTraverser(r.env)
		if len(args) == 0 {
			trav.RenderDotAll(os.Stdout)
		} else {
			//trav.RenderDot(os.Stdout, args[0])
		}
		return true, false
	case ".quiet":
		r.quiet = !r.quiet
		fmt.Printf("Quiet mode: %t\n", r.quiet)
		return true, false
	case ".compile": // .compile <ident>
		if len(args) == 1 {
			if obj, have := r.env.Get(args[0]); !have {
				fmt.Printf("invalid identifier: '%s'", args[0])
			} else if node, ok := obj.(*object.Node); ok {
				if g, err := r.evaluator.CompileObject(node); err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
				} else if b, err := json.Marshal(g); err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
				} else {
					os.Stdout.Write(b)
					fmt.Println("")
				}
			} else {
				fmt.Printf("ERROR: '%s' is not a node object", args[0])
			}
			return true, false
		}
	case ".run": // .run <ident>
	case ".help":
		printReplHelp()
		return true, false
	}
	return false, false
}

func (r *Repl) ExecuteCode(reader io.Reader) error {
	prog := stitch.NewProgram(reader)
	if prog == nil {
		fmt.Printf("Parser error(s):\n")
		for _, e := range prog.Errors() {
			fmt.Printf("   %s\n", e)
		}
	} else {
		//if obj, err := r.evaluator.EvalProgram(prog, r.env); err != nil {
		//	fmt.Printf("ERROR: %s\n", err.Error())
		//} else if obj != nil && !r.quiet {
		//	fmt.Printf("%s\n", obj.Inspect())
		//}
	}
	return nil
}

func (r *Repl) LoadFile(path string) error {
	if f, err := os.Open(path); err != nil {
		return err
	} else {
		return r.ExecuteCode(f)
	}
}

func (r *Repl) ListEnv(args []string) {
}

func (r *Repl) Highlight(text string) (string, error) {
	lex := lexing.NewLexer(strings.NewReader(text))
	keywords := color.FgGreen.Render
	strs := color.FgCyan.Render
	literals := color.FgCyan.Render

	var buffer strings.Builder

	offset := 0
	for tok, _ := lex.NextToken(); tok.Type != lexing.EOF; {
		plainLength := tok.Position.Column - offset
		if plainLength > 0 {
			buffer.WriteString(text[offset:tok.Position.Column])
			offset += plainLength
		}
		switch {
		case tok.Type == lexing.L_INTEGER || tok.Type == lexing.K_TRUE || tok.Type == lexing.K_FALSE:
			buffer.WriteString(literals(tok.Text))
			offset += len(tok.Text)
		case tok.Type == lexing.L_STRING:
			buffer.WriteString(strs(fmt.Sprintf("\"%s\"", tok.Text)))
			offset += len(tok.Text) + 2
		case tok.Type >= lexing.K_LET && tok.Type <= lexing.K_IN:
			buffer.WriteString(keywords(tok.Text))
			offset += len(tok.Text)
		}
		tok, _ = lex.NextToken()
	}
	if offset < len(text) {
		buffer.WriteString(text[offset:len(text)])
	}

	return buffer.String(), nil
}

// We're given the full line, and the offset of where the cursor is.
// It's up to use to then split the words, figure out which word the cursor
// is on, and then return all words that can match.. while also determining
// what would result in the head, and tail, of the new line should one of those
// words be substituted..
func (r *Repl) completeWord(line string, pos int) (string, []string, string) {
	return line, nil, ""
}

func (r *Repl) Resolve(name string) (object.Object, error) {
	if fn, ok := internal.HostedFuncs[name]; ok {
		return &object.InternalFunction{Fn: fn, Env: nil}, nil
	} else if nodeType, ok := internal.HostedNodeTypes[name]; ok {
		return nodeType, nil
	} else {
		return nil, fmt.Errorf("unknown internal \"%s\"", name)
	}
}

func printReplHelp() {
	fmt.Printf(`Stitch REPL Help
   Commands:
      .ls [pkg]    - List named variables, and unnamed nodes in global scope, or package.
      .dot [var]   - Render the current graph (or graph rooted by var) in dot syntax.
      .quiet       - Turn off auto-inspect when evaluating expressions.
		.compile <ident> - Compile a given node to its gaufre graph.
		.run <ident> - Compile a node, and run it.
		.stop        - Stop running the current graph.
`)
}
