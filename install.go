package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func installFromMeta(meta AddonMeta) ([]string, error) {
	ret := []string{}
	zipfile, err := getZipfile(meta.Name, meta.URL, meta.Source)
	if err != nil {
		return ret, err
	}
	reader, err := zip.OpenReader(zipfile)
	if err != nil {
		return ret, err
	}
	defer reader.Close()

	for _, f := range reader.Reader.File {
		zipped, err := f.Open()
		if err != nil {
			return ret, err
		}

		defer zipped.Close()
		elPath := filepath.Join(addonDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(elPath, f.Mode())
			dirname := strings.TrimRight(f.Name, "/")
			if !strings.Contains(dirname, "/") {
				seen := false
				for _, s := range ret {
					if s == dirname {
						seen = true
					}
				}
				if !seen {
					ret = append(ret, dirname)
				}
			}
		} else {
			// Some zip files don't include the containing directory
			// None of this is terribly efficient, but CPU isn't our limit
			slash := strings.LastIndex(f.Name, "/")
			dirname := f.Name[:slash]
			if !strings.Contains(dirname, "/") {
				seen := false
				for _, s := range ret {
					if s == dirname {
						seen = true
					}
				}
				if !seen {
					ret = append(ret, dirname)
				}
			}
			containingDir := filepath.Dir(elPath)
			_, err = os.Stat(containingDir)
			if err != nil {
				_ = os.MkdirAll(containingDir, 0755)

			}
			writer, err := os.OpenFile(elPath, os.O_WRONLY|os.O_CREATE, f.Mode())
			if err != nil {
				return ret, err
			}
			defer writer.Close()

			_, err = io.Copy(writer, zipped)
			if err != nil {
				return ret, err
			}
		}
	}
	return ret, nil
}

func installAddon(name string, source string, verb string) error {
	failed := color.New(color.FgRed).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()
	meta, err := downloadURL(name, source)
	if err != nil {
		fmt.Printf("%s: failed: %s\n", failed(name), err.Error())
	} else {
		subdirs, err := installFromMeta(meta)

		if err == nil {
			fmt.Printf("%s: %s\n", success(name), verb)
			addonCompat := 999999
			for _, dir := range subdirs {
				cf, err := readToc(dir)
				if err == nil {
					ver, ok := cf["interface"]
					if ok {
						nver, err := strconv.Atoi(ver)
						if err == nil && nver < addonCompat {
							addonCompat = nver
						}
					}
				}
			}
			if addonCompat == 999999 {
				addonCompat = 0
			}
			config.Addons[name] = Addon{
				Source:    meta.Source,
				Version:   meta.Version,
				Folders:   subdirs,
				Interface: addonCompat,
			}
			err = writeConfig()
			if err != nil {
				failed := color.New(color.FgRed).SprintFunc()
				fmt.Printf("%s: failed to write configuration file '%s': %s\n", failed(name), configFile, err.Error())
				return err
			}
		} else {
			fmt.Printf("%s: failed: %s\n", failed(name), err.Error())
		}
	}
	return nil
}

func install(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return cli.NewExitError(fmt.Sprintf("Usage: %s install <addon_name>...", c.App.Name), 1)
	}

	for _, name := range c.Args() {
		err := installAddon(name, addonSource, "installed")
		if err != nil {
			return err
		}
	}
	return nil
}

func reinstall(c *cli.Context) error {
	for name, meta := range config.Addons {
		err := installAddon(name, meta.Source, "reinstalled")
		if err != nil {
			return err
		}
	}
	return nil
}
