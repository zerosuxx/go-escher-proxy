package httphelper

import "net/http"

type ResponseBody struct {
	*http.Response
}

func (response ResponseBody) readWithoutClear() []byte {
	bodyBytes := GetBytesFromBody(response.Body)
	response.Body = GetBodyFromBytes(bodyBytes)

	return bodyBytes
}
