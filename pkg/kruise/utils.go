package kruise

import "os"

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
