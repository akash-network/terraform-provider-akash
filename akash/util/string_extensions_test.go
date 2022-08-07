package util

import "testing"

func TestContains(t *testing.T) {
	testSlice := []string{"A", "B", "C", "D", "test", "F"}

	t.Run("should return true if slice contains value", func(t *testing.T) {
		if !Contains(testSlice, "test") {
			t.Errorf("Expected slice to contain %s", testSlice)
		}
	})

	t.Run("should return false if slice contains value", func(t *testing.T) {
		if Contains(testSlice, "Non-existent test") {
			t.Errorf("Expected slice not to contain %s", testSlice)
		}
	})
}
