package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func bootstrap(c *cli.Context) error {
	return bootstrapConfig()
}

func bootstrapConfig() error {
	success := color.New(color.FgGreen).SprintFunc()
	err := loadCatalog()
	if err != nil {
		return err
	}

	dirs := map[string]bool{}

	files, err := ioutil.ReadDir(addonDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		// fmt.Printf("dir %s found\n", f.Name())
		dirs[f.Name()] = false
	}

	for _, source := range catalog.Preference {
		// fmt.Printf("catalog source: %s\n", source)
		for name, addon := range catalog.Sources[source].Addons {
			if len(addon.Folder) == 0 {
				continue
			}
			possible := true
			someunused := false
			for _, dir := range addon.Folder {
				// fmt.Printf("Checking for %s\n", dir)
				used, ok := dirs[dir]
				if !ok {
					possible = false
					break
				}
				if !used {
					someunused = true
				}
			}
			if possible {
				if !someunused {
					fmt.Printf("Multiple addons (%s) use directories %s - check configuration\n", name, strings.Join(addon.Folder, ", "))
				} else {
					// fmt.Printf("Found addon %s\n", name)
					config.Addons[name] = Addon{
						Source:  source,
						Folders: addon.Folder,
					}
					populateVersion(name)
					for _, dir := range addon.Folder {
						dirs[dir] = true
					}
				}
			}
		}
	}
	err = writeConfig()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", success("Configuration created"))
	return nil
}
