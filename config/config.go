package config

import (
    "os"
    "strconv"
)

type Config struct {
    GitHubToken       string
    Port              int
    RedisURL          string
    Env               string
    RateLimitReqPerMin int
}

func LoadConfig() *Config {
    port := 8080
    if p := os.Getenv("PORT"); p != "" {
        if v, err := strconv.Atoi(p); err == nil {
            port = v
        }
    }

    rateLimit := 120 // default 120 reqs per min
    if r := os.Getenv("RATE_LIMIT_PER_MIN"); r != "" {
        if v, err := strconv.Atoi(r); err == nil {
            rateLimit = v
        }
    }

    return &Config{
        GitHubToken:       os.Getenv("GITHUB_TOKEN"),
        Port:              port,
        RedisURL:          os.Getenv("REDIS_URL"),
        Env:               os.Getenv("APP_ENV"),
        RateLimitReqPerMin: rateLimit,
    }
}
