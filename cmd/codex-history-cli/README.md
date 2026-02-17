# codex-history-cli

Local AI conversation history manager for multiple providers.

## Features

- Multi-provider import: `codex`, `ollama`, `grok`, `claude`, `gemini`
- SQLite + FTS5 search
- Auto summary and tag generation (heuristic, optional Ollama local model)
- Prompt injection risk scoring
- Topic trend graph (ASCII)
- Interactive dashboard/TUI mode

## Build

```bash
go build -o codex-history-cli ./cmd/codex-history-cli
```

## Quick Start

```bash
./codex-history-cli init
./codex-history-cli import -provider codex -file ./export.jsonl
./codex-history-cli search -q "prompt injection"
./codex-history-cli analyze
./codex-history-cli security -threshold 30
./codex-history-cli topics -days 30
./codex-history-cli dashboard
```

## Optional local LLM summary

```bash
./codex-history-cli analyze -ollama-model qwen2.5:3b
```

Set DB path with `-db` or `CODEX_HISTORY_DB`.

## TUI Commands

Inside `dashboard`:

- `r`: refresh
- `s <query>`: search
- `days <n>`: change lookback days
- `th <n>`: change risk threshold
- `q`: quit
