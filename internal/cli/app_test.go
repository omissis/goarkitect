package cli_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	flag "github.com/spf13/pflag"
	"golang.org/x/exp/slices"

	"github.com/omissis/goarkitect/internal/cli"
)

func Test_NewApp_Fail(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc     string
		name     string
		commands []cli.Command
		wantErr  error
	}{
		{
			desc:    "no name",
			wantErr: cli.ErrEmptyAppName,
		},
		{
			desc:    "no commands",
			name:    "test",
			wantErr: cli.ErrEmptyCommands,
		},
	}
	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			_, err := cli.NewApp(tC.name, tC.commands, nil)
			if !errors.Is(err, tC.wantErr) {
				t.Errorf("wantErr = %v, got = %v", tC.wantErr, err)
			}
		})
	}
}

func Test_App_Run_NoArgs(t *testing.T) {
	t.Parallel()

	app, err := cli.NewApp("test", []cli.Command{&cmd{}}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err = app.Run(); !errors.Is(err, cli.ErrNoCommandSpecified) {
		t.Errorf("wantErr = %v, got = %v", cli.ErrNoCommandSpecified, err)
	}
}

func Test_App_Run_Args_And_Flags(t *testing.T) {
	t.Parallel()

	osArgs := slices.Clone(os.Args)

	testCases := []struct {
		desc    string
		args    []string
		wantErr string
	}{
		{
			desc: "one arg, no flags",
			args: []string{"test"},
		},
		{
			desc: "two args, no flags",
			args: []string{"test", "example"},
		},
		{
			desc:    "one arg, one unknown global flag",
			args:    []string{"--foo=bar", "test"},
			wantErr: "unknown flag: --foo",
		},
		{
			desc: "one arg, one known global flag",
			args: []string{"--output=json", "test"},
		},
		{
			desc: "one arg, one known global flag specified as command flag",
			args: []string{"test", "--output=json"},
		},
	}
	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			os.Args = append(os.Args, tC.args...)
			defer func() {
				os.Args = slices.Clone(osArgs)
			}()

			out := "text"
			flagSet := flag.NewFlagSet("global", flag.ContinueOnError)
			flagSet.StringVar(&out, "output", "text", "format of the output")

			app, err := cli.NewApp("test", []cli.Command{&cmd{}}, flagSet)
			if err != nil {
				t.Fatal(err)
			}

			if err = app.Run(); err != nil && !strings.Contains(err.Error(), tC.wantErr) {
				t.Errorf("wanted err = %v, got = %v", tC.wantErr, err)
			}
		})
	}
}

type cmd struct{}

func (c *cmd) Name() string {
	return "test"
}

func (c *cmd) Help() string {
	return "test"
}

func (c *cmd) Run(args []string) error {
	return nil
}

func (c *cmd) Synopsis() string {
	return "test"
}
