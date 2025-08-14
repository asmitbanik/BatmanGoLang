package api

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "new implementation/internal"
)


type analyzeRequest struct {
    Filename string `json:"filename" binding:"required"`
    Code     string `json:"code" binding:"required"`
    Language string `json:"language"` // optional user hint
}

func RegisterRoutes(r *gin.Engine, gh *internal.GitHubClient) {
    r.POST("/analyze", makeAnalyzeHandler(gh, zap.NewExample()))
}

func makeAnalyzeHandler(gh *internal.GitHubClient, logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req analyzeRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // --- Modular pipeline ---
        // 1. Language detection (static + ML)
        lang, langConf := internal.DetectLanguageWithConfidence(req.Filename, []byte(req.Code))
        // TODO: If ambiguous, call ML-based n-gram classifier and merge results

        // 2. Framework/library detection (static + ML)
        framework, fwConf := internal.DetectFrameworkWithConfidence(req.Code)
        // TODO: AST scan, ML idiom detection, weighted merge

        // 3. Security fingerprint (CVE/deprecated API scan)
        security := internal.AnalyzeSecurity(req.Code, lang, framework)

        // 4. Purpose guess (LLM/ML stub)
        purpose := internal.GuessPurpose(req.Code, lang, framework)

        // 5. Complexity & style metrics
        complexity, styleScore := internal.AnalyzeComplexityAndStyle(req.Code, lang)

        // 6. Similar code search (with caching, filtering)
        ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
        defer cancel()
        repos, err := gh.SearchCodeSmart(ctx, req.Code, lang, framework)
        if err != nil {
            logger.Sugar().Warnw("github search failed", "err", err)
        }

        // --- Response ---
        c.JSON(http.StatusOK, gin.H{
            "language": gin.H{
                "name": lang,
                "confidence": langConf,
            },
            "framework": gin.H{
                "name": framework,
                "confidence": fwConf,
            },
            "security": security,
            "purpose_guess": purpose,
            "complexity": gin.H{
                "cyclomatic": complexity,
                "style_score": styleScore,
            },
            "similar_repos": repos,
        })
    }
}
