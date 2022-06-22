package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

const CONFIG_FILE_PATH = "/.config/personal-website/config.json"

type Config struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Port   int    `json:"port"`

	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

func ReadConfigFile() (*Config, error) {
	h, err := os.UserHomeDir()

	if err != nil {
		log.Panic("Problem getting home directory: " + err.Error())
	}

	b, err := ioutil.ReadFile(h + CONFIG_FILE_PATH)

	if errors.Is(err, os.ErrNotExist) {
		log.Panic("Config file does not exist at ~" + CONFIG_FILE_PATH + "\n" + err.Error())
	}
	if err != nil {
		log.Panic("Problem opening ~" + CONFIG_FILE_PATH + "\n" + err.Error())
	}

	if err != nil {
		log.Panic("Problem reading config file\n" + err.Error())
	}

	var c Config
	err = json.Unmarshal(b, &c)

	if err != nil {
		log.Panic("Problem unmarshaling data from config\n" + err.Error())
	}

	return &c, nil
}
