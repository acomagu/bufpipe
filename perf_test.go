package bufpipe

import (
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkReadOnly(b *testing.B) {
	length := 2
	buf := make([]byte, b.N, b.N)
	r, w := New(buf, 0)
	w.Close()
	data := make([]byte, length)
	b.ResetTimer()
	for {
		_, err := io.ReadFull(r, data)
		if err != nil {
			if math.Mod(float64(b.N), float64(length)) == 0 {
				assert.EqualError(b, err, "EOF")
			} else {
				assert.EqualError(b, err, "unexpected EOF")
			}
			break
		}
	}
}
