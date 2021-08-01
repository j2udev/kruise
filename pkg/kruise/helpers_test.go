package kruise

import (
	"testing"
)

func TestContains(t *testing.T) {
	truthyResult := contains([]string{"this", "is", "a", "test"}, "test")
	if truthyResult != true {
		t.Fatalf("Expected true")
	}
	falsyResult := contains([]string{"this", "is", "a", "test"}, "doesnthaveme")
	if falsyResult != false {
		t.Fatalf("Expected false")
	}
}
