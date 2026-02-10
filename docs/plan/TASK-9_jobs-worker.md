# TASK-9 ジョブワーカー（SPEC_PENDING/EXEC_PENDING）と再実行

## Goal
Publish/Approve/差戻し操作で非同期ジョブが走り、完了でレビュー状態へ進む“閉ループ”を成立させる。

## Behavior (MVP)
- `SPEC_PENDING`:
  - ProviderでSpec生成
  - Specを `rev+1` で保存
  - 状態を `SPEC_REVIEW` へ
- `EXEC_PENDING`:
  - Providerで実行
  - Executionを `rev+1` で保存
  - 状態を `EXEC_REVIEW` へ
- Provider失敗:
  - 状態を `ERROR` にして `last_error` を保存

## Restart Policy (simple)
- 起動時に `SPEC_PENDING` / `EXEC_PENDING` が残っていたら `ERROR` に落とす
- ユーザーが `Retry` を押して再投入

## Steps
1. `internal/jobs/worker.go` を実装（in-memory queue）
2. enqueue API（Spec生成、Execute）を公開
3. app起動時にpendingタスクの回復方針を実装

## Acceptance Criteria
- Publish → Spec生成 → Specレビュー表示まで自動で進む
- Spec差戻しでrevが増え、再生成される
- Execute差戻しでrevが増え、再実行される
- 失敗時にERRORになり、Retryできる

## Notes
- MVPでは「同時実行数=1」でも良い（シンプル優先）。
