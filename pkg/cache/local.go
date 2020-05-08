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
	objectFile      = "objects.json"
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
	cachedObjectsMetadataFile := buildCachedObjectsMetadataFile(rootDir)
	metadata, err := readCacheMetadata(cachedObjectsMetadataFile)
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
		return nil, storeCacheMetadata(cachedObjectsMetadataFile, metadata)
	} else if !ok {
		return nil, nil
	}

	content, err := ioutil.ReadFile(cachedObjectPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	return content, err
}

// Set stores the given content in provided cache path for the requested duration.
func Set(relpath string, content []byte, expiresIn time.Duration) error {
	rootDir, err := getRootDir()
	if err != nil {
		return err
	}

	return set(rootDir, relpath, content, expiresIn)
}

func set(rootDir, relpath string, content []byte, expiresIn time.Duration) error {
	cachedObjectsMetadataFile := buildCachedObjectsMetadataFile(rootDir)
	metadata, err := readCacheMetadata(cachedObjectsMetadataFile)
	if err != nil {
		return err
	}

	objectPath := filepath.Join(rootDir, relObjectsDir, relpath)
	if err := os.MkdirAll(filepath.Dir(objectPath), 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(objectPath, content, 0600); err != nil {
		return err
	}

	metadata[relpath] = time.Now().Add(expiresIn)
	return storeCacheMetadata(cachedObjectsMetadataFile, metadata)
}

func readCacheMetadata(cachedObjectsMetadataFile string) (map[string]time.Time, error) {
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

func storeCacheMetadata(cachedObjectsMetadataFile string, metadata map[string]time.Time) error {
	bytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(cachedObjectsMetadataFile), 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(cachedObjectsMetadataFile, bytes, 0600)
}

func getRootDir() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userCacheDir, relCacheRootDir), nil
}

func buildCachedObjectsMetadataFile(rootDir string) string {
	return filepath.Join(rootDir, relMetadataDir, objectFile)
}
