package category

import "testing"

func TestNewCategory(t *testing.T) {
	c := Category(ID: "123")

	if c == nil {
		t.Fatal("category should not be nil")
	}
}

