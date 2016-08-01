package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func dlurl(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return cli.NewExitError(fmt.Sprintf("Usage: %s dlurl <addon_name>...", c.App.Name), 1)
	}
	for _, name := range c.Args() {
		meta, err := downloadURL(name, addonSource)
		if err != nil {
			fmt.Printf("%s: %s\n", name, err.Error())
		} else {
			fmt.Printf("%s: %d %s\n", name, meta.Version, meta.URL)
		}
	}
	return nil
}
