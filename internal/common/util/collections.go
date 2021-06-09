package util

func Contains(s []string, e interface{}) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
