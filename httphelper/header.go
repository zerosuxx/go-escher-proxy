package httphelper

import (
	"net/http"
)

func AssignHeaders(header http.Header, additionalHeaders [][2]string) {
	for _, escherRequestHeader := range additionalHeaders {
		header.Set(escherRequestHeader[0], escherRequestHeader[1])
	}
}

func ExtractHeaders(header http.Header) [][2]string {
	var headers [][2]string
	for name, values := range header {
		for _, value := range values {
			headers = append(headers, [2]string{name, value})
		}
	}

	return headers
}
