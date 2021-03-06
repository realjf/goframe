package conv

import (
	"testing"

	"github.com/realjf/goframe/pkg/encoding/binary"
)

var value = binary.Encode(123456789)

func BenchmarkBytesToStringNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = string(value)
	}
}
