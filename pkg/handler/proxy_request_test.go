package handler_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zerosuxx/go-escher-proxy/pkg/config"
	"github.com/zerosuxx/go-escher-proxy/pkg/escherhelper"
	"github.com/zerosuxx/go-escher-proxy/pkg/handler"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProxyRequest_HandleWithNonLocalHostSchemeForcedToHttps(t *testing.T) {
	request := httptest.NewRequest("GET", "http://example.com", nil)

	appConfig := config.NewAppConfig(map[string]config.SiteConfig{}, "localhost:1234", false)

	proxyRequest := handler.ProxyRequest{
		AppConfig: appConfig,
	}

	proxyRequest.Handle(request, nil)
	assert.Equal(t, "https", request.URL.Scheme)
}

func TestProxyRequest_HandleWithLocalHostSchemeNotForcedToHttps(t *testing.T) {
	request := httptest.NewRequest("GET", "http://localhost", nil)

	appConfig := config.NewAppConfig(map[string]config.SiteConfig{}, "localhost:1234", false)

	proxyRequest := handler.ProxyRequest{
		AppConfig: appConfig,
	}

	proxyRequest.Handle(request, nil)
	assert.Equal(t, "http", request.URL.Scheme)
}

func TestProxyRequest_HandleWithDisableForceHttpsHeaderNotForcedToHttps(t *testing.T) {
	request := httptest.NewRequest("GET", "http://example.com", nil)
	request.Header.Set("X-Disable-Force-Https", "1")

	appConfig := config.NewAppConfig(map[string]config.SiteConfig{}, "localhost:1234", false)

	proxyRequest := handler.ProxyRequest{
		AppConfig: appConfig,
	}

	proxyRequest.Handle(request, nil)
	assert.Equal(t, "http", request.URL.Scheme)
}

func TestProxyRequest_HandleSetRequestHostHeader(t *testing.T) {
	request := httptest.NewRequest("GET", "http://localhost", nil)

	appConfig := config.NewAppConfig(map[string]config.SiteConfig{}, "localhost:1234", false)

	proxyRequest := handler.ProxyRequest{
		AppConfig: appConfig,
	}

	proxyRequest.Handle(request, nil)
	assert.Equal(t, "localhost", request.Header.Get("Host"))
}

func TestProxyRequest_HandleGetWithEscherAuthenticationSetProperHeaders(t *testing.T) {
	request := httptest.NewRequest("GET", "http://escher.url", nil)

	credentials := escherhelper.CredentialsConfig{
		AccessKeyID:     "key",
		APISecret:       "secret",
		CredentialScope: "eu/test/scope",
		Date:            "2011-03-11T15:59:01.888888Z",
	}
	sites := map[string]config.SiteConfig{}
	sites["escher.url"] = config.SiteConfig{
		EscherCredentials: &credentials,
	}
	appConfig := config.NewAppConfig(sites, "localhost:1234", false)

	proxyRequest := handler.ProxyRequest{
		AppConfig: appConfig,
	}

	proxyRequest.Handle(request, nil)
	assert.Equal(t, "20110311T155901Z", request.Header.Get("x-ems-date"))
	assert.Equal(
		t,
		"EMS-HMAC-SHA256 Credential=key/20110311/eu/test/scope, SignedHeaders=host;x-ems-date, Signature=f00d5a438853852dc74fe963735a8de4a7e4a819e430e3471bb4271badb4f4cf",
		request.Header.Get("x-ems-auth"),
	)
}

func TestProxyRequest_HandlePostWithEscherAuthenticationSetProperHeaders(t *testing.T) {
	body := strings.NewReader("sample body")
	request := httptest.NewRequest("POST", "http://escher.url", body)

	credentials := escherhelper.CredentialsConfig{
		AccessKeyID:     "key",
		APISecret:       "secret",
		CredentialScope: "eu/test/scope",
		Date:            "2011-03-11T15:59:01.888888Z",
	}
	sites := map[string]config.SiteConfig{}
	sites["escher.url"] = config.SiteConfig{
		EscherCredentials: &credentials,
	}
	appConfig := config.NewAppConfig(sites, "localhost:1234", false)

	proxyRequest := handler.ProxyRequest{
		AppConfig: appConfig,
	}

	proxyRequest.Handle(request, nil)
	assert.Equal(t, "20110311T155901Z", request.Header.Get("x-ems-date"))
	assert.Equal(
		t,
		"EMS-HMAC-SHA256 Credential=key/20110311/eu/test/scope, SignedHeaders=host;x-ems-date, Signature=7f1ec2f2dcd2b43409dcb0c7e2374945af2fa3d9f74bf9a8514f6987cf0f85fd",
		request.Header.Get("x-ems-auth"),
	)
}

func TestProxyRequest_HandlePostWithEscherAuthenticationAndDisabledBodyCheckSetProperHeaders(t *testing.T) {
	body := strings.NewReader("sample body")
	request := httptest.NewRequest("POST", "http://escher.url", body)

	credentials := escherhelper.CredentialsConfig{
		AccessKeyID:      "key",
		APISecret:        "secret",
		CredentialScope:  "eu/test/scope",
		Date:             "2011-03-11T15:59:01.888888Z",
		DisableBodyCheck: true,
	}
	sites := map[string]config.SiteConfig{}
	sites["escher.url"] = config.SiteConfig{
		EscherCredentials: &credentials,
	}
	appConfig := config.NewAppConfig(sites, "localhost:1234", false)

	proxyRequest := handler.ProxyRequest{
		AppConfig: appConfig,
	}

	proxyRequest.Handle(request, nil)
	assert.Equal(t, "20110311T155901Z", request.Header.Get("x-ems-date"))
	assert.Equal(
		t,
		"EMS-HMAC-SHA256 Credential=key/20110311/eu/test/scope, SignedHeaders=host;x-ems-date, Signature=a515018564868eee1a9e6ef2b56ca638b154a227ac6a27b4ab70191f8b8a0bea",
		request.Header.Get("x-ems-auth"),
	)
}
