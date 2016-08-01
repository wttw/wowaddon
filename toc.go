package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func readToc(subdir string) (map[string]string, error) {
	ret := map[string]string{}
	fieldre := regexp.MustCompile(`##\s*([a-zA-Z0-9-]+)\s*:\s*(.*)$`)
	filename := filepath.Join(addonDir, subdir, subdir+".toc")
	file, err := os.Open(filename)
	if err != nil {
		return ret, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		match := fieldre.FindStringSubmatch(scanner.Text())
		if match != nil {
			ret[strings.ToLower(match[1])] = match[2]
		}
	}
	err = scanner.Err()
	if err != nil {
		return ret, err
	}
	return ret, nil
}
