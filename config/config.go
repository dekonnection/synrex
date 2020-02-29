package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config is our config
type Config struct {
	DbHost          string            `yaml:"db_host"`
	DbName          string            `yaml:"db_name"`
	DbUser          string            `yaml:"db_user"`
	DbPassword      string            `yaml:"db_password"`
	Rooms           map[string]string `yaml:"rooms"`
	OutputDirectory string            `yaml:"output_directory"`
	LastTsFile      string            `yaml:"last_ts_file"`
	DaemonInterval  int               `yaml:"daemon_interval"`
	LogLevel        int               `yaml:"log_level"`
}

// Load config file and return Config struct
func Load(cfgFilePath string) (cfg Config, err error) {
	cfgFile, err := ioutil.ReadFile(cfgFilePath)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return
	}

	err = yaml.Unmarshal(cfgFile, &cfg)
	if err != nil {
		fmt.Printf("Error while parsing YAML config file: %s\n", err)
	}

	return
}
