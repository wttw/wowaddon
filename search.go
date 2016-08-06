package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var useRegex bool

func search(c *cli.Context) error {
	failed := color.New(color.FgRed).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()
	warn := color.New(color.FgYellow).SprintFunc()
	wowV := wowVersion()

	if len(c.Args()) == 0 {
		fmt.Printf("You must provide a string to search for\n")
		os.Exit(1)
	}
	loadCatalog()

	words := []string{}

	for _, w := range c.Args() {
		if useRegex {
			_, err := regexp.Compile(w)
			if err != nil {
				fmt.Printf("%s: '%s' isn't a valid regex: %s\n", failed("Error"), w, err.Error())
				os.Exit(1)
			}
			words = append(words, w)
		} else {
			words = append(words, regexp.QuoteMeta(w))
		}
	}

	pattern := "(?i:" + strings.Join(words, ")|(?i:") + ")"

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("%s: Failed to compile '%s': %s\n", failed("Internal Error"), pattern, err.Error())
	}

	for _, source := range catalog.Preference {
		for name, addon := range catalog.Sources[source].Addons {
			namematch := re.MatchString(name)
			titlematch := re.MatchString(addon.Title)
			notes := []string{}
			for _, n := range addon.Notes {
				match := re.FindAllStringIndex(n, -1)
				if match != nil {
					r := ""
					idx := 0
					for _, m := range match {
						if m[0] > idx {
							r += n[idx:m[0]]
						}
						r += warn(n[m[0]:m[1]])
						idx = m[1]
					}
					r += n[idx:]
					notes = append(notes, r)

				}
			}
			if namematch || titlematch || len(notes) > 0 {
				if addon.Interface >= wowV {
					fmt.Printf("%s: %s\n", success(name), addon.Title)
				} else {
					fmt.Printf("%s: %s (out of date)\n", failed(name), addon.Title)
				}
				for _, note := range notes {
					fmt.Printf("  %s\n", note)
				}
			}
		}
	}

	return nil
}
