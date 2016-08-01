package main

import (
	"fmt"
	"runtime"

	"github.com/urfave/cli"
)

func environment(c *cli.Context) error {
	fmt.Printf("%s version:  %s %s/%s\n", c.App.Name, c.App.Version, runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Configuration:     %s\n", configFile)
	fmt.Printf("WoW directory:     %s\n", wowDir)
	cf, err := readWowConfig()
	if err != nil {
		fmt.Printf("Interface version: (failed to read configuration)\n")
	} else {
		version, ok := cf["lastAddonVersion"]
		if !ok {
			fmt.Printf("Interface version: (not found)\n")
		} else {
			fmt.Printf("Interface version: %s\n", version)
		}
	}
	return nil
}
