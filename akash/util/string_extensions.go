package util

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsAny(strings []string, e []string) bool {
	ch := make(chan bool, len(e))

	for _, str := range e {
		go func(s string) {
			ch <- Contains(strings, s)
		}(str)
	}

	for range e {
		contains := <-ch
		if contains {
			return true
		}
	}

	return false
}

func FindAll(strings []string, e []string) []string {
	ch := make(chan string, len(e))

	for _, str := range e {
		go func(s string) {
			if Contains(strings, s) {
				ch <- s
			} else {
				ch <- ""
			}
		}(str)
	}

	finds := make([]string, 0)

	for range e {
		s := <-ch
		if s != "" {
			finds = append(finds, s)
		}
	}

	return finds
}
