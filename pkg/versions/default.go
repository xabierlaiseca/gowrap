package versions

import (
	"github.com/xabierlaiseca/gowrap/pkg/config"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
)

func SetDefaultVersion(gowrapHome, version string) error {
	_, err := FindLatestInstalledForPrefix(gowrapHome, version)
	if customerrors.IsNotFound(err) {
		if _, err = InstallLatestIfNotInstalled(gowrapHome, version); customerrors.IsNotFound(err) {
			return customerrors.Errorf("%s is not a valid go version", version)
		} else if err != nil {
			return err
		}
	}

	configuration, err := config.Load(gowrapHome)
	if err != nil {
		return err
	}

	configuration.DefaultVersion = version
	return configuration.Save()
}
