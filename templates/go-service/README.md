# go-service template

The canonical Duke-ECE Go service skeleton. Copy it, don't re-derive it.

## Instantiate

```sh
cp -r templates/go-service <name>
cd <name>
go mod edit -module github.com/Duke-ECE/<name>
GOTOOLCHAIN=local go mod tidy   # go directive must stay on the 1.25 line
# replace YOURSERVICE in k8s.yaml and .github/workflows/ci-cd.yml
```

Then: `gh repo create Duke-ECE/<name> --public --source . --push`, set the
`KUBE_CONFIG` repo secret, and the first push to `main` deploys.

## What's inside

- `cmd/server/main.go` — env, wiring, graceful shutdown. Nothing else.
- `internal/server/` — server construction (gRPC + reflection) and the
  bufconn test pattern. HTTP/Gin variant: build an explicit `http.Server`
  with `ReadHeaderTimeout` (never a blanket `WriteTimeout` — it kills SSE)
  and put the endpoint table in `internal/server/routes.go`.
- `internal/store/` — the PostgREST persistence seam (delete if stateless).
- `Dockerfile`, `k8s.yaml`, `.github/workflows/ci-cd.yml` — the standard
  build → ghcr → rollout pipeline.
- `Makefile` — `make build vet test run tidy` (with `GOTOOLCHAIN=local`).

Gates before every push: `gofmt -l .`, `go vet ./...`, `go test ./...`.
