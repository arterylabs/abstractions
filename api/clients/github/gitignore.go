package github_client

import (
	"path/filepath"
	"strings"
)

// GitIgnoreMatcher applies gitignore-style path rules (subset sufficient for indexing).
// Supported: comments, negation (!), directory patterns (trailing /), ** and * globs,
// and leading-slash anchor to repo root.
type GitIgnoreMatcher struct {
	rules []ignoreRule
}

type ignoreRule struct {
	negated bool
	dirOnly bool
	pattern string
}

// ParseGitIgnore parses gitignore file contents into a matcher.
func ParseGitIgnore(content string) *GitIgnoreMatcher {
	m := &GitIgnoreMatcher{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimRight(line, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		negated := false
		if strings.HasPrefix(trimmed, "!") {
			negated = true
			trimmed = strings.TrimSpace(trimmed[1:])
			if trimmed == "" {
				continue
			}
		}
		dirOnly := strings.HasSuffix(trimmed, "/")
		if dirOnly {
			trimmed = strings.TrimSuffix(trimmed, "/")
		}
		// Unescape leading \# and \!
		if strings.HasPrefix(trimmed, `\`) && len(trimmed) > 1 {
			trimmed = trimmed[1:]
		}
		m.rules = append(m.rules, ignoreRule{
			negated: negated,
			dirOnly: dirOnly,
			pattern: trimmed,
		})
	}
	return m
}

// Merge combines multiple matchers; later rules override earlier ones for the same path.
func (m *GitIgnoreMatcher) Merge(other *GitIgnoreMatcher) *GitIgnoreMatcher {
	if other == nil {
		return m
	}
	if m == nil {
		return other
	}
	combined := &GitIgnoreMatcher{
		rules: make([]ignoreRule, 0, len(m.rules)+len(other.rules)),
	}
	combined.rules = append(combined.rules, m.rules...)
	combined.rules = append(combined.rules, other.rules...)
	return combined
}

// Match returns true if path should be ignored.
// path uses forward slashes relative to the repository root.
func (m *GitIgnoreMatcher) Match(path string) bool {
	if m == nil || len(m.rules) == 0 {
		return false
	}
	path = strings.TrimPrefix(path, "./")
	ignored := false
	for _, rule := range m.rules {
		if rule.dirOnly {
			// Directory rules match the directory itself and anything under it.
			if matchGitPattern(rule.pattern, path) || hasPathPrefixMatch(rule.pattern, path) {
				if rule.negated {
					ignored = false
				} else {
					ignored = true
				}
			}
			continue
		}
		if matchGitPattern(rule.pattern, path) {
			if rule.negated {
				ignored = false
			} else {
				ignored = true
			}
		}
	}
	return ignored
}

func hasPathPrefixMatch(pattern, path string) bool {
	// "build" dir-only should match "build/foo"
	clean := strings.TrimPrefix(pattern, "/")
	if strings.HasPrefix(path, clean+"/") {
		return true
	}
	// Pattern may be a glob for a directory segment.
	parts := strings.Split(path, "/")
	for i := range parts {
		prefix := strings.Join(parts[:i+1], "/")
		if matchGitPattern(pattern, prefix) {
			return true
		}
	}
	return false
}

func matchGitPattern(pattern, path string) bool {
	pattern = strings.TrimPrefix(pattern, "/")
	// "**/" prefix: match anywhere
	if strings.HasPrefix(pattern, "**/") {
		rest := pattern[3:]
		if matchGitPattern(rest, path) {
			return true
		}
		// any prefix path
		for i := 0; i < len(path); i++ {
			if path[i] == '/' && matchGitPattern(rest, path[i+1:]) {
				return true
			}
		}
		return false
	}
	// "/**" suffix
	if strings.HasSuffix(pattern, "/**") {
		prefix := strings.TrimSuffix(pattern, "/**")
		return path == prefix || strings.HasPrefix(path, prefix+"/")
	}
	// Pattern with slash: match from root (or full path)
	if strings.Contains(pattern, "/") {
		ok, _ := filepath.Match(pattern, path)
		if ok {
			return true
		}
		// Also try matching trailing segments
		return false
	}
	// Basename-only pattern matches any path segment / basename
	base := filepath.Base(path)
	if ok, _ := filepath.Match(pattern, base); ok {
		return true
	}
	// Match full path last segment-style: "foo" matches "a/foo" basename already covered
	if ok, _ := filepath.Match(pattern, path); ok {
		return true
	}
	return false
}
