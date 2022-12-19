package helpers

func Contains(haystack []string, needle string) bool {
	if needle == "" || len(haystack) == 0 {
		return false
	}
	for _, c := range haystack {
		if c == needle {
			return true
		}
	}
	return false
}
