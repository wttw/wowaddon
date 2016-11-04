package main

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/urfave/cli"
	"github.com/wttw/wowaddon/output"
)

func wowVersion() int {
	cf, err := readWowConfig()
	wowV := 0
	if err == nil {
		version, ok := cf["lastAddonVersion"]
		if ok {
			nver, err := strconv.Atoi(version)
			if err == nil {
				wowV = nver
			}
		}
	}
	return wowV
}

func list(c *cli.Context) error {
	wowV := wowVersion()

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

	for name, addon := range config.Addons {
		installed := true
		for _, d := range addon.Folders {
			_, ok := dirs[d]
			if !ok {
				installed = false
			}
			dirs[d] = true
		}
		if !installed {
			output.Printf("%s: not installed\n", failed(name))
			continue
		}
		locked := ""
		if addon.Locked {
			locked = "(locked) "
		}
		if wowV != 0 && addon.Interface != 0 && addon.Interface < wowV {
			output.Printf("%s: %s(out of date) version %s installed\n", warn(name), locked, addon.Version)
		} else {
			output.Printf("%s: %sversion %s installed\n", success(name), locked, addon.Version)
		}
	}
	orphans := []string{}
	for dirname, seen := range dirs {
		if !seen {
			orphans = append(orphans, dirname)
		}
	}
	if len(orphans) > 0 {
		output.Printf("%s: %s\n", warn("Unmanaged addon directories"), strings.Join(orphans, ", "))
	}
	return nil
}

func fullinfo(c *cli.Context) error {
	addons := c.Args()
	if len(addons) == 0 {
		for name := range config.Addons {
			addons = append(addons, name)
		}
	}
	for _, name := range addons {
		addon, ok := config.Addons[name]
		if !ok {
			output.Printf("%s: not installed\n", failed(name))
			continue
		}
		output.Printf("%s: version %s\n", success(name), addon.Version)
		for _, dir := range addon.Folders {
			toc, err := readToc(dir)
			if err != nil {
				output.Printf("  %s: failed to read toc: %s\n", failed(dir), err.Error())
				continue
			}
			output.Printf("  %s:\n", success(dir))
			for k, v := range toc {
				output.Printf("    %s: %s\n", k, v)
			}
		}
	}
	return nil
}

func info(c *cli.Context) error {
	addons := c.Args()
	if len(addons) == 0 {
		for name := range config.Addons {
			addons = append(addons, name)
		}
	}
	for _, name := range addons {
		addon, ok := config.Addons[name]
		if !ok {
			output.Printf("%s: not installed\n", failed(name))
			continue
		}
		output.Printf("%s: version %s\n", success(name), addon.Version)
		for _, dir := range addon.Folders {

			toc, err := readToc(dir)
			if err != nil {
				output.Printf("  %s: failed to read toc: %s\n", failed(dir), err.Error())
				continue
			}

			output.Printf("  %s:", success(dir))
			ver, ok := toc["version"]
			if ok {
				output.Printf(" version: %s", ver)
			}
			iface, ok := toc["interface"]
			if ok {
				output.Printf(" compatible: %s", iface)
			}
			output.Printf("\n")
			notes, ok := toc["notes"]
			if ok {
				output.Printf("    %s\n", notes)
			}
		}
	}
	return nil
}
