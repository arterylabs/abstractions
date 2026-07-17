package github_client

import (
	"path/filepath"
	"strings"
)

const maxFileSizeBytes = 500 * 1024 // 500 KB

// ignoredDirs are directory prefixes that should never be indexed.
var ignoredDirs = []string{
	"node_modules/",
	"vendor/",
	"dist/",
	"build/",
	".git/",
	"__pycache__/",
	".next/",
	"out/",
	"coverage/",
	"target/", // Rust/Java build dir
	".gradle/",
	".idea/",
	".vscode/",
	"public/static/", // bundled static assets
}

// ignoredExtensions are file extensions that produce no indexable signal.
var ignoredExtensions = []string{
	".lock", // package-lock.json, yarn.lock, Cargo.lock, go.sum analog
	".sum",  // go.sum
	".min.js",
	".min.css",
	".map", // source maps
	".png",
	".jpg",
	".jpeg",
	".gif",
	".svg",
	".ico",
	".webp",
	".woff",
	".woff2",
	".ttf",
	".eot",
	".otf",
	".pdf",
	".zip",
	".tar",
	".gz",
	".tgz",
	".mp4",
	".mp3",
	".mov",
	".avi",
	".pyc",
	".class", // Java bytecode
	".o",     // C/C++ object files
	".a",     // static libraries
	".so",    // shared libraries
	".dylib",
	".exe",
	".dll",
}

// ignoredFileNames are exact file base names to skip.
var ignoredFileNames = []string{
	".env",
	".env.local",
	".env.production",
	".env.staging",
	".DS_Store",
	"Thumbs.db",
	"package-lock.json",
	"yarn.lock",
	"pnpm-lock.yaml",
	"go.sum",
	"Pipfile.lock",
	"poetry.lock",
	"Gemfile.lock",
	"composer.lock",
}

// ignoredMigrationSuffix matches SQL migration files (they are schema, not domain logic).
var ignoredMigrationSuffix = []string{
	".up.sql",
	".down.sql",
}

// ShouldIgnore returns true when the file at path should be excluded from indexing.
// sizeBytes is the file's size in bytes as reported by the GitHub Trees API.
func ShouldIgnore(path string, sizeBytes int) bool {
	if int64(sizeBytes) > maxFileSizeBytes {
		return true
	}

	// Check directory prefixes.
	lowerPath := strings.ToLower(path)
	for _, dir := range ignoredDirs {
		if strings.HasPrefix(lowerPath, dir) || strings.Contains(lowerPath, "/"+dir) {
			return true
		}
	}

	base := filepath.Base(path)
	lowerBase := strings.ToLower(base)

	// Exact file name match.
	for _, name := range ignoredFileNames {
		if lowerBase == strings.ToLower(name) {
			return true
		}
	}

	// Migration SQL files.
	for _, suffix := range ignoredMigrationSuffix {
		if strings.HasSuffix(lowerBase, suffix) {
			return true
		}
	}

	// Extension match — handle compound extensions like .min.js first.
	fullLower := strings.ToLower(path)
	for _, ext := range ignoredExtensions {
		if strings.HasSuffix(fullLower, ext) {
			return true
		}
	}

	return false
}
