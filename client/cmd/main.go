package main

import (
	"fmt"
	"os"
)

var usage = `version: 0.0.1-SNAPSHOT
usage: fsec [-hncmHbfs] [-a apiUrl]
`

func main() {
	args := parseArguments()

	if args.HasArg("login") {
		// TODO: opt name conflicts with group name

	}
}

func parseArguments() (args *Args) {
	parser := NewArgParser()
	parser.
		Options(Option{
			opt:          "h",
			help:         "help",
			required:     false,
			defaultValue: false,
		}).
		Group("login",
			Option{
				opt:          "s",
				help:         "server address",
				required:     true,
				defaultValue: "",
			},
			Option{
				opt:          "c",
				help:         "credential",
				required:     true,
				defaultValue: "",
			})
	parser.Usage(func() {
		fmt.Print(usage)
		parser.PrintDefaultUsage()
	})

	args = parser.Parse()

	if args.HasArg("h") {
		parser.PrintUsage()
		os.Exit(0)
	}
	return args
}
