package main

import (
	"log"
	"os"

	"goarkitect/cmd"

	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("app", "0.1.0-dev")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"validate": cmd.ValidateFactory,
		"verify":   cmd.VerifyFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
