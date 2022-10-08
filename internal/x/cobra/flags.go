package cobra

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitEnvs(envPrefix string) *viper.Viper {
	v := viper.New()

	v.SetEnvPrefix(envPrefix)

	v.AutomaticEnv()

	return v
}

func BindFlags(cmd *cobra.Command, v *viper.Viper, logger func(v error), envPrefix string) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))

			env := envSuffix
			if envPrefix != "" {
				env = fmt.Sprintf("%s_%s", envPrefix, envSuffix)
			}

			if err := v.BindEnv(f.Name, env); err != nil {
				logger(err)
			}
		}

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)

			if err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
				logger(err)
			}
		}
	})
}
