package semver

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_IsLessThan(t *testing.T) {
	testCases := map[string]struct {
		semver1  string
		semver2  string
		expected bool
	}{
		"SameMayorAndMinorAndPatch": {
			semver1:  "1.2.4",
			semver2:  "1.2.4",
			expected: false,
		},
		"SameMayorAndMinorButPatchLessThan": {
			semver1:  "1.2.3",
			semver2:  "1.2.4",
			expected: true,
		},
		"SameMayorAndMinorButPatchGreaterThan": {
			semver1:  "1.2.5",
			semver2:  "1.2.4",
			expected: false,
		},
		"SameMayorAndMinorButPatchMissingInFirst": {
			semver1:  "1.2",
			semver2:  "1.2.4",
			expected: true,
		},
		"SameMayorAndMinorButPatchMissingInSecond": {
			semver1:  "1.2.4",
			semver2:  "1.2",
			expected: false,
		},
		"SameMayorAndMinor": {
			semver1:  "1.2",
			semver2:  "1.2",
			expected: false,
		},
		"SameMayorButMinorLessThan": {
			semver1:  "1.1.7",
			semver2:  "1.2.4",
			expected: true,
		},
		"SameMayorButMinorGreaterThan": {
			semver1:  "1.3.2",
			semver2:  "1.2.4",
			expected: false,
		},
		"SameMayorButMinorMissingInFirst": {
			semver1:  "1",
			semver2:  "1.2",
			expected: true,
		},
		"SameMayorButMinorMissingInSecond": {
			semver1:  "1.2",
			semver2:  "1",
			expected: false,
		},
		"SameMayor": {
			semver1:  "1",
			semver2:  "1",
			expected: false,
		},
		"MayorLessThan": {
			semver1:  "1.4.7",
			semver2:  "2.2.4",
			expected: true,
		},
		"MayorGreaterThan": {
			semver1:  "2.1.2",
			semver2:  "1.2.4",
			expected: false,
		},
		"ComparesMayorAsNumber": {
			semver1:  "2",
			semver2:  "10",
			expected: true,
		},
		"ComparesMinorAsNumber": {
			semver1:  "1.2",
			semver2:  "1.10",
			expected: true,
		},
		"ComparesPatchAsNumber": {
			semver1:  "1.1.2",
			semver2:  "1.1.10",
			expected: true,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			actual := IsLessThan(testCase.semver1, testCase.semver2)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func Test_SliceStableComparatorFor_ValidVersions(t *testing.T) {
	semvers := []string{"2", "1.20.1", "1.3", "1.20.4", "1.19.1"}
	expected := []string{"1.3", "1.19.1", "1.20.1", "1.20.4", "2"}

	comparator, err := SliceStableComparatorFor(semvers)
	require.NoError(t, err)

	sort.SliceStable(semvers, comparator)

	assert.Equal(t, expected, semvers)
}

func Test_SliceStableComparatorFor_InvalidVersions(t *testing.T) {
	semvers := []string{"2", "1.20.a"}

	_, err := SliceStableComparatorFor(semvers)
	assert.EqualError(t, err, "invalid semantic version: 1.20.a")
}

func Test_Latest(t *testing.T) {
	testCases := map[string]struct {
		input    []string
		expected string
	}{
		"OneVersion": {
			input:    []string{"1.14.2"},
			expected: "1.14.2",
		},
		"TwoVersionsFirstIsLatest": {
			input:    []string{"1.14.2", "1.3.2"},
			expected: "1.14.2",
		},
		"TwoVersionsSecondIsLatest": {
			input:    []string{"1.3.2", "1.14.2"},
			expected: "1.14.2",
		},
		"MultipleVersions": {
			input:    []string{"1.3.2", "1.14.2", "1.13.2", "1.14.1"},
			expected: "1.14.2",
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			actual, err := Latest(testCase.input)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func Test_Latest_NoVersionsProvided(t *testing.T) {
	_, err := Latest([]string{})
	assert.Error(t, err)
}

func Test_HasPrefix(t *testing.T) {
	testCases := map[string]struct {
		version  string
		prefix   string
		expected bool
	}{
		"MayorPrefix": {
			version:  "1.14",
			prefix:   "1",
			expected: true,
		},
		"ComparesMayorAsNumber": {
			version:  "10",
			prefix:   "1",
			expected: false,
		},
		"MayorAndMinorPrefix": {
			version:  "1.14.2",
			prefix:   "1.14",
			expected: true,
		},
		"ComparesMinorAsNumber": {
			version:  "1.14.2",
			prefix:   "1.1",
			expected: false,
		},
		"ComparesPatchAsNumber": {
			version:  "1.14.12",
			prefix:   "1.14.1",
			expected: false,
		},
		"SameVersionOnlyMayor": {
			version:  "1",
			prefix:   "1",
			expected: true,
		},
		"SameVersionMayorAndMinor": {
			version:  "1.12",
			prefix:   "1.12",
			expected: true,
		},
		"SameVersionMayorMinorAndPatch": {
			version:  "1.12.1",
			prefix:   "1.12.1",
			expected: true,
		},
		"PrefixWithDot": {
			version:  "1.14",
			prefix:   "1.",
			expected: true,
		},

		"LongerPrefixThanVersion": {
			version:  "1.14",
			prefix:   "1.14.1",
			expected: false,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			actual := HasPrefix(testCase.version, testCase.prefix)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
