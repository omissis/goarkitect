package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/omissis/goarkitect/cmd/cmdutil"
	"github.com/omissis/goarkitect/cmd/validate"
	"github.com/omissis/goarkitect/internal/schema/santhosh"
)

func NewValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate the configuration file(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			output := cmd.Flag("output").Value.String()

			basePath := cmdutil.GetWd()
			schema, err := santhosh.LoadSchema(basePath)
			if err != nil {
				return fmt.Errorf("failed to load schema: %w", err)
			}

			if len(args) == 0 {
				args = append(args, filepath.Join(basePath, ".goarkitect.yaml"))
			}

			cfs := cmdutil.ListConfigFiles(args)
			if len(cfs) == 0 {
				return cmdutil.ErrNoConfigFileFound
			}

			hasErrors := error(nil)
			for _, cf := range cfs {
				conf := cmdutil.LoadConfig[any](cf)

				if err := schema.ValidateInterface(conf); err != nil {
					validate.PrintResults(output, err, conf, cf)

					hasErrors = validate.ErrHasValidationErrors
				}
			}

			validate.PrintSummary(output, hasErrors != nil)

			return hasErrors
		},
	}
}
