package cli

import (
	"flag"
)

func FlagHandling() (bool, error) {
	initiate := flag.Bool("init", false, "adds premade data for db. To use it: --init")
	flag.Parse()

	return *initiate, nil
}
