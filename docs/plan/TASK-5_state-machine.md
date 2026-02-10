# TASK-5 状態機械（遷移/ガード）実装

## Goal
「ボタン操作のみ」で状態が遷移することをコード上で保証し、UI/HTTPがドメイン不整合を起こさないようにする。

## Transitions (MVP)
- `DRAFT` → (Publish) → `SPEC_PENDING`
- `SPEC_PENDING` → (AI完了) → `SPEC_REVIEW`
- `SPEC_REVIEW` → (Approve) → `EXEC_PENDING`
- `SPEC_REVIEW` → (Request changes) → `SPEC_PENDING`（一段前の *_PENDING に戻してAIに再作業させる）
- `EXEC_PENDING` → (AI完了) → `EXEC_REVIEW`
- `EXEC_REVIEW` → (Approve) → `DONE`
- `EXEC_REVIEW` → (Request changes) → `EXEC_PENDING`（一段前の *_PENDING に戻してAIに再作業させる）
- `*` → (AI失敗) → `ERROR`
- `ERROR` → (Reset) → `DRAFT`

## Implementation Notes
- 遷移は `workflow.Can(action, from)` のような純関数に集約
- ハンドラ層は必ず状態機械を通してから永続化
- 許可されない操作は 400/409 で返す
- Request changes は「一段前の *_PENDING に戻す」ことを保証し、その後のAI処理（Spec/Executionの生成）はワーカーが行う
- Reset は ERROR を解消して DRAFT に戻す（MVPでは *_PENDING へ直接戻さない）

## Steps
1. `internal/workflow/state_machine.go` を作成
2. `Action`（Publish/Approve/RequestChanges/Reset/AIComplete/AIFail）を定義
3. テスト（テーブル駆動）で全遷移を固定

## Acceptance Criteria
- 許可遷移がテストで固定される
- 不正遷移が必ず拒否される

## Out of scope
- `EXEC_REVIEW → SPEC_REVIEW` の強制巻き戻し（MVP外）
