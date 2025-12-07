# CLAUDE.md

必ず日本語で返答してください

## プロジェクト概要

Minecraft Modの日本語翻訳を管理するGoライブラリ + CLI。
SQLite DBで翻訳データを一元管理し、高品質で一貫性のある翻訳を実現する。

## 基本原則

1. **完全性優先**: 1つのModを100%翻訳してから次へ進む
2. **公式翻訳尊重**: JARのja_jp.jsonがあれば必ず優先インポート
3. **データ品質**: 翻訳JSONに空文字("")を含めない、翻訳できないキーは含めない

## 標準ワークフロー

```bash
# 1. インポート
./moddict import -jar workspace/imports/mods/[mod.jar]

# 2. 公式翻訳確認・適用（あれば）
unzip -l [mod.jar] | grep -i ja_jp
unzip -p [mod.jar] 'assets/*/lang/ja_jp.json' > /tmp/official.json
./moddict translate -mod [mod_id] -official /tmp/official.json

# 3. 翻訳ループ（Pending: 0 になるまで）
./moddict translate -mod [mod_id] -status
./moddict translate -mod [mod_id] -export /tmp/pending.json -limit 20
# 翻訳生成後
./moddict translate -mod [mod_id] -json /tmp/translated.json

# 4. 完了確認 → Progress: 100.0%
./moddict translate -mod [mod_id] -status
```

## サブエージェント活用（重要）

**翻訳作業は必ずサブエージェント（Haiku）に委任する**

理由：
- 翻訳JSONがメインコンテキストを消費しない
- 定型作業なのでHaikuで十分、コスト効率が良い
- 複数カテゴリを並列処理できる

```
# サブエージェント起動例
Task tool (subagent_type: general-purpose, model: haiku)
→ 「[Mod名]の翻訳を実行してください」
→ docs/subagent-workflow.md のテンプレートに従う
```

## 翻訳ルール

- Minecraft公式訳に準拠（バニラアイテム・ブロック名）
- フォーマットコード保持: `§`, `$(...)`, `%s`, `%d`, `\n`
- 用語優先順位: global < category < mod

## プロジェクト構造

```
cmd/moddict/     # CLI
internal/        # 内部パッケージ
data/patterns/   # ファイルパターン定義
data/terms/      # 用語辞書
workspace/       # 作業ディレクトリ（Git対象外）
docs/            # 詳細ドキュメント
```

## 詳細ドキュメント

タスク固有の詳細は以下を参照:
- `docs/cli-reference.md` - 全コマンドリファレンス
- `docs/db-schema.md` - DBスキーマ詳細
- `docs/subagent-workflow.md` - サブエージェント運用ガイド
