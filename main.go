package main

import (
	"github.com/omissis/goarkitect/cmd"
	"github.com/omissis/goarkitect/internal/cli"
	"github.com/omissis/goarkitect/internal/logx"

	flag "github.com/spf13/pflag"
)

var (
	Version   string = "unknown"
	GitCommit string = "unknown"
	BuildTime string = "unknown"
	GoVersion string = "unknown"
	OsArch    string = "unknown"
)

func main() {
	out := "text"

	flagSet := flag.NewFlagSet("global", flag.ContinueOnError)
	flagSet.StringVar(&out, "output", "text", "format of the output")

	app, err := cli.NewApp(
		"goarkitect",
		[]cli.Command{
			cmd.NewValidateCommand(&out),
			cmd.NewVerifyCommand(&out),
			cmd.NewVersionCommand(&out, map[string]string{
				"version":   Version,
				"gitCommit": GitCommit,
				"buildTime": BuildTime,
				"goVersion": GoVersion,
				"osArch":    OsArch,
			}),
		},
		flagSet,
	)
	if err != nil {
		logx.Fatal(err)
	}

	if err := app.Run(); err != nil {
		logx.Fatal(err)
	}
}
