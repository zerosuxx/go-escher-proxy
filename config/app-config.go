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
	KeyDB         *[]escherhelper.CredentialConfig
	ListenAddress *string
	Verbose       *bool
	ForcedHTTPS   *bool
}

func NewAppConfig(
	keyDB []escherhelper.CredentialConfig,
	listenAddress string,
	verbose bool,
	forceHTTPS bool,
) AppConfig {
	appConfig := AppConfig{}
	appConfig.KeyDB = &keyDB
	appConfig.ListenAddress = &listenAddress
	appConfig.Verbose = &verbose
	appConfig.ForcedHTTPS = &forceHTTPS

	return appConfig
}

func (appConfig *AppConfig) FindCredentialConfigByHost(host string) *escherhelper.CredentialConfig {
	for _, credentialConfig := range *appConfig.KeyDB {
		if host == credentialConfig.Host {
			return &credentialConfig
		}
	}

	return nil
}

func (appConfig *AppConfig) LoadFromJSONFile(jsonFile string) {
	if _, err := os.Stat(jsonFile); err == nil {
		jsonData := readFromFile(jsonFile)
		jsonError := json.Unmarshal(jsonData, appConfig)
		if jsonError != nil {
			log.Println(jsonError)
		}
	}
}

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
