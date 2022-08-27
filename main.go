package main

import (
	"log"

	"github.com/omissis/goarkitect/cmd"
)

var (
	Version   string = "unknown"
	GitCommit string = "unknown"
	BuildTime string = "unknown"
	GoVersion string = "unknown"
	OsArch    string = "unknown"
)

func main() {
	versions := map[string]string{
		"version":   Version,
		"gitCommit": GitCommit,
		"buildTime": BuildTime,
		"goVersion": GoVersion,
		"osArch":    OsArch,
	}

	if err := cmd.NewRootCommand(versions).Execute(); err != nil {
		log.Fatal(err)
	}
}
