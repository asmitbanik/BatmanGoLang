package internal

import (
    "bytes"
    "strings"
    "log"
)

enry "github.com/src-d/enry/v2"

func DetectLanguage(filename string, content []byte) string {
    lang := enry.GetLanguage(filename, content)
    if lang == "" || lang == "Unknown" {
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
