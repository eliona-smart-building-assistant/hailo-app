package hailo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCodingBase64(t *testing.T) {
	const s = "0abc1defghij2klmnop3qrstuvwxyz4ABCDEF5GHIJKLM6NOPQ7RSTUVW8XY9Z"
	enc := encodeBase64(s)
	dec := decodeBase64(enc)
	assert.Equal(t, dec, s)
}
