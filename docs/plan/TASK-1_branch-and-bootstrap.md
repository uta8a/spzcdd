# TASK-1 ブランチ作成と初期化

## Goal
MVP実装用の作業ブランチを作り、以降のタスクが迷わず進む最小の下準備を整える。

## Scope
- gitブランチ作成
- 最低限のディレクトリ（docs/plan は作成済み）方針を決める
- AI readyなリポジトリ設定（AGENTS.md / Agent Skills / PR template）

## Out of scope
- 実装コードの追加（TASK-2以降）

## Steps
1. ブランチ作成
   - `git switch -c feat/mvp-closed-loop`
2. AI readyなメタ情報を追加
   - `AGENTS.md` を追加（このリポジトリのGoals / Source of Truth / 品質ルール）
   - `.agent/skills/` に Agent Skills を追加
     - commit messageに関するSkillsを追加（例: `.agent/skills/commit-messages.md`）
3. Pull Request templateを追加
   - `.github/pull_request_template.md` を追加（Summary / Related Tasks / How To Test / Checklist）
4. 以降の作業方針をメモ（このdocs/plan配下のTASKを単位に進める）

## Acceptance Criteria
- 作業ブランチに切り替わっている
- 以降のタスクがこのブランチ上で進められる
- `AGENTS.md` が存在し、作業ルールが明文化されている
- `.agent/skills/commit-messages.md` が存在し、TASK-IDを含むcommit message規約が定義されている
- `.github/pull_request_template.md` が存在し、PR作成時に必要項目が揃っている

## Notes
- このリポジトリは現状ほぼ空なので、MVPの設計判断はdocs/planを単一の“仕様の正”として扱う。
