package sha1

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"

	"github.com/realjf/goframe/pkg/utils/conv"
)

func Encrypt(v interface{}) string {
	r := sha1.Sum(conv.Bytes(v))
	return hex.EncodeToString(r[:])
}

func EncryptFile(path string) (encrypt string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha1.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
