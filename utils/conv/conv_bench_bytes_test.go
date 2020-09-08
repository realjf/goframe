package conv

import (
	"github.com/realjf/goframe/encoding/binary"
	"testing"
)

var value = binary.Encode(123456789)

func BenchmarkBytesToStringNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = string(value)
	}
}
