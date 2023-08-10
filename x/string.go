package x

import (
	"regexp"
)

var regInputCheck = regexp.MustCompile(`^[a-zA-Z0-9].$`)

func InputCheck(input string) bool {
	return regInputCheck.MatchString(input)
}
