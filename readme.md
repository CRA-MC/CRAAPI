# CRAAPI — Yggdrasil-Compatible Auth API (Go + fasthttp)

CRAAPI is a lightweight authentication API service built in **Go** using **fasthttp**.
It provides **Yggdrasil-compatible** endpoints for the Minecraft ecosystem and can be extended for
custom project authorization workflows.

This repo focuses on **reliability-oriented service design**:
- Non-blocking request handling
- Failure isolation for SMTP email delivery
- Concurrency-safe logging (single-writer)

---

## Highlights

- **High-performance HTTP stack**: `fasthttp` + router-based routing
- **CLI + config**: `cobra` for commands, `viper` for config loading (TOML)
- **SMTP worker pool**: buffered channel + fixed goroutines to decouple email I/O from core auth requests
- **Concurrency-safe logging**: channel-based logger + single goroutine writer; log files named by date
- **MongoDB integration** (auth/account-related persistence)

---

## Quick Start

### Requirements
- Go (recent version recommended)
- MongoDB (if enabled by your config)
- An SMTP server (if you enable email verification / delivery)

### Build
```bash
go mod tidy
go build -o craapi .
````

### Run

CRAAPI loads configuration from a TOML file (see the example below).

```bash
./craapi --config ./craapi.config.example
```

> Tip: If you use a different config path/name, update the flag accordingly.

---

## Configuration

CRAAPI uses a TOML config (via `viper`). A minimal example:

```toml
[server]
port = 30021

[mongo]
dsn = "mongodb://127.0.0.1:27017"
db  = "craapi"

[yggdrasil]
# Your deployment-specific options
# (keep consistent with your routing and frontend expectations)

[panel]
# Optional admin/panel config

[smtp]
host = "smtp.example.com"
port = 587
username = "user@example.com"
password = "your-password"
from = "noreply@example.com"
connections = 4
```

### Notes on Email Delivery

CRAAPI sends emails through a **worker pool** pattern:

* A buffered channel acts as a queue
* A fixed number of goroutines consume messages and send via SMTP
* This prevents SMTP latency/outages from blocking auth endpoints

If an email fails to send, CRAAPI records the failure (logging) and returns a failure response to the client,
so the client can decide whether/when to retry.

---

## API Overview

CRAAPI includes **Yggdrasil-style** routes under `/mcauth/...` (see `cmd/root.go` for the full list).
Common categories include:

* Authentication/session flows
* Account/profile-related endpoints
* Auxiliary endpoints for integrations

> For the definitive list, check the router setup in `cmd/root.go`.

---

## Logging

The logger is designed for concurrency correctness:

* Requests push log events into a channel
* A dedicated goroutine serializes writes to disk (prevents interleaving / “log corruption” under concurrency)
* Log file names are date-based: `YYYY-MM-DD.log`

Logs are written under the configured `log/` directory (see `packages/log/log.go`).

---

## Repo Structure

* `cmd/` — CLI entry, router setup, config loading
* `packages/api/` — API handlers and route logic
* `packages/email/` — SMTP pool and email delivery pipeline
* `packages/log/` — channel-based logger (single-writer)
* `packages/mongo/` — MongoDB connection and persistence utilities
* `packages/yggdrasil/` — Yggdrasil compatibility layer / helpers

---

## Development

### Run in dev mode (example)

```bash
go run . --config ./craapi.config.example
```

### Common improvements (PR-friendly)

* Add a `docs/` folder with endpoint documentation and examples
* Add a healthcheck endpoint and Prometheus metrics
* Add integration tests for SMTP failure scenarios and request flows

---

## Security Notes

* Do **not** commit real SMTP credentials or production MongoDB DSNs.
* Use environment variables or separate secret management for production deployments.

---

## License

(TODO)

```

---

### 中文描述

> CRAAPI 是一个基于 Go + fasthttp 的 Minecraft 账号体系认证 API（兼容 Yggdrasil 路由风格）。项目重点是工程化与可靠性：通过 SMTP worker pool 把邮件发送从主请求链路剥离，并用 channel + 单写协程保证高并发下日志不会错位。
