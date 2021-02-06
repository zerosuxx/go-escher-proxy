package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/zerosuxx/go-escher-proxy/config"
	"github.com/zerosuxx/go-escher-proxy/escherhelper"
	"github.com/zerosuxx/go-escher-proxy/handler"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ClientMock struct {
	request  http.Request
	response http.Response
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	c.request = *req

	return &c.response, nil
}

func TestWebRequestWithoutXTargetUrlHeaderReturns500(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()

	appConfig := config.AppConfig{
		KeyDB:         nil,
		ListenAddress: "localhost:1234",
		Verbose:       false,
		ForcedHTTPS:   false,
	}

	webRequest := handler.WebRequest{
		AppConfig: appConfig,
	}

	webRequest.Handle(req, w)

	assert.Equal(t, 500, w.Code)
}

func TestWebRequestWithXTargetUrlHeaderReturns200(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.Header.Add("X-Target-Url", "http://escher.url")

	w := httptest.NewRecorder()
	clientMock := &ClientMock{
		response: http.Response{StatusCode: 200},
	}

	appConfig := config.AppConfig{}
	appConfig.Verbose = false

	var credentials []escherhelper.CredentialConfig
	credentials = append(credentials, escherhelper.CredentialConfig{
		Host:            "escher.url",
		AccessKeyID:     "key",
		APISecret:       "secret",
		CredentialScope: "eu/test/scope",
		Date:            "2011-03-11T15:59:01.888888Z",
	})
	appConfig.KeyDB = credentials

	webRequest := handler.WebRequest{
		AppConfig: appConfig,
		Client:    clientMock,
	}

	webRequest.Handle(req, w)

	requestHeader := clientMock.request.Header

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "escher.url", clientMock.request.Host)
	assert.Equal(t, "http://escher.url", clientMock.request.URL.String())
	assert.Equal(t, "20110311T155901Z", requestHeader.Get("x-ems-date"))
	assert.Equal(
		t,
		"EMS-HMAC-SHA256 Credential=key/20110311/eu/test/scope, SignedHeaders=host;x-ems-date, Signature=f00d5a438853852dc74fe963735a8de4a7e4a819e430e3471bb4271badb4f4cf",
		requestHeader.Get("x-ems-auth"),
	)
}
