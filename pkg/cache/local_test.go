package cache

import (
	"testing"
	"time"
  "os"

  "io/ioutil"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
)

type Verification func(*testing.T, string, string)

func Test_Cache(t *testing.T) {
	testCases := map[string]struct {
		cache       map[string]time.Duration
		verifyFiles map[string]Verification
	}{
		"SetAndGetOne": {
		  cache: map[string]time.Duration{
        "test.txt": (60 * time.Minute),
      },
      verifyFiles: map[string]Verification {
        "test.txt": verifyExists,
      },
		},
    "GetNotExists": {
		  cache: map[string]time.Duration{},
      verifyFiles: map[string]Verification {
        "test.txt": verifyNotExists,
      },
		},
    "SetAndGetUntilExpires": {
		  cache: map[string]time.Duration{
        "test.txt": (50 * time.Millisecond),
      },
      verifyFiles: map[string]Verification {
        "test.txt": verifyExpires,
      },
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
      tmpDir, err := ioutil.TempDir(os.TempDir(), "test-cache-")
      require.NoError(t, err)

			for relpath, expiresIn := range testCase.cache {
        err := set(tmpDir, relpath, []byte(relpath), expiresIn)
        assert.NoError(t, err)
			}

      for relpath, verify := range testCase.verifyFiles {
        verify(t, tmpDir, relpath)
      }
		})
	}
}

func verifyExists(t *testing.T, rootDir, relpath string) {
  actual, err := get(rootDir, relpath)
  assert.NoError(t, err)
  assert.Equal(t, []byte(relpath), actual)
}

func verifyNotExists(t *testing.T, rootDir, relpath string) {
  actual, err := get(rootDir, relpath)
  assert.NoError(t, err)
  assert.Nil(t, actual)
}

func verifyExpires(t *testing.T, rootDir, relpath string) {
  verifyExists(t, rootDir, relpath)
  cacheEntryNotExists := func() bool {
    content, err := get(rootDir, relpath)
    return content == nil && err == nil
  }

  assert.Eventually(t, cacheEntryNotExists, 1 * time.Second, 50 * time.Millisecond)
}
