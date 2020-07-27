package repl

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nirosys/stitch/cmd/stitch/subcmd/internal/shellcmd"
	"github.com/nirosys/stitch/object"
)

func (r *Repl) compileCommand() *shellcmd.Command {
	var compileCommand = &shellcmd.Command{
		Use:   "compile <ident>",
		Short: "Compile a node object to Gaufre graph",
		RunE:  r.compile,
	}
	return compileCommand
}

func (r *Repl) compile(cmd *shellcmd.Command, args []string) error {
	if len(args) == 1 {
		if obj, have := r.env.Get(args[0]); !have {
			fmt.Printf("invalid identifier: '%s'", args[0])
		} else if node, ok := obj.(*object.Node); ok {
			if g, err := r.evaluator.CompileObject(node); err != nil {
				return err
			} else if b, err := json.Marshal(g); err != nil {
				return err
			} else {
				os.Stdout.Write(b)
				fmt.Println("")
			}
		}
	}
	return nil
}
