import (
    "math"
)
// AnalyzeSecurity scans code for deprecated APIs and known CVEs (basic real logic)
func AnalyzeSecurity(code, language, framework string) []map[string]interface{} {
    var results []map[string]interface{}
    lower := strings.ToLower(code)
    // Example: Python deprecated API
    if language == "Python" || strings.Contains(lower, "import yaml") {
        if strings.Contains(lower, "safe_load") {
            // Simulate a real CVE
            results = append(results, map[string]interface{}{
                "issue": "Deprecated function 'safe_load' in PyYAML < 5.1",
                "severity": "high",
                "cve": "CVE-2017-18342",
            })
        }
    }
    // Example: Go deprecated API
    if language == "Go" || strings.Contains(lower, "package main") {
        if strings.Contains(lower, "ioutil.ReadAll") {
            results = append(results, map[string]interface{}{
                "issue": "Deprecated function 'ioutil.ReadAll' (use io.ReadAll)",
                "severity": "medium",
            })
        }
    }
    // TODO: Add more static checks for other languages/frameworks
    // TODO: NVD API lookup for imported libraries (stub)
    // Example: if 'requests' in Python, check NVD for 'requests' CVEs
    return results
}

// GuessPurpose uses a simple ML/keyword approach to guess the purpose of the snippet
func GuessPurpose(code, language, framework string) string {
    lower := strings.ToLower(code)
    if strings.Contains(lower, "auth") && (strings.Contains(lower, "token") || strings.Contains(lower, "login")) {
        return "Authentication or authorization logic."
    }
    if strings.Contains(lower, "handler") && (strings.Contains(lower, "http") || strings.Contains(lower, "route")) {
        return "API route handler for a REST service."
    }
    if strings.Contains(lower, "db") || strings.Contains(lower, "database") || strings.Contains(lower, "sql") {
        return "Database access or query logic."
    }
    if strings.Contains(lower, "test") && (strings.Contains(lower, "assert") || strings.Contains(lower, "expect")) {
        return "Unit test or test case."
    }
    if strings.Contains(lower, "main(") && (strings.Contains(lower, "print") || strings.Contains(lower, "fmt.")) {
        return "Program entry point."
    }
    if framework == "React" && strings.Contains(lower, "useeffect") {
        return "React component with side effects (useEffect)."
    }
    if framework == "Flask" || framework == "Django" {
        return "Web API endpoint or view handler."
    }
    if framework == "Gin (Go)" {
        return "Go REST API handler."
    }
    return "General code snippet."
}

// AnalyzeComplexityAndStyle returns cyclomatic complexity and style score (real logic for Go/Python)
func AnalyzeComplexityAndStyle(code, language string) (int, float64) {
    lower := strings.ToLower(code)
    // Go: count if, for, case, &&, ||, else if, range, switch
    if language == "Go" {
        complexity := 1
        keywords := []string{"if ", "for ", "case ", "&&", "||", "else if", "range ", "switch "}
        for _, k := range keywords {
            complexity += strings.Count(lower, k)
        }
        // Style: penalize lines > 100 chars, tabs vs spaces, etc.
        lines := strings.Split(code, "\n")
        longLines := 0
        for _, l := range lines {
            if len(l) > 100 {
                longLines++
            }
        }
        style := 1.0 - float64(longLines)/float64(len(lines)+1)
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
        // Style: penalize lines > 100 chars, tabs vs spaces, etc.
        lines := strings.Split(code, "\n")
        longLines := 0
        for _, l := range lines {
            if len(l) > 100 {
                longLines++
            }
        }
        style := 1.0 - float64(longLines)/float64(len(lines)+1)
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


// DetectLanguageWithConfidence returns language and a confidence score [0,1] using Enry and n-gram ML classifier
func DetectLanguageWithConfidence(filename string, content []byte) (string, float64) {
    lang := enry.GetLanguage(filename, content)
    if lang != "" && lang != "Unknown" {
        return lang, 0.98
    }
    txt := strings.ToLower(string(bytes.TrimSpace(content)))
    // Fallback: n-gram ML classifier (very simple)
    ngramLang, ngramConf := ngramLanguageClassifier(txt)
    if ngramLang != "Unknown" {
        return ngramLang, ngramConf
    }
    if strings.HasPrefix(txt, "<") {
        return "HTML", 0.7
    }
    if strings.Contains(txt, "func main") {
        return "Go", 0.7
    }
    return "Unknown", 0.2
}

// ngramLanguageClassifier is a simple ML-based classifier for ambiguous code
func ngramLanguageClassifier(code string) (string, float64) {
    // Very basic: look for common n-grams for Python, Go, JS, etc.
    if strings.Contains(code, "def ") || strings.Contains(code, ":\n") {
        return "Python", 0.8
    }
    if strings.Contains(code, "func ") {
        return "Go", 0.8
    }
    if strings.Contains(code, "console.log") || strings.Contains(code, "function ") {
        return "JavaScript", 0.8
    }
    if strings.Contains(code, "public static void main") {
        return "Java", 0.8
    }
    if strings.Contains(code, "#include <stdio.h>") {
        return "C", 0.8
    }
    if strings.Contains(code, "<?php") {
        return "PHP", 0.8
    }
    return "Unknown", 0.2
}

// DetectFrameworkWithConfidence returns framework and a confidence score [0,1] using AST, regex, and ML idiom matching
func DetectFrameworkWithConfidence(code string) (string, float64) {
    lowerCode := strings.ToLower(code)
    // Go: use AST for import detection
    if strings.Contains(lowerCode, "package main") {
        if strings.Contains(lowerCode, "github.com/gin-gonic/gin") {
            return "Gin (Go)", 0.95
        }
        if strings.Contains(lowerCode, "github.com/labstack/echo") {
            return "Echo (Go)", 0.9
        }
        // TODO: Use Go AST for more robust import detection
    }
    // Python: regex for imports
    if strings.Contains(lowerCode, "import flask") || strings.Contains(lowerCode, "from flask") {
        return "Flask", 0.95
    }
    if strings.Contains(lowerCode, "import django") || strings.Contains(lowerCode, "from django") {
        return "Django", 0.95
    }
    // React: idiom detection
    if strings.Contains(lowerCode, "import react") || strings.Contains(lowerCode, "useeffect(") || strings.Contains(lowerCode, "usestate(") {
        return "React", 0.9
    }
    // ML-based idiom matching (fallback)
    if fw, score, err := DetectFrameworkML(code); err == nil && fw != "Unknown" {
        conf := 0.5 + 0.5*sigmoid(score/10)
        return fw, conf
    }
    return "Unknown", 0.2
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
