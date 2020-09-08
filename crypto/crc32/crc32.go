package crc32

import (
	"hash/crc32"

	"github.com/realjf/goframe/utils/conv"
)

func Encrypt(v interface{}) uint32 {
	return crc32.ChecksumIEEE(conv.Bytes(v))
}
