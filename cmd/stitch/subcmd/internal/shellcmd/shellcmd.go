package shellcmd

import (
	"flag"
	"strings"
)

type Command struct {
	// Command usage (first word is used as the actual command)
	Use string

	// A short description for the command
	Short string

	// A long description of the command.
	Long string

	RunE func(cmd *Command, args []string) error

	flagSet *flag.FlagSet
}

func (c *Command) Flags() *flag.FlagSet {
	if c.flagSet == nil {
		toks := strings.SplitN(strings.TrimSpace(c.Use), " ", 2)
		c.flagSet = flag.NewFlagSet(toks[0], flag.ContinueOnError)
	}
	return c.flagSet
}

type Parser struct {
	Prefix     string
	CommandMap map[string]*Command
}

func NewParser() *Parser {
	return &Parser{
		Prefix:     "",
		CommandMap: map[string]*Command{},
	}
}

func (p *Parser) Run(line string) error {
	toks := strings.Split(line, " ")
	if len(toks) > 0 {
		c := strings.TrimPrefix(toks[0], p.Prefix)
		if cmd, ok := p.CommandMap[c]; ok {
			return cmd.RunE(cmd, toks[1:])
		}
	}
	return nil
}

func (p *Parser) AddCommand(cmd *Command) {
	toks := strings.SplitN(strings.TrimSpace(cmd.Use), " ", 2)
	if len(toks) > 0 {
		p.CommandMap[toks[0]] = cmd
	}
}
