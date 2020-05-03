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
	OS     string
	ARCH   string
	URL    string
	SHA256 string
}

func Load() (map[string]GoArchive, error) {
	content, err := cache.Get(localVersionsCachedFile)
	if err != nil {
		logrus.Warningf("failed to get cached go versions file: %v", err)
	}

	if content == nil {
		rvf, err := download()
		if err != nil {
			return nil, err
		}

		archivesForPlatform := rvf.getArchivesFor(runtime.GOARCH, runtime.GOOS)
		toCache, err := json.Marshal(archivesForPlatform)
		if err == nil {
			cache.Set(localVersionsCachedFile, toCache, 24*time.Hour)
		} else {
			logrus.Warningf("failed to serialise archives for caching: %v", err)
		}

		return archivesForPlatform, nil
	}

	var archivesForPlatform map[string]GoArchive
	err = json.Unmarshal(content, &archivesForPlatform)
	return archivesForPlatform, err
}
