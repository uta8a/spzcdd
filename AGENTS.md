# AGENTS.md

このリポジトリで作業する（人間/エージェント）向けの最低限ルールです。

## Goals
- MVPとして「Spec Drivenを成立させる最小の閉ループ（Draft→Spec→Exec→Done）」を作る
- 単一Goバイナリ + embedded HTML + bbolt 永続化

## Source of Truth
- 実装タスク計画は docs/plan が正
  - 作業は docs/plan/TASK.md の順で進める

## Documentation Policy
- `docs/spec/` 配下のストック仕様書を追加/更新した場合は、必ずこの `AGENTS.md` にも参照（または方針）を追記して、正の所在が追える状態を維持する

## Architecture
- ディレクトリ構造と責務分割は `docs/spec/architecture.md` を正とする
  - TASKでコードを追加する際は、このドキュメントの構造に沿って配置する

## Domain Model
- ドメインモデル（Task/Spec/Execution・ULID方針）は `docs/spec/domain-model.md` を正とする

## Coding Guidelines
- 追加依存は最小（stdlib優先）
- 状態遷移は workflow 層に集約し、ハンドラは必ずそれを経由
- Spec/Execution は revision ごとに履歴を残す（上書き禁止）

## Tests / Quality
- Goの導入/バージョン固定/タスク実行は **mise** を唯一の入口にする
  - 変更後は原則 `mise run check`（= `mise run tidy` + `mise run test`）を実行
  - 必要に応じて `mise run tidy` / `mise run test` / `gofmt -w .` を使う

## Commit Messages
- `docs/agent/skills/commit-messages/SKILL.md` を参照

## Pull Requests
- PR作成時は `.github/pull_request_template.md` の項目を埋める
