package password_utils_test

import (
	"testing"

	"abstractions/utils/password_utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHashAndCompare(t *testing.T) {
	plain := []byte("correct-horse-battery-staple")
	hash := password_utils.GeneratePasswordHash(plain)
	require.NotEmpty(t, hash)
	assert.NotEqual(t, string(plain), hash)

	assert.True(t, password_utils.ComparePasswords(hash, plain))
	assert.False(t, password_utils.ComparePasswords(hash, []byte("wrong-password")))
	assert.False(t, password_utils.ComparePasswords("not-a-bcrypt-hash", plain))
}
