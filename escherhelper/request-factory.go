package escherhelper

import (
	"github.com/emartech/escher-go"
	"github.com/zerosuxx/go-escher-proxy/httphelper"
	"net/http"
)

func RequestFactory(request *http.Request) escher.EscherRequest {
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
		Body:    string(httphelper.ReadBodyWithoutClear(httphelper.RequestBody{Request: request})),
	}
}
