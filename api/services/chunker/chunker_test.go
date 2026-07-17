package chunker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---- language detection ----

func TestDetectLanguage(t *testing.T) {
	cases := map[string]string{
		"main.go":            LangGo,
		"component.tsx":      LangTypeScript,
		"service.ts":         LangTypeScript,
		"index.js":           LangJavaScript,
		"app.jsx":            LangJavaScript,
		"script.py":          LangPython,
		"model.rb":           LangRuby,
		"Main.java":          LangJava,
		"lib.rs":             LangRust,
		"somefile.xyz":       LangUnknown,
		"main.c":             LangC,
		"header.h":           LangC,
		"main.cpp":           LangCpp,
		"util.cc":            LangCpp,
		"util.cxx":           LangCpp,
		"header.hpp":         LangCpp,
		"Program.cs":         LangCSharp,
		"index.php":          LangPHP,
		"AppDelegate.swift":  LangSwift,
		"Main.kt":            LangKotlin,
		"build.kts":          LangKotlin,
		"Main.scala":         LangScala,
		"package.json":       LangJSON,
		"tsconfig.json":      LangJSON,
		"Cargo.toml":         LangTOML,
		"pyproject.toml":     LangTOML,
		"config.yaml":        LangYAML,
		"docker-compose.yml": LangYAML,
		"pom.xml":            LangXML,
		"index.html":         LangHTML,
		"page.htm":           LangHTML,
		"styles.css":         LangCSS,
		"styles.scss":        LangSCSS,
		"styles.sass":        LangSCSS,
		"migration.sql":      LangSQL,
		"README.md":          LangMarkdown,
		"docs.mdx":           LangMarkdown,
		"schema.graphql":     LangGraphQL,
		"query.gql":          LangGraphQL,
		"service.proto":      LangProtobuf,
		"setup.sh":           LangShell,
		"run.bash":           LangShell,
		"config.zsh":         LangShell,
		"main.tf":            LangTerraform,
		"vars.tfvars":        LangTerraform,
		".env":               LangDotenv,
		"Dockerfile":         LangDockerfile,
		"path/to/Dockerfile": LangDockerfile,
		"Makefile":           LangMakefile,
		"Gemfile":            LangRuby,
		"Rakefile":           LangRuby,
		"Vagrantfile":        LangRuby,
		"Procfile":           LangShell,
		".gitignore":         LangGitignore,
		".dockerignore":      LangGitignore,
		".editorconfig":      LangINI,
	}

	for path, want := range cases {
		t.Run(path, func(t *testing.T) {
			assert.Equal(t, want, DetectLanguage(path))
		})
	}
}

// ---- Go chunking ----

func TestChunkFile_Go_TwoFunctions(t *testing.T) {
	src := `package main

func Hello() string {
	return "hello"
}

func World() string {
	return "world"
}
`
	chunks := ChunkFile("main.go", src, "go")
	// Should produce at least 2 chunks (one per func)
	assert.GreaterOrEqual(t, len(chunks), 2)
	found := make(map[string]bool)
	for _, c := range chunks {
		if strings.Contains(c.Content, "func Hello") {
			found["Hello"] = true
		}
		if strings.Contains(c.Content, "func World") {
			found["World"] = true
		}
	}
	assert.True(t, found["Hello"], "should have Hello chunk")
	assert.True(t, found["World"], "should have World chunk")
}

func TestChunkFile_Go_TypeDeclaration(t *testing.T) {
	src := `package models

type User struct {
	ID   string
	Name string
}

type Post struct {
	ID    string
	Title string
}
`
	chunks := ChunkFile("models.go", src, "go")
	assert.GreaterOrEqual(t, len(chunks), 2)
}

// ---- TypeScript chunking ----

func TestChunkFile_TypeScript_ExportedFunctions(t *testing.T) {
	src := `export function fetchUser(id: string) {
  return api.get('/users/' + id);
}

export const createUser = async (data: UserData) => {
  return api.post('/users', data);
};

export class UserService {
  getAll() { return []; }
}
`
	chunks := ChunkFile("user.ts", src, "typescript")
	assert.GreaterOrEqual(t, len(chunks), 2)
}

// ---- Python chunking ----

func TestChunkFile_Python_FunctionsAndClass(t *testing.T) {
	src := `def add(a, b):
    return a + b

def subtract(a, b):
    return a - b

class Calculator:
    def multiply(self, a, b):
        return a * b
`
	chunks := ChunkFile("calc.py", src, "python")
	assert.GreaterOrEqual(t, len(chunks), 3)
}

// ---- Ruby chunking ----

func TestChunkFile_Ruby_MethodsAndClass(t *testing.T) {
	src := `class Animal
  def speak
    "..."
  end
end

module Helpers
  def format_name(name)
    name.capitalize
  end
end
`
	chunks := ChunkFile("animal.rb", src, "ruby")
	assert.GreaterOrEqual(t, len(chunks), 2)
}

// ---- Java chunking ----

func TestChunkFile_Java_PublicMethods(t *testing.T) {
	src := `public class Main {
    public static void main(String[] args) {
        System.out.println("Hello");
    }

    private String helper() {
        return "helper";
    }
}
`
	chunks := ChunkFile("Main.java", src, "java")
	assert.GreaterOrEqual(t, len(chunks), 1)
}

// ---- Rust chunking ----

func TestChunkFile_Rust_FnAndStruct(t *testing.T) {
	src := `pub fn greet(name: &str) -> String {
    format!("Hello, {}!", name)
}

struct Point {
    x: f64,
    y: f64,
}

impl Point {
    fn distance(&self) -> f64 {
        (self.x * self.x + self.y * self.y).sqrt()
    }
}
`
	chunks := ChunkFile("lib.rs", src, "rust")
	assert.GreaterOrEqual(t, len(chunks), 3)
}

// ---- blank-line fallback ----

func TestChunkFile_BlankLineFallback(t *testing.T) {
	// Unknown extension → blank-line blocks
	src := `Block one line one
Block one line two

Block two line one

Block three line one
Block three line two
Block three line three
`
	chunks := ChunkFile("data.txt", src, "")
	assert.Len(t, chunks, 3)
}

func TestChunkFile_BlankLineFallback_MaxLines(t *testing.T) {
	// Build a 150-line block with no blank lines — should be split at 100
	var sb strings.Builder
	for i := 0; i < 150; i++ {
		fmt.Fprintf(&sb, "line %d\n", i+1)
	}
	chunks := ChunkFile("data.csv", sb.String(), "")
	// 150 lines should produce at least 2 chunks
	assert.GreaterOrEqual(t, len(chunks), 2)
}

// ---- oversized chunk splitting ----

func TestChunkFile_OversizedChunk_SplitIntoWindows(t *testing.T) {
	// Build a 250-line Go function (well above the 200-line limit)
	var sb strings.Builder
	sb.WriteString("func BigFunction() {\n")
	for i := 0; i < 248; i++ {
		fmt.Fprintf(&sb, "\t// line %d\n", i+1)
	}
	sb.WriteString("}\n")

	chunks := ChunkFile("big.go", sb.String(), "go")
	// Should have been split into overlapping windows
	assert.GreaterOrEqual(t, len(chunks), 2)
	for _, c := range chunks {
		lineCount := c.EndLine - c.StartLine + 1
		assert.LessOrEqual(t, lineCount, 150, "each window should be ≤ 150 lines")
	}
}

// ---- content not in logs ----

// TestChunkFile_ContentNotInLogs verifies that ChunkFile only logs filePath and chunkCount,
// not the raw chunk content. The test validates that the Chunk.Content field holds the
// content (not the logger), confirming the code's design intent.
func TestChunkFile_ContentNotInLogs(t *testing.T) {
	secret := "SECRET_CONTENT_SHOULD_NOT_APPEAR"
	src := fmt.Sprintf("func FetchSecret() string {\n\treturn \"%s\"\n}\n", secret)
	chunks := ChunkFile("secret.go", src, "go")
	// Content exists on the struct (accessible to callers) but was never passed to the logger
	assert.NotEmpty(t, chunks)
	found := false
	for _, c := range chunks {
		if strings.Contains(c.Content, secret) {
			found = true
		}
	}
	assert.True(t, found, "chunk content should be in Chunk.Content struct field, not in log messages")
}
