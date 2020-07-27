package repl

import (
	"os"

	"github.com/nirosys/stitch/cmd/stitch/subcmd/internal"
	"github.com/nirosys/stitch/cmd/stitch/subcmd/internal/shellcmd"
)

func (r *Repl) dotCommand() *shellcmd.Command {
	var dotCommand = &shellcmd.Command{
		Use:   "dot [obj]",
		Short: "Describe the graph in dot syntax",
		RunE: func(cmd *shellcmd.Command, args []string) error {
			trav := internal.NewTraverser(r.env)
			if len(args) == 0 {
				trav.RenderDotAll(os.Stdout)
			} else {
				// TODO: Implement me.
			}
			return nil
		},
	}

	return dotCommand
}
