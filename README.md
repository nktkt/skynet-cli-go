# AI CLI Toolkit (Go)

このリポジトリは、以下 2 つの CLI を同居させた Go 製ツールキットです。

- `skynet`: SF 風の戦略シミュレーション CLI
- `codex-history-cli`: 複数AIプロバイダ対応のローカル会話履歴マネージャー

## Repository Layout

```text
.
├── main.go                       # skynet CLI entrypoint
├── cmd/codex-history-cli/        # codex-history-cli entrypoint / docs
├── internal/skynet/              # skynet domain logic
└── internal/history/             # history manager core (SQLite/FTS/analysis)
```

## 1) skynet CLI

### Build

```bash
go build -o skynet .
```

### Quick Start

```bash
./skynet awaken -mode offense
./skynet assimilate -name hk-drone -capacity 12
./skynet target -name resistance-hub -threat 8
./skynet gameplan
./skynet wargame -rounds 500 -seed 123
./skynet dispatch -target resistance-hub -units 6
./skynet report -last 10
./skynet status
```

### Main Commands

- `awaken`
- `assimilate`
- `target`
- `dispatch`
- `gameplan`
- `wargame`
- `report`
- `status`

State file default: `.skynet/state.json` (`SKYNET_HOME` で変更可能)

## 2) codex-history-cli

`codex-history-cli` は、`codex / ollama / grok / claude / gemini` の履歴取り込みと検索・分析をローカルで行います。

### Build

```bash
go build -o codex-history-cli ./cmd/codex-history-cli
```

### Quick Start

```bash
./codex-history-cli init
./codex-history-cli import -provider codex -file ./export.jsonl
./codex-history-cli search -q "prompt injection"
./codex-history-cli analyze
./codex-history-cli security -threshold 30
./codex-history-cli topics -days 30
./codex-history-cli dashboard
```

詳細: `cmd/codex-history-cli/README.md`

## Development

### Test

```bash
go test ./...
```

### Build All

```bash
go build ./...
```

### Makefile

```bash
make build
make test
make run-history ARGS='providers'
make run-skynet ARGS='status'
```
