package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var Conf Config

type Config struct {
	API struct {
		Url string `json:"url"`
	} `json:"api"`
}

func Load(path string) (Config, error) {
	conf := Config{}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return conf, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, err
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		return conf, err
	}
	log.Println("config loaded .. ", conf)
	return conf, nil
}
