package kruise

import (
	"log"
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
