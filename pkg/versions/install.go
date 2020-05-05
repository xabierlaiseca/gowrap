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
	"strings"

	"github.com/mholt/archiver/v3"
	"github.com/xabierlaiseca/gowrap/pkg/versionsfile"
)

var versionsDirRelToHome = filepath.Join(".gowrap", "versions")

func Install(version string) error {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	installableVersions, err := versionsfile.Load()
	if err != nil {
		return err
	}

	archive, found := installableVersions[version]
	if !found {
		return fmt.Errorf("version %s is not available", version)
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
	if archive.IsSHA256Checksum() {
		hasher = sha256.New()
	} else if archive.IsSHA1Checksum() {
		hasher = sha1.New()
	} else {
		return fmt.Errorf("Unsupported checksum algorithm: %s", archive.ChecksumAlgorithm)
	}

	totalMB := response.ContentLength / (1024 * 1024)
	fmt.Printf("Downloading go from %s...\n", archive.URL)
	fmt.Printf("[>%s] 0MB/%dMB", strings.Repeat(" ", 49), totalMB)

	bytes := make([]byte, 64*1024)
	totalRead := int64(0)
	finished := false

	for !finished {
		readCount, err := response.Body.Read(bytes)
		totalRead += int64(readCount)
		downloadedPercentage := int(100 * totalRead / response.ContentLength)

		if err == io.EOF {
			finished = true
		} else if err != nil {
			return err
		}

		writeCount, err := dst.Write(bytes[:readCount])
		hasher.Write(bytes[:readCount])

		equalSigns := downloadedPercentage / 2
		greaterSigns := 1
		if equalSigns == 50 {
			greaterSigns = 0
		}
		spaces := 50 - equalSigns - greaterSigns
		fmt.Printf("\r\033[K[%s%s%s] %dMB/%dMB", strings.Repeat("=", equalSigns),
			strings.Repeat(">", greaterSigns), strings.Repeat(" ", spaces),
			totalRead/(1024*1024), totalMB)

		if err != nil {
			return err
		}

		if writeCount != readCount {
			return fmt.Errorf("failed to write file to disk")
		}
	}

	fmt.Println()
	checksum := hex.EncodeToString(hasher.Sum(nil))
	if checksum != archive.Checksum {
		return fmt.Errorf("failed to download file, checksums don't match")
	}

	if err := archiver.Unarchive(dstPath, downloadsDir); err != nil {
		return err
	}

	versionsDir := filepath.Join(userHome, versionsDirRelToHome)
	if err := os.MkdirAll(versionsDir, 0755); err != nil {
		return err
	}

	destinationDir := filepath.Join(versionsDir, version)
	if err := os.Rename(filepath.Join(downloadsDir, "go"), destinationDir); err != nil {
		return err
	}

	fmt.Printf("Successfully installed version %s\n", version)
	return nil
}
