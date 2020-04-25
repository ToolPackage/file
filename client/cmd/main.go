package main

import (
	"flag"
	"fmt"
	"os"
)

var usage = `Version: 0.0.1-SNAPSHOT
Usage: fsec [-hncmHbfs] [-a apiUrl]
Options:
`

func main() {
	parseArguments()
}

func parseArguments() {
	help := flag.Bool("h", false, "help")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, usage)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	cmd := flag.Arg(0)
	fmt.Printf(cmd)
}