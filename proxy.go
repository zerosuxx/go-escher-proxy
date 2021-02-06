package main

import (
	"github.com/elazarl/goproxy"
	"github.com/zerosuxx/go-escher-proxy/config"
	"github.com/zerosuxx/go-escher-proxy/handler"
	"log"
	"net/http"
)

const VERSION = "0.2.0"
const ConfigFile = ".proxy-config.json"

func main() {
	appConfig := config.AppConfig{}
	appConfig.LoadFromArgument()
	appConfig.LoadFromJsonFile(ConfigFile)

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *appConfig.Verbose

	proxy.NonproxyHandler = http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		webRequestHandler := handler.WebRequest{
			AppConfig: appConfig,
		}

		webRequestHandler.Handle(request, responseWriter)
	})

	proxy.OnRequest().DoFunc(
		func(request *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			proxy := handler.ProxyRequest{
				AppConfig: appConfig,
			}

			return proxy.Handle(request, ctx)
		})

	log.Println("Escher Pr0xy " + VERSION + " | Listening on: " + *appConfig.Host)
	log.Fatal(http.ListenAndServe(*appConfig.Host, proxy))
}
