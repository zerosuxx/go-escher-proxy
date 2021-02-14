package handler

import (
	"github.com/elazarl/goproxy"
	"github.com/emartech/escher-go"
	"github.com/zerosuxx/go-escher-proxy/pkg/config"
	"github.com/zerosuxx/go-escher-proxy/pkg/escherhelper"
	"github.com/zerosuxx/go-escher-proxy/pkg/httphelper"
	"log"
	"net/http"
)

type ProxyRequest struct {
	AppConfig config.AppConfig
}

func (proxy *ProxyRequest) Handle(request *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if *proxy.AppConfig.ForcedHTTPS && request.Header.Get("X-Disable-Force-Https") != "1" {
		request.URL.Scheme = "https"
	}

	request.Header.Set("Host", request.Host)

	credentialConfig := proxy.AppConfig.FindCredentialConfigByHost(request.Host)
	if credentialConfig == nil {
		log.Println("Escher config not found for given host: " + request.Host)

		return request, nil
	}

	escherRequestFactory := escherhelper.RequestFactory{}
	escherRequest := escherRequestFactory.CreateFromCredentialConfig(request, credentialConfig)
	escherSigner := escher.Escher(credentialConfig.GetEscherConfig())
	signedEscherRequest := escherSigner.SignRequest(escherRequest, []string{"host"})
	httphelper.AssignHeaders(request.Header, signedEscherRequest.Headers)

	if *proxy.AppConfig.Verbose {
		log.Println("Headers", request.Header)
	}

	return request, nil
}
