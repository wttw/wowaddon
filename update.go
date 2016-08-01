package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func update(c *cli.Context) error {
	failed := color.New(color.FgRed).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()
	var addons []string
	if len(c.Args()) == 0 {
		addons = make([]string, len(config.Addons))
		i := 0
		for k := range config.Addons {
			addons[i] = k
			i++
		}
	} else {
		addons = c.Args()
	}

	updated := 0
	for _, name := range addons {
		addon, ok := config.Addons[name]
		if !ok {
			fmt.Printf("%s: isn't installed\n", failed(name))
			continue
		}
		meta, err := downloadURL(name, addon.Source)
		if err != nil {
			fmt.Printf("%s: failed to retrieve metadata: %s\n", failed(name), err.Error())
			continue
		}
		if meta.Version <= addon.Version {
			fmt.Printf("%s: up to date at version %d\n", success(name), addon.Version)
			continue
		}
		err = installAddon(name, addon.Source, "updated")
		if err != nil {
			updated++
		}
	}
	fmt.Printf("%d addons updated\n", updated)
	return nil
}

func checkupdate(c *cli.Context) error {
	failed := color.New(color.FgRed).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()
	updated := 0
	for name, addon := range config.Addons {
		meta, err := downloadURL(name, addon.Source)
		if err != nil {
			fmt.Printf("%s: failed to retrieve metadata: %s\n", failed(name), err.Error())
			continue
		}
		if meta.Version > addon.Version {
			fmt.Printf("%s: can be updated from %d to %d\n", success(name), addon.Version, meta.Version)
			updated++
		}
	}
	fmt.Printf("%d addons can be updated\n", updated)
	return nil
}
