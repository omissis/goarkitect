package main

import (
	"log"

	"github.com/omissis/goarkitect/cmd"
)

var (
	version   = "unknown"
	gitCommit = "unknown"
	buildTime = "unknown"
	goVersion = "unknown"
	osArch    = "unknown"
)

func main() {
	versions := map[string]string{
		"version":   version,
		"gitCommit": gitCommit,
		"buildTime": buildTime,
		"goVersion": goVersion,
		"osArch":    osArch,
	}

	if err := cmd.NewRootCommand(versions).Execute(); err != nil {
		log.Fatal(err)
	}
}
