# TASK-2 Goモジュール/依存関係の確定

## Goal
Go単一バイナリのMVPに必要な最小依存を確定し、`go test ./...` が通る土台を作る。

このリポジトリでは、Goの導入/依存関係/タスク実行は **すべてmiseで管理** する（ローカル・CIともに同一手順）。

## Decisions (MVP)
- Web: `net/http`, `html/template`, `embed`（標準）
- DB: bbolt https://github.com/etcd-io/bbolt
- Config: `config.yaml` 必須
- AI: Codex CLI と JSON-RPC で通信する Provider
- Task ID: ULID（時系列ソート可能な一意ID）

## Steps
1. mise設定を追加
   - `mise.toml` を更新して `go` をピン留め（バージョンはプロジェクトで固定）
   - miseのtask runnerでチェック用コマンドを定義
     - `mise run tidy`: `go mod tidy`
     - `mise run test`: `go test ./...`
     - `mise run check`: `tidy` と `test` をまとめて実行（CIとローカルの統一入口）
2. `go.mod` を作成（mise経由で実行）
   - `mise run init`（taskとして用意するか、初回のみ手動）
   - `go mod init github.com/uta8a/spzcdd`（module pathは必要なら変更）
3. 依存追加（最小）
   - bbolt: `go.etcd.io/bbolt`
   - YAML: yaml/go-yaml（`go.yaml.in/yaml/v4`）（config.yaml読取用）
   - ULID: `github.com/oklog/ulid/v2`
4. CIにmiseを導入
   - `.github/workflows/ci.yml` を追加
   - `jdx/mise-action` を使い、`mise run check` を実行
5. `.gitignore` を確認（DBファイルなどが除外されていること）

## Acceptance Criteria
- `go mod tidy` が成功する
- `go test ./...` が（まだテストがなくても）実行可能な状態になる
- `mise run tidy` / `mise run test` / `mise run check` が動作する
- CIが `jdx/mise-action` 経由で `mise run check` を実行できる

## Notes
- ルータは `http.ServeMux` に固定（追加依存を避ける）。
- Task IDはULIDに固定する（生成方法は `ulid.Monotonic` を用いた時刻+乱数でよい）。

## Repo Policy
- Goの導入/バージョン固定/タスク実行は mise を唯一の入口にする
   - ローカル: `mise run check`
   - CI: `jdx/mise-action` + `mise run check`
