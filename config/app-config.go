package config

import (
	"encoding/json"
	"flag"
	"github.com/zerosuxx/go-escher-proxy/escherhelper"
	"io/ioutil"
	"log"
	"os"
)

type AppConfig struct {
	KeyDB       []escherhelper.CredentialConfig
	Host        *string
	Verbose     *bool
	ForcedHttps *bool
}

func (appConfig *AppConfig) FindCredentialConfigByHost(host string) *escherhelper.CredentialConfig {
	for _, credentialConfig := range appConfig.KeyDB {
		if host == credentialConfig.Host {
			return &credentialConfig
		}
	}

	return nil
}

func (appConfig *AppConfig) LoadFromJsonFile(jsonFile string) {
	if _, err := os.Stat(jsonFile); err == nil {
		jsonData := readFromFile(jsonFile)
		jsonError := json.Unmarshal(jsonData, appConfig)
		if jsonError != nil {
			log.Println(jsonError)
		}
	}
}

func (appConfig *AppConfig) LoadFromArgument() {
	appConfig.Host = flag.String("host", "0.0.0.0:8181", "Proxy server listen address")
	appConfig.ForcedHttps = flag.Bool("https", true, "Force Https")
	appConfig.Verbose = flag.Bool("v", false, "Verbose")

	flag.Parse()
}

func readFromFile(file string) []byte {
	jsonFile, err := os.Open(file)

	if err != nil {
		return []byte("")
	}

	jsonData, _ := ioutil.ReadAll(jsonFile)

	return jsonData
}
