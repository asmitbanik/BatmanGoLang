package internal

import (
    "bytes"
    "strings"

    enry "github.com/src-d/enry/v2"
)

// DetectLanguage guesses language using enry
func DetectLanguage(filename string, content []byte) string {
    // enry expects filename and content
    lang := enry.GetLanguage(filename, content)
    if lang == "" || lang == "Unknown" {
        // fallback to content-based heuristics
        txt := strings.ToLower(string(bytes.TrimSpace(content)))
        if strings.HasPrefix(txt, "<") {
            return "HTML"
        }
        if strings.Contains(txt, "func main") {
            return "Go"
        }
    }
    return lang
}

// DetectFramework looks for heuristics in the source
func DetectFramework(code string, lang string) string {
    lower := strings.ToLower(code)
    switch lang {
    case "JavaScript", "TypeScript":
        if strings.Contains(lower, "react") || strings.Contains(lower, "from 'react'") || strings.Contains(lower, "from \"react\"") {
            return "React"
        }
        if strings.Contains(lower, "@angular") || strings.Contains(lower, "angular.module") {
            return "Angular"
        }
    case "Python":
        if strings.Contains(lower, "django") {
            return "Django"
        }
        if strings.Contains(lower, "flask") {
            return "Flask"
        }
    case "Go":
        if strings.Contains(lower, "github.com/gin-gonic/gin") || strings.Contains(lower, "gin") && strings.Contains(lower, "router") {
            return "Gin"
        }
    case "Java":
        if strings.Contains(lower, "springframework") {
            return "Spring"
        }
    }
    return "Unknown"
}
