package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/wttw/wowaddon/output"
)

func uninstall(c *cli.Context) error {
	for _, name := range c.Args() {
		addon, ok := config.Addons[name]
		if !ok {
			output.Printf("%s: %s: wasn't installed by me\n", name, failed("failed"))
		} else {
			for _, d := range addon.Folders {
				// For each folder in the addon we're removing
				unused := true
				var usedby string
				for cname, caddon := range config.Addons {
					if cname == name {
						continue
					}
					for _, cdir := range caddon.Folders {
						if cdir == d {
							unused = false
							usedby = cname
						}
					}
				}
				if unused {
					dir := filepath.Join(addonDir, d)
					err := os.RemoveAll(dir)
					if err != nil {
						output.Printf("%s: %s: failed to remove directory %s: %s\n", name, failed("failed"), dir, err.Error())
					} else {
						output.Printf("%s: directory %s %s\n", name, d, success("removed"))
					}
				} else {
					output.Printf("%s: directory %s not removed, also used by %s\n", name, d, usedby)
				}
			}
			delete(config.Addons, name)
			err := writeConfig()
			if err != nil {
				output.Printf("%s: Failed to update configuration file %s: %s", failed("failed"), configFile, err.Error())
			} else {
				output.Printf("%s: %s\n", name, success("uninstalled"))
			}
		}
	}
	if !config.KeepCache {
		purgeCache()
	}
	return nil
}
