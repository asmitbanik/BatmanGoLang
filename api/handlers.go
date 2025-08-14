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
        lang := internal.DetectLanguage(req.Filename, []byte(req.Code))
        framework, mlScore, _ := internal.DetectFrameworkML(req.Code)
        if framework == "" || framework == "Unknown" {
            framework = internal.DetectFramework(req.Code)
        }
        ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
        defer cancel()
        repos, err := gh.SearchCode(ctx, req.Code)
        if err != nil {
            logger.Sugar().Warnw("github search failed", "err", err)
        }
        c.JSON(http.StatusOK, gin.H{
            "language":      lang,
            "framework":     framework,
            "ml_score":      mlScore,
            "similar_repos": repos,
        })
    }
}
