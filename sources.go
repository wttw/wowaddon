package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"strings"

	"github.com/urfave/cli"
)

// AddonMeta holds the metadata for an addon
type AddonMeta struct {
	Name        string
	URL         string
	Version     string
	Description string
	Source      string
}

// downloadURL returns the metadata for an addon, including the download location
func downloadURL(name string, source string) (AddonMeta, error) {
	switch source {
	case "curse":
		return curseDownloadURL(name)
	case "tukui":
		return tukuiDownloadURL(name)
	case "":
		if name == "tukui" || name == "elvui" {
			return tukuiDownloadURL(name)
		}
		return curseDownloadURL(name)
	default:
		return AddonMeta{}, cli.NewExitError("Bad source '%s'. Must be one of curse or tukui", 1)
	}
}

func tukString(project map[string]interface{}, name string) (string, error) {
	iface, ok := project[name]
	if !ok {
		return "", fmt.Errorf("Tukui response didn't include %s", name)
	}
	str, ok := iface.(string)
	if !ok {
		return "", fmt.Errorf("Tukui response for %s was unexpected type", name)
	}
	return str, nil
}

// tukuiDownloadURL gets the download URL and version for an addon from tukui
func tukuiDownloadURL(name string) (AddonMeta, error) {
	url := fmt.Sprintf("http://www.tukui.org/api.php?project=%s", name)
	resp, err := Get(url)
	if err != nil {
		return AddonMeta{}, err
	}
	defer resp.Body.Close()
	var projects []map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		return AddonMeta{}, err
	}
	if len(projects) < 1 {
		return AddonMeta{}, fmt.Errorf("Tukui didn't return metadata for %s", name)
	}
	ret := AddonMeta{
		Name:   name,
		Source: "tukui",
	}

	ret.URL, err = tukString(projects[0], "url")
	if err != nil {
		return AddonMeta{}, err
	}

	version, err := tukString(projects[0], "version")
	if err != nil {
		return AddonMeta{}, err
	}
	// parts := strings.Split(version, ".")
	// nver := 0
	// for _, part := range parts {
	// 	np, err := strconv.Atoi(part)
	// 	if err != nil {
	// 		return AddonMeta{}, fmt.Errorf("Tukui returned version '%s', which I couldn't parse: %s", version, err.Error())
	// 	}
	// 	nver = nver*1000 + np
	// }
	ret.Version = version
	return ret, nil
}

// curseDownloadURL gets the download URL and version for an addon from curse
func curseDownloadURL(name string) (AddonMeta, error) {
	url := fmt.Sprintf("http://www.curse.com/addons/wow/%s/download", name)
	resp, err := Get(url)
	if err != nil {
		return AddonMeta{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return AddonMeta{}, err
	}

	downloadRe := regexp.MustCompile(`data-href="(http:\/\/addons\.curse\.cursecdn\.com\/files\/[^\n"]*)"`)
	versionRe := regexp.MustCompile(`data-file="([0-9]+)"`)

	dlmatch := downloadRe.FindSubmatch(body)
	if dlmatch == nil {
		// No match for download URL
		return AddonMeta{}, fmt.Errorf("%s not found at Curse", name)
	}

	ret := AddonMeta{
		Name:   name,
		URL:    string(dlmatch[1]),
		Source: "curse",
	}

	vermatch := versionRe.FindSubmatch(body)
	if vermatch == nil {
		// No match for version
		return AddonMeta{}, fmt.Errorf("Version for %s not found at Curse", name)
	}

	ret.Version = string(vermatch[1])

	return ret, nil
}

//wowinterfaceURL gets the download URL and version for an addon from wowinterfaceUrl
func wowinterfaceURL(name string) (AddonMeta, error) {
	return AddonMeta{}, fmt.Errorf("wowinterface not implemented")
}

func getZipfile(name string, url string, source string) (string, error) {
	slash := strings.LastIndex(url, "/")
	if slash == -1 {
		return "", fmt.Errorf("Invalid URL for %s: %s", url, name)
	}
	zipname := fmt.Sprintf("%s__%s__%s", source, name, url[slash+1:])
	zipfile := filepath.Join(cacheDir, zipname)
	_, err := os.Stat(zipfile)
	if err == nil {
		return zipfile, nil
	}
	fmt.Printf("%s: Fetching from %s\n", name, url)

	tempfile, err := ioutil.TempFile(cacheDir, fmt.Sprintf("downloading-%s", name))
	if err != nil {
		return "", err
	}
	defer os.Remove(tempfile.Name())

	resp, err := Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(tempfile, resp.Body)
	if err != nil {
		return "", err
	}
	err = tempfile.Close()
	if err != nil {
		return "", err
	}

	err = os.Rename(tempfile.Name(), zipfile)
	if err != nil {
		return "", err
	}
	return zipfile, nil
}
