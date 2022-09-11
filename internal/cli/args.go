package cli

import (
	"flag"
	"strings"
)

func GetArgs(args []string, low int) []string {
	if flag.Lookup("test.v") != nil {
		nargs := normalizeArgs(args)
		cnt := countTestFlags(nargs)

		if len(nargs) > low+cnt {
			return nargs[low+cnt:]
		}

		// Getting here means the tests are invoked using unexpected flags.
		return []string{}
	}

	return args[low:]
}

func normalizeArgs(args []string) []string {
	nargs := make([]string, 0)

	for _, arg := range args {
		if len(arg) > 0 && arg[0] == '-' {
			nargs = append(nargs, strings.Split(arg, "=")...)
		} else {
			nargs = append(nargs, arg)
		}
	}

	return nargs
}

func countTestFlags(args []string) int {
	count := 0

	for i, arg := range args {
		if len(arg) >= 6 && arg[0:6] == "-test." {
			count++

			if i+1 <= len(args) && args[i+1][0] != '-' {
				count++
			}
		}
	}

	return count
}
