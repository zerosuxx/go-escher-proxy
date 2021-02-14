package httphelper

type BodyReader interface {
	readWithoutClear() []byte
}

func ReadBodyWithoutClear(bodyReader BodyReader) []byte {
	return bodyReader.readWithoutClear()
}
