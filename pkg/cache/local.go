package cache

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	relCacheRootDir = "gowrap"
	relMetadataDir  = "metadata"
	relObjectsDir   = "objects"
)

var (
	relCachedObjectsMetadataFile = filepath.Join(relMetadataDir, "objects.json")
)

// Get returns the content of previously cached if not expired.
// If no error is returned and bytes is set to nil, it means that either no
// content was never cached or the content expired.
func Get(relpath string) ([]byte, error) {
	rootDir, err := getRootDir()
	if err != nil {
		return nil, err
	}

	return get(rootDir, relpath)
}

func get(rootDir, relpath string) ([]byte, error) {
	metadata, err := readCacheMetadata(rootDir)
	if err != nil {
		return nil, err
	}

	cachedObjectPath := filepath.Join(rootDir, relObjectsDir, relpath)
	if expiresAt, ok := metadata[relpath]; ok && expiresAt.Before(time.Now()) {
		err := os.Remove(cachedObjectPath)
		if err != nil {
			logrus.Error(err.Error())
		}

		delete(metadata, relpath)
		return nil, storeCacheMetadata(rootDir, metadata)
	} else if !ok {
		return nil, nil
	}

	content, err := ioutil.ReadFile(cachedObjectPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	return content, err
}

// Set stores the given content in provided cache path for the requested duration
func Set(relpath string, content []byte, expiresIn time.Duration) error {
	rootDir, err := getRootDir()
	if err != nil {
		return err
	}

	return set(rootDir, relpath, content, expiresIn)
}

func set(rootDir, relpath string, content []byte, expiresIn time.Duration) error {
	metadata, err := readCacheMetadata(rootDir)
	if err != nil {
		return err
	}

	objectPath := filepath.Join(rootDir, relObjectsDir, relpath)
	if err := os.MkdirAll(filepath.Dir(objectPath), 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(objectPath, content, 0644); err != nil {
		return err
	}

	metadata[relpath] = time.Now().Add(expiresIn)
	return storeCacheMetadata(rootDir, metadata)
}

func readCacheMetadata(rootDir string) (map[string]time.Time, error) {
	cachedObjectsMetadataFile := filepath.Join(rootDir, relCachedObjectsMetadataFile)
	content, err := ioutil.ReadFile(cachedObjectsMetadataFile)
	if os.IsNotExist(err) {
		return make(map[string]time.Time), nil
	} else if err != nil {
		return nil, err
	}

	metadata := make(map[string]time.Time)
	err = json.Unmarshal(content, &metadata)
	return metadata, err
}

func storeCacheMetadata(rootDir string, metadata map[string]time.Time) error {
	bytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	cachedObjectsMetadataFile := filepath.Join(rootDir, relCachedObjectsMetadataFile)
	if err := os.MkdirAll(filepath.Dir(cachedObjectsMetadataFile), 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(cachedObjectsMetadataFile, bytes, 0644)
}

func getRootDir() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userCacheDir, relCacheRootDir), nil
}
