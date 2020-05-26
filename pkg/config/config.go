package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	configFileName = "config.json"
)

const (
	AutoInstallEnabled  string = "enabled"
	AutoInstallMissing  string = "missing"
	AutoInstallDisabled string = "disabled"

	SelfUpgradesEnabled  = "enabled"
	SelfUpgradesDisabled = "disabled"
)

type Configuration struct {
	gowrapHome     string
	DefaultVersion string `json:"defaultVersion,omitempty"`
	AutoInstall    string `json:"autoInstall,omitempty"`
	SelfUpgrade    string `json:"selfUpgrade,omitempty"`
}

func Load(gowrapHome string) (*Configuration, error) {
	cfg := &Configuration{
		gowrapHome:  gowrapHome,
		AutoInstall: AutoInstallMissing,
		SelfUpgrade: SelfUpgradesDisabled,
	}
	configFilePath := getConfigFilePath(gowrapHome)
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return cfg, nil
	} else if err != nil {
		return nil, err
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	content, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(content, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Configuration) Save() error {
	content, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	configFilePath := getConfigFilePath(c.gowrapHome)

	configFile, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer configFile.Close()

	_, err = configFile.Write(content)
	return err
}

func getConfigFilePath(gowrapHome string) string {
	return filepath.Join(gowrapHome, configFileName)
}
