package chunker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSupportedLanguages_AllHavePatterns ensures every entry in SupportedLanguages
// has a corresponding regex in topLevelPatterns.
func TestSupportedLanguages_AllHavePatterns(t *testing.T) {
	for _, lang := range SupportedLanguages {
		_, ok := topLevelPatterns[lang]
		assert.True(t, ok, "SupportedLanguages entry %q has no topLevelPatterns entry", lang)
	}
}

// TestSupportedLanguages_NoDuplicates ensures there are no duplicate entries.
func TestSupportedLanguages_NoDuplicates(t *testing.T) {
	seen := make(map[string]bool)
	for _, lang := range SupportedLanguages {
		assert.False(t, seen[lang], "duplicate language in SupportedLanguages: %q", lang)
		seen[lang] = true
	}
}

// ---- Per-language pattern tests ----

func TestTopLevelPattern_Go_Matches(t *testing.T) {
	p := topLevelPatterns[LangGo]
	require.NotNil(t, p)

	matches := []string{
		"func Hello() {}",
		"type User struct {}",
		"var x = 1",
		"const Pi = 3.14",
	}
	for _, line := range matches {
		assert.True(t, p.MatchString(line), "Go pattern should match: %q", line)
	}
}

func TestTopLevelPattern_Go_NoMatch(t *testing.T) {
	p := topLevelPatterns[LangGo]
	require.NotNil(t, p)

	noMatches := []string{
		"\tfunc indented() {}",
		"  var indented = 1",
		"// comment",
		"package main",
		"import \"fmt\"",
	}
	for _, line := range noMatches {
		assert.False(t, p.MatchString(line), "Go pattern should NOT match: %q", line)
	}
}

func TestTopLevelPattern_TypeScript_Matches(t *testing.T) {
	p := topLevelPatterns[LangTypeScript]
	require.NotNil(t, p)

	matches := []string{
		"function hello() {}",
		"class UserService {}",
		"const fetchUser = async () => {}",
		"let count = 0",
		"export function greet() {}",
		"export default class Foo {}",
		"export const bar = 1",
		"interface IUser {}",
		"type Alias = string",
		"enum Direction {}",
		"async function load() {}",
	}
	for _, line := range matches {
		assert.True(t, p.MatchString(line), "TypeScript pattern should match: %q", line)
	}
}

func TestTopLevelPattern_TypeScript_NoMatch(t *testing.T) {
	p := topLevelPatterns[LangTypeScript]
	require.NotNil(t, p)

	noMatches := []string{
		"  function indented() {}",
		"\tclass indented {}",
		"// comment",
		"import { foo } from 'bar'",
	}
	for _, line := range noMatches {
		assert.False(t, p.MatchString(line), "TypeScript pattern should NOT match: %q", line)
	}
}

func TestTopLevelPattern_JavaScript_Matches(t *testing.T) {
	p := topLevelPatterns[LangJavaScript]
	require.NotNil(t, p)

	matches := []string{
		"function hello() {}",
		"class Foo {}",
		"const x = 1",
		"let y = 2",
		"var z = 3",
		"export function greet() {}",
		"export const bar = () => {}",
		"async function load() {}",
	}
	for _, line := range matches {
		assert.True(t, p.MatchString(line), "JavaScript pattern should match: %q", line)
	}
}

func TestTopLevelPattern_Python_Matches(t *testing.T) {
	p := topLevelPatterns[LangPython]
	require.NotNil(t, p)

	matches := []string{
		"def my_function():",
		"class MyClass:",
		"async def fetch_data():",
	}
	for _, line := range matches {
		assert.True(t, p.MatchString(line), "Python pattern should match: %q", line)
	}
}

func TestTopLevelPattern_Python_NoMatch(t *testing.T) {
	p := topLevelPatterns[LangPython]
	require.NotNil(t, p)

	noMatches := []string{
		"    def indented():",
		"\tclass indented:",
		"# comment",
		"import os",
		"from os import path",
	}
	for _, line := range noMatches {
		assert.False(t, p.MatchString(line), "Python pattern should NOT match: %q", line)
	}
}

func TestTopLevelPattern_Ruby_Matches(t *testing.T) {
	p := topLevelPatterns[LangRuby]
	require.NotNil(t, p)

	matches := []string{
		"def my_method",
		"class MyClass",
		"module MyModule",
	}
	for _, line := range matches {
		assert.True(t, p.MatchString(line), "Ruby pattern should match: %q", line)
	}
}

func TestTopLevelPattern_Java_Matches(t *testing.T) {
	p := topLevelPatterns[LangJava]
	require.NotNil(t, p)

	matches := []string{
		"public class Main {",
		"private String helper() {",
		"protected void doSomething() {",
		"class Foo {",
		"interface Bar {",
		"enum Color {",
		"@interface Annotation {",
	}
	for _, line := range matches {
		assert.True(t, p.MatchString(line), "Java pattern should match: %q", line)
	}
}

func TestTopLevelPattern_Rust_Matches(t *testing.T) {
	p := topLevelPatterns[LangRust]
	require.NotNil(t, p)

	matches := []string{
		"fn main() {}",
		"pub fn greet() {}",
		"pub(crate) fn internal() {}",
		"struct Point {",
		"enum Direction {",
		"impl Point {",
		"trait Animal {",
		"mod config {",
		"type Alias = String;",
		"const MAX: u32 = 100;",
		"static GREETING: &str = \"hello\";",
	}
	for _, line := range matches {
		assert.True(t, p.MatchString(line), "Rust pattern should match: %q", line)
	}
}

func TestTopLevelPattern_Rust_NoMatch(t *testing.T) {
	p := topLevelPatterns[LangRust]
	require.NotNil(t, p)

	noMatches := []string{
		"    fn indented() {}",
		"\tstruct Indented {",
		"// comment",
		"use std::io;",
	}
	for _, line := range noMatches {
		assert.False(t, p.MatchString(line), "Rust pattern should NOT match: %q", line)
	}
}

// ---- importLinePatterns tests ----

func TestImportLinePatterns_CoverSupportedLanguages(t *testing.T) {
	for _, lang := range SupportedLanguages {
		_, ok := importLinePatterns[lang]
		assert.True(t, ok, "importLinePatterns missing entry for supported language %q", lang)
	}
}

func TestIsImportChunk_GoImportBlock(t *testing.T) {
	chunk := Chunk{
		Content: `package main

import (
	"fmt"
	"os"
)`,
	}
	assert.True(t, IsImportChunk(chunk, LangGo))
}

func TestIsImportChunk_GoFunctionIsNotImport(t *testing.T) {
	chunk := Chunk{
		Content: `func main() {
	fmt.Println("hello")
}`,
	}
	assert.False(t, IsImportChunk(chunk, LangGo))
}

func TestIsImportChunk_TypeScriptImport(t *testing.T) {
	chunk := Chunk{
		Content: `import React from 'react';
import { useState } from 'react';`,
	}
	assert.True(t, IsImportChunk(chunk, LangTypeScript))
}

func TestIsImportChunk_UnknownLanguage(t *testing.T) {
	chunk := Chunk{Content: "some content"}
	assert.False(t, IsImportChunk(chunk, LangUnknown))
}
