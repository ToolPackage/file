package utils

func OrString(str ...string) string {
	for _, s := range str {
		if s != "" {
			return s
		}
	}
	return ""
}

func Min(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
