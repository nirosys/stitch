package subcmd

import (
	"errors"
	//"fmt"
	//"os"

	"github.com/spf13/cobra"
)

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format a stitch program",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("invalid arguments")
		}
		return nil
	},
	RunE: format,
}

func init() {
	RootCmd.AddCommand(fmtCmd)
}

func format(cmd *cobra.Command, args []string) error {
	/*
		if f, err := os.Open(args[0]); err != nil {
			fmt.Printf("ERROR Opening File: %s\n", err.Error())
			return nil
		} else {
			//parser := parsing.NewParser(f)
			//prog := parser.ParseProgram()
			//os.Stdout.WriteString(prog.String())
		}
	*/
	return nil
}
