package main

import (
	"runtime"

	"github.com/pdbogen/wowaddon/output"
	"github.com/urfave/cli"
)

func environment(c *cli.Context) error {
	output.Printf("%s version:     %s %s/%s\n", c.App.Name, c.App.Version, runtime.GOOS, runtime.GOARCH)
	output.Printf("Configuration:        %s\n", configFile)
	output.Printf("WoW directory:        %s\n", wowDir)
	output.Printf("Catalog:              %s\n", catalogFile)
	output.Printf("Catalog fetched:      %s\n", config.CatalogDownloaded)
	output.Printf("Next catalog refresh: %s\n", config.NextCatalogUpdate)
	output.Printf("Cache directory:      %s\n", cacheDir)
	cf, err := readWowConfig()
	if err != nil {
		output.Printf("Interface version:    (failed to read configuration)\n")
	} else {
		version, ok := cf["lastAddonVersion"]
		if !ok {
			output.Printf("Interface version:    (not found)\n")
		} else {
			output.Printf("Interface version:    %s\n", version)
		}
	}
	return nil
}
