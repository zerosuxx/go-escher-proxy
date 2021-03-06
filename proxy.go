package main

import (
	"github.com/elazarl/goproxy"
	"github.com/zerosuxx/go-escher-proxy/pkg/config"
	"github.com/zerosuxx/go-escher-proxy/pkg/handler"
	"log"
	"net/http"
	"os"
)

func main() {
	const VERSION = "0.6.4"
	const configFileName = "proxy-config.json"

	appConfig := config.AppConfig{}
	appConfig.LoadFromArgument()
	path, _ := os.Getwd()
	configFile := path + "/" + configFileName
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

	log.Println("Escher Pr0xy " + VERSION + " | Listening on: " + appConfig.ListenAddress)
	log.Fatal(http.ListenAndServe(appConfig.ListenAddress, proxy))
}
