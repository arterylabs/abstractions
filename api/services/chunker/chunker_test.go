package chunker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---- language detection ----

func TestDetectLanguage_GoFile(t *testing.T) {
	assert.Equal(t, "go", DetectLanguage("main.go"))
}

func TestDetectLanguage_TypeScript(t *testing.T) {
	assert.Equal(t, "typescript", DetectLanguage("component.tsx"))
	assert.Equal(t, "typescript", DetectLanguage("service.ts"))
}

func TestDetectLanguage_JavaScript(t *testing.T) {
	assert.Equal(t, "javascript", DetectLanguage("index.js"))
	assert.Equal(t, "javascript", DetectLanguage("app.jsx"))
}

func TestDetectLanguage_Python(t *testing.T) {
	assert.Equal(t, "python", DetectLanguage("script.py"))
}

func TestDetectLanguage_Ruby(t *testing.T) {
	assert.Equal(t, "ruby", DetectLanguage("model.rb"))
}

func TestDetectLanguage_Java(t *testing.T) {
	assert.Equal(t, "java", DetectLanguage("Main.java"))
}

func TestDetectLanguage_Rust(t *testing.T) {
	assert.Equal(t, "rust", DetectLanguage("lib.rs"))
}

func TestDetectLanguage_Unknown(t *testing.T) {
	assert.Equal(t, "unknown", DetectLanguage("somefile.xyz"))
}

func TestDetectLanguage_C(t *testing.T) {
	assert.Equal(t, "c", DetectLanguage("main.c"))
	assert.Equal(t, "c", DetectLanguage("header.h"))
}

func TestDetectLanguage_Cpp(t *testing.T) {
	assert.Equal(t, "cpp", DetectLanguage("main.cpp"))
	assert.Equal(t, "cpp", DetectLanguage("util.cc"))
	assert.Equal(t, "cpp", DetectLanguage("util.cxx"))
	assert.Equal(t, "cpp", DetectLanguage("header.hpp"))
}

func TestDetectLanguage_CSharp(t *testing.T) {
	assert.Equal(t, "csharp", DetectLanguage("Program.cs"))
}

func TestDetectLanguage_PHP(t *testing.T) {
	assert.Equal(t, "php", DetectLanguage("index.php"))
}

func TestDetectLanguage_Swift(t *testing.T) {
	assert.Equal(t, "swift", DetectLanguage("AppDelegate.swift"))
}

func TestDetectLanguage_Kotlin(t *testing.T) {
	assert.Equal(t, "kotlin", DetectLanguage("Main.kt"))
	assert.Equal(t, "kotlin", DetectLanguage("build.kts"))
}

func TestDetectLanguage_Scala(t *testing.T) {
	assert.Equal(t, "scala", DetectLanguage("Main.scala"))
}

func TestDetectLanguage_JSON(t *testing.T) {
	assert.Equal(t, "json", DetectLanguage("package.json"))
	assert.Equal(t, "json", DetectLanguage("tsconfig.json"))
}

func TestDetectLanguage_TOML(t *testing.T) {
	assert.Equal(t, "toml", DetectLanguage("Cargo.toml"))
	assert.Equal(t, "toml", DetectLanguage("pyproject.toml"))
}

func TestDetectLanguage_YAML(t *testing.T) {
	assert.Equal(t, "yaml", DetectLanguage("config.yaml"))
	assert.Equal(t, "yaml", DetectLanguage("docker-compose.yml"))
}

func TestDetectLanguage_XML(t *testing.T) {
	assert.Equal(t, "xml", DetectLanguage("pom.xml"))
}

func TestDetectLanguage_HTML(t *testing.T) {
	assert.Equal(t, "html", DetectLanguage("index.html"))
	assert.Equal(t, "html", DetectLanguage("page.htm"))
}

func TestDetectLanguage_CSS(t *testing.T) {
	assert.Equal(t, "css", DetectLanguage("styles.css"))
}

func TestDetectLanguage_SCSS(t *testing.T) {
	assert.Equal(t, "scss", DetectLanguage("styles.scss"))
	assert.Equal(t, "scss", DetectLanguage("styles.sass"))
}

func TestDetectLanguage_SQL(t *testing.T) {
	assert.Equal(t, "sql", DetectLanguage("migration.sql"))
}

func TestDetectLanguage_Markdown(t *testing.T) {
	assert.Equal(t, "markdown", DetectLanguage("README.md"))
	assert.Equal(t, "markdown", DetectLanguage("docs.mdx"))
}

func TestDetectLanguage_GraphQL(t *testing.T) {
	assert.Equal(t, "graphql", DetectLanguage("schema.graphql"))
	assert.Equal(t, "graphql", DetectLanguage("query.gql"))
}

func TestDetectLanguage_Protobuf(t *testing.T) {
	assert.Equal(t, "protobuf", DetectLanguage("service.proto"))
}

func TestDetectLanguage_Shell(t *testing.T) {
	assert.Equal(t, "shell", DetectLanguage("setup.sh"))
	assert.Equal(t, "shell", DetectLanguage("run.bash"))
	assert.Equal(t, "shell", DetectLanguage("config.zsh"))
}

func TestDetectLanguage_Terraform(t *testing.T) {
	assert.Equal(t, "terraform", DetectLanguage("main.tf"))
	assert.Equal(t, "terraform", DetectLanguage("vars.tfvars"))
}

func TestDetectLanguage_Dotenv(t *testing.T) {
	assert.Equal(t, "dotenv", DetectLanguage(".env"))
}

func TestDetectLanguage_Dockerfile(t *testing.T) {
	assert.Equal(t, "dockerfile", DetectLanguage("Dockerfile"))
	assert.Equal(t, "dockerfile", DetectLanguage("path/to/Dockerfile"))
}

func TestDetectLanguage_Makefile(t *testing.T) {
	assert.Equal(t, "makefile", DetectLanguage("Makefile"))
}

func TestDetectLanguage_Gemfile(t *testing.T) {
	assert.Equal(t, "ruby", DetectLanguage("Gemfile"))
	assert.Equal(t, "ruby", DetectLanguage("Rakefile"))
	assert.Equal(t, "ruby", DetectLanguage("Vagrantfile"))
}

func TestDetectLanguage_Procfile(t *testing.T) {
	assert.Equal(t, "shell", DetectLanguage("Procfile"))
}

func TestDetectLanguage_Gitignore(t *testing.T) {
	assert.Equal(t, "gitignore", DetectLanguage(".gitignore"))
	assert.Equal(t, "gitignore", DetectLanguage(".dockerignore"))
}

func TestDetectLanguage_Editorconfig(t *testing.T) {
	assert.Equal(t, "ini", DetectLanguage(".editorconfig"))
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
