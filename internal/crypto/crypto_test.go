package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrypto(t *testing.T) {
	t.Run("positive: success encode, decode and match source and decoded", func(t *testing.T) {
		private, public, err := GenerateCrypto()
		assert.NoError(t, err)
		require.NotEmpty(t, private)
		require.NotEmpty(t, public)

		sourceMessage := "test 1234567890-=@#$%^&*()XCVBNM<>?QWERTYUIOP{}:LKJHGASDFGHJKL:"

		encodedMessage, err := Encode(public, sourceMessage)
		assert.NoError(t, err)
		assert.NotEmpty(t, encodedMessage)

		decodedMessage, err := Decode(private, encodedMessage)
		assert.NoError(t, err)
		assert.NotEmpty(t, decodedMessage)

		assert.Equal(t, sourceMessage, decodedMessage)
	})
	t.Run("negative: private and public key not matched", func(t *testing.T) {
		_, public, err := GenerateCrypto()
		assert.NoError(t, err)
		require.NotEmpty(t, public)

		wrongPrivate, _, err := GenerateCrypto()
		assert.NoError(t, err)
		require.NotEmpty(t, wrongPrivate)

		sourceMessage := "test 1234567890-=@#$%^&*()XCVBNM<>?QWERTYUIOP{}:LKJHGASDFGHJKL:"

		encodedMessage, err := Encode(public, sourceMessage)
		assert.NoError(t, err)
		assert.NotEmpty(t, encodedMessage)

		decodedMessage, err := Decode(wrongPrivate, encodedMessage)
		assert.Error(t, err)
		assert.Empty(t, decodedMessage)

		assert.NotEqual(t, sourceMessage, decodedMessage)
	})
	t.Run("negative: public key empty", func(t *testing.T) {
		sourceMessage := "test 1234567890-=@#$%^&*()XCVBNM<>?QWERTYUIOP{}:LKJHGASDFGHJKL:"
		encodedMessage, err := Encode("", sourceMessage)
		assert.Error(t, err)
		assert.Empty(t, encodedMessage)
	})
	t.Run("negative: source message empty", func(t *testing.T) {
		_, public, err := GenerateCrypto()
		assert.NoError(t, err)
		require.NotEmpty(t, public)

		sourceMessage := ""
		encodedMessage, err := Encode(public, sourceMessage)
		assert.Error(t, err)
		assert.Empty(t, encodedMessage)
	})
	t.Run("negative: private key empty", func(t *testing.T) {
		_, public, err := GenerateCrypto()
		assert.NoError(t, err)
		require.NotEmpty(t, public)

		sourceMessage := "test 1234567890-=@#$%^&*()XCVBNM<>?QWERTYUIOP{}:LKJHGASDFGHJKL:"

		encodedMessage, err := Encode(public, sourceMessage)
		assert.NoError(t, err)
		assert.NotEmpty(t, encodedMessage)

		decodedMessage, err := Decode("", encodedMessage)
		assert.Error(t, err)
		assert.Empty(t, decodedMessage)
	})
}
