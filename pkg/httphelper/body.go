package httphelper

import (
	"bytes"
	"io"
	"io/ioutil"
)

func GetBytesFromBody(body io.ReadCloser) []byte {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = ioutil.ReadAll(body)
	}

	return bodyBytes
}

func GetBodyFromBytes(data []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBuffer(data))
}
