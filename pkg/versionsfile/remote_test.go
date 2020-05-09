package versionsfile

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	thisARCH = runtime.GOARCH
	thisOS   = runtime.GOOS

	otherARCH string
	otherOS   string
)

func init() {
	if runtime.GOARCH == "amd64" {
		otherARCH = "386"
	} else {
		otherARCH = "amd64"
	}

	if runtime.GOOS == "linux" {
		otherOS = "darwin"
	} else {
		otherOS = "linux"
	}
}

func Test_RemoteVersionsFile_GetArchivesFor(t *testing.T) {
	testCases := map[string]struct {
		availableVersions  map[string][]platformGoArchive
		inputARCH          string
		inputOS            string
		expectedGoArchives map[string]platformGoArchive
	}{
		"NoVersionsForCurrentPlatform": {
			availableVersions: map[string][]platformGoArchive{
				"1.2.3": {
					{ARCH: otherARCH, OS: otherOS},
				},
			},
			inputARCH:          thisARCH,
			inputOS:            thisOS,
			expectedGoArchives: make(map[string]platformGoArchive),
		},
		"OneVersionsForCurrentPlatform": {
			availableVersions: map[string][]platformGoArchive{
				"1.2.3": {
					{ARCH: otherARCH, OS: otherOS},
					{ARCH: thisARCH, OS: thisOS},
				},
			},
			inputARCH: thisARCH,
			inputOS:   thisOS,
			expectedGoArchives: map[string]platformGoArchive{
				"1.2.3": {ARCH: thisARCH, OS: thisOS},
			},
		},
		"ArchivesForOtherPlatform": {
			availableVersions: map[string][]platformGoArchive{
				"1.2.3": {
					{ARCH: otherARCH, OS: otherOS},
					{ARCH: thisARCH, OS: thisOS},
				},
			},
			inputARCH: otherARCH,
			inputOS:   otherOS,
			expectedGoArchives: map[string]platformGoArchive{
				"1.2.3": {ARCH: otherARCH, OS: otherOS},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			underTest := remoteVersionsFile{
				versions: testCase.availableVersions,
			}

			result := underTest.getArchivesFor(testCase.inputARCH, testCase.inputOS)
			assert.Equal(t, testCase.expectedGoArchives, result)
		})
	}
}
