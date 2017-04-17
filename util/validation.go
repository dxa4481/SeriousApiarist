package util

import "regexp"

// ValidName validates letters, numbers, dash, period, and underscore 0 to 8 chars
func ValidName(str string) bool {
	var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9\-._]{0,80}$`)
	return nameRegex.MatchString(str)
}

// ValidRef validates letters, numbers, dash, and forward slash 3 to 40 chars
func ValidRef(str string) bool {
	var refRegex = regexp.MustCompile(`^[a-zA-Z0-9\-/]{3,40}$`)
	return refRegex.MatchString(str)
}

// ValidInt validates integers from 1 to 8 chars
func ValidInt(str string) bool {
	var intRegex = regexp.MustCompile(`^[0-9]{1,8}$`)
	return intRegex.MatchString(str)
}
