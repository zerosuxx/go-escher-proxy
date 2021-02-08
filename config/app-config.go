package config

import (
	"encoding/json"
	"flag"
	"github.com/zerosuxx/go-escher-proxy/escherhelper"
	"io/ioutil"
	"log"
	"os"
)

// AppConfig base application config
type AppConfig struct {
	KeyDB         []escherhelper.CredentialConfig
	ListenAddress *string
	Verbose       *bool
	ForcedHTTPS   *bool
}

// FindCredentialConfigByHost Find Escher Credential by host
func (appConfig *AppConfig) FindCredentialConfigByHost(host string) *escherhelper.CredentialConfig {
	for _, credentialConfig := range appConfig.KeyDB {
		if host == credentialConfig.Host {
			return &credentialConfig
		}
	}

	return nil
}

// LoadFromJSONFile Load Config from JSON file
func (appConfig *AppConfig) LoadFromJSONFile(jsonFile string) {
	if _, err := os.Stat(jsonFile); err == nil {
		jsonData := readFromFile(jsonFile)
		jsonError := json.Unmarshal(jsonData, appConfig)
		if jsonError != nil {
			log.Println(jsonError)
		}
	}
}

// LoadFromArgument Load Config from argument (ListenAddress, ForcedHTTPS, Verbose)
func (appConfig *AppConfig) LoadFromArgument() {
	appConfig.ListenAddress = flag.String("addr", "0.0.0.0:8181", "Proxy server listen address")
	appConfig.Verbose = flag.Bool("v", false, "Verbose")
	appConfig.ForcedHTTPS = flag.Bool("https", true, "Force Https")

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
