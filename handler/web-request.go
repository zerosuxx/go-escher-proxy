package handler

import (
	"github.com/emartech/escher-go"
	"github.com/zerosuxx/go-escher-proxy/config"
	"github.com/zerosuxx/go-escher-proxy/escherhelper"
	"github.com/zerosuxx/go-escher-proxy/httphelper"
	"log"
	"net/http"
)

type WebRequest struct {
	AppConfig config.AppConfig
	Client    httphelper.Client
}

func (web *WebRequest) Handle(request *http.Request, responseWriter http.ResponseWriter) {
	targetURL := request.Header.Get("X-Target-Url")

	if targetURL == "" {
		responseWriter.WriteHeader(500)

		return
	}

	url, parseErr := request.URL.Parse(targetURL)
	if parseErr != nil {
		panic(parseErr)
	}
	newRequest, requestErr := http.NewRequest(request.Method, url.String(), nil)
	if requestErr != nil {
		panic(requestErr)
	}

	newRequest.Header = request.Header
	newRequest.Header.Del("X-Target-Url")
	newRequest.Header.Set("Host", url.Host)
	newRequest.Body = request.Body

	escherConfig := web.AppConfig.FindCredentialConfigByHost(newRequest.Host)
	if escherConfig != nil {
		escherSigner := escher.Escher(escherConfig.GetEscherConfig())
		signedEscherRequest := escherSigner.SignRequest(escherhelper.RequestFactory(newRequest), []string{"host"})

		httphelper.AssignHeaders(newRequest.Header, signedEscherRequest.Headers)
	} else {
		log.Println("Escher config not found for given host: " + newRequest.Host)
	}

	if web.AppConfig.Verbose {
		log.Println("Host", newRequest.Host)
		log.Println("Headers", newRequest.Header)
	}

	clientResponse, clientErr := web.Client.Do(newRequest)
	if clientErr != nil {
		panic(clientErr)
	}

	responseWriter.WriteHeader(clientResponse.StatusCode)
	for _, value := range httphelper.ExtractHeaders(newRequest.Header) {
		responseWriter.Header().Add(value[0], value[1])
	}

	_, responseError := responseWriter.Write(
		httphelper.ReadBodyWithoutClear(
			httphelper.ResponseBody{
				Response: clientResponse,
			},
		),
	)
	if responseError != nil {
		panic(responseError)
	}
}
