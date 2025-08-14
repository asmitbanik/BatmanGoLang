package tests

import (
    "testing"
    // "github.com/rohanair/shazam-for-code/internal"
)

func TestDetectLanguageGo(t *testing.T) {
    code := []byte("package main\nfunc main(){}")
    lang := internal.DetectLanguage("main.go", code)
    if lang != "Go" {
        t.Fatalf("expected Go got %s", lang)
    }
}
