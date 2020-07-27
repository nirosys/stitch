package repl

import (
	"github.com/nirosys/stitch/cmd/stitch/subcmd/internal/shellcmd"
)

func (r *Repl) quitCommand() *shellcmd.Command {
	var quitCommand = &shellcmd.Command{
		Use:   "quit",
		Short: "Quit the current session",
		RunE: func(cmd *shellcmd.Command, args []string) error {
			r.quit = true
			return nil
		},
	}
	return quitCommand
}
