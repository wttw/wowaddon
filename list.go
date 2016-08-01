package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func list(c *cli.Context) error {
	failed := color.New(color.FgRed).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()
	warn := color.New(color.FgYellow).SprintFunc()
	cf, err := readWowConfig()
	wowVersion := 0
	if err == nil {
		version, ok := cf["lastAddonVersion"]
		if ok {
			nver, err := strconv.Atoi(version)
			if err == nil {
				wowVersion = nver
			}
		}
	}
	for name, addon := range config.Addons {
		installed := true
		for _, d := range addon.Folders {
			dirpath := filepath.Join(addonDir, d)
			fi, err := os.Stat(dirpath)
			if err != nil {
				installed = false
			} else {
				if !fi.IsDir() {
					installed = false
				}
			}
		}
		if !installed {
			fmt.Printf("%s: not installed\n", failed(name))
			continue
		}
		if wowVersion != 0 && addon.Interface != 0 && addon.Interface < wowVersion {
			fmt.Printf("%s: (out of date) version %d installed in %s\n", warn(name), addon.Version, strings.Join(addon.Folders, ", "))
		} else {
			fmt.Printf("%s: version %d installed in %s\n", success(name), addon.Version, strings.Join(addon.Folders, ", "))
		}
	}
	return nil
}

func fullinfo(c *cli.Context) error {
	failed := color.New(color.FgRed).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()
	addons := c.Args()
	if len(addons) == 0 {
		for name := range config.Addons {
			addons = append(addons, name)
		}
	}
	for _, name := range addons {
		addon, ok := config.Addons[name]
		if !ok {
			fmt.Printf("%s: not installed\n", failed(name))
			continue
		}
		fmt.Printf("%s: version %d\n", success(name), addon.Version)
		for _, dir := range addon.Folders {
			toc, err := readToc(dir)
			if err != nil {
				fmt.Printf("  %s: err\n", failed(dir))
				continue
			}
			fmt.Printf("  %s:\n", success(dir))
			for k, v := range toc {
				fmt.Printf("    %s: %s\n", k, v)
			}
		}
	}
	return nil
}

func info(c *cli.Context) error {
	failed := color.New(color.FgRed).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()
	addons := c.Args()
	if len(addons) == 0 {
		for name := range config.Addons {
			addons = append(addons, name)
		}
	}
	for _, name := range addons {
		addon, ok := config.Addons[name]
		if !ok {
			fmt.Printf("%s: not installed\n", failed(name))
			continue
		}
		fmt.Printf("%s: version %d\n", success(name), addon.Version)
		for _, dir := range addon.Folders {
			toc, err := readToc(dir)
			if err != nil {
				fmt.Printf("  %s: err\n", failed(dir))
				continue
			}

			fmt.Printf("  %s:", success(dir))
			ver, ok := toc["version"]
			if ok {
				fmt.Printf(" version: %s", ver)
			}
			iface, ok := toc["interface"]
			if ok {
				fmt.Printf(" compatible: %s", iface)
			}
			fmt.Printf("\n")
			notes, ok := toc["notes"]
			if ok {
				fmt.Printf("    %s\n", notes)
			}
		}
	}
	return nil
}
