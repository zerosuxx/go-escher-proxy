package main

import (
	"bytes"
	"flag"
	"github.com/elazarl/goproxy"
	"github.com/emartech/escher-go"
	"io"
	"log"
	"net/http"
	"os"
)

const VERSION = "0.0.3"

func main() {
	addr := flag.String("addr", "0.0.0.0:8181", "Proxy listen address")
	accessKeyId := flag.String("key", "", "Key name (required)")
	apiSecret := flag.String("secret", "", "Secret key (required)")
	credentialScope := flag.String("scope", "eu/suite/ems_request", "Credential scope")
	isVerbose := flag.Bool("v", false, "Verbose")
	flag.Parse()

	if *accessKeyId == "" || *apiSecret == "" {
		flag.Usage()
		os.Exit(1)
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *isVerbose

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			escherSigner := escher.Escher(getEscherConfig(accessKeyId, apiSecret, credentialScope))
			escherRequest := getEscherRequest(r)

			signedEscherRequest := escherSigner.SignRequest(escherRequest, []string{"host"})
			assignHeaders(r, signedEscherRequest.Headers)

			if *isVerbose {
				log.Println(signedEscherRequest)
			}

			return r, nil
		})

	log.Println("Escher Pr0xy " + VERSION + " | Listening on: " + *addr)
	log.Fatal(http.ListenAndServe(*addr, proxy))
}

func getEscherConfig(accessKeyId *string, apiSecret *string, credentialScope *string) escher.EscherConfig {
	return escher.EscherConfig{
		VendorKey:       "Escher",
		AlgoPrefix:      "EMS",
		HashAlgo:        "SHA256",
		AuthHeaderName:  "X-Ems-Auth",
		DateHeaderName:  "X-Ems-Date",
		AccessKeyId:     *accessKeyId,
		ApiSecret:       *apiSecret,
		CredentialScope: *credentialScope,
	}
}

func extractHeaders(header http.Header) [][2]string {
	var headers [][2]string
	var headersToSign []string
	for name, values := range header {
		headersToSign = append(headersToSign, name)
		for _, value := range values {
			headers = append(headers, [2]string{name, value})
		}
	}

	return headers
}

func getBodyAsString(body io.ReadCloser) string {
	bodyBuffer := new(bytes.Buffer)
	_, err := bodyBuffer.ReadFrom(body)
	if err != nil {
		log.Fatalln(err)
	}

	return bodyBuffer.String()
}

func getEscherRequest(r *http.Request) escher.EscherRequest {
	return escher.EscherRequest{
		Method:  r.Method,
		Url:     r.URL.String(),
		Headers: extractHeaders(r.Header),
		Body:    getBodyAsString(r.Body),
	}
}

func assignHeaders(r *http.Request, headers escher.EscherRequestHeaders) {
	for _, header := range headers {
		r.Header.Set(header[0], header[1])
	}
}
