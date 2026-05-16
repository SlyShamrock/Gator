package config

import (
	"encoding/json"
	"os"
	"fmt"
	"path/filepath"		
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL string	`json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get filepath: %s", err)
	}
	
	file, err := os.Open(fullPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open file: %s", err)
	}

	defer file.Close()

	var cfg Config
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read file: %s", err)
	}

	return cfg, nil
}

func (u *Config) SetUser(username string) error {
	u.CurrentUserName = username
	err := write(*u)
	return err
}

func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("failed to get filepath: %s", err)
	}
	
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(cfg)
	if err != nil {
		return fmt.Errorf("failed to encode file: %s", err)
	}
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to obtain filepath: %s", err)
	}
	return filepath.Join(home, configFileName), nil
}