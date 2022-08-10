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

func TestContainsAny(t *testing.T) {
	testSlice := []string{"A", "B", "C", "D", "test", "F"}

	t.Run("should return true if slice contains any value", func(t *testing.T) {
		if !ContainsAny(testSlice, []string{"yes", "test", "hello", "there"}) {
			t.Errorf("Expected slice to contain %s", testSlice)
		}
	})

	t.Run("should return false if slice contains any value", func(t *testing.T) {
		if ContainsAny(testSlice, []string{"Non-existent test"}) {
			t.Errorf("Expected slice not to contain %s", testSlice)
		}
	})
}

func TestFindAlly(t *testing.T) {
	testSlice := []string{"A", "B", "C", "D", "test", "F"}

	t.Run("should return a slice containing finds", func(t *testing.T) {
		finds := FindAll(testSlice, []string{"yes", "test", "A", "there"})
		if len(finds) != 2 {
			t.Errorf("Expected to find %d but found %d (%v)", 2, len(finds), finds)
		}

		if !Contains(testSlice, finds[0]) || !Contains(testSlice, finds[1]) {
			t.Fail()
		}
	})
}
