# TASK-10 Web UI（一覧/詳細）+ HTTPハンドラ/ルート

## Goal
embedded HTMLで、Taskの作成→レビュー→差戻し→完了までの操作をボタンのみで実行できるUIを提供する。

## Pages (MVP)
### Task一覧（簡易カンバン）
- 列: `DRAFT / SPEC_REVIEW / EXEC_REVIEW / DONE`
- `*_PENDING` は同じ列で「処理中」表示
- カード: `id, title, state`

### Task詳細
- Draft表示/編集（DRAFTのみ編集可）
- Spec表示（SPEC_REVIEW以降）
- Exec表示（EXEC_REVIEW以降）
- レビュー入力（Spec用/Exec用のテキスト1つずつでOK）
- 状態に応じたボタンのみ表示
  - Publish / Approve / Request changes / Retry

## Routes (proposal)
- GET `/`（一覧）
- GET `/tasks/{id}`（詳細）
- POST `/tasks`（新規作成）
- POST `/tasks/{id}/draft`
- POST `/tasks/{id}/publish`
- POST `/tasks/{id}/spec/approve`
- POST `/tasks/{id}/spec/request_changes`
- POST `/tasks/{id}/exec/approve`
- POST `/tasks/{id}/exec/request_changes`
- POST `/tasks/{id}/retry`

## Steps
1. `internal/web/` にハンドラを実装（`http.ServeMux`）
2. `web/templates/` を作成し `embed` で読み込む
3. POST後はPRG（Post/Redirect/Get）に統一
4. 遷移時は必ずworkflowを通し、ストア更新→必要ならジョブenqueue

## Acceptance Criteria
- ブラウザ操作だけでMVP閉ループが完走できる
- 状態に応じて表示されるボタンが正しい
- 不正操作が行われても状態が壊れない（409等で拒否）

## Notes
- CSS/JSビルドは入れない。見た目は最小（可読性重視）。
