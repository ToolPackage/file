package internal

import (
	"os"
	"runtime"
)

func getUserHomeDir() string {
	home := os.Getenv("HOME")
	if home != "" {
		return home
	}

	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
		if home != "" {
			return home
		}

		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home != "" {
			return home
		}
	}

	panic("could not detect home directory")
}
