package kruise

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
)

// contains is used to generically determine whether an object is contained
// within a slice of other objects
func contains[T comparable](list []T, t T) bool {
	for _, l := range list {
		if l == t {
			return true
		}
	}
	return false
}

// containsAny is used to generically determine whether any object given is
// contained within a slice of other objects
func containsAny[T comparable](list []T, any ...T) bool {
	for _, l := range list {
		for _, t := range any {
			if l == t {
				return true
			}
		}
	}
	return false
}

// exists is used to determine whether a file or directory already exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// captureStdout is used to captureStdout from another function and return it
// in a string; this is useful for testing
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		Logger.Error(err)
	}
	return buf.String()
}

// trimDeployPrefix is used to trim the absolute path prefix from any
// underlying commands used by Deploy
func trimDeployPrefix(actual string) string {
	actuals := strings.Split(actual, "\n")
	reg := regexp.MustCompile(`^\/.+(helm |kubectl )`)
	for i, a := range actuals {
		actuals[i] = reg.ReplaceAllString(a, "${1}")
	}
	return strings.Join(actuals, "\n")
}

// trimDeployStdoutPrefix is used to trim the absolute path prefix from the
// stdout of the given function
func trimDeployStdoutPrefix(f func()) string {
	out := captureStdout((f))
	return trimDeployPrefix(out)
}
