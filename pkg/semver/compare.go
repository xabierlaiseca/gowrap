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
		return "", customerrors.New("no versions provided")
	}

	latest := versions[0]
	for i := 1; i < len(versions); i++ {
		if isOlder(latest, versions[i]) {
			latest = versions[i]
		}
	}
	return latest, nil
}

func isOlder(semver1, semver2 string) bool {
	split1 := strings.SplitN(semver1, ".", 2)
	split2 := strings.SplitN(semver2, ".", 2)
	currentSegment1, _ := strconv.Atoi(split1[0])
	currentSegment2, _ := strconv.Atoi(split2[0])

	if currentSegment1 != currentSegment2 {
		return currentSegment1 < currentSegment2
	}

	switch {
	case len(split1) == 1 && len(split2) == 1:
		return false
	case len(split1) == 1:
		return true
	case len(split2) == 1:
		return false
	default:
		return isOlder(split1[1], split2[1])
	}
}
