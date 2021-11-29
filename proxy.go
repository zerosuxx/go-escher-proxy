package main

import (
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

func main() {
	const configFileName = "proxy-config.json"

	appConfig := config.AppConfig{}
	appConfig.LoadFromArgument()

	var configFile string
	if isDebugMode(Version) {
		currentWorkingPath, _ := os.Getwd()
		configFile = currentWorkingPath + "/" + configFileName
	} else {
		currentScriptPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		configFile = currentScriptPath + "/" + configFileName
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

	log.Println("Escher Pr0xy " + Version + " | Listening on: " + appConfig.ListenAddress)
	log.Fatal(http.ListenAndServe(appConfig.ListenAddress, proxy))
}
