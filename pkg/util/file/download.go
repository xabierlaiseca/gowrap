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
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
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

	progressBar, err := newProgressBar(response.ContentLength)
	if err != nil {
		return err
	}

	_, err = io.Copy(io.MultiWriter(dst, hasher, progressBar), response.Body)
	if err != nil {
		return err
	}

	checksum := hex.EncodeToString(hasher.Sum(nil))
	if checksum != expectedChecksum {
		return customerrors.Error("failed to download file, checksums don't match")
	}

	return nil
}

const progressBarRefreshThrottle = 65 * time.Millisecond

func newProgressBar(max int64) (*progressbar.ProgressBar, error) {
	stdout := ansi.NewAnsiStdout()
	bar := progressbar.NewOptions64(
		max,
		progressbar.OptionSetWriter(stdout),
		progressbar.OptionShowBytes(true),
		progressbar.OptionThrottle(progressBarRefreshThrottle),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(stdout, "\n")
		}),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	return bar, bar.RenderBlank()
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
