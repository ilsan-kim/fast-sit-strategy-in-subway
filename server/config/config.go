package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var Conf Config

type Config struct {
	SkApi struct {
		AppKey    string `json:"app_key"`
		SecretKey string `json:"secret_key"`
	} `json:"sk_api"`
}

type MaskedConfig struct {
	Config
}

func (c Config) String() string {
	var masked MaskedConfig
	jsonData, _ := json.Marshal(c)
	_ = json.Unmarshal(jsonData, &masked)
	masked.SkApi.AppKey = "***********"
	masked.SkApi.SecretKey = "***********"
	jsonData, _ = json.Marshal(masked)
	return string(jsonData)
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
