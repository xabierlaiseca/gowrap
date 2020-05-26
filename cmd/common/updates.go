package common

import (
	"bufio"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xabierlaiseca/gowrap/pkg/cache"
	"github.com/xabierlaiseca/gowrap/pkg/config"

	"github.com/google/go-github/github"
	"github.com/mholt/archiver/v3"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/util/file"
)

const (
	oneDay = 24 * time.Hour

	selfUpgradesFile        = "gowrap.upgrade.attempt"
	selfUpgradesFileContent = "yes"
)

func SelfUpgrade(gowrapHome, currentVersion string) {
	err := trySelfUpgrade(gowrapHome, currentVersion)
	if customerrors.IsNotFound(err) {
		err = customerrors.Errorf("upgrade for version %s not found", currentVersion)
	}

	if err != nil {
		logrus.Warningf("Failed to upgrade gowrap: %s", err.Error())
	}
}

func trySelfUpgrade(gowrapHome, currentVersion string) error {
	content, err := cache.Get(selfUpgradesFile)
	if err != nil {
		return err
	} else if content != nil {
		return nil
	}

	c, err := config.Load(gowrapHome)
	if err != nil || c.SelfUpgrade == config.SelfUpgradesDisabled {
		return nil
	}

	if err := cache.Set(selfUpgradesFile, []byte(selfUpgradesFileContent), oneDay); err != nil {
		return err
	}

	client := github.NewClient(nil)
	ctx := context.Background()

	release, _, err := client.Repositories.GetLatestRelease(ctx, "xabierlaiseca", "gowrap")
	if err != nil {
		return err
	}

	releaseSemver := strings.TrimPrefix(*release.Name, "v")
	currentSemver := strings.Split(currentVersion, "-")[0]

	if !semver.IsLessThan(currentSemver, releaseSemver) {
		return nil
	}

	if err := upgrade(release); err != nil {
		return err
	}

	return nil
}

func upgrade(release *github.RepositoryRelease) error {
	gowrapAsset, checksum, err := findAsset(release)
	if err != nil {
		return err
	}

	downloadsDir, err := ioutil.TempDir(os.TempDir(), "gowrap-download-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(downloadsDir)

	gowrapArchivePath := filepath.Join(downloadsDir, gowrapAsset.GetName())
	if err := file.DownloadTo("gowrap", gowrapArchivePath, gowrapAsset.GetBrowserDownloadURL(), checksum, "sha256"); err != nil {
		return err
	}

	archiveContentDir := filepath.Join(downloadsDir, "unarchived")
	if err := os.Mkdir(archiveContentDir, 0700); err != nil {
		return err
	}

	if err := archiver.Unarchive(gowrapArchivePath, archiveContentDir); err != nil {
		return err
	}

	files, err := ioutil.ReadDir(archiveContentDir)
	if err != nil {
		return err
	}

	binariesDir := filepath.Dir(os.Args[0])
	backupsDir := filepath.Join(downloadsDir, "backups")
	if err := os.MkdirAll(backupsDir, 0700); err != nil {
		return err
	}

	backedUp, err := moveAll(files, binariesDir, backupsDir, true)
	if err != nil {
		_, _ = moveAll(backedUp, backupsDir, binariesDir, false)
		return err
	}

	if _, err := moveAll(files, archiveContentDir, binariesDir, true); err != nil {
		_, _ = moveAll(backedUp, backupsDir, binariesDir, false)
		return err
	}

	return nil
}

func moveAll(files []os.FileInfo, srcDir, dstDir string, failFast bool) ([]os.FileInfo, error) {
	moved := make([]os.FileInfo, 0, len(files))
	var lastErr error

	for _, f := range files {
		src := filepath.Join(srcDir, f.Name())
		if _, err := os.Stat(src); os.IsNotExist(err) {
			continue
		} else if err != nil && failFast {
			return moved, err
		}

		dst := filepath.Join(dstDir, f.Name())

		err := os.Rename(src, dst)
		switch {
		case err != nil && failFast:
			return moved, err
		case err != nil:
			lastErr = err
		default:
			moved = append(moved, f)
		}
	}

	return moved, lastErr
}

func findAsset(release *github.RepositoryRelease) (*github.ReleaseAsset, string, error) {
	var checksumsAsset *github.ReleaseAsset
	var gowrapAsset *github.ReleaseAsset

	for i := range release.Assets {
		asset := release.Assets[i]

		if asset.GetName() == "checksums.txt" {
			checksumsAsset = &asset
			continue
		}

		segments := strings.Split(asset.GetName(), "_")

		if len(segments) >= 4 && segments[2] == runtime.GOOS && strings.HasPrefix(segments[3], runtime.GOARCH+".") {
			gowrapAsset = &asset
		}
	}

	if checksumsAsset == nil || gowrapAsset == nil {
		return nil, "", customerrors.NotFound()
	}

	checksum, err := fetchChecksumFor(gowrapAsset.GetName(), checksumsAsset)
	return gowrapAsset, checksum, err
}

var checksumLineRegex = regexp.MustCompile(`^([0-9a-f]+)\s+([^\s]+)\s*$`)

func fetchChecksumFor(assetName string, checksumsAsset *github.ReleaseAsset) (string, error) {
	response, err := http.Get(checksumsAsset.GetBrowserDownloadURL())
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	reader := bufio.NewReader(response.Body)

	eof := false
	for !eof {
		line, err := reader.ReadString('\n')
		eof = err == io.EOF
		if err != nil && !eof {
			return "", err
		}

		matches := checksumLineRegex.FindStringSubmatch(line)
		if len(matches) == 3 && strings.TrimSpace(matches[2]) == assetName {
			return strings.TrimSpace(matches[1]), nil
		}
	}

	return "", customerrors.NotFound()
}
