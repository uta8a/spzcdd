# TASK-8 AI Provider抽象 + Codex JSON-RPC Provider

## Goal
Spec生成・Executeを差し替え可能にするため、Provider抽象を定義し、Codex CLIのJSON-RPC実装を追加する。

## Interfaces
- `GenerateSpec(ctx, input) -> specBody`
- `Execute(ctx, input) -> execBody`

## Inputs (MVP)
### GenerateSpec input
- `title`
- `draft_body`
- `previous_review_feedback`（任意）

### Execute input
- `spec_body`
- `previous_review_feedback`（任意）

## Codex JSON-RPC (implementation outline)
- `ai.command` / `ai.args` でプロセス起動
- stdin/stdout でJSON-RPCリクエスト/レスポンス
- タイムアウト/キャンセル（context）を尊重

## Steps
1. `internal/ai/provider.go` にインターフェース定義
2. `internal/ai/codexrpc/` に
   - JSON-RPCクライアント
   - Provider実装
3. Providerのユニットテスト（プロセスをモックできる形にする、またはインターフェース分離）

## Acceptance Criteria
- Providerが `GenerateSpec` / `Execute` を満たす
- 失敗時にエラーを返し、呼び出し側で `ERROR` へ落とせる

## Notes
- JSON-RPCのmethod名・payloadはCodex CLI側仕様に合わせて確定する必要がある。
  - ここが未確定なら、まずは「コマンド実行Provider（stdinにJSON、stdoutをそのまま採用）」のように薄くして吸収する。
