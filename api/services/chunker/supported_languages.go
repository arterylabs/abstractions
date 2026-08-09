package chunker

import "regexp"

// Language name constants for every language that detectLanguage can return.
const (
	// Fully-supported languages: declaration-based chunking with topLevelPatterns.
	LangGo         = "go"
	LangTypeScript = "typescript"
	LangJavaScript = "javascript"
	LangPython     = "python"
	LangRuby       = "ruby"
	LangJava       = "java"
	LangRust       = "rust"
	LangC          = "c"
	LangCpp        = "cpp"
	LangCSharp     = "csharp"
	LangPHP        = "php"
	LangSwift      = "swift"
	LangKotlin     = "kotlin"
	LangScala      = "scala"
	LangShell      = "shell"
	LangTerraform  = "terraform"
	LangGraphQL    = "graphql"
	LangProtobuf   = "protobuf"
	LangSQL        = "sql"

	// Detected languages with specialized or blank-line chunking.
	LangJSON       = "json"
	LangTOML       = "toml"
	LangYAML       = "yaml"
	LangXML        = "xml"
	LangHTML       = "html"
	LangCSS        = "css"
	LangSCSS       = "scss"
	LangMarkdown   = "markdown"
	LangDotenv     = "dotenv"
	LangDockerfile = "dockerfile"
	LangMakefile   = "makefile"
	LangGitignore  = "gitignore"
	LangINI        = "ini"
	LangUnknown    = "unknown"
)

// SupportedLanguages lists every language that has declaration-based chunking via
// topLevelPatterns (or specialized heading chunking for markdown).
var SupportedLanguages = []string{
	LangGo,
	LangTypeScript,
	LangJavaScript,
	LangPython,
	LangRuby,
	LangJava,
	LangRust,
	LangC,
	LangCpp,
	LangCSharp,
	LangPHP,
	LangSwift,
	LangKotlin,
	LangScala,
	LangShell,
	LangTerraform,
	LangGraphQL,
	LangProtobuf,
	LangSQL,
	// LangMarkdown uses heading-based chunking (handled specially in ChunkFile), not topLevelPatterns.
}

// topLevelPatterns maps language names to the regex that matches the start of a
// top-level declaration. The match anchors to the beginning of a line.
var topLevelPatterns = map[string]*regexp.Regexp{
	LangGo:         regexp.MustCompile(`^(func |type |var |const )`),
	LangTypeScript: regexp.MustCompile(`^(export\s+)?(default\s+)?(async\s+)?(function\b|class\b|const\b|let\b|var\b|interface\b|type\b|enum\b)`),
	LangJavaScript: regexp.MustCompile(`^(export\s+)?(default\s+)?(async\s+)?(function\b|class\b|const\b|let\b|var\b)`),
	LangPython:     regexp.MustCompile(`^(def |class |async def |@)`),
	LangRuby:       regexp.MustCompile(`^(def |class |module )`),
	LangJava:       regexp.MustCompile(`^(public|private|protected|class|interface|enum|@interface|@)\b`),
	LangRust:       regexp.MustCompile(`^(pub(\([^)]*\))?\s+)?(fn |struct |enum |impl |trait |mod |type |const |static |#\[)`),
	LangC:          regexp.MustCompile(`^([a-zA-Z_][\w\s\*]+\s+)?[a-zA-Z_]\w*\s*\([^;]*\)\s*\{?$|^(struct |enum |typedef |#define |#include )`),
	LangCpp:        regexp.MustCompile(`^(class |struct |namespace |template\b|enum |typedef |using |#include )|^([a-zA-Z_][\w:<>,\s\*\&]+)\s+[a-zA-Z_]\w*\s*\(`),
	LangCSharp:     regexp.MustCompile(`^(public|private|protected|internal|static|partial|class|interface|struct|enum|namespace|record)\b`),
	LangPHP:        regexp.MustCompile(`^(function |class |interface |trait |namespace |abstract |final |public function |private function |protected function )`),
	LangSwift:      regexp.MustCompile(`^(public |private |internal |open |fileprivate |final |static )?(func |class |struct |enum |protocol |extension |actor |typealias )`),
	LangKotlin:     regexp.MustCompile(`^(public |private |internal |protected |open |abstract |data |sealed )?(fun |class |interface |object |enum class |typealias )`),
	LangScala:      regexp.MustCompile(`^(private |protected |override |final |abstract )?(def |class |object |trait |type |val |var |enum )`),
	LangShell:      regexp.MustCompile(`^[a-zA-Z_][\w]*\s*\(\)\s*\{|^function\s+[a-zA-Z_]`),
	LangTerraform:  regexp.MustCompile(`^(resource |data |module |variable |output |locals |provider |terraform )\b`),
	LangGraphQL:    regexp.MustCompile(`^(type |input |enum |interface |union |scalar |extend |query |mutation |subscription )\b`),
	LangProtobuf:   regexp.MustCompile(`^(message |service |enum |rpc |extend |option |import |syntax |package )\b`),
	LangSQL:        regexp.MustCompile(`^(?i)(CREATE |ALTER |DROP |INSERT |UPDATE |DELETE |SELECT |WITH |GRANT |REVOKE )`),
}

// importLinePatterns maps language names to prefixes that indicate a line is part of
// an import/package declaration block (no executable logic).
var importLinePatterns = map[string][]string{
	LangGo:         {"package ", "import "},
	LangTypeScript: {"import ", "export * from", "export {", "export type {"},
	LangJavaScript: {"import ", "require(", "export * from", "export {"},
	LangPython:     {"import ", "from "},
	LangRuby:       {"require ", "require_relative "},
	LangJava:       {"import ", "package "},
	LangRust:       {"use ", "extern crate ", "mod "},
	LangC:          {"#include ", "#import "},
	LangCpp:        {"#include ", "using namespace ", "using "},
	LangCSharp:     {"using "},
	LangPHP:        {"use ", "namespace ", "require ", "require_once ", "include ", "include_once "},
	LangSwift:      {"import "},
	LangKotlin:     {"import ", "package "},
	LangScala:      {"import ", "package "},
	LangProtobuf:   {"import ", "syntax ", "package ", "option "},
	LangGraphQL:    {"#import ", "import "},
	LangTerraform:  {},
	LangShell:      {},
	LangSQL:        {},
}

// commentLinePrefixes are language-agnostic comment line starts used when classifying import chunks.
var commentLinePrefixes = []string{"//", "/*", " *", "*/", "#", "--", "<!--"}
