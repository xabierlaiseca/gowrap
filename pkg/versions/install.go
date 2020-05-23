package versions

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/mholt/archiver/v3"
	"github.com/xabierlaiseca/gowrap/pkg/util/console"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versionsfile"
)

const downloadBufferSize = 64 * 1024

// InstallLatestIfNotInstalled installs latest version for given prefix if not already installed.
// If no error, `true` will be returned if the version was installed or `false` if the version
// was already available.
func InstallLatestIfNotInstalled(gowrapHome, prefix string) (bool, error) {
	versionToInstall, err := FindLatestAvailable(prefix)
	if err != nil {
		return false, err
	}

	return InstallIfNotInstalled(gowrapHome, versionToInstall)
}

// InstallIfNotInstalled installs the requested version if not already installed.
// If no error, `true` will be returned if the version was installed or `false` if the version
// was already available.
func InstallIfNotInstalled(gowrapHome, version string) (bool, error) {
	versionsDir, err := GetVersionsDir(gowrapHome)
	if err != nil {
		return false, err
	}

	if alreadyInstalled, err := isVersionInstalled(versionsDir, version); err != nil {
		return false, err
	} else if alreadyInstalled {
		return false, nil
	}

	installableVersions, err := versionsfile.Load()
	if err != nil {
		return false, err
	}

	archive, found := installableVersions[version]
	if !found {
		return false, customerrors.Errorf("version %s is not available", version)
	}

	filename := path.Base(archive.URL)
	downloadsDir, downloadsDirSet := os.LookupEnv("GOWRAP_DOWNLOADS_DIR")
	if !downloadsDirSet {
		downloadsDir, err = ioutil.TempDir(os.TempDir(), "gowrap-download-")
		if err != nil {
			return false, err
		}
	}

	archivePath := filepath.Join(downloadsDir, filename)
	if !downloadsDirSet || !exists(archivePath) {
		response, err := http.Get(archive.URL)
		if err != nil {
			return false, err
		}
		defer response.Body.Close()

		if err := storeDownload(response, archivePath, archive); err != nil {
			return false, err
		}
	}

	if err := archiver.Unarchive(archivePath, downloadsDir); err != nil {
		return false, err
	}

	destinationDir := filepath.Join(versionsDir, version)
	if err := os.Rename(filepath.Join(downloadsDir, "go"), destinationDir); err != nil {
		return false, err
	}

	fmt.Printf("Successfully installed version %s\n", version)
	return true, nil
}

func Uninstall(gowrapHome, version string) error {
	versionsDir, err := GetVersionsDir(gowrapHome)
	if err != nil {
		return err
	}

	versionDir := filepath.Join(versionsDir, version)

	if _, err = os.Stat(versionDir); os.IsNotExist(err) {
		return customerrors.Errorf("version %s was not previously installed", version)
	} else if err != nil {
		return err
	}

	return os.RemoveAll(versionDir)
}

func storeDownload(response *http.Response, dstPath string, archive versionsfile.GoArchive) error {
	var hasher hash.Hash
	switch {
	case archive.IsSHA256Checksum():
		hasher = sha256.New()
	case archive.IsSHA1Checksum():
		sha1.New()
	default:
		return customerrors.Errorf("unsupported checksum algorithm: %s", archive.ChecksumAlgorithm)
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	fmt.Printf("Downloading go from %s...\n", archive.URL)
	progressBar := console.NewProgressBar(response.ContentLength, sizeToMBString)

	bytes := make([]byte, downloadBufferSize)
	finished := false

	for !finished {
		readCount, err := response.Body.Read(bytes)

		if err == io.EOF {
			finished = true
		} else if err != nil {
			return err
		}

		writeCount, err := dst.Write(bytes[:readCount])
		if err != nil {
			return err
		}

		if _, err = hasher.Write(bytes[:readCount]); err != nil {
			return err
		}

		progressBar.Increment(int64(writeCount))

		if writeCount != readCount {
			return customerrors.Errorf("failed to write file to disk")
		}
	}

	progressBar.Done()

	checksum := hex.EncodeToString(hasher.Sum(nil))
	if checksum != archive.Checksum {
		return customerrors.Error("failed to download file, checksums don't match")
	}

	return nil
}

const oneMB = 1024 * 1024

func sizeToMBString(size int64) string {
	return fmt.Sprintf("%dMB", size/oneMB)
}

func isVersionInstalled(versionsDir, version string) (bool, error) {
	versionDir := filepath.Join(versionsDir, version)

	stat, err := os.Stat(versionDir)
	switch {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, err
	case !stat.IsDir():
		return false, customerrors.Errorf("unexpected file in %s, it should be removed for proper functioning of this tool", versionDir)
	}

	goBinPath := filepath.Join(versionDir, "bin", "go")
	if stat, err := os.Stat(goBinPath); err == nil && stat.Mode().IsRegular() {
		return true, nil
	}

	return false, customerrors.Errorf("unexpected content in %s, it should be removed for proper functioning of this tool", versionDir)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
