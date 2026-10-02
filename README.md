# 🐙 Inkwell Bank

Banking with all eight arms.

Inkwell Bank is a demo app for Chainguard Workspaces and Checks: a customer
dashboard for Linky, Chainguard's octopus mascot, with features gated behind
flags in [`flags.json`](./flags.json).

Live at https://inkwell-546601746857.us-central1.run.app.

| Flag | What customers see |
|---|---|
| `camouflage_mode` | Balances blend into the reef until tapped |
| `ink_cloud_card_freeze` | A button that inks (freezes) a card |

## Develop

```bash
chainctl develop github.com/priyawadhwa/inkwell
```

or locally:

```bash
go run .            # http://localhost:8080
go test ./...
```

## CI and deploy

Pull requests run the checks in [`.chainguard/ci.yaml`](./.chainguard/ci.yaml) on
Chainguard Checks. Merging to `main` deploys to Cloud Run with
[`build.yaml`](./build.yaml). Open dashboards reload themselves when
a new revision goes live.
