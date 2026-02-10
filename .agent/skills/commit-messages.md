# Skill: Commit Messages

このリポジトリでは、コミットメッセージを「後から見て作業単位が追える」ことを最優先にします。

## Rule (MVP)
- 1コミット = 1つのTASKを進める（小さければTASK内で複数コミット可）
- 先頭に必ず TASK-ID を入れる
- 変更の種類は Conventional Commits 風の接頭辞で揃える

## Format
```
<TYPE>: <short summary> (TASK-<ID>)

<optional body>
```

### TYPE
- `feat`: 仕様追加/機能追加
- `fix`: バグ修正
- `refactor`: 仕様変更なしの整理
- `test`: テスト追加/調整
- `docs`: ドキュメント
- `chore`: 雑務（依存更新、整形など）

## Examples
- `docs: add MVP task plans (TASK-1)`
- `feat: add bbolt store layer (TASK-6)`
- `test: cover state transitions (TASK-5)`

## Notes
- レビューで追いやすいように、サマリは「何をしたか」を短く具体的に書く。
- TASKをまたぐ巨大コミットは避ける（差分が読めない）。
