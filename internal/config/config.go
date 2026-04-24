// reads and writes "~/.gatorconfig.json"

package config

// json file structure with struct tags
type Config struct {

}

// reads "~/.gatorconfig.json"
func Read() Config {
return Config{}
}

// sets user field then writes to "~/.gatorconfig.json"
func (c Config) SetUser() {
	configPath := getConfigFilePath()
	currentConfig := Read(configPath)
	
}


// helpers
func getConfigFilePath() (string, error) {

}

func write(c Config) error {

}

const configFileName = ".gatorconfig.json"
