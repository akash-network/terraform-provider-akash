package praetor

import "testing"

func TestPraetor_UnwrapAttributes(t *testing.T) {

	t.Run("should unwrap all attributes", func(t *testing.T) {
		attributes := []Attribute{
			{Key: "test1", Value: "value1"},
			{Key: "test2", Value: "value2"},
			{Key: "test3", Value: "value3"},
		}

		unwrapped := unwrapAttributes(attributes)

		if len(unwrapped) != len(attributes) {
			t.Errorf("Expected unwrapped to have the same length as attributes slice. Expected %d got %d", len(attributes), len(unwrapped))
		}

		for _, attr := range attributes {
			if unwrapped[attr.Key] == "" {
				t.Errorf("Key %s was not found in unwrapped attributes", attr.Key)
			}

			if unwrapped[attr.Key] != attr.Value {
				t.Errorf("Expected key %s to have value %s got %s instead", attr.Key, attr.Value, unwrapped[attr.Key])
			}
		}
	})

	t.Run("should overwrite attributes with the same key", func(t *testing.T) {
		attributes := []Attribute{
			{Key: "test1", Value: "value1"},
			{Key: "test2", Value: "value2"},
			{Key: "test1", Value: "value-over"},
		}

		unwrapped := unwrapAttributes(attributes)

		if len(unwrapped) != 2 {
			t.Errorf("Expected unwrapped attributes to have length %d got %d", 2, len(unwrapped))
		}

		if unwrapped["test1"] != "value-over" {
			t.Errorf("Expected key %s to have value %s got %s", "test1", "value-over", unwrapped["test1"])
		}
	})
}
