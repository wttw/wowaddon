package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func lock(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return fmt.Errorf("You must provide an addon to lock")
	}
	addons := c.Args()
	for _, name := range addons {
		addon, ok := config.Addons[name]
		if !ok {
			fmt.Printf("%s: isn't installed\n", failed(name))
			continue
		}
		if addon.Locked {
			fmt.Printf("%s: already locked\n", warn(name))
		} else {
			addon.Locked = true
			config.Addons[name] = addon
			fmt.Printf("%s: locked\n", success(name))
		}
	}
	return writeConfig()
}

func unlock(c *cli.Context) error {
	addons := []string{}
	if len(c.Args()) == 0 {
		for name, addon := range config.Addons {
			if addon.Locked {
				addons = append(addons, name)
			}
		}
	} else {
		addons = c.Args()
	}

	for _, name := range addons {
		addon, ok := config.Addons[name]
		if !ok {
			fmt.Printf("%s: isn't installed\n", failed(name))
			continue
		}
		if addon.Locked {
			addon.Locked = false
			config.Addons[name] = addon
			fmt.Printf("%s: unlocked\n", success(name))
		} else {
			fmt.Printf("%s: already unlocked\n", warn(name))
		}
	}
	return writeConfig()
}
