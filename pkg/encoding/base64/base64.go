package base64

import (
	"encoding/base64"
	"io/ioutil"

	"github.com/realjf/goframe/pkg/utils/conv"
)

func Encode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

func Decode(dst []byte) ([]byte, error) {
	src := make([]byte, base64.StdEncoding.DecodedLen(len(dst)))
	n, err := base64.StdEncoding.Decode(src, dst)
	return src[:n], err
}

func EncodeString(src string) string {
	return EncodeToString([]byte(src))
}

func EncodeToString(src []byte) string {
	return conv.UnsafeBytesToString(Encode(src))
}

// base64加密文件
func EncodeFile(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Encode(content), nil
}

func EncodeFileToString(path string) (string, error) {
	content, err := EncodeFile(path)
	if err != nil {
		return "", err
	}
	return conv.UnsafeBytesToString(content), nil
}

func DecodeString(str string) ([]byte, error) {
	return Decode([]byte(str))
}

func DecodeToString(str string) (string, error) {
	b, err := DecodeString(str)
	return conv.UnsafeBytesToString(b), err
}
