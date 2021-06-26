package utils

import (
	"testing"

	u "github.com/j2udevelopment/kruise/pkg/utils"
)

func TestContainsTruthy(t *testing.T) {
	result := u.Contains([]string{"this", "is", "a", "test"}, "test")
	if result != true {
		t.Fatalf("Expected true")
	}
}

func TestContainsFalsy(t *testing.T) {
	result := u.Contains([]string{"this", "is", "a", "test"}, "doesnthaveme")
	if result != false {
		t.Fatalf("Expected false")
	}
}
