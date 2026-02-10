# TASK-11 結合・検証（テスト/手動動作確認）

## Goal
MVPが「壊れない閉ループ」になっていることを、最低限の自動テスト + 手動手順で確認できる状態にする。

## Automated Tests (minimum)
- workflow: 全遷移テスト（許可/不許可）
- store: Task/Spec/Execution のCRUD、rev増分と履歴保持
- （可能なら）ハンドラ: 不正遷移時のHTTPステータス

## Manual Verification
1. `config.yaml` を用意
2. `go run ./cmd/spzcdd` で起動
3. Task作成 → Draft編集 → Publish
4. Specレビュー画面で
   - Approve → Executeへ
   - Request changes → Spec再生成（rev+1）
5. Execレビュー画面で
   - Approve → DONE
   - Request changes → 再実行（rev+1）
6. Provider失敗を意図的に起こして ERROR → Retry

## Acceptance Criteria
- `go test ./...` が通る
- 上の手動手順で閉ループが回る
- 再起動試験（pendingが残ってもERRORへ落ち、Retryで回復できる）

## Notes
- MVPでは速度よりも整合性（状態・rev・履歴）を優先。
