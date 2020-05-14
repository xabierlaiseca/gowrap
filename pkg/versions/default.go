package versions

import (
	"github.com/xabierlaiseca/gowrap/pkg/config"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
)

func SetDefaultVersion(version string) error {
	_, err := FindLatestInstalledForPrefix(version)
	if customerrors.IsNotFound(err) {
		if err = InstallLatestForPrefix(version); customerrors.IsNotFound(err) {
			return customerrors.Errorf("%s is not a valid go version", version)
		} else if err != nil {
			return err
		}
	}

	configuration, err := config.Load()
	if err != nil {
		return err
	}

	configuration.DefaultVersion = version
	return configuration.Save()
}
