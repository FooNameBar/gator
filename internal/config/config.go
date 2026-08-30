package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
)

var configName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Error: json.Marshal: %v, %v\n", c, err)
	}

	ConfigLocation, err := getConfigFilePath()
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigLocation, data, 0644)
}

func Read() (config Config, err error) {

	ConfigLocation, err := getConfigFilePath()
	if err != nil {
		return config, err
	}

	f, err := os.Open(ConfigLocation)
	if err != nil {
		return config, fmt.Errorf("os.Open: %s: %v\n", ConfigLocation, err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("json.Unmarshal: %v\n", err)
	}

	return config, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error: os.UserHomeDir: %v\n", err)
	}

	return path.Join(homeDir, configName), nil
}
