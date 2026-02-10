# TASK-7 config.yaml（必須）読み込みとバリデーション

## Goal
`config.yaml` が存在しない場合は起動失敗にし、必要な設定（db/listen/ai）を確実に受け取る。

## Config Fields (MVP)
```yaml
db_path: "/path/to/spzcdd.db"
listen_addr: "127.0.0.1:8080"
ai:
  # Codex CLI JSON-RPC
  command: "codex"   # 例
  args: ["--jsonrpc"]
  timeout_seconds: 300
```

## Validation Rules
- `db_path` 必須
- `listen_addr` 必須
- `ai.command` 必須
- `ai.args` は任意（空でもOK）

## Steps
1. `internal/config/config.go` を実装
2. `cmd/spzcdd/main.go` からロードし、エラー時は終了
3. サンプル config を `config.example.yaml` として追加（任意だが運用上便利）

## Dependency
- YAML: `go.yaml.in/yaml/v4`

## Acceptance Criteria
- configなしで起動すると明示的にエラー
- 必須項目が欠けているとエラー

## Notes
- flags/env はMVPでは入れない（必要になったら後続タスクで追加）。
