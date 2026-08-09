package security_utils_test

import (
	"bytes"
	"strings"
	"testing"

	"abstractions/utils/security_utils"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeStripsScriptTags(t *testing.T) {
	input := `<p>Hello</p><script>alert("xss")</script><img src=x onerror=alert(1)>`
	out := security_utils.Sanitize(input)

	assert.Contains(t, out, "Hello")
	assert.NotContains(t, strings.ToLower(out), "<script")
	assert.NotContains(t, strings.ToLower(out), "onerror")
	assert.NotContains(t, out, "alert(")
}

func TestSanitizeBytesStripsScriptTags(t *testing.T) {
	input := []byte(`<b>ok</b><script>evil()</script>`)
	out := security_utils.SanitizeBytes(input)
	assert.Contains(t, string(out), "ok")
	assert.NotContains(t, strings.ToLower(string(out)), "<script")
}

func TestSanitizeReaderStripsScriptTags(t *testing.T) {
	input := bytes.NewBufferString(`safe <script>document.cookie</script> text`)
	out := security_utils.SanitizeReader(input)
	assert.NotNil(t, out)
	s := out.String()
	assert.Contains(t, s, "safe")
	assert.NotContains(t, strings.ToLower(s), "<script")
}

func TestSanitizeAllowsSafeMarkup(t *testing.T) {
	input := `<p><strong>Bold</strong> and <em>italic</em></p>`
	out := security_utils.Sanitize(input)
	assert.Contains(t, out, "Bold")
	assert.Contains(t, out, "italic")
}
