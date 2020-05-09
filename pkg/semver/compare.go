package semver

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
)

func SliceStableComparatorFor(semvers []string) (func(int, int) bool, error) {
	for _, semver := range semvers {
		if !IsValid(semver) {
			return nil, customerrors.Errorf("invalid semantic version: %s", semver)
		}
	}

	return func(i1, i2 int) bool {
		return isOlder(semvers[i1], semvers[i2])
	}, nil
}

var validSemVerRegex = regexp.MustCompile(`^[0-9]+(\.[0-9]+){0,2}$`)

func IsValid(semver string) bool {
	return validSemVerRegex.MatchString(semver)
}

func Latest(versions []string) (string, error) {
	if len(versions) == 0 {
		return "", customerrors.Error("no versions provided")
	}

	latest := versions[0]
	for i := 1; i < len(versions); i++ {
		if isOlder(latest, versions[i]) {
			latest = versions[i]
		}
	}
	return latest, nil
}

func HasPrefix(version, prefix string) bool {
	switch {
	case len(prefix) == 0:
		return true
	case len(version) == 0:
		return false
	}

	firstSegmentVersion, restVersion := splitFirstSegment(version)
	firstSegmentPrefix, restPrefix := splitFirstSegment(prefix)

	if firstSegmentVersion != firstSegmentPrefix {
		return false
	}

	return HasPrefix(restVersion, restPrefix)
}

func isOlder(semver1, semver2 string) bool {
	firstSegment1, rest1 := splitFirstSegment(semver1)
	firstSegment2, rest2 := splitFirstSegment(semver2)

	if firstSegment1 != firstSegment2 {
		return firstSegment1 < firstSegment2
	}

	switch {
	case len(rest1) == 0 && len(rest2) == 0:
		return false
	case len(rest1) == 0:
		return true
	case len(rest2) == 0:
		return false
	default:
		return isOlder(rest1, rest2)
	}
}

func splitFirstSegment(version string) (int, string) {
	split := strings.SplitN(version, ".", 2)
	segment, _ := strconv.Atoi(split[0])

	if len(split) == 1 {
		return segment, ""
	}

	return segment, split[1]
}
