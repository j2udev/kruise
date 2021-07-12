package utils

import (
	"testing"
)

func TestContainsTruthy(t *testing.T) {
	result := Contains([]string{"this", "is", "a", "test"}, "test")
	if result != true {
		t.Fatalf("Expected true")
	}
}

func TestContainsFalsy(t *testing.T) {
	result := Contains([]string{"this", "is", "a", "test"}, "doesnthaveme")
	if result != false {
		t.Fatalf("Expected false")
	}
}
