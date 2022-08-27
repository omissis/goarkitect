package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/omissis/goarkitect/internal/logx"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

var ErrNoOutputFormat = errors.New("output cannot be nil")
var ErrNoConfigFileFound = errors.New("no config files found")

func getWd() string {
	cwd, err := os.Getwd()
	if err != nil {
		logx.Fatal(err)
	}

	return cwd
}

func listConfigFiles(cfs []string) []string {
	configFiles := make([]string, 0)

	for _, cf := range cfs {
		fileInfo, err := os.Stat(cf)
		if err != nil {
			logx.Fatal(err)
		}

		if !fileInfo.IsDir() {
			configFiles = append(configFiles, cf)
			continue
		}

		if err := filepath.Walk(cf, visitConfigFolder(&configFiles)); err != nil {
			logx.Fatal(err)
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
		logx.Fatal(err)
	}

	if err := yaml.Unmarshal(configData, &conf); err != nil {
		logx.Fatal(err)
	}

	return conf
}
