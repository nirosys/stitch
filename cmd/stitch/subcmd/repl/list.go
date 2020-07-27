package repl

import (
	"fmt"
	"os"
	"strings"

	"github.com/nirosys/stitch/cmd/stitch/subcmd/internal"
	"github.com/nirosys/stitch/cmd/stitch/subcmd/internal/shellcmd"
)

func (r *Repl) listCommand() *shellcmd.Command {
	var listCommand = &shellcmd.Command{
		Use:   "ls [pkg]",
		Short: "List all objects in a scope.",
		RunE:  r.list,
	}
	listCommand.Flags().Bool("all", false, "List all objects, even unbounded")

	return listCommand
}

func (r *Repl) list(cmd *shellcmd.Command, args []string) error {
	table := internal.NewTableWriter()
	table.SetColumnTitles([]string{"Name", "Type"})

	all := false

	idx := -1
	for idx = 0; idx < len(args) && strings.HasPrefix(args[idx], "-"); idx++ {
		switch args[idx] {
		case "-a":
			all = true
		default:
			fmt.Printf("unknown argument: %s\n", args[idx])
			break
		}
	}

	args = args[idx:]

	scope := r.env

	if len(args) > 0 { // assume we have a package name
		pkg := args[0]
		fmt.Printf("Listing package: %s\n", pkg)
		if p, err := scope.GetPackage(pkg); err != nil {
			return err
		} else {
			scope = p.Environment
		}
	}

	bound := scope.GetNames()
	if len(bound) > 0 {
		for _, n := range scope.GetNames() {
			obj, _ := scope.Get(n)
			table.AddRow([]string{n, string(obj.Type())})
		}
		fmt.Printf("Scope:\n")
		table.Write(os.Stdout)
	}

	table.ClearData()

	if all {
		unbound := scope.GetUnboundNodes()
		if len(unbound) > 0 {
			for _, n := range scope.GetUnboundNodes() {
				obj, _ := scope.Get(n)
				table.AddRow([]string{n, string(obj.Type())})
			}
			fmt.Printf("\nUnbounded Nodes\n")
			table.Write(os.Stdout)
		}
	}
	return nil
}
