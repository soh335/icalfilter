package main

import (
	"flag"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/soh335/icalparser"
	"github.com/soh335/sliceflag"
)

var (
	filterComponents = sliceflag.String(flag.CommandLine, "filter-component", []string{}, "filter component")
	filterLines      = sliceflag.String(flag.CommandLine, "filter-line", []string{}, "filter line")
)

func main() {
	flag.Parse()
	if err := _main(os.Stdin, os.Stdout, *filterComponents, *filterLines); err != nil {
		log.Fatal(err)
	}
}

func _main(r io.Reader, w io.Writer, filterComponents, filterLines []string) error {
	obj, err := icalparser.NewParser(r).Parse()
	if err != nil {
		return err
	}

	filterComponentRegexpList, filterLineRegexpList, err := setupflag(filterComponents, filterLines)
	if err != nil {
		return err
	}

	components := []*icalparser.Component{}
OUTER:
	for _, component := range obj.Components {
		properties := []*icalparser.ContentLine{}
	INNER:
		for _, property := range component.PropertiyLines {
			str := property.String()
			if isMatchInList(str, filterComponentRegexpList) {
				continue OUTER
			}
			if isMatchInList(str, filterLineRegexpList) {
				continue INNER
			}
			properties = append(properties, property)
		}

		component.PropertiyLines = properties
		components = append(components, component)
	}

	obj.Components = components

	icalparser.NewPrinter(obj).WriteTo(w)

	return nil
}

func setupflag(filterComponents, filterLines []string) (filterComponentRegexpList []*regexp.Regexp, filterLineRegexpList []*regexp.Regexp, err error) {

	for _, s := range filterComponents {
		var r *regexp.Regexp
		r, err = regexp.Compile(s)
		if err != nil {
			return
		}
		filterComponentRegexpList = append(filterComponentRegexpList, r)
	}

	for _, s := range filterLines {
		var r *regexp.Regexp
		r, err = regexp.Compile(s)
		if err != nil {
			return
		}
		filterLineRegexpList = append(filterLineRegexpList, r)
	}

	return
}

func isMatchInList(str string, regexpList []*regexp.Regexp) bool {
	for _, r := range regexpList {
		if r.MatchString(str) {
			return true
		}
	}
	return false
}
