# Security & Production Checklist

- Use GitHub App or fine-grained PATs with minimal scopes (code search scope).
- Store tokens in secrets manager (AWS Secrets Manager / Vault / Kubernetes secrets).
- Enable Transport Layer Security (TLS) in front of the service (ingress/ALB).
- Use rate limiting and request validation.
- Sanitize any user-provided strings sent to external APIs.
- Limit GitHub request size and escape queries to avoid injection-like issues.
- Configure logging to avoid writing secrets to logs.
- Monitor rate-limit headers returned by GitHub and backoff.
- Use health checks and liveness/readiness probes in k8s.
- Use automated vulnerability scanning for images and dependencies.
