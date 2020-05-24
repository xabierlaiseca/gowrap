package file

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/xabierlaiseca/gowrap/pkg/util/console"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
)

const (
	downloadBufferSize = 64 * 1024
)

func DownloadTo(packageName, dst, url, checksum, algorithm string) error {
	var hasher hash.Hash
	switch strings.ToLower(algorithm) {
	case "sha256":
		hasher = sha256.New()
	case "sha1":
		hasher = sha1.New()
	case "":
		hasher = noopHasher{}

	default:
		return customerrors.Errorf("unsupported checksum algorithm: %s", algorithm)
	}

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	fmt.Printf("Downloading %s from %s...\n", packageName, url)
	return storeDownload(response, dst, checksum, hasher)
}

func storeDownload(response *http.Response, dstPath, expectedChecksum string, hasher hash.Hash) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

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
	if checksum != expectedChecksum {
		return customerrors.Error("failed to download file, checksums don't match")
	}

	return nil
}

const oneMB = 1024 * 1024

func sizeToMBString(size int64) string {
	return fmt.Sprintf("%dMB", size/oneMB)
}

type noopHasher struct{}

func (n2 noopHasher) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (n2 noopHasher) Sum(b []byte) []byte {
	return []byte{}
}

func (n2 noopHasher) Reset() {}

func (n2 noopHasher) Size() int {
	return 0
}

func (n2 noopHasher) BlockSize() int {
	return 0
}
