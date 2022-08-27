package cli

type Command interface {
	Help() string
	Name() string
	Run(args []string) error
	Synopsis() string
}
