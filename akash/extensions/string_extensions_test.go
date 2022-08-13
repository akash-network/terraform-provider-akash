package extensions

import "testing"

func TestFindAll(t *testing.T) {
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
