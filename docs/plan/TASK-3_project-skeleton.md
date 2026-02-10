# TASK-3 パッケージ構成とエントリポイント作成

## Goal
単一バイナリの起動経路（main → app wiring）と、今後の追加が自然に収まるディレクトリ構成を作る。

## Proposed Structure
- `cmd/spzcdd/main.go`（エントリポイント）
- `internal/app/`（DI・起動・HTTPサーバ起動）
- `internal/config/`（config.yamlロード）
- `internal/domain/`（エンティティ）
- `internal/workflow/`（状態機械）
- `internal/store/`（bbolt）
- `internal/ai/`（Provider抽象、Codex実装）
- `internal/jobs/`（ワーカー）
- `internal/web/`（HTTP handlers + templates）
- `web/templates/`（embedするテンプレ）

## Steps
1. `cmd/spzcdd/main.go` を作成
   - configロード
   - app生成
   - HTTPサーバ起動
2. `internal/app/app.go` を作成
   - Store / AI Provider / Worker / Handlers を組み立て

## Acceptance Criteria
- `go run ./cmd/spzcdd` が（実装未完でも）起動経路まで到達する
- `internal` 配下が責務ごとに分割されている

## Notes
- この段階ではルーティングやDB詳細に踏み込まない（TASK-6,10へ）。
