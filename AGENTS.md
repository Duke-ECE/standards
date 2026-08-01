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

**Vertical slices inside a hexagonal core.** Use `templates/go-service/`.
The layout, and what belongs in each file:

```
cmd/server/main.go            # the ONLY assembly point (composition root)
internal/
  transport/
    grpc|http/                # server construction, thin handlers, single
                              # error→status mapper (errors.go); HTTP: routes.go
  <slice>/                    # one package per aggregate:
    <name>.go                 #   types
    service.go                #   business rules — imports nothing above itself
    store.go                  #   the Store port, OWNED BY THE SLICE
    errors.go                 #   domain errors
  infrastructure/
    postgrest/                # port implementations over Supabase PostgREST
    memory/                   # in-memory implementations — REQUIRED
  config/                     # env vars with defaults
```

Rules:

- Dependencies point inward: transport → slices ← infrastructure.
  Ports are declared in the slice that consumes them, never next to
  their implementation.
- `cmd/server`, not `cmd/api`. Handlers do zero business logic; domain
  errors map to statuses in exactly one file per transport.
- **Zero-dep local runs**: every port has a memory implementation; it is
  the reference for behavior parity (same IDs, ordering, errors).
- **No orphans**: every package is wired into the running binary.
- Slices sit flat under `internal/` until there are more than five; then
  group them under `internal/domain/` in one move.
- Construction is explicit (`NewServer()`), not bare `http.ListenAndServe`.
- **Timeouts are SSE-aware**: set `ReadHeaderTimeout`; never set a blanket
  `WriteTimeout` on servers that stream (it kills long-lived SSE/gRPC streams).
- Every long-running service implements **graceful shutdown**
  (signal.NotifyContext → Shutdown/GracefulStop).
- `/health` reports dependency status and must never `log.Fatal` — a sick
  dependency is a 200 with `"status":"down"` detail, not a dead process.
- **Escape valve**: a service whose slice would be a pure pass-through
  (a proxy with no rules of its own) may collapse transport and slice
  into one package. The moment a rule appears (ownership, validation,
  orchestration of two backends), it moves into a slice.

## 3. Go — testing, gates, toolchain

- **stdlib `testing` only.** No testify, no gomock. Seams are small
  interfaces (e.g. `Store`, `Fetcher`) faked in tests; gRPC is tested over
  `bufconn`; HTTP backends of external APIs are faked with `httptest`.
- **Two test levels.** Unit tests live next to the code they cover.
  A top-level `test/` holds **integration tests**: assemble the whole
  service exactly like `cmd/server/main.go` does, drive its public API
  over a real connection (TCP listener, not just in-package calls), and
  fake only true external boundaries (Supabase HTTP, the k8s API).
- Gates: `gofmt -l .` clean, `go vet ./...`, `go test ./...` — all must pass
  before every push. Tests live next to the code they cover.
- `./scripts/check.sh` (also `make check`, wired into CI) mechanically
  verifies what a script can: go-directive on the 1.25 line, gofmt, no
  committed `.env`, timestamped migration names, no secret literals in
  manifests. Keep it passing; extend it when a rule is checkable.
- **Toolchain trap (macOS dev machines auto-upgrade Go)**: after any
  `go get` / `go mod tidy`, the `go` directive in `go.mod` must stay on the
  1.25 line (Docker builds with `golang:1.25-alpine`, `GOTOOLCHAIN=local`).
  Verify everything as `GOTOOLCHAIN=local go build/vet/test ./...`.
  If tidy writes `go 1.26`, pin dependencies back (e.g. k8s.io v0.34.x).

## 3a. Logging, dependencies, review

- **Logging**: `log/slog` for anything new (structured, levels);
  stdlib `log` is acceptable in `main` for boot lines. Never log secrets;
  log startup config as booleans/paths, not values.
- **Dependencies**: stdlib first. Every new module dependency needs a
  one-line justification in the commit message; toolchain-sensitive dep
  families (k8s.io) stay pinned to the versions in the template's go.sum.
- **Review**: solo work may go straight to `main` once gates are green
  locally; cross-cutting changes (contracts, shared CI, another repo's
  domain) go through a PR. CI must be green before any deploy happens —
  never bypass a red check.

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

## 9. New service checklist

1. Copy `templates/go-service`, `go mod edit -module`, rename
   `YOURSERVICE` in `k8s.yaml` + CI, rename `internal/example`.
2. Contract first if it has an API: proto in `Duke-ECE/protos` → tag →
   pin in go.mod.
3. `gh repo create Duke-ECE/<name> --public --source . --push`; set the
   `KUBE_CONFIG` repo secret.
4. Tables needed? Timestamped migration in `supabase/migrations/`,
   `supabase link --project-ref <ref>` + `supabase db push`.
5. Secrets via `kubectl create secret` (never committed); wire with
   `secretKeyRef` in `k8s.yaml`.
6. First push to `main` deploys; verify with the rollout + a grpcurl /
   curl smoke check (see the docs Operations page).
