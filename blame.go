package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func blame(c *cli.Context) error {
	dirs := map[string]struct{}{}
	for _, d := range c.Args() {
		dirs[d] = struct{}{}
	}
	for name, meta := range config.Addons {
		for _, dir := range meta.Folders {
			_, ok := dirs[dir]
			if ok {
				fmt.Printf("%s: %s\n", dir, name)
			}
		}
	}

	return nil
}
