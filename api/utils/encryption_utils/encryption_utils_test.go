package encryption_utils_test

import (
	"testing"

	"abstractions/utils/encryption_utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AES-256 requires a 32-byte key.
const testKey = "12345678901234567890123456789012"

func TestEncryptDecryptRoundTrip(t *testing.T) {
	plaintext := "super-secret-api-key-value"

	ciphertext, err := encryption_utils.EncryptString(testKey, plaintext)
	require.NoError(t, err)
	require.NotEmpty(t, ciphertext)
	assert.NotEqual(t, plaintext, ciphertext)

	decrypted, err := encryption_utils.DecryptString(testKey, ciphertext)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestDecryptString_WrongKeyFails(t *testing.T) {
	ciphertext, err := encryption_utils.EncryptString(testKey, "payload")
	require.NoError(t, err)

	wrongKey := "abcdefghijklmnopqrstuvwxyz123456" // also 32 bytes
	_, err = encryption_utils.DecryptString(wrongKey, ciphertext)
	assert.Error(t, err)
}

func TestDecryptString_InvalidHexFails(t *testing.T) {
	_, err := encryption_utils.DecryptString(testKey, "not-hex!!")
	assert.Error(t, err)
}

func TestEncryptString_InvalidKeyLengthFails(t *testing.T) {
	_, err := encryption_utils.EncryptString("short", "payload")
	assert.Error(t, err)
}
