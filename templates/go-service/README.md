# go-service template

The canonical Duke-ECE Go service skeleton: **vertical slices inside a
hexagonal core**. Copy it, don't re-derive it.

## Instantiate

```sh
cp -r templates/go-service <name>
cd <name>
go mod edit -module github.com/Duke-ECE/<name>
GOTOOLCHAIN=local go mod tidy   # go directive must stay on the 1.25 line
# replace YOURSERVICE in k8s.yaml and .github/workflows/ci-cd.yml
# rename internal/example to your first real slice
```

Then: `gh repo create Duke-ECE/<name> --public --source . --push`, set the
`KUBE_CONFIG` repo secret, and the first push to `main` deploys.

## Layout and rules

```
cmd/server/main.go              # the ONLY assembly point (composition root)
internal/
  transport/
    grpc/                       # server construction, thin handlers, single
                                # error→status mapper; bufconn tests
  example/                      # a domain slice — rename it:
    example.go                  #   types
    service.go                  #   business rules (nothing above imported)
    store.go                    #   the Store port, OWNED BY THE SLICE
    errors.go                   #   domain errors
  infrastructure/
    postgrest/                  # port impl over Supabase PostgREST (+ httptest tests)
    memory/                     # port impl in memory — REQUIRED (zero-dep runs,
                                #   slice tests, reference behavior for parity)
test/                           # whole-service integration tests: assemble like
                                # main, drive the public API over real TCP, fake
                                # only true external boundaries
```

- Dependencies point inward: transport → slices ← infrastructure.
- Ports are declared in the slice that consumes them, never next to their
  implementation.
- Handlers do zero business logic; domain errors map to statuses in
  exactly one file per transport.
- Both backends behave identically (IDs, ordering, errors).
- No orphans: every package is wired into the running binary.
- Slices sit flat under `internal/` until there are more than five; then
  group them under `internal/domain/` in one move.

HTTP/Gin variant: `internal/transport/http/` with `server.go` (explicit
`http.Server`, `ReadHeaderTimeout` — never a blanket `WriteTimeout`, it
kills SSE), `routes.go`, `middleware.go`, `errors.go`.

Gates before every push: `gofmt -l .`, `go vet ./...`, `go test ./...`
(Makefile pins `GOTOOLCHAIN=local`).
