package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/osext"
	"github.com/urfave/cli"
)

const configFilename = "addons.json"
const cacheDirname = "ZipFiles"
const catalogFilename = "addoncatalog.json"

var wowDir string
var addonDir string
var configFile string
var addonSource string
var cacheDir string
var catalogFile string

// Addon holds the configuration and state for a single addon
type Addon struct {
	Source    string   `json:"source"`
	Version   string   `json:"version"`
	Folders   []string `json:"folders"`
	Interface int      `json:"interface"`
	Zip       string   `json:"zipfile"`
}

// Config holds the configuration file
type Config struct {
	KeepCache         bool             `json:"keepcache"`
	NextCatalogUpdate time.Time        `json:"next_catalog_update"`
	CatalogDownloaded time.Time        `json:"catalog_retrieved"`
	Addons            map[string]Addon `json:"addons"`
}

var config = Config{
	Addons: map[string]Addon{},
}

func main() {
	app := cli.NewApp()
	app.Name = "wowaddon"
	app.Usage = "Install WoW addons"
	app.Version = "0.2.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Steve Atkins",
			Email: "steve@blighty.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "wowdir, dir, d",
			Usage:       "WoW base directory",
			EnvVar:      "WOWDIR",
			Destination: &wowDir,
		},
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Use an alternate configuration file",
			EnvVar:      "WOW_ADDON_CONFIG",
			Destination: &configFile,
		},
		cli.StringFlag{
			Name:        "cache",
			Usage:       "Use an alternate cache directory",
			EnvVar:      "WOW_ADDON_CACHE",
			Destination: &cacheDir,
		},
	}
	app.Before = setup
	app.Commands = []cli.Command{
		{
			Name:   "install, i",
			Usage:  "Install addon `NAME`",
			Action: install,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "source, s",
					Usage:       "Install from `SOURCE`",
					Destination: &addonSource,
				},
			},
		},
		{
			Name:   "update, u",
			Usage:  "Update all addons",
			Action: update,
		},
		{
			Name:   "search, s",
			Usage:  "Search for new addons",
			Action: search,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "pattern, regexp, regex, r",
					Usage:       "use regular expressions",
					Destination: &useRegex,
				},
			},
		},
		{
			Name:   "uninstall",
			Usage:  "Uninstall addon `NAME`",
			Action: uninstall,
		},
		{
			Name:   "reinstall",
			Usage:  "Reinstall all addons",
			Action: reinstall,
		},
		{
			Name:   "checkupdate",
			Usage:  "List addons that can be updated",
			Action: checkupdate,
		},
		{
			Name:    "folders",
			Aliases: []string{"list", "ls"},
			Usage:   "List addons and their folders",
			Action:  list,
		},
		{
			Name:   "blame",
			Usage:  "Show which addon created a folder",
			Action: blame,
		},
		{
			Name:    "environment",
			Aliases: []string{"env"},
			Usage:   "Show environment",
			Action:  environment,
		},
		{
			Name:   "info",
			Usage:  "Show information about installed addons",
			Action: info,
		},
		{
			Name:   "fullinfo",
			Usage:  "Show toc metadata about installed addons",
			Action: fullinfo,
		},
		{
			Name:   "bootstrap",
			Usage:  "Create a configuration file from existing addons",
			Action: bootstrap,
		},
		{
			Name:   "dlurl",
			Usage:  "Find a download URL",
			Action: dlurl,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "source, s",
					Usage:       "Show URLs for `SOURCE`",
					Destination: &addonSource,
				},
			},
		},
		{
			Name:   "releasetag",
			Usage:  "display release tag",
			Hidden: true,
			Action: func(c *cli.Context) error {
				fmt.Printf("v%s\n", app.Version)
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func setup(*cli.Context) error {
	err := findBaseDir()
	if err != nil {
		return err
	}
	addonDir = filepath.Join(wowDir, "Interface", "Addons")
	err = os.MkdirAll(addonDir, 0755)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Can't create directory '%s': %s", addonDir, err.Error()), 1)
	}
	if configFile == "" {
		configFile = filepath.Join(wowDir, configFilename)
	}

	if cacheDir == "" {
		cacheDir = filepath.Join(wowDir, "Interface", cacheDirname)
	}

	if catalogFile == "" {
		catalogFile = filepath.Join(wowDir, catalogFilename)
	}

	cf, err := os.Open(configFile)
	if err == nil {
		defer cf.Close()
		jsonParser := json.NewDecoder(cf)
		err = jsonParser.Decode(&config)
		if err != nil {
			fmt.Printf("I couldn't parse configuration file '%s': %s\nMaybe fix it up, or delete it and start over?\n", configFile, err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Printf("Setting up your configuration...\n")
		err = bootstrapConfig()
		if err != nil {
			fmt.Printf("Failed to set up configuration: %s\n", err.Error())
		}
	}

	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		fmt.Printf("Can't create directory '%s': %s\n", cacheDir, err.Error())
		os.Exit(1)
	}

	return nil
}

func writeConfig() error {
	tmpfile := fmt.Sprintf("%s.tmp", configFile)
	backupfile := fmt.Sprintf("%s.bak", configFile)
	out, err := os.OpenFile(tmpfile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	j, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		out.Close()
		return err
	}

	_, err = out.Write(j)
	out.Close()
	if err != nil {
		return err
	}

	_, err = os.Stat(configFile)
	if err == nil {
		// Current config file exists, probably
		_ = os.Remove(backupfile)
		err = os.Rename(configFile, backupfile)
		if err != nil {
			return err
		}
	}
	err = os.Rename(tmpfile, configFile)
	return err
}

func checkDir(dir string) bool {
	iface := filepath.Join(dir, "Interface")
	fi, err := os.Stat(iface)
	if err != nil {
		return false
	}
	if !fi.IsDir() {
		return false
	}
	wowDir = dir
	return true
}

func findBaseDir() error {
	if wowDir != "" {
		return nil
	}
	for _, dir := range installDirs {
		if checkDir(dir) {
			return nil
		}
	}
	wd, err := os.Getwd()
	if err == nil {
		if checkDir(wd) {
			return nil
		}
	}
	bind, err := osext.ExecutableFolder()
	if err == nil {
		if checkDir(bind) {
			return nil
		}
	}
	return cli.NewExitError("I can't find the World Of Warcraft folder\nSet the environment variable WOWDIR or use --wowdir", 1)
}
