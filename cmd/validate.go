package cmd

import (
	"path/filepath"

	"github.com/omissis/goarkitect/cmd/validate"
	"github.com/omissis/goarkitect/internal/schema/santhosh"
	"github.com/spf13/cobra"
)

func NewValidateCommand(output *string) *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate the configuration file(s)",
		RunE: func(_ *cobra.Command, args []string) error {
			if output == nil {
				return ErrNoOutputFormat
			}

			basePath := getWd()
			schema, err := santhosh.LoadSchema(basePath)
			if err != nil {
				return err
			}

			if len(args) == 0 {
				args = append(args, filepath.Join(basePath, ".goarkitect.yaml"))
			}

			cfs := listConfigFiles(args)
			if len(cfs) == 0 {
				return ErrNoConfigFileFound
			}

			hasErrors := error(nil)
			for _, cf := range cfs {
				conf := loadConfig[any](cf)

				if err := schema.ValidateInterface(conf); err != nil {
					validate.PrintResults(*output, err, conf, cf)

					hasErrors = validate.ErrHasValidationErrors
				}
			}

			return hasErrors
		},
	}
}
