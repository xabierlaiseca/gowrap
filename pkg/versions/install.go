package versions

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/mholt/archiver/v3"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/util/file"
	"github.com/xabierlaiseca/gowrap/pkg/versionsfile"
)

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

	destinationDir := filepath.Join(versionsDir, version)
	if err := unarchiveRemoteFile(archive, destinationDir); err != nil {
		return false, err
	}

	fmt.Printf("Successfully installed version %s\n", version)
	return true, nil
}

func unarchiveRemoteFile(archive versionsfile.GoArchive, destinationDir string) error {
	filename := path.Base(archive.URL)
	downloadsDir, downloadsDirSet := os.LookupEnv("GOWRAP_DOWNLOADS_DIR")
	if !downloadsDirSet {
		var err error
		downloadsDir, err = ioutil.TempDir(os.TempDir(), "go-download-")
		if err != nil {
			return err
		}

		defer os.RemoveAll(downloadsDir)
	}

	archiveDst := filepath.Join(downloadsDir, filename)
	if !downloadsDirSet || !exists(archiveDst) {
		if err := file.DownloadTo("go", archiveDst, archive.URL, archive.Checksum, archive.ChecksumAlgorithm); err != nil {
			return err
		}
	}

	if err := archiver.Unarchive(archiveDst, downloadsDir); err != nil {
		return err
	}

	return os.Rename(filepath.Join(downloadsDir, "go"), destinationDir)
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
