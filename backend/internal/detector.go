import (
    "math"
)
// AnalyzeSecurity scans code for deprecated APIs, insecure patterns, hardcoded secrets, and known CVEs (expanded)
func AnalyzeSecurity(code, language, framework string) []map[string]interface{} {
    var results []map[string]interface{}
    lower := strings.ToLower(code)
    // Python: deprecated APIs, insecure patterns
    if language == "Python" || strings.Contains(lower, "import yaml") {
        if strings.Contains(lower, "safe_load") {
            results = append(results, map[string]interface{}{
                "issue": "Deprecated function 'safe_load' in PyYAML < 5.1",
                "severity": "high",
                "cve": "CVE-2017-18342",
            })
        }
        if strings.Contains(lower, "pickle.load") {
            results = append(results, map[string]interface{}{
                "issue": "Use of 'pickle.load' is insecure (arbitrary code execution)",
                "severity": "high",
            })
        }
        if strings.Contains(lower, "os.system") || strings.Contains(lower, "subprocess.call") {
            results = append(results, map[string]interface{}{
                "issue": "Use of os.system/subprocess.call can be insecure",
                "severity": "medium",
            })
        }
        if strings.Contains(lower, "import requests") {
            // Simulate NVD lookup for requests
            results = append(results, map[string]interface{}{
                "issue": "requests < 2.20.0 has known CVEs",
                "severity": "medium",
                "cve": "CVE-2018-18074",
            })
        }
    }
    // Go: deprecated APIs, insecure patterns
    if language == "Go" || strings.Contains(lower, "package main") {
        if strings.Contains(lower, "ioutil.ReadAll") {
            results = append(results, map[string]interface{}{
                "issue": "Deprecated function 'ioutil.ReadAll' (use io.ReadAll)",
                "severity": "medium",
            })
        }
        if strings.Contains(lower, "os/exec") {
            results = append(results, map[string]interface{}{
                "issue": "Use of os/exec can be insecure",
                "severity": "medium",
            })
        }
    }
    // JavaScript/Node: insecure patterns
    if language == "JavaScript" || language == "TypeScript" {
        if strings.Contains(lower, "eval(") {
            results = append(results, map[string]interface{}{
                "issue": "Use of eval() is dangerous and should be avoided",
                "severity": "high",
            })
        }
        if strings.Contains(lower, "child_process.exec") {
            results = append(results, map[string]interface{}{
                "issue": "Use of child_process.exec can be insecure",
                "severity": "medium",
            })
        }
    }
    // PHP: insecure patterns
    if language == "PHP" {
        if strings.Contains(lower, "eval(") {
            results = append(results, map[string]interface{}{
                "issue": "Use of eval() is dangerous and should be avoided",
                "severity": "high",
            })
        }
        if strings.Contains(lower, "mysql_query") {
            results = append(results, map[string]interface{}{
                "issue": "mysql_query is deprecated, use PDO or MySQLi",
                "severity": "medium",
            })
        }
    }
    // Hardcoded secrets/tokens (generic)
    if strings.Contains(lower, "password=") || strings.Contains(lower, "api_key=") || strings.Contains(lower, "secret=") {
        results = append(results, map[string]interface{}{
            "issue": "Possible hardcoded secret or credential",
            "severity": "high",
        })
    }
    // TODO: Add more static checks for other languages/frameworks
    // TODO: NVD API lookup for imported libraries (stub)
    return results
}

// GuessPurpose uses ML/keyword/context rules to guess the purpose of the snippet (expanded)
func GuessPurpose(code, language, framework string) string {
    lower := strings.ToLower(code)
    // Authentication/authorization
    if strings.Contains(lower, "auth") && (strings.Contains(lower, "token") || strings.Contains(lower, "login") || strings.Contains(lower, "jwt")) {
        return "Authentication or authorization logic."
    }
    // API handler
    if strings.Contains(lower, "handler") && (strings.Contains(lower, "http") || strings.Contains(lower, "route") || strings.Contains(lower, "endpoint")) {
        return "API route handler for a REST service."
    }
    // Database
    if strings.Contains(lower, "db") || strings.Contains(lower, "database") || strings.Contains(lower, "sql") || strings.Contains(lower, "orm") {
        return "Database access or query logic."
    }
    // CLI tool
    if strings.Contains(lower, "flag.") || strings.Contains(lower, "argparse") || strings.Contains(lower, "sys.argv") {
        return "Command-line tool or script."
    }
    // Data processing
    if strings.Contains(lower, "pandas") || strings.Contains(lower, "numpy") || strings.Contains(lower, "dataframe") {
        return "Data processing or analysis."
    }
    // Web scraping
    if strings.Contains(lower, "beautifulsoup") || strings.Contains(lower, "requests.get") {
        return "Web scraping or HTTP client."
    }
    // Test
    if strings.Contains(lower, "test") && (strings.Contains(lower, "assert") || strings.Contains(lower, "expect") || strings.Contains(lower, "unittest")) {
        return "Unit test or test case."
    }
    // Config
    if strings.Contains(lower, "config") || strings.Contains(lower, "yaml") || strings.Contains(lower, "ini") {
        return "Configuration file or loader."
    }
    // Logging
    if strings.Contains(lower, "log.") || strings.Contains(lower, "logger") {
        return "Logging or monitoring logic."
    }
    // Entry point
    if strings.Contains(lower, "main(") && (strings.Contains(lower, "print") || strings.Contains(lower, "fmt.") || strings.Contains(lower, "console.log")) {
        return "Program entry point."
    }
    // React
    if framework == "React" && strings.Contains(lower, "useeffect") {
        return "React component with side effects (useEffect)."
    }
    // Flask/Django
    if framework == "Flask" || framework == "Django" {
        return "Web API endpoint or view handler."
    }
    // Gin
    if framework == "Gin (Go)" {
        return "Go REST API handler."
    }
    // Data science
    if strings.Contains(lower, "sklearn") || strings.Contains(lower, "tensorflow") || strings.Contains(lower, "keras") {
        return "Machine learning or data science code."
    }
    // Fallback
    return "General code snippet."
}

// AnalyzeComplexityAndStyle returns cyclomatic complexity and style score (expanded for more languages and metrics)
func AnalyzeComplexityAndStyle(code, language string) (int, float64) {
    lower := strings.ToLower(code)
    lines := strings.Split(code, "\n")
    longLines := 0
    commentLines := 0
    funcCount := 0
    for _, l := range lines {
        if len(l) > 100 {
            longLines++
        }
        if strings.HasPrefix(strings.TrimSpace(l), "//") || strings.HasPrefix(strings.TrimSpace(l), "#") {
            commentLines++
        }
        if strings.Contains(l, "func ") || strings.Contains(l, "def ") || strings.Contains(l, "function ") {
            funcCount++
        }
    }
    // Go: count if, for, case, &&, ||, else if, range, switch
    if language == "Go" {
        complexity := 1
        keywords := []string{"if ", "for ", "case ", "&&", "||", "else if", "range ", "switch "}
        for _, k := range keywords {
            complexity += strings.Count(lower, k)
        }
        // Style: penalize long lines, low comment density, few functions
        style := 1.0 - float64(longLines)/float64(len(lines)+1)
        if funcCount < 1 { style -= 0.1 }
        if float64(commentLines)/float64(len(lines)+1) < 0.05 { style -= 0.1 }
        if style < 0.5 { style = 0.5 }
        return complexity, style
    }
    // Python: count if, for, while, elif, and, or
    if language == "Python" {
        complexity := 1
        keywords := []string{"if ", "for ", "while ", "elif ", " and ", " or "}
        for _, k := range keywords {
            complexity += strings.Count(lower, k)
        }
        style := 1.0 - float64(longLines)/float64(len(lines)+1)
        if funcCount < 1 { style -= 0.1 }
        if float64(commentLines)/float64(len(lines)+1) < 0.05 { style -= 0.1 }
        if style < 0.5 { style = 0.5 }
        return complexity, style
    }
    // JavaScript/TypeScript: count if, for, while, function, =>
    if language == "JavaScript" || language == "TypeScript" {
        complexity := 1
        keywords := []string{"if ", "for ", "while ", "function ", "=>"}
        for _, k := range keywords {
            complexity += strings.Count(lower, k)
        }
        style := 1.0 - float64(longLines)/float64(len(lines)+1)
        if funcCount < 1 { style -= 0.1 }
        if float64(commentLines)/float64(len(lines)+1) < 0.05 { style -= 0.1 }
        if style < 0.5 { style = 0.5 }
        return complexity, style
    }
    // Fallback for other languages
    return 1, 0.9
}
package internal

import (
    "bytes"
    "strings"
    "log"
)

enry "github.com/src-d/enry/v2"


// DetectLanguageWithConfidence returns language and a confidence score [0,1].
// Uses Enry, filename hints, n-gram/keyword rules, and snippet heuristics.
func DetectLanguageWithConfidence(filename string, content []byte) (string, float64) {
    txt := strings.ToLower(string(bytes.TrimSpace(content)))
    ext := getFileExtension(filename)

    // 1. Enry (static analysis)
    lang := enry.GetLanguage(filename, content)
    if lang != "" && lang != "Unknown" {
        return lang, 0.98
    }

    // 2. Filename/extension hints
    if ext != "" {
        if l, conf := languageFromExtension(ext); l != "Unknown" {
            return l, conf
        }
    }

    // 3. N-gram/keyword rules for many languages
    if l, conf := ngramLanguageClassifier(txt); l != "Unknown" {
        return l, conf
    }

    // 4. Heuristics for markup, config, shell, SQL, etc.
    if strings.HasPrefix(txt, "<") {
        return "HTML", 0.7
    }
    if strings.HasPrefix(txt, "#!/bin/bash") || strings.Contains(txt, "echo ") {
        return "Shell", 0.7
    }
    if strings.Contains(txt, "SELECT ") && strings.Contains(txt, " FROM ") {
        return "SQL", 0.7
    }
    if strings.HasPrefix(txt, "{") && strings.HasSuffix(txt, "}") && strings.Contains(txt, ":") {
        return "JSON", 0.7
    }
    if strings.HasPrefix(txt, "---") || strings.Contains(txt, ": ") {
        return "YAML", 0.7
    }

    // 5. Fallback
    return "Unknown", 0.2
}

// getFileExtension extracts the file extension from a filename
func getFileExtension(filename string) string {
    idx := strings.LastIndex(filename, ".")
    if idx == -1 || idx == len(filename)-1 {
        return ""
    }
    return strings.ToLower(filename[idx+1:])
}

// --- Add more language extension mappings for broad coverage ---

// languageFromExtension maps file extensions to languages and confidence
func languageFromExtension(ext string) (string, float64) {
    switch ext {
    case "go": return "Go", 0.95
    case "py": return "Python", 0.95
    case "js": return "JavaScript", 0.95
    case "ts": return "TypeScript", 0.95
    case "rb": return "Ruby", 0.95
    case "rs": return "Rust", 0.95
    case "java": return "Java", 0.95
    case "c": return "C", 0.95
    case "cpp", "cc", "cxx": return "C++", 0.95
    case "cs": return "C#", 0.95
    case "php": return "PHP", 0.95
    case "html", "htm": return "HTML", 0.95
    case "css": return "CSS", 0.95
    case "json": return "JSON", 0.95
    case "yaml", "yml": return "YAML", 0.95
    case "sh", "bash": return "Shell", 0.95
    case "sql": return "SQL", 0.95
    case "swift": return "Swift", 0.95
    case "kt": return "Kotlin", 0.95
    case "scala": return "Scala", 0.95
    case "dart": return "Dart", 0.95
    case "lua": return "Lua", 0.95
    case "m": return "Objective-C", 0.95
    case "r": return "R", 0.95
    case "pl": return "Perl", 0.95
    case "hs": return "Haskell", 0.95
    case "erl": return "Erlang", 0.95
    case "ex", "exs": return "Elixir", 0.95
    case "jl": return "Julia", 0.95
    case "groovy": return "Groovy", 0.95
    case "f90", "f95", "f03": return "Fortran", 0.95
    case "v": return "Verilog", 0.95
    case "vhd", "vhdl": return "VHDL", 0.95
    case "tsx": return "TypeScript", 0.93
    case "jsx": return "JavaScript", 0.93
    // Additional languages
    case "ml", "mli": return "OCaml", 0.95
    case "fs", "fsi": return "F#", 0.95
    case "elm": return "Elm", 0.95
    case "nim": return "Nim", 0.95
    case "cr": return "Crystal", 0.95
    case "re": return "Reason", 0.95
    case "vala": return "Vala", 0.95
    case "zig": return "Zig", 0.95
    case "sol": return "Solidity", 0.95
    case "graphql": return "GraphQL", 0.95
    case "dockerfile": return "Dockerfile", 0.95
    case "makefile": return "Makefile", 0.95
    case "cmake": return "CMake", 0.95
    case "ini": return "INI", 0.95
    case "toml": return "TOML", 0.95
    case "proto": return "Protobuf", 0.95
    }
    return "Unknown", 0.2
}

// ngramLanguageClassifier is a modular ML-inspired classifier for ambiguous code
// Add more rules for more languages and idioms as needed
func ngramLanguageClassifier(code string) (string, float64) {
    // Python
    if strings.Contains(code, "def ") || strings.Contains(code, ":\n") || strings.Contains(code, "import sys") {
        return "Python", 0.8
    }
    // Go
    if strings.Contains(code, "func ") || strings.Contains(code, "package main") {
        return "Go", 0.8
    }
    // JavaScript/TypeScript
    if strings.Contains(code, "console.log") || strings.Contains(code, "function ") || strings.Contains(code, "let ") || strings.Contains(code, "const ") {
        if strings.Contains(code, ":") && strings.Contains(code, ";") {
            return "TypeScript", 0.8
        }
        return "JavaScript", 0.8
    }
    // Java
    if strings.Contains(code, "public static void main") || strings.Contains(code, "System.out.println") {
        return "Java", 0.8
    }
    // C/C++
    if strings.Contains(code, "#include <stdio.h>") || strings.Contains(code, "#include <iostream>") {
        if strings.Contains(code, "std::") {
            return "C++", 0.8
        }
        return "C", 0.8
    }
    // Rust
    if strings.Contains(code, "fn main()") || strings.Contains(code, "println!") {
        return "Rust", 0.8
    }
    // Ruby
    if strings.Contains(code, "def ") && strings.Contains(code, "end") {
        return "Ruby", 0.8
    }
    // PHP
    if strings.Contains(code, "<?php") {
        return "PHP", 0.8
    }
    // C#
    if strings.Contains(code, "using System;") || strings.Contains(code, "namespace ") {
        return "C#", 0.8
    }
    // Swift
    if strings.Contains(code, "import Foundation") || strings.Contains(code, "func main(") {
        return "Swift", 0.8
    }
    // Kotlin
    if strings.Contains(code, "fun main(") || strings.Contains(code, "val ") {
        return "Kotlin", 0.8
    }
    // Scala
    if strings.Contains(code, "object ") && strings.Contains(code, "extends App") {
        return "Scala", 0.8
    }
    // Dart
    if strings.Contains(code, "void main()") && strings.Contains(code, "import 'dart:") {
        return "Dart", 0.8
    }
    // OCaml
    if strings.Contains(code, "let rec ") || strings.Contains(code, ";;") {
        return "OCaml", 0.8
    }
    // F#
    if strings.Contains(code, "let mutable ") || strings.Contains(code, "open System") {
        return "F#", 0.8
    }
    // Elm
    if strings.Contains(code, "module Main exposing") || strings.Contains(code, "Html.text") {
        return "Elm", 0.8
    }
    // Nim
    if strings.Contains(code, "proc ") || strings.Contains(code, "echo ") {
        return "Nim", 0.8
    }
    // Crystal
    if strings.Contains(code, "def initialize") || strings.Contains(code, "Crystal::") {
        return "Crystal", 0.8
    }
    // Reason
    if strings.Contains(code, "let make = (_self) =>") {
        return "Reason", 0.8
    }
    // Vala
    if strings.Contains(code, "public static int main") {
        return "Vala", 0.8
    }
    // Zig
    if strings.Contains(code, "pub fn main()") {
        return "Zig", 0.8
    }
    // Solidity
    if strings.Contains(code, "pragma solidity") {
        return "Solidity", 0.8
    }
    // GraphQL
    if strings.Contains(code, "type Query {") {
        return "GraphQL", 0.8
    }
    // Dockerfile
    if strings.HasPrefix(code, "from ") && strings.Contains(code, "docker") {
        return "Dockerfile", 0.8
    }
    // Makefile
    if strings.HasPrefix(code, "all:") || strings.Contains(code, ".PHONY") {
        return "Makefile", 0.8
    }
    // CMake
    if strings.Contains(code, "cmake_minimum_required") {
        return "CMake", 0.8
    }
    // INI
    if strings.Contains(code, "[section]") || strings.Contains(code, "=") {
        return "INI", 0.8
    }
    // TOML
    if strings.Contains(code, "[package]") && strings.Contains(code, "version = ") {
        return "TOML", 0.8
    }
    // Protobuf
    if strings.Contains(code, "syntax = \"proto3\"") {
        return "Protobuf", 0.8
    }
    // YAML
    if strings.HasPrefix(code, "---") || strings.Contains(code, ": ") {
        return "YAML", 0.8
    }
    // JSON
    if strings.HasPrefix(code, "{") && strings.HasSuffix(code, "}") && strings.Contains(code, ":") {
        return "JSON", 0.8
    }
    // HTML
    if strings.HasPrefix(code, "<") && strings.Contains(code, ">") {
        return "HTML", 0.8
    }
    return "Unknown", 0.2
}

// DetectFrameworkWithConfidence returns framework and a confidence score [0,1].
// Uses modular static rules, idiom/ML matching, and weighted scoring.
func DetectFrameworkWithConfidence(code string) (string, float64) {
    lower := strings.ToLower(code)
    var candidates []struct{ name string; score float64 }

    // --- Go frameworks ---
    if strings.Contains(lower, "package main") {
        if strings.Contains(lower, "github.com/gin-gonic/gin") {
            candidates = append(candidates, struct{ name string; score float64 }{"Gin (Go)", 0.98})
        }
        if strings.Contains(lower, "github.com/labstack/echo") {
            candidates = append(candidates, struct{ name string; score float64 }{"Echo (Go)", 0.95})
        }
        if strings.Contains(lower, "github.com/gofiber/fiber") {
            candidates = append(candidates, struct{ name string; score float64 }{"Fiber (Go)", 0.93})
        }
        if strings.Contains(lower, "github.com/astaxie/beego") {
            candidates = append(candidates, struct{ name string; score float64 }{"Beego (Go)", 0.92})
        }
    }

    // --- Python frameworks ---
    if strings.Contains(lower, "import flask") || strings.Contains(lower, "from flask") {
        candidates = append(candidates, struct{ name string; score float64 }{"Flask", 0.98})
    }
    if strings.Contains(lower, "import django") || strings.Contains(lower, "from django") {
        candidates = append(candidates, struct{ name string; score float64 }{"Django", 0.98})
    }
    if strings.Contains(lower, "import fastapi") {
        candidates = append(candidates, struct{ name string; score float64 }{"FastAPI", 0.97})
    }
    if strings.Contains(lower, "import tornado") {
        candidates = append(candidates, struct{ name string; score float64 }{"Tornado", 0.95})
    }
    if strings.Contains(lower, "import pyramid") {
        candidates = append(candidates, struct{ name string; score float64 }{"Pyramid", 0.93})
    }

    // --- JavaScript/TypeScript frameworks ---
    if strings.Contains(lower, "import react") || strings.Contains(lower, "useeffect(") || strings.Contains(lower, "usestate(") {
        candidates = append(candidates, struct{ name string; score float64 }{"React", 0.97})
    }
    if strings.Contains(lower, "import vue") || strings.Contains(lower, "new vue(") {
        candidates = append(candidates, struct{ name string; score float64 }{"Vue", 0.95})
    }
    if strings.Contains(lower, "@component(") || strings.Contains(lower, "import { component }") {
        candidates = append(candidates, struct{ name string; score float64 }{"Angular", 0.95})
    }
    if strings.Contains(lower, "require('express')") || strings.Contains(lower, "const app = express()") {
        candidates = append(candidates, struct{ name string; score float64 }{"Express", 0.95})
    }
    if strings.Contains(lower, "import next") || strings.Contains(lower, "getstaticprops") {
        candidates = append(candidates, struct{ name string; score float64 }{"Next.js", 0.93})
    }
    if strings.Contains(lower, "@nestjs/") {
        candidates = append(candidates, struct{ name string; score float64 }{"NestJS", 0.93})
    }

    // --- Java frameworks ---
    if strings.Contains(lower, "org.springframework") || strings.Contains(lower, "@controller") {
        candidates = append(candidates, struct{ name string; score float64 }{"Spring Boot", 0.97})
    }
    if strings.Contains(lower, "micronaut") {
        candidates = append(candidates, struct{ name string; score float64 }{"Micronaut", 0.93})
    }
    if strings.Contains(lower, "quarkus") {
        candidates = append(candidates, struct{ name string; score float64 }{"Quarkus", 0.93})
    }

    // --- Ruby frameworks ---
    if strings.Contains(lower, "rails") || strings.Contains(lower, "activerecord") {
        candidates = append(candidates, struct{ name string; score float64 }{"Rails", 0.97})
    }
    if strings.Contains(lower, "sinatra") {
        candidates = append(candidates, struct{ name string; score float64 }{"Sinatra", 0.93})
    }

    // --- PHP frameworks ---
    if strings.Contains(lower, "laravel") {
        candidates = append(candidates, struct{ name string; score float64 }{"Laravel", 0.97})
    }
    if strings.Contains(lower, "symfony") {
        candidates = append(candidates, struct{ name string; score float64 }{"Symfony", 0.95})
    }
    if strings.Contains(lower, "codeigniter") {
        candidates = append(candidates, struct{ name string; score float64 }{"CodeIgniter", 0.93})
    }

    // --- .NET/C# frameworks ---
    if strings.Contains(lower, "asp.net") || strings.Contains(lower, "microsoft.aspnetcore") {
        candidates = append(candidates, struct{ name string; score float64 }{"ASP.NET", 0.97})
    }

    // --- Rust frameworks ---
    if strings.Contains(lower, "rocket::") {
        candidates = append(candidates, struct{ name string; score float64 }{"Rocket (Rust)", 0.95})
    }
    if strings.Contains(lower, "actix::") {
        candidates = append(candidates, struct{ name string; score float64 }{"Actix (Rust)", 0.93})
    }

    // --- ML-based idiom matching (fallback, lower confidence) ---
    if fw, score, err := DetectFrameworkML(code); err == nil && fw != "Unknown" {
        conf := 0.5 + 0.5*sigmoid(score/10)
        candidates = append(candidates, struct{ name string; score float64 }{fw, conf * 0.8})
    }

    // --- Weighted scoring and conflict resolution ---
    if len(candidates) == 0 {
        return "Unknown", 0.2
    }
    // Penalize if multiple high-confidence conflicting frameworks detected
    if len(candidates) > 1 && candidates[0].score > 0.9 && candidates[1].score > 0.9 {
        return candidates[0].name, 0.6 // lower confidence
    }
    // Return highest scoring candidate
    best := candidates[0]
    for _, c := range candidates {
        if c.score > best.score {
            best = c
        }
    }
    return best.name, best.score
}

func sigmoid(x float64) float64 {
    return 1.0 / (1.0 + math.Exp(-x))
}

func DetectFramework(code string) string {
    if fw, score, err := DetectFrameworkML(code); err == nil {
        if score > -1000 {
            log.Printf("ML framework detection: %s (score %.2f)\n", fw, score)
            return fw
        }
    } else {
        log.Println("ML detection error:", err)
    }
    lowerCode := strings.ToLower(code)
    if strings.Contains(lowerCode, "import react") || strings.Contains(lowerCode, "useeffect(") || strings.Contains(lowerCode, "usestate(") {
        return "React"
    }
    if strings.Contains(lowerCode, "github.com/gin-gonic/gin") || strings.Contains(lowerCode, "gin.") {
        return "Gin (Go)"
    }
    if strings.Contains(lowerCode, "django") {
        return "Django"
    }
    if strings.Contains(lowerCode, "from flask") || strings.Contains(lowerCode, "flask import") {
        return "Flask"
    }
    return "Unknown"
}
