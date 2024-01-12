package versionsfile

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xabierlaiseca/gowrap/pkg/cache"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
)

const localVersionsCachedFile = "goversions.json"

type GoArchive struct {
	URL               string `json:"url,omitempty"`
	Checksum          string `json:"checksum,omitempty"`
	ChecksumAlgorithm string `json:"checksumAlgorithm,omitempty"`
}

const oneDay = 24 * time.Hour

func Load() (map[string]GoArchive, error) {
	content, err := cache.Get(localVersionsCachedFile)
	if err != nil {
		logrus.Warningf("failed to get cached go versions file: %v", err)
	}

	if content != nil {
		archivesForPlatform := make(map[string]GoArchive)
		err = json.Unmarshal(content, &archivesForPlatform)
		return archivesForPlatform, err
	}

	return DownloadToCache()
}

func DownloadToCache() (map[string]GoArchive, error) {
	rvf, err := download()
	if err != nil {
		return nil, err
	}

	archivesForPlatform := make(map[string]GoArchive)
	for version, pga := range rvf.getArchivesFor(runtime.GOARCH, runtime.GOOS) {
		if !semver.IsFullVersion(version) {
			version = fmt.Sprintf("%s.0", version)
		}
		archivesForPlatform[version] = pga.GoArchive
	}

	toCache, err := json.Marshal(archivesForPlatform)
	if err == nil {
		err = cache.Set(localVersionsCachedFile, toCache, oneDay)
		if err != nil {
			logrus.Warningf("failed to store local versions file: %v", err)
		}
	} else {
		logrus.Warningf("failed to serialise archives for caching: %v", err)
	}

	return archivesForPlatform, err
}
