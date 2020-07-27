package subcmd

import (
	"fmt"

	"github.com/nirosys/stitch/cmd/stitch/subcmd/repl"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "stitch",
	Short: "A PoC flow language for SNITCH",
	Args:  cobra.MaximumNArgs(1),
	RunE:  do_repl,
}

func init() {
	RootCmd.Flags().StringP("init-with", "i", "", "Specify a script to run at the start of the session")
}

func do_repl(cmd *cobra.Command, args []string) error {
	repl := repl.NewRepl()

	if v, err := cmd.Flags().GetString("init-with"); err == nil && v != "" {
		if err := repl.LoadFile(v); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
	}

	if len(args) == 1 {
		if err := repl.LoadFile(args[0]); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
	} else {
		if err := repl.Run(); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
	}

	return nil
}
