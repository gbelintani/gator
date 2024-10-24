package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	configFileName = ".gatorconfig.json"
)

func fullPath() (string, error) {
	h, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get homedir: %w", err)
	}
	return fmt.Sprintf("%s/%s", h, configFileName), nil
}

func Read() (Config, error) {
	fullPath, err := fullPath()
	if err != nil {
		return Config{}, fmt.Errorf("could not get homedir: %w", err)
	}

	file, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, fmt.Errorf("could read file in path %s: %w", fullPath, err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return Config{}, fmt.Errorf("could not unmarshal: %w", err)
	}

	return config, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user

	fullPath, err := fullPath()
	if err != nil {
		return fmt.Errorf("could not get homedir: %w", err)
	}

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("could not marshal: %w", err)
	}
	err = os.WriteFile(fullPath, data, 0666)
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}
