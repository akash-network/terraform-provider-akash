package extensions

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
