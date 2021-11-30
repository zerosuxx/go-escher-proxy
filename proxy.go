package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"github.com/zerosuxx/go-escher-proxy/pkg/config"
	"github.com/zerosuxx/go-escher-proxy/pkg/handler"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var Version = "development"

func isDebugMode(version string) bool {
	return version == "development"
}

func getConfigPath(version string) string {
	if isDebugMode(version) {
		currentWorkingPath, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		return currentWorkingPath
	} else {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}

		return filepath.Dir(ex)
	}
}

func main() {
	const configFileName = "proxy-config.json"

	appConfig := config.AppConfig{}
	appConfig.LoadFromArgument()

	configFile := getConfigPath(Version) + "/" + configFileName
	if appConfig.Verbose {
		log.Println("Try to loading config file: " + configFile)
	}
	appConfig.LoadFromJSONFile(configFile)

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = appConfig.Verbose

	proxy.NonproxyHandler = http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		webRequestHandler := handler.WebRequest{
			AppConfig: appConfig,
			Client:    &http.Client{},
		}

		webRequestHandler.Handle(request, responseWriter)
	})

	proxy.OnRequest().DoFunc(func(request *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		proxy := handler.ProxyRequest{
			AppConfig: appConfig,
		}

		return proxy.Handle(request, ctx)
	})

	fmt.Println("Escher Pr0xy " + Version + " | Listening on: " + appConfig.ListenAddress)
	log.Fatal(http.ListenAndServe(appConfig.ListenAddress, proxy))
}
