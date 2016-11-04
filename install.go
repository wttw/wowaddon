package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
	"github.com/wttw/wowaddon/output"
)

func installFromMeta(meta AddonMeta) (Addon, error) {
	ret := Addon{
		Source:  meta.Source,
		Version: meta.Version,
	}
	subdirs := []string{}
	zipfile, err := getZipfile(meta.Name, meta.URL, meta.Source)
	if err != nil {
		return ret, err
	}
	ret.Zip = filepath.Base(zipfile)
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

		filename := strings.Replace(f.Name, `\`, `/`, -1)

		elPath := filepath.Join(addonDir, filepath.FromSlash(filename))
		if f.FileInfo().IsDir() {
			os.MkdirAll(elPath, f.Mode())
			dirname := strings.TrimRight(filename, "/")
			if !strings.Contains(dirname, "/") {
				seen := false
				for _, s := range subdirs {
					if s == dirname {
						seen = true
					}
				}
				if !seen {
					subdirs = append(subdirs, dirname)
				}
			}
		} else {
			// Some zip files don't include the containing directory
			// None of this is terribly efficient, but CPU isn't our limit
			containingDir := filepath.Dir(elPath)
			_, err = os.Stat(containingDir)
			if err != nil {
				_ = os.MkdirAll(containingDir, 0755)
			}
			pathParts := strings.Split(filename, "/")
			if len(pathParts) > 0 {
				dirname := pathParts[0]
				seen := false
				for _, s := range subdirs {
					if s == dirname {
						seen = true
					}
				}
				if !seen {
					subdirs = append(subdirs, dirname)
				}
			}
			writer, err := os.OpenFile(elPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return ret, err
			}
			defer writer.Close()

			_, err = io.Copy(writer, zipped)
			if err != nil {
				return ret, err
			}
			err = os.Chtimes(elPath, time.Now(), f.ModTime())
			if err != nil {
				return ret, err
			}
		}
	}
	ret.Folders = subdirs
	return ret, nil
}

func populateVersion(name string) {
	ao, ok := config.Addons[name]
	if ok {
		addonCompat := 999999
		for _, dir := range ao.Folders {
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
		ao.Interface = addonCompat
		config.Addons[name] = ao
	}
}

func installAddon(name string, source string, verb string) error {
	meta, err := downloadURL(name, source)
	if err != nil {
		output.Printf("%s: failed: %s\n", failed(name), err.Error())
	} else {
		ao, err := installFromMeta(meta)

		if err == nil {
			config.Addons[name] = ao
			populateVersion(name)

			err = writeConfig()
			if err != nil {
				output.Printf("%s: failed to write configuration file '%s': %s\n", failed(name), configFile, err.Error())
				return err
			}
		} else {
			output.Printf("%s: failed: %s\n", failed(name), err.Error())
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
	if !config.KeepCache {
		purgeCache()
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
	if !config.KeepCache {
		purgeCache()
	}
	return nil
}

func purgeCache() {
	zipFiles := map[string]struct{}{}
	for _, addon := range config.Addons {
		zipFiles[addon.Zip] = struct{}{}
	}
	cacheFiles, err := ioutil.ReadDir(cacheDir)
	if err != nil {
		output.Printf("Error reading cache directory: %s\n", err.Error())
		return
	}
	for _, file := range cacheFiles {
		if !file.IsDir() {
			_, ok := zipFiles[file.Name()]
			if !ok {
				_ = os.Remove(filepath.Join(cacheDir, file.Name()))
			}
		}
	}
}
