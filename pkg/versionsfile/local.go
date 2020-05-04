package versionsfile

import (
	"encoding/json"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xabierlaiseca/gowrap/pkg/cache"
)

const localVersionsCachedFile = "goversions.json"

type GoArchive struct {
	URL    string `json:"url,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
}

func Load() (map[string]GoArchive, error) {
	content, err := cache.Get(localVersionsCachedFile)
	if err != nil {
		logrus.Warningf("failed to get cached go versions file: %v", err)
	}

	archivesForPlatform := make(map[string]GoArchive)

	if content == nil {
		rvf, err := download()
		if err != nil {
			return nil, err
		}

		for version, pga := range rvf.getArchivesFor(runtime.GOARCH, runtime.GOOS) {
			archivesForPlatform[version] = pga.GoArchive
		}

		toCache, err := json.Marshal(archivesForPlatform)
		if err == nil {
			cache.Set(localVersionsCachedFile, toCache, 24*time.Hour)
		} else {
			logrus.Warningf("failed to serialise archives for caching: %v", err)
		}

	} else {
		err = json.Unmarshal(content, &archivesForPlatform)
	}

	return archivesForPlatform, err
}
