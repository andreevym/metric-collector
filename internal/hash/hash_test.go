package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeHash(t *testing.T) {
	sourceMsg := "source message"
	key := "testpassword"
	encodeHash := EncodeHash([]byte(sourceMsg), key)
	require.NotEmpty(t, sourceMsg)
	resultMsg := DecodeHash(encodeHash, key)
	require.NotEmpty(t, resultMsg)
	require.Equal(t, sourceMsg, resultMsg)
}
