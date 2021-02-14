package httphelper

import "net/http"

type RequestBody struct {
	*http.Request
}

func (request RequestBody) readWithoutClear() []byte {
	bodyBytes := GetBytesFromBody(request.Body)
	request.Body = GetBodyFromBytes(bodyBytes)

	return bodyBytes
}
