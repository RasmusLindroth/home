package home

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`

	IKEA struct {
		Gateway  string `yaml:"gateway"`
		ClientID string `yaml:"clientID"`
		PSK      string `yaml:"psk"`
	} `yaml:"ikea"`
}

func (c Config) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

func (c Config) GetAddressWithoutHost() string {
	return fmt.Sprintf(":%s", c.Server.Port)
}

func ParseConfig() (Config, error) {
	confFile, ok := os.LookupEnv("XDG_CONFIG_HOME")
	conf := Config{}

	if !ok {
		homeDir, err := os.UserHomeDir()

		if err != nil {
			return conf, err
		}

		confFile = homeDir + "/.config"
	}

	confFile += "/gohome.yaml"

	_, err := os.Stat(confFile)

	if os.IsNotExist(err) {
		return conf, fmt.Errorf("Config file %s doesn't exist", confFile)
	}

	f, err := os.Open(confFile)
	defer f.Close()

	if err != nil {
		return conf, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return conf, err
	}

	if conf.Server.Host == "" || conf.Server.Port == "" {
		return conf, errors.New("Host and/or port can't be empty")
	}

	return conf, nil
}
