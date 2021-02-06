package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/elazarl/goproxy"
	"github.com/emartech/escher-go"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const VERSION = "0.1.0"

func main() {
	jsonData := readFromFile(".proxy-config.json")
	jsonConfig := getJsonConfig(jsonData)

	addr := flag.String("addr", "0.0.0.0:8181", "Proxy listen address")
	isHttpsForced := flag.Bool("https", true, "Force Https")
	isVerbose := flag.Bool("v", false, "Verbose")
	flag.Parse()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *isVerbose

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if *isHttpsForced && r.Header.Get("X-Disable-Force-Https") != "1" {
				r.URL.Scheme = "https"
			}
			r.Header.Set("Host", r.Host)
			r.Header.Del("Proxy-Connection")

			escherConfig := jsonConfig.getEscherConfigByHost(r.Host)
			if escherConfig == nil {
				return r, nil
			}

			escherSigner := escher.Escher(
				getEscherConfig(
					&escherConfig.AccessKeyId,
					&escherConfig.ApiSecret,
					escherConfig.GetCredentialScope(),
				),
			)
			signedEscherRequest := escherSigner.SignRequest(getEscherRequest(r), []string{"host"})
			assignHeaders(r.Header, signedEscherRequest.Headers)

			if *isVerbose {
				log.Println("Headers", r.Header)
			}

			return r, nil
		})

	log.Println("Escher Pr0xy " + VERSION + " | Listening on: " + *addr)
	log.Fatal(http.ListenAndServe(*addr, proxy))
}

type EscherConfig struct {
	Host            string
	AccessKeyId     string
	ApiSecret       string
	CredentialScope string
}

func (e *EscherConfig) GetCredentialScope() *string {
	if e.CredentialScope == "" {
		credentialScope := "eu/suite/ems_request"

		return &credentialScope
	}

	return &e.CredentialScope
}

type JsonConfig struct {
	KeyDB []EscherConfig
}

func (j *JsonConfig) getEscherConfigByHost(host string) *EscherConfig {
	for _, escherConfig := range j.KeyDB {
		if host == escherConfig.Host {
			return &escherConfig
		}
	}

	log.Panicln("Escher config not found for given host: " + host)
	return nil
}

func readFromFile(file string) []byte {
	jsonFile, err := os.Open(file)

	if err != nil {
		return []byte("")
	}

	jsonData, _ := ioutil.ReadAll(jsonFile)

	return jsonData
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

func getJsonConfig(jsonData []byte) JsonConfig {
	var jsonConfig JsonConfig

	err := json.Unmarshal(jsonData, &jsonConfig)

	if err != nil {
		log.Println("Invalid json config file: " + err.Error())
	}

	return jsonConfig
}

func getEscherRequest(r *http.Request) escher.EscherRequest {
	url := r.URL.Path
	query := r.URL.RawQuery
	if query != "" {
		url += "?" + query
	}

	return escher.EscherRequest{
		Method:  r.Method,
		Url:     url,
		Headers: extractHeaders(r.Header),
		Body:    getBodyAsString(r.Body),
	}
}

func extractHeaders(header http.Header) [][2]string {
	var headers [][2]string
	for name, values := range header {
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

func assignHeaders(header http.Header, escherRequestHeaders escher.EscherRequestHeaders) {
	for _, escherRequestHeader := range escherRequestHeaders {
		header.Set(escherRequestHeader[0], escherRequestHeader[1])
	}
}
