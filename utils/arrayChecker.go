package utils

func ContainsString(list []string, target string) bool {
	for _, element := range list {
		if element == target {
			return true
		}
	}
	return false
}
