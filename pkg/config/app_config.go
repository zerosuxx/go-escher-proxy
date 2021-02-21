package config

import (
	"encoding/json"
	"flag"
	"github.com/zerosuxx/go-escher-proxy/pkg/escherhelper"
	"io/ioutil"
	"log"
	"os"
)

type AppConfig struct {
	Sites         map[string]SiteConfig
	ListenAddress string
	Verbose       bool
}

type SiteConfig struct {
	EscherCredentials *escherhelper.CredentialsConfig
}

func NewAppConfig(
	sites map[string]SiteConfig,
	listenAddress string,
	verbose bool,
) AppConfig {
	appConfig := AppConfig{}
	appConfig.Sites = sites
	appConfig.ListenAddress = listenAddress
	appConfig.Verbose = verbose

	return appConfig
}

func (appConfig *AppConfig) FindCredentialConfigByHost(host string) *escherhelper.CredentialsConfig {
	if appConfig.Sites == nil {
		return nil
	}

	if val, exists := appConfig.Sites[host]; exists {
		return val.EscherCredentials
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
	flag.StringVar(&appConfig.ListenAddress, "addr", "0.0.0.0:8181", "Proxy server listen address")
	flag.BoolVar(&appConfig.Verbose, "v", false, "Verbose")

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
