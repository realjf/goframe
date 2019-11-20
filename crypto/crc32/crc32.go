package crc32

import (
	"goframe/utils/conv"
	"hash/crc32"
)

func Encrypt(v interface{}) uint32 {
	return crc32.ChecksumIEEE(conv.Bytes(v))
}


