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
		availableVersions  map[string][]GoArchive
		expectedGoArchives map[string]GoArchive
	}{
		"NoVersionsForCurrentPlatform": {
			availableVersions: map[string][]GoArchive{
				"1.2.3": []GoArchive{
					{ARCH: otherARCH, OS: otherOS},
				},
			},
			expectedGoArchives: make(map[string]GoArchive),
		},
		"OneVersionsForCurrentPlatform": {
			availableVersions: map[string][]GoArchive{
				"1.2.3": []GoArchive{
					{ARCH: otherARCH, OS: otherOS},
					{ARCH: thisARCH, OS: thisOS},
				},
			},
			expectedGoArchives: map[string]GoArchive{
				"1.2.3": GoArchive{ARCH: thisARCH, OS: thisOS},
			},
		},
		// // "VersionExistsButNoArchivesForPlatform": {
		// //
		// // }
		// "ArchiveExistsForPlatform": {
		// 	availableVersions: map[string][]GoArchive{"1.3.4": []GoArchive{
		// 		{OS: runtime.GOOS, ARCH: runtime.GOARCH},
		// 	}},
		// 	expectedGoArchive: &GoArchive{OS: runtime.GOOS, ARCH: runtime.GOARCH},
		// },
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			underTest := RemoteVersionsFile{
				versions: testCase.availableVersions,
			}

			result := underTest.GetArchivesFor(thisARCH, thisOS)
			assert.Equal(t, testCase.expectedGoArchives, result)
		})
	}
}
