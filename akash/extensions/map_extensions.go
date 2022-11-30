package extensions

import (
	"errors"
)

// IsSubset checks whether a given map is a subset of another. Return false if an entry in the subset is found that is not
// present in the source map. An empty map is a subset of any map.
func IsSubset[K, V comparable](source, subset map[K]V) bool {
	for k, v := range subset {
		if source[k] != v {
			return false
		}
	}

	return true
}

func SafeCastMapValues[K, V comparable](m map[K]interface{}) (out map[K]V, err error) {
	out = make(map[K]V, len(m))

	for k, v := range m {
		out[k] = v.(V)
	}

	defer func() {
		if r := recover(); r != nil {
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			out = nil
		}
	}()

	return out, nil
}
