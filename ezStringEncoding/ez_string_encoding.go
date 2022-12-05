package DataTransfer

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
)

func GBKToUTF8(data []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	return io.ReadAll(reader)
}
