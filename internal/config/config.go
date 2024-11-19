package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBurl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigPath() (string, error) {
	home_path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home_path + "/" + configFileName, nil
}

func Read() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	conf := Config{}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}

func SetUser(username string, conf Config) error {
	conf.CurrentUserName = username
	write(conf)
	return nil
}

func write(conf Config) error {
	string_data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, string_data, 0644)
	if err != nil {
		return err
	}

	return nil
}
