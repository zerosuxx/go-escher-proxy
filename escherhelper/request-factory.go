package escherhelper

import (
	"github.com/emartech/escher-go"
	"github.com/zerosuxx/go-escher-proxy/httphelper"
	"net/http"
)

type RequestFactory struct {
}

func (factory *RequestFactory) Create(request *http.Request) escher.EscherRequest {
	return createEscherRequest(request, string(httphelper.ReadBodyWithoutClear(httphelper.RequestBody{Request: request})))
}

func (factory *RequestFactory) CreateWithEmptyBody(request *http.Request) escher.EscherRequest {
	return createEscherRequest(request, "")
}

func (factory *RequestFactory) CreateFromCredentialConfig(request *http.Request, config *CredentialConfig) escher.EscherRequest {
	var escherRequest escher.EscherRequest
	if config.DisableBodyCheck {
		escherRequest = factory.CreateWithEmptyBody(request)
	} else {
		escherRequest = factory.Create(request)
	}

	return escherRequest
}

func createEscherRequest(request *http.Request, body string) escher.EscherRequest {
	path := request.URL.Path
	if path == "" {
		path = "/"
	}

	query := request.URL.RawQuery
	if query != "" {
		path += "?" + query
	}

	return escher.EscherRequest{
		Method:  request.Method,
		Url:     path,
		Headers: httphelper.ExtractHeaders(request.Header),
		Body:    body,
	}
}
