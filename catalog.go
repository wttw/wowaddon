package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/wttw/wowaddon/output"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// IAddon contains published data about a single addon
type IAddon struct {
	Folder    []string `json:"f"`
	Interface int      `json:"i"`
	Title     string   `json:"t"`
	Notes     []string `json:"n"`
}

// ISource contains all addons from a single source
type ISource struct {
	Addons map[string]IAddon `json:"addons"`
}

// Index contains all published data
type Index struct {
	Version    int                `json:"app_version"`
	Abort      string             `json:"abort"`
	Created    time.Time          `json:"created"`
	Preference []string           `json:"preference"`
	Sources    map[string]ISource `json:"source"`
	Downloaded time.Time          `json:"downloaded"`
}

var catalog Index

var catalogLoaded = false

func loadCatalogFromJSON() error {
	file, err := os.Open(catalogFile)
	if err == nil {
		defer file.Close()
		jsonParser := json.NewDecoder(file)
		err = jsonParser.Decode(&catalog)
		if err != nil {
			output.Printf("I couldn't parse catalog file '%s': %s\nMaybe fix it up, or delete it and start over?\n", configFile, err.Error())
			os.Exit(1)
		}
		if catalog.Abort != "" {
			output.Printf("%s\n", failed(catalog.Abort))
			os.Exit(1)
		}
		if catalog.Version > numericVersion {
			output.Printf("%s: There is an update to wowaddon available\nSee https://github.com/wttw/wowaddon/releases/latest\n", warn("Out of date"))
		}
		return nil
	}
	return err
}

func loadCatalogFromZip() error {
	z, err := zip.OpenReader(catalogFile + ".zip")
	if err != nil {
		return err
	}
	defer z.Close()

	for _, f := range z.File {
		if f.Name == "addoncatalog.json" {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			of, err := os.OpenFile(catalogFile, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				output.Printf("Couldn't create catalog file %s: %s\n", catalogFile, err.Error())
				os.Exit(1)
			}
			_, err = io.Copy(of, rc)
			if err != nil {
				output.Printf("Couldn't create catalog file %s: %s\n", catalogFile, err.Error())
				os.Exit(1)
			}
			err = of.Close()
			if err != nil {
				output.Printf("Couldn't create catalog file %s: %s\n", catalogFile, err.Error())
				os.Exit(1)
			}
			config.CatalogDownloaded = time.Now()
			config.NextCatalogUpdate = time.Now().Add(6 * time.Hour)
			err = writeConfig()
			if err != nil {
				output.Printf("Failed to save configuration: %s\n", err.Error())
				os.Exit(1)
			}
			err = loadCatalogFromJSON()
			if err != nil {
				return err
			}

			// mapping := bleve.NewIndexMapping()
			// index, err := bleve.New(catalogFile+".bleve", mapping)
			// if err != nil {
			// 	return err
			// }

			// for _, source := range catalog.Preference {
			// 	for name, data := range catalog.Sources[source].Addons {
			// 		err = index.Index(name, data)
			// 		if err != nil {
			// 			output.Printf("Problem indexing '%+v': %s\n", data, err.Error())
			// 		}
			// 	}
			// }

			// err = index.Close()
			// if err != nil {
			// 	output.Printf("Problem creating index: %s\n", err.Error())
			// }

			return nil
		}
	}
	return fmt.Errorf("No catalog zipfile found")
}

func fetchCatalog(current time.Time) error {
	url := "https://api.github.com/repos/wttw/wowaddon/releases/latest"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	// github ask that api users use a specific user-agent
	req.Header.Set("User-Agent", fmt.Sprintf("wttw/wowaddon (%s)", Version))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var release map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return err
	}

	assetsmap, ok := release["assets"]
	if !ok {
		return fmt.Errorf("No assets found in response")
	}

	// output.Printf("%T\n", assetsmap)
	// output.Printf("%+v\n", assetsmap)
	assets, ok := assetsmap.([]interface{})
	if !ok {
		return fmt.Errorf("Couldn't decode assets")
	}
	for _, assety := range assets {
		asset, ok := assety.(map[string]interface{})
		if !ok {
			continue
		}
		// Will panic if json isn't in the format we expect
		name := asset["name"].(string)
		if name != "addoncatalog.json.zip" {
			continue
		}
		dlurl := asset["browser_download_url"].(string)
		updatedAt := asset["updated_at"].(string)
		update, err := time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp '%s': %s", updatedAt, err.Error())
		}
		if current.After(update) {
			// not stale
			return nil
		}
		output.Printf("Fetching catalog from %s\n", dlurl)
		resp, err := Get(dlurl)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		tempfile, err := ioutil.TempFile(wowDir, "downloading-addoncatalog.json.zip")
		if err != nil {
			return err
		}
		defer os.Remove(tempfile.Name())
		_, err = io.Copy(tempfile, resp.Body)
		if err != nil {
			return err
		}
		err = tempfile.Close()
		if err != nil {
			return err
		}

		err = os.Rename(tempfile.Name(), catalogFile+".zip")
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("asset not found")
}

func loadCatalog() error {
	if catalogLoaded {
		return nil
	}

	err := loadCatalogFromJSON()
	if err != nil {
		err = loadCatalogFromZip()
		if err != nil {
			err = fetchCatalog(time.Unix(0, 0))
			if err != nil {
				output.Printf("Failed to fetch catalog: %s\n", err.Error())
				os.Exit(1)
			}
			err = loadCatalogFromZip()
			if err != nil {
				output.Printf("Failed to load catalog just fetched: %s\n", err.Error())
				os.Exit(1)
			}
		}
	}
	if config.NextCatalogUpdate.Before(time.Now()) {
		err = fetchCatalog(config.CatalogDownloaded)
		if err != nil {
			output.Printf("Failed to fetch catalog: %s\n", err.Error())
			os.Exit(1)
		}
		err = loadCatalogFromZip()
		if err != nil {
			output.Printf("Failed to load catalog just fetched: %s\n", err.Error())
			os.Exit(1)
		}
	}
	catalogLoaded = true
	return nil
}
