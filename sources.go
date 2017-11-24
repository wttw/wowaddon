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
	"github.com/wttw/wowaddon/output"
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
			return uiDownloadURL(name)
		}
		return curseDownloadURL(name)
	default:
		return AddonMeta{}, cli.NewExitError("Bad source '%s'. Must be one of curse or tukui", 1)
	}
}

// Corner Case: tukui and elvui cannot use regular tukui API
func uiDownloadURL(name string) (AddonMeta, error) {
	url := fmt.Sprintf("https://www.tukui.org/download.php?ui=%s",name)
	resp, err := Get(url)
	if err != nil {
		return AddonMeta{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return AddonMeta{}, err
	}

	downloadRe := regexp.MustCompile(fmt.Sprintf(`/downloads/%s-([0-9]+\.[0-9]+)\.zip`,name))

	dlmatch := downloadRe.FindSubmatch(body)
	if dlmatch == nil {
		return AddonMeta{}, fmt.Errorf("%s not found at Tukui", name)
	}

	ret := AddonMeta{
		Name:    name,
		URL:     fmt.Sprintf("https://www.tukui.org%s",string(dlmatch[0])),
		Source:  "tukui",
		Version: string(dlmatch[1]),
	}

	return ret, nil
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
	resp, err := Get("https://www.tukui.org/api.php?addons=all")
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
	var project map[string]interface{}

	for _, proj := range projects {
		addonName, err := tukString(proj, "name")
		if err != nil {
			continue
		}
		if addonName == name {
			project = proj
			break
		}
	}

	if project == nil {
		return AddonMeta{}, fmt.Errorf("Tukui could not find addon %s", name)
	}

	ret := AddonMeta{
		Name:   name,
		Source: "tukui",
	}

	ret.URL, err = tukString(project, "url")
	if err != nil {
		return AddonMeta{}, err
	}

	ret.Version, err = tukString(project, "version")
	if err != nil {
		return AddonMeta{}, err
	}

	return ret, nil
}

// curseDownloadURL gets the download URL and version for an addon from curse
func curseDownloadURL(name string) (AddonMeta, error) {
	url := fmt.Sprintf("https://www.curseforge.com/wow/addons/%s/download", name)
	resp, err := Get(url)
	if err != nil {
		return AddonMeta{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return AddonMeta{}, err
	}

	downloadRe := regexp.MustCompile(fmt.Sprintf(`%s\/download/([0-9]+)/file`,name))

	dlmatch := downloadRe.FindSubmatch(body)
	if dlmatch == nil {
		// No match for download URL
		return AddonMeta{}, fmt.Errorf("%s not found at Curse", name)
	}

	ret := AddonMeta{
		Name:    name,
		URL:     fmt.Sprintf("https://www.curseforge.com/wow/addons/%s",string(dlmatch[0])),
		Source:  "curse",
		Version: string(dlmatch[1]),
	}

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
	output.Printf("%s: Fetching from %s\n", name, url)

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
