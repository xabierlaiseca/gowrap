package semver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IsValid(t *testing.T) {
	testCases := map[string]struct {
		semver   string
		expected bool
	}{
		"Mayor": {
			semver:   "1",
			expected: true,
		},
		"MayorWithMultipleDigits": {
			semver:   "21",
			expected: true,
		},
		"MayorEndsWithDot": {
			semver:   "1.",
			expected: false,
		},
		"MayorWithInvalidChar": {
			semver:   "2a",
			expected: false,
		},
		"MayorAndMinor": {
			semver:   "1.3",
			expected: true,
		},
		"MayorAndMinorWithMultipleDigits": {
			semver:   "1.34",
			expected: true,
		},
		"MayorAndMinorEndsWithDot": {
			semver:   "1.2.",
			expected: false,
		},
		"MayorAndMinorWithInvalidChar": {
			semver:   "2.2a",
			expected: false,
		},
		"MayorAndMinorAndPatch": {
			semver:   "1.1.3",
			expected: true,
		},
		"MayorAndMinorAndPatchWithMultipleDigits": {
			semver:   "1.1.34",
			expected: true,
		},
		"MayorAndMinorAmdPatchEndsWithDot": {
			semver:   "1.2.4.",
			expected: false,
		},
		"MayorAndMinorAndPatchWithInvalidChar": {
			semver:   "2.1.2a",
			expected: false,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			actual := IsValid(testCase.semver)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func Test_IsFullVersion(t *testing.T) {
	testCases := map[string]struct {
		semver   string
		expected bool
	}{
		"Major": {
			semver:   "1",
			expected: false,
		},
		"MajorAndMinor": {
			semver:   "1.2",
			expected: false,
		},
		"MajorMinorAndPatch": {
			semver:   "1.2.3",
			expected: true,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			actual := IsFullVersion(testCase.semver)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
