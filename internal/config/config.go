// reads and writes "~/.gatorconfig.json"

package config

import (
	"fmt"
	"os"
	"encoding/json"
	"path/filepath"
	"io/fs"
)

// json file structure with struct tags
type Config struct {
	Db_url				string	`json:"db_url"`
	Current_user_name	string	`json:"current_user_name"`
}

// reads "~/.gatorconfig.json"
func ReadConfig() Config {

	configPath, err := getConfigFilePath()
	if err != nil {
		fmt.Errorf("home path error: %v", err)
	}

	// fmt.Printf("config path: %v\n", configPath)

	body, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Errorf("%v", err)
	}

	// os.Stdout.Write(body)

	var config Config
	err = json.Unmarshal(body, &config)
	if err != nil {
		fmt.Errorf("%v", err)
	}

	return config
}

// sets user field then writes to "~/.gatorconfig.json"
func (c Config) SetUser(username string) {
	c.Current_user_name = username
	err := writeConfig(c)
	if err != nil {
		fmt.Errorf("error writing config to file: %v", err)
	}
}


// helpers
func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configFileName), nil
}

func writeConfig(c Config) error {

	// fmt.Printf("config in write: %v\n", c)

	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}
	// fmt.Printf("json: %s\n", jsonData)

	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, jsonData, fs.ModeType)
	if err != nil {
		return err
	}

	return nil
}

const configFileName = ".gatorconfig.json"
