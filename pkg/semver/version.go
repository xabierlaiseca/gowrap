package semver

import (
	"regexp"
	"strings"
)

var validSemVerRegex = regexp.MustCompile(`^[0-9]+(\.[0-9]+){0,2}$`)

func IsValid(semver string) bool {
	return validSemVerRegex.MatchString(semver)
}

// IsFullVersion returns true if the given version at least contains major, minor and patch segments.
func IsFullVersion(version string) bool {
	return IsValid(version) && strings.Count(version, ".") >= 2
}
