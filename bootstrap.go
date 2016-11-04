package main

import (
	"io/ioutil"
	"strings"

	"github.com/urfave/cli"
	"github.com/wttw/wowaddon/output"
)

func bootstrap(c *cli.Context) error {
	return bootstrapConfig()
}

func bootstrapConfig() error {
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
		// output.Printf("dir %s found\n", f.Name())
		dirs[f.Name()] = false
	}

	for _, source := range catalog.Preference {
		// output.Printf("catalog source: %s\n", source)
		for name, addon := range catalog.Sources[source].Addons {
			if len(addon.Folder) == 0 {
				continue
			}
			possible := true
			someunused := false
			for _, dir := range addon.Folder {
				// output.Printf("Checking for %s\n", dir)
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
					output.Printf("Multiple addons (%s) use directories %s - check configuration\n", name, strings.Join(addon.Folder, ", "))
				} else {
					// output.Printf("Found addon %s\n", name)
					config.Addons[name] = Addon{
						Source:  source,
						Folders: addon.Folder,
						Version: "unknown",
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
	output.Printf("%s\n", success("Configuration created"))
	return nil
}
