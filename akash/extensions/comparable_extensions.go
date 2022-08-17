package extensions

func Union[T comparable](a []T, b []T) []T {
	ch := make(chan T, len(b))

	for _, bElem := range b {
		go func(t T) {
			if Contains(a, t) {
				ch <- t
			} else {
				var zeroValue T
				ch <- zeroValue
			}
		}(bElem)
	}

	finds := make([]T, 0)
	var zeroValue T

	for range b {
		s := <-ch
		if s != zeroValue {
			finds = append(finds, s)
		}
	}

	return finds
}

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsAny[T comparable](a []T, b []T) bool {
	ch := make(chan bool, len(b))

	for _, str := range b {
		go func(s T) {
			ch <- Contains(a, s)
		}(str)
	}

	for range b {
		contains := <-ch
		if contains {
			return true
		}
	}

	return false
}
