# ドメインモデル（MVP）

このドキュメントは、MVPで必要な永続化・表示・状態遷移のための最小データモデルを定義する。

## ID方針

- TaskのIDは **ULID** を採用する
  - 理由: 文字列として扱いやすく、時系列で概ねソート可能
  - 表現: 26文字のCrockford Base32（例: `01ARZ3NDEKTSV4RRFFQ69G5FAV`）

- Taskには **通し番号（task number）** を付与する
  - 目的: ユーザーが `HOGE-1` のように短い識別子で指定しやすくする
  - `id (ULID)` は内部の主キー、`number` は人間向けの参照用
  - `number` は作成時に採番され、以後不変（再利用しない）
  - 表示用キーは `PREFIX-number` の形式（例: `HOGE-1`）

## エンティティ

### Task

1タスクの“現在地”と、参照する最新リビジョン（Spec/Execution）を保持する。

- `id` (string, ULID)
- `number` (int; 通し番号、1から開始)
- `title` (string)
- `draft_body` (string)
- `state` (string; 状態機械で定義された値のみ)
- `current_spec_rev` (int)
- `current_exec_rev` (int)
- `last_error` (string, optional)
- `created_at` (timestamp)
- `updated_at` (timestamp)

#### 制約
- `draft_body` の編集は `state == DRAFT` のときのみ許可
- `current_spec_rev` / `current_exec_rev` は「最新の承認対象（表示対象）」を指す
- `number` は一意で単調増加（欠番は許容、重複は不可）

#### 採番（永続化）方針
- `number` は永続層が採番する（例: bbolt の `meta` にシーケンスを持ち、トランザクション内で increment）
- UI/URL/ユーザー入力で使う “短い識別子” は基本 `PREFIX-number` を表示する
  - `PREFIX` は `config.yaml` で設定できるものとする（例: `task_id_prefix: HOGE`）

### Spec

Taskに紐づく実行計画。差戻しのたびに revision を増やして履歴を残す。

- `task_id` (string, ULID)
- `rev` (int; 1から開始)
- `body` (string)
- `status` (string; `proposed|approved`)
- `review_feedback` (string, optional; 差戻し理由の最新1つ)
- `created_at` (timestamp)

#### 制約
- **上書き禁止**: `(task_id, rev)` は不変
- `status == approved` は `SPEC_REVIEW -> EXEC_PENDING` の承認操作でのみセット

### Execution

Specに基づく実行結果。差戻しのたびに revision を増やして履歴を残す。

- `task_id` (string, ULID)
- `rev` (int; 1から開始)
- `body` (string)
- `status` (string; `proposed|approved`)
- `review_feedback` (string, optional; 差戻し理由の最新1つ)
- `created_at` (timestamp)

#### 制約
- **上書き禁止**: `(task_id, rev)` は不変
- `status == approved` は `EXEC_REVIEW -> DONE` の承認操作でのみセット

## 状態（State）

状態は Task が持つ。状態名は文字列として固定し、workflow層で遷移を強制する。

- `DRAFT`: 下書き（編集可）
- `SPEC_PENDING`: Spec生成待ち/生成中
- `SPEC_REVIEW`: Specレビュー待ち
- `EXEC_PENDING`: 実行待ち/実行中
- `EXEC_REVIEW`: 実行結果レビュー待ち
- `DONE`: 完了
- `ERROR`: 自動処理失敗（手動リトライ）

## Revision（履歴）

- Spec/Execution は revision ごとに履歴として保存し、後から参照できる
- Task は `current_*_rev` で「いまレビュー/表示に使う revision」を指す

## 直近コメントの扱い（MVP）

- ReviewCommentのスレッド化はMVPでは行わない
- 差戻し理由は Spec/Execution の `review_feedback` に「最新1つ」だけ持つ

## 参照

- ディレクトリ構造: `docs/spec/architecture.md`
- 実装タスク: `docs/plan/TASK-4_domain-model.md`
