# CLAUDE.md

必ず日本語で返答してください

## 概要

Minecraft Modの翻訳辞書とファイル構造パターンを管理するGoライブラリ。
高品質で一貫性のある日本語翻訳を実現する。

---

## 基本原則

### 完全性優先（最重要）

**1つのModを100%翻訳してから次に進む。**

```
❌ NG: 6つのModを並列で60%ずつ翻訳
✅ OK: 1つのModを100%翻訳 → 次のModへ
```

### 公式翻訳の尊重（必須）

```
⚠️ 公式翻訳がある場合は必ず尊重する
- JARのja_jp.jsonを優先的にインポート（translator=official, status=verified）
- 公式翻訳を勝手に上書きしない
- 不足キーのみ翻訳対象とする
- 公式にしかないキーがあれば検出漏れとして調査
```

### 禁止事項

- サンプル翻訳だけ作成して「完了」とすること
- 翻訳率60%で次のModに進むこと
- 主要エントリだけ翻訳して残りを放置すること

---

## 翻訳ワークフロー

### 作業手順

1. **CurseForgeからModファイル（JAR）を取得**
   - 保存先: `workspace/imports/mods/`

2. **公式翻訳の確認**
   - `unzip -l [mod.jar] | grep ja_jp` で確認
   - **公式翻訳がある場合は必ず優先する**

3. **DBにインポート**
   - `moddict import` でen_us.jsonをDBに投入
   - 公式ja_jp.jsonがあれば同時にインポート（status=verified）

4. **不足分のみ翻訳**
   - DBからpendingをクエリ
   - **20件ずつ**Haikuサブエージェントで翻訳
   - **100件ごとにエージェントをリセット**（コンテキスト肥大化防止）
   - **並列実行可**

5. **検証**
   - 翻訳率100%を確認
   - フォーマットコードの保持確認

6. **次のModへ**

### 厳密ワークフロー（コマンド）

**⚠️ 全てのエージェントはこのワークフローに厳密に従うこと**

#### ステップ1: 初期化

```bash
# 1-1. Modのインポート
./moddict import -jar workspace/imports/mods/[mod.jar]
./moddict repair

# 1-2. 公式翻訳の確認と適用
unzip -l workspace/imports/mods/[mod.jar] | grep -i ja_jp
# ja_jp.jsonが存在する場合のみ:
unzip -p workspace/imports/mods/[mod.jar] 'assets/*/lang/ja_jp.json' > /tmp/official_ja.json
./moddict translate -mod [mod_id] -official /tmp/official_ja.json

# 1-3. 初期ステータス確認
./moddict translate -mod [mod_id] -status
```

#### ステップ2: 翻訳ループ

```bash
# 2-1. pendingエクスポート
./moddict translate -mod [mod_id] -export /tmp/pending.json -limit 20

# 2-2. 翻訳実行（エージェント内で翻訳を生成）
# ⚠️ 翻訳JSONには空文字("")を含めてはならない
# ⚠️ 翻訳できないキーはJSONに含めない（空文字で埋めない）

# 2-3. 翻訳インポート
./moddict translate -mod [mod_id] -json /tmp/translated.json

# 2-4. 進捗確認（Pending: 0 になるまで繰り返し）
./moddict translate -mod [mod_id] -status
```

#### ステップ3: 完了検証

```bash
# 最終ステータス確認 → Progress: 100.0% を確認
./moddict translate -mod [mod_id] -status
```

### バージョン更新時の動作

- 新バージョンをインポート → 既存と同じsource_textなら同じSourceを再利用
- 翻訳が自動的に継承される（継承コマンド不要）
- source_textが変わったキーのみ新規ソースを作成（pending状態）

### 翻訳JSON生成ルール（厳守）

```
✅ 正しい翻訳JSON:
{
  "item.mod.example": "サンプルアイテム",
  "block.mod.example": "サンプルブロック"
}

❌ 禁止パターン1: 空文字
{
  "item.mod.example": "",      ← データ破損の原因
  "block.mod.example": "サンプルブロック"
}

❌ 禁止パターン2: 翻訳できないキーを含める
{
  "item.mod.example": "item.mod.example",  ← 原文のままは無意味
  "block.mod.example": "サンプルブロック"
}

✅ 翻訳できないキーは含めない:
{
  "block.mod.example": "サンプルブロック"
}
```

---

## サブエージェント運用

### エージェント構成

| エージェント | 役割 | 並列実行 |
|-------------|------|----------|
| 翻訳Agent | Mod単位・ファイル単位の翻訳実行 | 可 |
| 用語抽出Agent | 新規用語の抽出・辞書提案 | 可 |
| レビュアーAgent | 翻訳品質・一貫性の検証 | 最終段階 |
| 調査Agent | Mod構造・既存翻訳の調査 | 可 |

### 並列実行の原則

1. **Mod単位では順次実行**
   - 1つのModを100%完了してから次へ
   - 中途半端な状態で複数Modを並列しない

2. **同一Mod内のカテゴリは並列可**
   - block / item / advancement 等のカテゴリ分割
   - lang / patchouli の同時処理

3. **レビューは各Mod完了時に実施**
   - 翻訳率100%を確認
   - 用語の一貫性、フォーマット保持を検証

### コンテキスト効率化

**重要: サブエージェントにDB操作を委任する**

```
❌ NG: メインでJSON出力 → 内容確認 → サブエージェントに渡す
       （メインコンテキストにJSONが展開されてしまう）

✅ OK: サブエージェント起動 → エージェント内でDB操作 → 完了報告のみ返す
       （メインコンテキストを消費しない）
```

### サブエージェント指示テンプレート

```
[Mod名]の翻訳を行ってください。

【厳守事項】
- 翻訳JSONに空文字("")を含めない
- 翻訳できないキーはJSONに含めない
- 各ステップでステータスを確認する

1. JARをインポート:
   ./moddict import -jar workspace/imports/mods/[mod.jar]
   ./moddict repair

2. 公式翻訳確認・インポート:
   unzip -l [mod.jar] | grep ja_jp
   （あれば）./moddict translate -mod [mod_id] -official [ja_jp.json]

3. ステータス確認:
   ./moddict translate -mod [mod_id] -status

4. pendingがあれば20件ずつ翻訳:
   ./moddict translate -mod [mod_id] -export pending.json -limit 20
   （翻訳後 - 空文字を含めない）
   ./moddict translate -mod [mod_id] -json translated.json
   ./moddict translate -mod [mod_id] -status  # 毎回確認
   （繰り返し）

5. 100%になったら報告

翻訳ルール: フォーマットコード保持、Minecraft公式用語準拠
最終報告: Mod ID、総キー数、翻訳完了数、ステータス
```

**Haikuを使う理由**: 翻訳タスクは定型的で、コスト効率が良い。

---

## CLIリファレンス

### コマンド一覧

| コマンド | 説明 |
|---------|------|
| `moddict import -jar [file]` | JARからインポート（既存ソース再利用、新バージョンをデフォルトに設定） |
| `moddict import-dir` | ディレクトリからインポート |
| `moddict translate -mod [id] -status` | 翻訳進捗確認 |
| `moddict translate -mod [id] -export [file] -limit N` | pendingをエクスポート |
| `moddict translate -mod [id] -json [file]` | 翻訳をインポート |
| `moddict translate -mod [id] -official [file]` | 公式翻訳をインポート |
| `moddict export -mod [id]` | 翻訳済みファイル出力 |
| `moddict repair` | データベース整合性の修復 |
| `moddict migrate` | スキーマ移行・バージョン情報修正 |

### 基本的な使い方

```bash
# Mod単位で翻訳を管理（バージョン指定不要）
moddict translate -mod [mod_id] -status              # 進捗確認
moddict translate -mod [mod_id] -export pending.json # pendingをエクスポート
moddict translate -mod [mod_id] -json translated.json # 翻訳をインポート
moddict export -mod [mod_id]                          # 翻訳済みファイル出力
```

**⚠️ 重要**: 翻訳結果をJSONファイルに保存しただけでは不十分。
必ず `-json` フラグでDBにインポートすること。

---

## DBスキーマ

### データ管理の原則

```
❌ NG: YAML/JSONファイルで翻訳を管理 → 整合性の問題、重複、漏れ
✅ OK: SQLite DBで一元管理 → 正確なカウント、クエリ、検証が可能
```

### テーブル構造

```
┌─────────────┐     ┌─────────────────┐     ┌──────────────────┐
│    mods     │     │  mod_versions   │     │ translation_     │
│             │     │                 │     │ sources          │
│ id (PK)     │◄────│ mod_id (FK)     │     │                  │
│ display_name│     │ id (PK)         │     │ id (PK)          │
└─────────────┘     │ version         │     │ mod_id (FK)      │
                    │ mc_version      │     │ key              │
                    │ loader          │     │ source_text      │
                    │ is_default ★    │     │ source_lang      │
                    └─────────────────┘     └──────────────────┘
                              │                      │
                              │ N:M                  │ 1:N
                              ▼                      ▼
                    ┌─────────────────┐     ┌──────────────────┐
                    │ source_versions │     │  translations    │
                    │                 │     │                  │
                    │ source_id (FK)  │     │ id (PK)          │
                    │ version_id (FK) │     │ source_id (FK)   │
                    └─────────────────┘     │ target_text      │
                                            │ target_lang      │
                                            │ status           │
                                            └──────────────────┘

★ is_default=true のバージョンが「現行バージョン」
```

### テーブル説明

| テーブル | 説明 |
|---------|------|
| `mods` | Modメタデータ |
| `mod_versions` | バージョン情報（MC版、ローダー含む） |
| `translation_sources` | ソーステキスト管理（mod_id + key + source_textでユニーク） |
| `source_versions` | ソースとバージョンのN:Mリンク |
| `translations` | 翻訳データ（source_idで紐付け） |

### スキーマの特徴

- 翻訳はMod単位で管理（バージョン非依存）
- 同じ`mod_id + key + source_text`は1つのTranslationSourceとして共有
- `source_versions`でソースとバージョンをN:Mリンク
- `mod_versions.is_default=true`で現行バージョンを指定
- クエリ時はデフォルトバージョンにリンクされたソースのみを取得

---

## 翻訳ガイドライン

### 用語の優先順位

```
global (低) → category (中) → mod (高)
```

### 必須ルール

1. **Minecraft公式訳に準拠** - バニラアイテム・ブロック名
2. **フォーマットコード保持** - `§`, `$(...)`, `%s`, `%d`
3. **改行・空白保持** - `\n`, `$(br)`, `$(br2)`

### Patchouliマクロ

| マクロ | 説明 |
|-------|------|
| `$(br)` | 改行 |
| `$(br2)` | 空行 |
| `$(item)`, `$(thing)` | ハイライト |
| `$(l:path)テキスト$(/)` | 内部リンク |

---

※ 新バージョンリリース時は `moddict import` で差分のみ追加翻訳

---

## プロジェクト構造

```
data/
├── patterns/           # ファイルパターン定義
│   ├── default.yaml    # json_lang
│   ├── patchouli.yaml  # Patchouli
│   └── mantle_book.yaml # Mantle Book
└── terms/              # 用語辞書
    ├── global.yaml     # 共通用語
    ├── category/       # カテゴリ別
    └── mod/            # Mod固有
```

---

## 注意事項

- `workspace/imports/` のJARファイルはGitにコミットしない
- 翻訳データはSQLite DB (`moddict.db`) で管理
- DBファイルのバックアップは `moddict.db.backup` として保持
- `data/` ディレクトリはパターン定義と用語辞書用
