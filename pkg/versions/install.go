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

func Install(version string) error {
	installableVersions, err := versionsfile.Load()
	if err != nil {
		return err
	}

	archive, found := installableVersions[version]
	if !found {
		return customerrors.Errorf("version %s is not available", version)
	}

	response, err := http.Get(archive.URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	downloadsDir, err := ioutil.TempDir(os.TempDir(), "gowrap-download-")
	if err != nil {
		return err
	}

	filename := path.Base(archive.URL)
	dstPath := filepath.Join(downloadsDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	var hasher hash.Hash
	switch {
	case archive.IsSHA256Checksum():
		hasher = sha256.New()
	case archive.IsSHA1Checksum():
		sha1.New()
	default:
		return customerrors.Errorf("unsupported checksum algorithm: %s", archive.ChecksumAlgorithm)
	}

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

		if err != nil {
			return err
		}

		if writeCount != readCount {
			return customerrors.Errorf("failed to write file to disk")
		}
	}

	progressBar.Done()

	checksum := hex.EncodeToString(hasher.Sum(nil))
	if checksum != archive.Checksum {
		return customerrors.New("failed to download file, checksums don't match")
	}

	if err := archiver.Unarchive(dstPath, downloadsDir); err != nil {
		return err
	}

	versionsDir, err := GetVersionsDir()
	if err != nil {
		return err
	}

	destinationDir := filepath.Join(versionsDir, version)
	if err := os.Rename(filepath.Join(downloadsDir, "go"), destinationDir); err != nil {
		return err
	}

	fmt.Printf("Successfully installed version %s\n", version)
	return nil
}

const oneMB = 1024 * 1024

func sizeToMBString(size int64) string {
	return fmt.Sprintf("%dMB", size/oneMB)
}
