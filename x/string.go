package x

import (
	"hash/fnv"
	"regexp"
)

var regInputCheck = regexp.MustCompile(`^[a-zA-Z0-9].$`)

func InputCheck(input string) bool {
	return regInputCheck.MatchString(input)
}

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
