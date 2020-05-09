package customerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsNotFound(t *testing.T) {
	testCases := map[string]struct {
		input    error
		expected bool
	}{
		"NotFoundFnErr": {input: NotFound(), expected: true},
		"OtherErr":      {input: Error("message"), expected: false},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			actual := IsNotFound(testCase.input)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
