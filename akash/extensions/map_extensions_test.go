package extensions

import (
	"fmt"
	"testing"
)

var testCases = []struct {
	source   map[string]string
	subset   map[string]string
	expected bool
}{
	{map[string]string{"a": "va", "b": "vb"}, map[string]string{"a": "va"}, true},
	{map[string]string{"a": "va", "b": "vb"}, map[string]string{"c": "vc"}, false},
	{map[string]string{}, map[string]string{"c": "vc"}, false},
	{map[string]string{}, map[string]string{}, true},
}

func TestIsSubset(t *testing.T) {
	t.Run("empty map is always subset", func(t *testing.T) {
		if !IsSubset(map[int32]int32{1: 1}, map[int32]int32{}) {
			t.Fatalf("An empty map is always a subset of any map")
		}
	})

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v is subset of %v should be %v", tc.subset, tc.source, tc.expected), func(t *testing.T) {
			if IsSubset(tc.source, tc.subset) != tc.expected {
				t.Fatalf("Expected %v got %v", tc.expected, IsSubset(tc.source, tc.subset))
			}
		})
	}
}
