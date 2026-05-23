package chunker

import "regexp"

// Language name constants for every language that detectLanguage can return.
// Using constants avoids scattered string literals and allows callers to reference
// a specific language without risking typos.
const (
	// Fully-supported languages: declaration-based chunking with topLevelPatterns.
	LangGo         = "go"
	LangTypeScript = "typescript"
	LangJavaScript = "javascript"
	LangPython     = "python"
	LangRuby       = "ruby"
	LangJava       = "java"
	LangRust       = "rust"

	// Detected languages that fall back to blank-line chunking.
	LangC          = "c"
	LangCpp        = "cpp"
	LangCSharp     = "csharp"
	LangPHP        = "php"
	LangSwift      = "swift"
	LangKotlin     = "kotlin"
	LangScala      = "scala"
	LangJSON       = "json"
	LangTOML       = "toml"
	LangYAML       = "yaml"
	LangXML        = "xml"
	LangHTML       = "html"
	LangCSS        = "css"
	LangSCSS       = "scss"
	LangSQL        = "sql"
	LangMarkdown   = "markdown"
	LangGraphQL    = "graphql"
	LangProtobuf   = "protobuf"
	LangShell      = "shell"
	LangTerraform  = "terraform"
	LangDotenv     = "dotenv"
	LangDockerfile = "dockerfile"
	LangMakefile   = "makefile"
	LangGitignore  = "gitignore"
	LangINI        = "ini"
	LangUnknown    = "unknown"
)

// SupportedLanguages lists every language that has declaration-based chunking via
// topLevelPatterns. Languages not in this list fall back to blank-line chunking.
var SupportedLanguages = []string{
	LangGo,
	LangTypeScript,
	LangJavaScript,
	LangPython,
	LangRuby,
	LangJava,
	LangRust,
}

// topLevelPatterns maps language names to the regex that matches the start of a
// top-level declaration. The match anchors to the beginning of a line.
// Each pattern is tested by supported_languages_test.go to guard against regressions.
var topLevelPatterns = map[string]*regexp.Regexp{
	// Go: func, type, var, const at col 0 (not indented)
	LangGo: regexp.MustCompile(`^(func |type |var |const )`),
	// TypeScript/JavaScript: function, class, const/let/var arrow-fn, export variants
	LangTypeScript: regexp.MustCompile(`^(export\s+)?(default\s+)?(async\s+)?(function\b|class\b|const\b|let\b|var\b|interface\b|type\b|enum\b)`),
	LangJavaScript: regexp.MustCompile(`^(export\s+)?(default\s+)?(async\s+)?(function\b|class\b|const\b|let\b|var\b)`),
	// Python: def or class at col 0
	LangPython: regexp.MustCompile(`^(def |class |async def )`),
	// Ruby: def, class, module at col 0
	LangRuby: regexp.MustCompile(`^(def |class |module )`),
	// Java: public/private/protected method or class
	LangJava: regexp.MustCompile(`^(public|private|protected|class|interface|enum|@interface)\b`),
	// Rust: fn, struct, enum, impl, trait, mod at col 0
	LangRust: regexp.MustCompile(`^(pub(\([^)]*\))?\s+)?(fn |struct |enum |impl |trait |mod |type |const |static )`),
}

// importLinePatterns maps language names to a set of prefixes that indicate a line
// is part of an import/package declaration block (no executable logic).
var importLinePatterns = map[string][]string{
	LangGo:         {"package ", "import ", "//", "/*", " *", "*/", "(", ")", "\""},
	LangTypeScript: {"import ", "//", "/*", " *", "*/"},
	LangJavaScript: {"import ", "require(", "//", "/*", " *", "*/"},
	LangPython:     {"import ", "from ", "#"},
	LangRuby:       {"require ", "require_relative ", "#"},
	LangJava:       {"import ", "package ", "//", "/*", " *", "*/"},
	LangRust:       {"use ", "extern crate ", "//", "/*", " *", "*/"},
}
