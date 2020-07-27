package subcmd

// TODO: This subcmd needs to be reworked..

import (
	//	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/nirosys/stitch"

	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile a stitch program to a gaufre graph.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("invalid arguments")
		}
		return nil
	},
	RunE: compile,
}

func init() {
	RootCmd.AddCommand(compileCmd)
}

func compile(cmd *cobra.Command, args []string) error {
	var source io.Reader

	if args[0] == "-" {
		source = os.Stdin
	} else if f, err := os.Open(args[0]); err != nil {
		fmt.Printf("ERROR Opening File: %s\n", err.Error())
		return nil
	} else {
		source = f
	}
	prog := stitch.NewProgram(source)
	if prog == nil {
		fmt.Printf("ERROR\n")
		for _, e := range prog.Errors() {
			fmt.Printf("   %s\n", e)
		}
	} else {
		//eval := stitch.NewEvaluator()
		//eval.Resolver = internal.NewResolver()
		//if graph, err := eval.Compile(prog); err != nil {
		//	fmt.Printf("ERROR: %s\n", err.Error())
		//	return nil
		//} else {
		//	b, err := json.Marshal(graph)
		//	if err != nil {
		//		fmt.Printf("ERROR: %s\n", err.Error())
		//		return nil
		//	}
		//	os.Stdout.Write(b)
		//}
	}

	return nil
}
