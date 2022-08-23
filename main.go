package main

import (
	"flag"
	"os"

	"github.com/omissis/goarkitect/cmd"
	"github.com/omissis/goarkitect/internal/logx"

	"github.com/mitchellh/cli"
)

func main() {
	out := "text"

	flagSet := flag.NewFlagSet("global", flag.ContinueOnError)
	flagSet.StringVar(&out, "output", "text", "format of the output")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		logx.Fatal(err)
	}

	logx.SetFormat(out)

	c := cli.NewCLI("app", "unknown")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"validate": func() (cli.Command, error) {
			return cmd.ValidateFactory(out)
		},
		"verify": func() (cli.Command, error) {
			return cmd.VerifyFactory(out)
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		logx.Fatal(err)
	}

	os.Exit(exitStatus)
}
