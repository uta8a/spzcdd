# TASK-6 bboltストア層（バケット設計 + CRUD）

## Goal
MVPが必要とする最小の永続化（Task/Spec/Executionの履歴）をbboltで実装する。

## Bucket Design (proposal)
- `meta`（schema version など）
- `tasks`
  - key: `task/{taskID}` → Task(JSON)
- `specs`
  - key: `spec/{taskID}/{rev}` → Spec(JSON)
- `execs`
  - key: `exec/{taskID}/{rev}` → Execution(JSON)

※revは順序安定のため `0000000001` のようにゼロパディングを推奨。

## Store API (MVP)
- Task
  - `CreateTask(title, draftBody)`
  - `UpdateDraft(taskID, draftBody)`（DRAFTのみ）
  - `SetState(taskID, state)`
  - `SetError(taskID, errMsg)`
  - `GetTask(taskID)` / `ListTasks()`
- Spec/Exec
  - `PutSpec(spec)` / `GetSpec(taskID, rev)` / `GetCurrentSpec(taskID)`
  - `PutExecution(exec)` / `GetExecution(taskID, rev)` / `GetCurrentExecution(taskID)`

## Steps
1. `internal/store/` にbbolt open/closeを実装
2. schema version を `meta` に保存（将来拡張のため）
3. CRUDを実装
4. 最小テスト（temp file + open + CRUD）

## Acceptance Criteria
- Task/Spec/Execution が永続化できる
- Spec/Execution はrevごとに上書きされない（履歴が残る）

## Notes
- 大きい成果物はMVPでは保存しない（exec_bodyは要約/ログ文字列）。
