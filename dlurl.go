package main

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/wttw/wowaddon/output"
)

func dlurl(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return cli.NewExitError(fmt.Sprintf("Usage: %s dlurl <addon_name>...", c.App.Name), 1)
	}
	for _, name := range c.Args() {
		meta, err := downloadURL(name, addonSource)
		if err != nil {
			output.Printf("%s: %s\n", name, err.Error())
		} else {
			output.Printf("%s: %s %s\n", name, meta.Version, meta.URL)
		}
	}
	return nil
}
