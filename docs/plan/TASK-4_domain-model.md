# TASK-4 ドメインモデル定義（Task/Spec/Execution）

## Goal
MVPの永続化/表示/状態遷移に必要なデータ構造を最小で固定し、後続タスクで迷わないようにする。

## Entities (MVP)
### Task
- `id`（ULID）
- `title`
- `draft_body`
- `state`（状態機械）
- `current_spec_rev`（int）
- `current_exec_rev`（int）
- `last_error`（string, optional）
- `created_at`, `updated_at`

### Spec
- `task_id`
- `rev`（int）
- `body`（spec_body）
- `status`（`proposed|approved`）
- `review_feedback`（差戻しコメント、最新1つ）
- `created_at`

### Execution
- `task_id`
- `rev`（int）
- `body`（exec_body）
- `status`（`proposed|approved`）
- `review_feedback`（差戻しコメント、最新1つ）
- `created_at`

## States (MVP)
- `DRAFT`
- `SPEC_PENDING`
- `SPEC_REVIEW`
- `EXEC_PENDING`
- `EXEC_REVIEW`
- `DONE`
- `ERROR`

## Steps
1. `internal/domain/` に型定義を追加
2. JSON/YAML保存用にシリアライズ方針を決める（例: JSON）
3. `updated_at` の更新タイミングを方針化

## Acceptance Criteria
- ドメイン型がコンパイル可能
- 状態名が文字列で一意に固定される

## Test Policy (MVP)
ドメイン層は「ルールの最小単位」なので、永続化/HTTPが無い段階でもユニットテストで不変条件を固定する。

### 対象
- `TaskState.IsValid()` が定義済み状態だけを許可する
- `ReviewStatus.IsValid()` が `proposed|approved` だけを許可する
- `NewTask(now, title, draft)`
	- `title` 必須
	- 初期状態が `DRAFT`
	- `CreatedAt/UpdatedAt` が `now`（UTC）になる
	- `id` がULIDとしてパース可能（形式の最低限チェック）
- `Task.SetDraftBody(now, body)`
	- `state == DRAFT` のときのみ成功
	- 成功時に `UpdatedAt` が更新される
- `NewSpec/NewExecution(now, taskID, rev, body)`
	- `taskID` 必須
	- `rev >= 1` 必須
	- 初期 `status == proposed`

### 非対象（TASK-4ではやらない）
- state machine（遷移/ガード）はTASK-5で別途テスト
- store層のシリアライズ互換やbucket設計はTASK-6でテスト

### 実装メモ
- テストファイル例: `internal/domain/*_test.go`
- ULIDのテストは「文字列長26 + `ulid.Parse` が通る」程度に留める（乱数の再現性は要求しない）

## Notes
- ReviewCommentは別エンティティにせず、`review_feedback` 1フィールドでMVP成立させる。

## Spec Doc
- ストックドキュメント: `docs/spec/domain-model.md`
