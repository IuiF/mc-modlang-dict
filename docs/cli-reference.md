# CLI リファレンス

## コマンド一覧

| コマンド | 説明 |
|---------|------|
| `moddict import -jar [file]` | JARからインポート（既存ソース・バージョン再利用） |
| `moddict import-dir` | ディレクトリからインポート |
| `moddict translate -mod [id] -status` | 翻訳進捗確認 |
| `moddict translate -mod [id] -export [file] -limit N` | pendingをエクスポート |
| `moddict translate -mod [id] -json [file]` | 翻訳をインポート |
| `moddict translate -mod [id] -official [file]` | 公式翻訳をインポート |
| `moddict export -mod [id]` | 翻訳済みファイル出力 |
| `moddict repair` | データベース整合性の修復 |
| `moddict migrate` | スキーマ移行・バージョン情報修正 |

## 使用例

```bash
# Mod単位で翻訳を管理（バージョン指定不要）
moddict translate -mod [mod_id] -status              # 進捗確認
moddict translate -mod [mod_id] -export pending.json # pendingをエクスポート
moddict translate -mod [mod_id] -json translated.json # 翻訳をインポート
moddict export -mod [mod_id]                          # 翻訳済みファイル出力
```

## 注意事項

- 翻訳結果をJSONファイルに保存しただけでは不十分
- 必ず `-json` フラグでDBにインポートすること

## データ品質保護（自動チェック）

CLIに以下の自動チェック機能が実装されています：

1. **翻訳インポート時** (`-json`, `-official`)
   - 空文字("")の翻訳は自動スキップ
   - source_text = target_text（原文と同一）の翻訳は自動スキップ
   - スキップされた件数がレポートされる

2. **修復コマンド** (`moddict repair`)
   - source=target問題を検出してレポート
   - 問題があれば手動修正用のSQLを表示

```bash
# 定期的にrepairを実行して品質チェック
./moddict repair -dry-run
```
