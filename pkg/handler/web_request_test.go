package handler

import (
	"github.com/stretchr/testify/assert"
	"github.com/zerosuxx/go-escher-proxy/pkg/config"
	"github.com/zerosuxx/go-escher-proxy/pkg/escherhelper"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestWebRequest_HandleWithoutTargetUrlReturns500(t *testing.T) {
	request := httptest.NewRequest("GET", "http://localhost", nil)
	responseRecorder := httptest.NewRecorder()

	webRequest := WebRequest{}
	webRequest.Handle(request, responseRecorder)

	assert.Equal(t, 500, responseRecorder.Code)
}

func TestWebRequest_HandleWithTargetUrlRequestUrlModifiedReturns200(t *testing.T) {
	request := httptest.NewRequest("GET", "http://localhost:1234", nil)
	request.Header.Add("X-Target-Url", "http://target.url")

	responseRecorder := httptest.NewRecorder()
	clientMock := &ClientMock{
		response: http.Response{StatusCode: 200},
	}

	appConfig := config.NewAppConfig(nil, "localhost:1234", false, false)

	webRequest := WebRequest{
		AppConfig: appConfig,
		Client:    clientMock,
	}
	webRequest.Handle(request, responseRecorder)

	assert.Equal(t, 200, responseRecorder.Code)
	assert.Equal(t, "target.url", clientMock.request.Host)
	assert.Equal(t, "http://target.url", clientMock.request.URL.String())
}

func TestWebRequest_HandlePostWithEscherAuthenticationSetProperHeaders(t *testing.T) {
	body := strings.NewReader("sample body")
	request := httptest.NewRequest("POST", "http://localhost", body)
	request.Header.Add("X-Target-Url", "http://escher.url")

	responseRecorder := httptest.NewRecorder()
	clientMock := &ClientMock{
		response: http.Response{StatusCode: 200},
	}

	var credentials []escherhelper.CredentialConfig
	credentials = append(credentials, escherhelper.CredentialConfig{
		Host:            "escher.url",
		AccessKeyID:     "key",
		APISecret:       "secret",
		CredentialScope: "eu/test/scope",
		Date:            "2011-03-11T15:59:01.888888Z",
	})

	appConfig := config.NewAppConfig(credentials, "localhost:1234", false, false)
	appConfig.KeyDB = &credentials

	webRequest := WebRequest{
		AppConfig: appConfig,
		Client:    clientMock,
	}

	webRequest.Handle(request, responseRecorder)

	requestHeader := clientMock.request.Header

	assert.Equal(t, "20110311T155901Z", requestHeader.Get("x-ems-date"))
	assert.Equal(
		t,
		"EMS-HMAC-SHA256 Credential=key/20110311/eu/test/scope, SignedHeaders=host;x-ems-date, Signature=7f1ec2f2dcd2b43409dcb0c7e2374945af2fa3d9f74bf9a8514f6987cf0f85fd",
		requestHeader.Get("x-ems-auth"),
	)
}

func TestWebRequest_HandlePostWithEscherAuthenticationAndDisabledBodyCheckSetProperHeaders(t *testing.T) {
	body := strings.NewReader("sample body")
	request := httptest.NewRequest("POST", "http://localhost:1234", body)
	request.Header.Add("X-Target-Url", "http://escher.url")

	w := httptest.NewRecorder()
	clientMock := &ClientMock{
		response: http.Response{StatusCode: 200},
	}

	var credentials []escherhelper.CredentialConfig
	credentials = append(credentials, escherhelper.CredentialConfig{
		Host:             "escher.url",
		AccessKeyID:      "key",
		APISecret:        "secret",
		CredentialScope:  "eu/test/scope",
		Date:             "2011-03-11T15:59:01.888888Z",
		DisableBodyCheck: true,
	})

	appConfig := config.NewAppConfig(credentials, "localhost:1234", false, false)
	appConfig.KeyDB = &credentials

	webRequest := WebRequest{
		AppConfig: appConfig,
		Client:    clientMock,
	}

	webRequest.Handle(request, w)

	requestHeader := clientMock.request.Header

	assert.Equal(t, "20110311T155901Z", requestHeader.Get("x-ems-date"))
	assert.Equal(
		t,
		"EMS-HMAC-SHA256 Credential=key/20110311/eu/test/scope, SignedHeaders=host;x-ems-date, Signature=a515018564868eee1a9e6ef2b56ca638b154a227ac6a27b4ab70191f8b8a0bea",
		requestHeader.Get("x-ems-auth"),
	)
}
