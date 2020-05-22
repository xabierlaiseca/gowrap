package config

import (
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
)

func SetDefaultVersion(gowrapHome, version string) error {
	_, err := versions.FindLatestInstalledForPrefix(gowrapHome, version)
	if customerrors.IsNotFound(err) {
		if _, err = versions.InstallLatestIfNotInstalled(gowrapHome, version); customerrors.IsNotFound(err) {
			return customerrors.Errorf("%s is not a valid go version", version)
		} else if err != nil {
			return err
		}
	}

	configuration, err := Load(gowrapHome)
	if err != nil {
		return err
	}

	configuration.DefaultVersion = version
	return configuration.Save()
}
