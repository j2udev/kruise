package kruise

import "os"

type void struct{}

func contains[T comparable](list []T, t T) bool {
	for _, l := range list {
		if l == t {
			return true
		}
	}
	return false
}

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
