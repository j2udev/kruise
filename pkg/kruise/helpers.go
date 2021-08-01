package kruise

import (
	"log"
	"strings"
)

// checkErr is used to generically handle error catching
func checkErr(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

// contains determines whether a string is contained within a slice of strings
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// collectValidArgs is used to filter a slice of human-readable options into a
// slice of strings to be used with the Cobra Command ValidArgs slice
func collectValidArgs(opts []Option) []string {
	var collector []string
	for _, opt := range opts {
		collector = append(collector, strings.Split(opt.Arguments, ", ")...)
	}
	return collector
}
