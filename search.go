package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli"
	"github.com/wttw/wowaddon/output"
)

var useRegex bool

func search(c *cli.Context) error {
	wowV := wowVersion()

	if len(c.Args()) == 0 {
		output.Printf("You must provide a string to search for\n")
		os.Exit(1)
	}
	loadCatalog()

	words := []string{}

	for _, w := range c.Args() {
		if useRegex {
			_, err := regexp.Compile(w)
			if err != nil {
				output.Printf("%s: '%s' isn't a valid regex: %s\n", failed("Error"), w, err.Error())
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
		output.Printf("%s: Failed to compile '%s': %s\n", failed("Internal Error"), pattern, err.Error())
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
						r += highlight(n[m[0]:m[1]])
						idx = m[1]
					}
					r += n[idx:]
					notes = append(notes, r)

				}
			}
			if namematch || titlematch || len(notes) > 0 {
				if addon.Interface >= wowV {
					output.Printf("%s: %s\n", success(name), addon.Title)
				} else {
					output.Printf("%s: %s (out of date)\n", failed(name), addon.Title)
				}
				for _, note := range notes {
					output.Printf("  %s\n", note)
				}
			}
		}
	}

	return nil
}
