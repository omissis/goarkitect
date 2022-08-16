package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/cli"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

// configs contains a list of files where ruleset are specified
type configFiles []string

func (i *configFiles) String() string {
	return strings.Join(*i, ",")
}

func (i *configFiles) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func VerifyFactory() (cli.Command, error) {
	return &verifyCommand{}, nil
}

func getWd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return cwd
}

func listConfigFiles(cfs []string) []string {
	configFiles := make([]string, 0)

	for _, cf := range cfs {
		fileInfo, err := os.Stat(cf)
		if err != nil {
			log.Fatal(err)
		}

		if !fileInfo.IsDir() {
			configFiles = append(configFiles, cf)
			continue
		}

		if err := filepath.Walk(cf, visitConfigFolder(&configFiles)); err != nil {
			log.Fatal(err)
		}
	}

	slices.Sort(configFiles)

	return slices.Compact(configFiles)
}

func visitConfigFolder(files *[]string) filepath.WalkFunc {
	return func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !file.IsDir() && filepath.Ext(path) == ".yaml" {
			*files = append(*files, path)
		}

		return nil
	}
}

func loadConfig[T any](file string) T {
	var conf T

	configData, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(configData, &conf); err != nil {
		log.Fatal(err)
	}

	return conf
}
