package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"

    "new implementation/api"
    "new implementation/config"
    "new implementation/internal"
)

func init() {
    _ = godotenv.Load()
}

func main() {
    cfg := config.LoadConfig()
    internal.InitOrLoadModel()
    logger := internal.NewLogger(cfg.Env)
    defer logger.Sync()
    cacheClient, err := internal.NewCache(cfg)
    if err != nil {
        logger.Sugar().Warnw("redis disabled, falling back to in-memory cache", "error", err)
    }
    gh := internal.NewGitHubClient(cfg.GitHubToken, cacheClient, logger)
    r := gin.New()
    r.Use(internal.GinZap(logger))
    r.Use(gin.Recovery())
    r.Use(internal.RateLimitMiddleware(cfg.RateLimitReqPerMin))
    api.RegisterRoutes(r, gh)
    srv := &http.Server{
        Addr:    fmt.Sprintf(":%d", cfg.Port),
        Handler: r,
    }
    go func() {
        logger.Sugar().Infow("starting server", "addr", srv.Addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Sugar().Fatalw("server error", "err", err)
        }
    }()
    quit := make(chan os.Signal, 1)
    <-quit
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        logger.Sugar().Fatalw("server forced to shutdown", "err", err)
    }
    logger.Sugar().Infow("server exiting")
}
