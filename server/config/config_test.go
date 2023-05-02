package config

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var testConfigPath string
var testConfig Config

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	testConfigPath = filepath.Join(dir, filepath.Base("test.json"))
	testConfig = Config{
		SkApi: struct {
			AppKey    string `json:"app_key"`
			SecretKey string `json:"secret_key"`
		}{
			AppKey:    "testAppKey",
			SecretKey: "testSecretKey",
		},
	}

	testConfigJson, err := json.Marshal(testConfig)
	if err != nil {
		log.Println(err)
		return
	}

	f, err := os.Create(testConfigPath)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write(testConfigJson)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("hello")
}

func TestLoad(t *testing.T) {
	c, err := Load(testConfigPath)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, c, testConfig)
}

func TestMaskingConfig(t *testing.T) {
	str := testConfig.String()
	assert.False(t, strings.Contains(str, testConfig.SkApi.SecretKey), "strings.Contains(str, testConfig.SkApi.SecretKey) returns true..")
	assert.False(t, strings.Contains(str, testConfig.SkApi.AppKey), "strings.Contains(str, testConfig.SkApi.AppKey) returns true..")
}
