package api

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    // "github.com/rohanair/shazam-for-code/internal"
)

type analyzeRequest struct {
    Filename string `json:"filename" binding:"required"`
    Code     string `json:"code" binding:"required"`
}

type analyzeResponse struct {
    Language string   `json:"language"`
    Framework string  `json:"framework"`
    SimilarRepos []string `json:"similar_repos"`
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

        // Detect language & framework
        lang := internal.DetectLanguage(req.Filename, []byte(req.Code))
        framework := internal.DetectFramework(req.Code, lang)

        // Call GitHub search with timeout
        ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
        defer cancel()
        repos, err := gh.SearchCode(ctx, req.Code)
        if err != nil {
            logger.Sugar().Warnw("github search failed", "err", err)
            // continue and return other info
        }

        res := analyzeResponse{
            Language: lang,
            Framework: framework,
            SimilarRepos: repos,
        }
        c.JSON(http.StatusOK, res)
    }
}
