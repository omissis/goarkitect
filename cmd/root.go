package cmd

import (
	"github.com/omissis/goarkitect/internal/cobrax"
	"github.com/spf13/cobra"
)

type rootConfig struct {
	Output string
}

type RootCommand struct {
	*cobra.Command
	config *rootConfig
}

func NewRootCommand(versions map[string]string) *RootCommand {
	const envPrefix = ""

	root := &RootCommand{
		Command: &cobra.Command{
			PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
				cobrax.BindFlags(cmd, cobrax.InitEnvs(envPrefix), envPrefix)

				return nil
			},
			Use:          "goarkitect",
			SilenceUsage: true,
		},
		config: &rootConfig{},
	}

	v := cobrax.InitEnvs(envPrefix)

	root.PersistentFlags().StringVar(
		&root.config.Output, "output", "text", "format to use for logs and console outputs",
	)

	cobrax.BindFlags(root.Command, v, envPrefix)

	// root.AddCommand()
	root.AddCommand(NewValidateCommand(&root.config.Output))
	root.AddCommand(NewVersionCommand(&root.config.Output, versions))

	return root
}
