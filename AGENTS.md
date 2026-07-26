# AGENTS.md — Duke-ECE engineering rules

Applies to every repo in the managed-agents platform. Service repos keep a
short `AGENTS.md` with repo-specific facts and a link to this file.

## 1. Repos and contracts

- One deployable per repo; each repo has its own CI/CD, Dockerfile, `k8s.yaml`.
- Cross-service contracts are **pinned versions**, never local paths:
  - gRPC: `github.com/Duke-ECE/protos` (buf v2). Additive changes go in the
    existing `vN` package; breaking changes get a new package version. After
    merging, tag (`vX.Y.Z`) and bump consumers. `gen/go/` is committed
    generated code — never hand-edit. Lint baseline: `STANDARD` (`*Service`
    names, `XxxResponse` replies).
  - TypeScript consumers (agent-runtime) vendor `.proto` files via
    `npm run sync-proto` pinned to a tag; never hand-edit vendored copies.
- New services get a repo under `Duke-ECE/`, public, with the standard CI
  workflow, and the `KUBE_CONFIG` repo secret for auto-deploy.

## 2. Go services — layout

Use `templates/go-service/`. The layout, and what belongs in each file:

```
cmd/server/main.go        # env parsing, wiring, graceful shutdown ONLY — no logic
internal/server/server.go # server construction: grpc.Server/http.Server, timeouts
internal/server/routes.go # HTTP services: the endpoint table (RegisterRoutes)
internal/<domain>/        # business logic + RPC handlers, one package per domain
internal/store/           # persistence behind an interface (delete if stateless)
supabase/migrations/      # schema, timestamped versions (see §4)
k8s.yaml                  # Deployment + Service; env from Secrets via secretKeyRef
Dockerfile                # golang:1.25-alpine → alpine:3.21, non-root uid 10001
Makefile                  # build / vet / test / run — thin wrappers, no magic
```

Rules:

- `cmd/server`, not `cmd/api`. Handlers never live in `main`.
- Construction is explicit (`NewServer()`), not bare `http.ListenAndServe`.
- **Timeouts are SSE-aware**: set `ReadHeaderTimeout`; never set a blanket
  `WriteTimeout` on servers that stream (it kills long-lived SSE/gRPC streams).
- Every long-running service implements **graceful shutdown**
  (signal.NotifyContext → Shutdown/GracefulStop).
- `/health` reports dependency status and must never `log.Fatal` — a sick
  dependency is a 200 with `"status":"down"` detail, not a dead process.

## 3. Go — testing, gates, toolchain

- **stdlib `testing` only.** No testify, no gomock. Seams are small
  interfaces (e.g. `Store`, `Fetcher`) faked in tests; gRPC is tested over
  `bufconn`; HTTP backends of external APIs are faked with `httptest`.
- Gates: `gofmt -l .` clean, `go vet ./...`, `go test ./...` — all must pass
  before every push. Tests live next to the code they cover.
- **Toolchain trap (macOS dev machines auto-upgrade Go)**: after any
  `go get` / `go mod tidy`, the `go` directive in `go.mod` must stay on the
  1.25 line (Docker builds with `golang:1.25-alpine`, `GOTOOLCHAIN=local`).
  Verify everything as `GOTOOLCHAIN=local go build/vet/test ./...`.
  If tidy writes `go 1.26`, pin dependencies back (e.g. k8s.io v0.34.x).

## 4. Data access (Supabase)

- Services talk to Postgres via **PostgREST** with a service key
  (`apikey` + `Authorization: Bearer` headers) — not pgx, not database/sql.
  One access pattern platform-wide.
- Tables: `enable row level security` with **no policies** for service-owned
  tables (service role bypasses RLS; anon/authenticated get nothing).
  API keys never reach the browser.
- Migrations: one shared Supabase project, several repos push migrations →
  **version files by timestamp** (`YYYYMMDDHHMMSS_name.sql`). Plain counters
  collide across repos and `db push` silently skips the loser.
- Never commit secrets. Migration files contain schema only; data with
  secrets is seeded by one-off commands.

## 5. Secrets and env

- Nothing secret in git — no `.env` files, no keys in manifests.
- Runtime secrets are k8s Secrets consumed via `secretKeyRef` in `k8s.yaml`.
- CI credentials live in GitHub repo secrets (`KUBE_CONFIG`, `GITHUB_TOKEN`).
- Document every env var (name, default, behavior when unset) in the repo
  README. Unset optional features must degrade gracefully (fail closed for
  auth, feature-off otherwise).

## 6. CI/CD, Docker, k8s

- Standard workflow (copy from the template): on PR/push → build + vet +
  test; on push to `main` also → docker buildx → `ghcr.io/duke-ece/<repo>
  :{sha,latest}` → `kubectl apply -f k8s.yaml` → `kubectl set image … :<sha>`
  → wait for rollout.
- Dockerfile: multi-stage `golang:1.25-alpine` build (`CGO_ENABLED=0`) →
  `alpine:3.21`, `USER 10001`, `EXPOSE` the service port.
- gRPC services also enable **server reflection** (grpcurl debugging).
- Readiness gates: a pod that cannot serve (e.g. JWKS not yet fetched) must
  not report ready.

## 7. TypeScript repos (quick rules)

- `tsc --strict` is the gate; ESM with `.js` extensions on relative imports.
- After any `npm install` on a China-network dev machine:
  `grep -c npmmirror package-lock.json` must be 0 before committing —
  regenerate with the default registry if not.
- gRPC gotcha (grpc-js): `call.destroy(err)` doesn't deliver status —
  use `call.emit("error", err)`.

## 8. Git hygiene

- Conventional-ish commit messages (`feat:`, `fix:`, `chore:`), English.
- Commit generated code (protos `gen/go/`) together with its source change.
- Don't commit: secrets, `supabase/.temp/`, build artifacts, `.DS_Store`.
