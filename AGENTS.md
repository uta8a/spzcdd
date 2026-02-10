# AGENTS.md

このリポジトリで作業する（人間/エージェント）向けの最低限ルールです。

## Goals
- MVPとして「Spec Drivenを成立させる最小の閉ループ（Draft→Spec→Exec→Done）」を作る
- 単一Goバイナリ + embedded HTML + bbolt 永続化

## Source of Truth
- 実装タスク計画は docs/plan が正
  - 作業は docs/plan/TASK.md の順で進める

## Coding Guidelines
- 追加依存は最小（stdlib優先）
- 状態遷移は workflow 層に集約し、ハンドラは必ずそれを経由
- Spec/Execution は revision ごとに履歴を残す（上書き禁止）

## Tests / Quality
- 変更後は原則 `gofmt -w .` と `go test ./...` を実行

## Commit Messages
- `.agent/skills/commit-messages.md` を参照

## Pull Requests
- PR作成時は `.github/pull_request_template.md` の項目を埋める
