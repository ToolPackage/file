package utils

func OrString(str ...string) string {
	for _, s := range str {
		if s != "" {
			return s
		}
	}
	return ""
}
