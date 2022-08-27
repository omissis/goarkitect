package cli

import (
	"errors"
	"os"

	flag "github.com/spf13/pflag"
)

var (
	ErrEmptyAppName       = errors.New("app name cannot be empty")
	ErrEmptyCommands      = errors.New("app must have at least one command")
	ErrNoCommandSpecified = errors.New("app command was not specified")
)

func NewApp(
	name string,
	commands []Command,
	flagSet *flag.FlagSet,
) (*App, error) {
	if name == "" {
		return nil, ErrEmptyAppName
	}

	if len(commands) < 1 {
		return nil, ErrEmptyCommands
	}

	if flagSet == nil {
		flagSet = flag.NewFlagSet("global", flag.ContinueOnError)
	}

	return &App{
		name:     name,
		commands: commands,
		flagSet:  flagSet,
	}, nil
}

type App struct {
	name     string
	commands []Command
	flagSet  *flag.FlagSet
}

func (a *App) Run() error {
	if err := a.flagSet.Parse(GetArgs(os.Args, 1)); err != nil {
		return err
	}

	args := a.flagSet.Args()
	if len(args) < 1 {
		return ErrNoCommandSpecified
	}

	for _, cmd := range a.commands {
		if cmd.Name() == args[0] {
			return cmd.Run(args)
		}
	}

	return nil
}
