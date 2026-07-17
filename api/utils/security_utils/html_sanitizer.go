package security_utils

import (
	"bytes"
	"io"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

func getPolicy() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowImages()
	p.AllowDataURIImages()
	p.AllowElements("iframe")
	p.AllowAttrs("width").Matching(bluemonday.Number).OnElements("iframe")
	p.AllowAttrs("height").Matching(bluemonday.Number).OnElements("iframe")
	p.AllowAttrs("src").OnElements("iframe")
	p.AllowAttrs("style").OnElements("iframe")
	p.AllowAttrs("frameborder").Matching(bluemonday.Number).OnElements("iframe")
	p.AllowAttrs("allow").Matching(regexp.MustCompile(`[a-z; -]*`)).OnElements("iframe")
	p.AllowAttrs("allowfullscreen").OnElements("iframe")
	return p
}

func Sanitize(value string) string {
	return getPolicy().Sanitize(value)
}

func SanitizeBytes(value []byte) []byte {
	return getPolicy().SanitizeBytes(value)
}

func SanitizeReader(value io.Reader) *bytes.Buffer {
	return getPolicy().SanitizeReader(value)
}
