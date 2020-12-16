package main

import (
	"github.com/ToolPackage/fse/client"
	"github.com/jessevdk/go-flags"
	"os"
)

func main() {
	_, err := flags.Parse(&internal.Opts)
	if err != nil {
		os.Exit(1)
	}
}
