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

    // "github.com/rohanair/shazam-for-code/api"
    // "github.com/rohanair/shazam-for-code/config"
    // "github.com/rohanair/shazam-for-code/internal"
)

func init() {
    // Try load .env for local dev, ignore errors in prod
    _ = godotenv.Load()
}

func main() {
    cfg := config.LoadConfig()

    // Initialize logger
    logger := internal.NewLogger(cfg.Env)
    defer logger.Sync()

    // Initialize cache (Redis optional)
    cacheClient, err := internal.NewCache(cfg)
    if err != nil {
        logger.Sugar().Warnw("redis disabled, falling back to in-memory cache", "error", err)
    }

    // Create GitHub client wrapper
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

    // graceful shutdown
    quit := make(chan os.Signal, 1)
    // signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        logger.Sugar().Fatalw("server forced to shutdown", "err", err)
    }

    logger.Sugar().Infow("server exiting")
}
