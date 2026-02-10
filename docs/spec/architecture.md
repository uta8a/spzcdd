# アーキテクチャ

## ゴール（MVP）
- 単一Goバイナリ
- embedded HTML（サーバーサイドレンダリング）
- bbolt による永続化
- 状態遷移（workflow/state-machine）の整合は workflow 層で担保する

## ディレクトリ構造

### エントリポイント
- `cmd/spzcdd/`
  - `main.go`: プロセスの入口（flags、シグナル処理）を担当し、`internal/app` に処理を委譲する

### アプリ配線（wiring）
- `internal/app/`
  - `app.go`: 実行時依存を組み立て、HTTPサーバを起動する

### ドメイン / workflow / アダプタ（導入予定）
以下はMVPコンポーネントの“置き場所”として想定しているパッケージです。TASKの進行に合わせて段階的に導入します。

- `internal/domain/`
  - エンティティ（Task/Spec/Execution）と値オブジェクト
- `internal/workflow/`
  - 状態機械とワークフロー操作
  - ハンドラは必ずこの層を経由して状態遷移させる
- `internal/store/`
  - bbolt 永続化（bucket/schema + CRUD）
- `internal/ai/`
  - AI Provider 抽象 + Codex JSON-RPC Provider
- `internal/jobs/`
  - `*_PENDING` フェーズやリトライを進めるプロセス内ワーカー
- `internal/web/`
  - HTTPハンドラと表示用モデル
- `web/templates/`
  - バイナリに埋め込むHTMLテンプレート（`embed` 経由）

## 実行時（Runtime）
- アプリはHTTPサーバを起動し、以下のエンドポイントを提供する:
  - `/`（TASK-3 時点では暫定のプレースホルダ）
  - `/healthz`

## 参照
- Task plan: `docs/plan/TASK-3_project-skeleton.md`
