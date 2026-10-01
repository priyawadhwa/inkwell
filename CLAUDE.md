# CLAUDE.md

Inkwell Bank is a small Go web app: one binary, no third-party dependencies,
the frontend embedded with `go:embed`.

## Layout

- `main.go` — HTTP server, routes, page rendering
- `bank.go` — the demo customer's accounts, transactions, and cards
- `flags.go`, `flags.json` — feature flags, compiled into the binary
- `templates/index.html`, `static/style.css` — the dashboard
- `.chainguard/ci/` — CI checks (build, vet, test, gofmt, lint)

## Feature flags

Flags live in `flags.json` and are read at startup. Turning a feature on means
setting its flag to `true` and updating any test that asserts the old
behavior. `TestHomePage` renders the dashboard with the shipped flags, so a
flag that changes what customers see changes what it expects.

## CI and pull requests

Run CI with the `run_ci` tool and get it green before calling a change done.
Don't push or open a pull request until asked. When asked to ship, push a
branch, open the PR with `gh pr create` (a one-line summary of what customers
will see), then `gh pr merge --auto --squash`.
