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

type Configuration struct {
	gowrapHome     string
	DefaultVersion string `json:"defaultVersion,omitempty"`
}

func Load(gowrapHome string) (*Configuration, error) {
	configFilePath := getConfigFilePath(gowrapHome)

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return &Configuration{}, nil
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

	config := &Configuration{}
	if err := json.Unmarshal(content, config); err != nil {
		return nil, err
	}

	config.gowrapHome = gowrapHome
	return config, nil
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
