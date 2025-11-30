# CLAUDE.md

必ず日本語で返答してください

## Project Overview

Minecraft Modの翻訳辞書とファイル構造パターンを管理するGoライブラリ。
高品質で一貫性のある日本語翻訳を実現する。

## 翻訳の基本原則

### 完全性優先（最重要）

**1つのModを100%翻訳してから次に進む。**

```
❌ NG: 6つのModを並列で60%ずつ翻訳
✅ OK: 1つのModを100%翻訳 → 次のModへ
```

### 翻訳作業の手順

1. **CurseForgeからModファイル（JAR）を取得**
   - 実際に翻訳対象となるModファイルを使用
   - 保存先: `workspace/imports/[mod_id]/`

2. **公式翻訳の確認（最重要）**
   - JARに`ja_jp.json`が含まれているか確認
   - **公式翻訳がある場合は必ず優先する**
   - `unzip -l [mod.jar] | grep ja_jp` で確認

3. **DBにインポート**
   - `moddict import` でen_us.jsonをDBに投入
   - **公式ja_jp.jsonがある場合は同時にインポート（status=verified）**
   - 公式翻訳がないキーのみpendingとして残る

4. **不足分のみ翻訳（Haikuサブエージェント）**
   - DBからpendingをクエリ（公式翻訳がないキーのみ）
   - **20件ずつ**Haikuサブエージェントで翻訳実行
   - **100件ごとにエージェントをリセット**（コンテキスト肥大化防止）
   - **並列実行可**（複数エージェントで同時処理）
   - **公式翻訳を上書きしない**

5. **検証**
   - 翻訳率100%を確認
   - フォーマットコードの保持確認

6. **次のModへ**

### 公式翻訳の尊重（必須）

```
⚠️ 公式翻訳がある場合は必ず尊重する
- JARのja_jp.jsonを優先的にインポート（translator=official, status=verified）
- 公式翻訳を勝手に上書きしない
- 不足キーのみ翻訳対象とする
- 公式にしかないキーがあれば検出漏れとして調査
```

### 翻訳エージェント実行パターン

```
┌─────────────────────────────────────────────────────────┐
│              DBから20件ずつクエリ                         │
└─────────────────────────────────────────────────────────┘
        │           │           │           │
        ▼           ▼           ▼           ▼
┌───────────┐ ┌───────────┐ ┌───────────┐ ┌───────────┐
│ Haiku #1  │ │ Haiku #2  │ │ Haiku #3  │ │ Haiku #4  │
│ 1-20件    │ │ 21-40件   │ │ 41-60件   │ │ 61-80件   │
└───────────┘ └───────────┘ └───────────┘ └───────────┘
        │           │           │           │
        ▼           ▼           ▼           ▼
┌─────────────────────────────────────────────────────────┐
│              100件完了 → エージェントリセット              │
└─────────────────────────────────────────────────────────┘
```

**Haikuを使う理由**: 翻訳タスクは定型的で、コスト効率が良い。
**20件ずつの理由**: 適度なバッチサイズで品質を維持。
**100件リセットの理由**: コンテキスト肥大化によるパフォーマンス低下を防止。

### データ管理はDBで行う（重要）

```
❌ NG: YAML/JSONファイルで翻訳を管理 → 整合性の問題、重複、漏れ
✅ OK: SQLite DBで一元管理 → 正確なカウント、クエリ、検証が可能
```

**翻訳ワークフロー:**
1. ソースファイルを静的に取得（リポジトリクローンまたはJAR）
2. `moddict import` または `moddict import-dir` でソースをDBに投入
3. `moddict translate -mod [mod_id] -export [file.json] -limit 1000` でpendingをエクスポート
4. Haikuサブエージェントで翻訳（JSON形式で出力）
5. **`moddict translate -mod [mod_id] -json [translated.json]` でDBにインポート（重要！）**
6. `moddict translate -mod [mod_id] -status` で進捗確認
7. `moddict export` で翻訳済みファイルを出力

**⚠️ 重要**: 翻訳結果をJSONファイルに保存しただけでは不十分。
必ず `-json` フラグでDBにインポートすること。

**CLIコマンド一覧:**
- `moddict import` - JARからインポート（既存ソース再利用、新バージョンをデフォルトに設定）
- `moddict import-dir` - ディレクトリからインポート
- `moddict translate` - 翻訳管理（status, export, import）
- `moddict export` - 翻訳済みファイル出力
- `moddict migrate` - スキーマ移行・バージョン情報修正
- `moddict repair` - データベース整合性の修復

**DBスキーマ:**

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

**テーブル説明:**
- `mods` - Modメタデータ
- `mod_versions` - バージョン情報（MC版、ローダー含む）
- `translation_sources` - ソーステキストのバリエーション管理（mod_id + key + source_textでユニーク）
- `source_versions` - どのバージョンでどのソースが有効か（N:M）
- `translations` - 翻訳データ（`source_id`でTranslationSourceに紐付け）
- `terms` - 用語辞書

**スキーマの特徴:**
- 翻訳はMod単位で管理（バージョン非依存）
- 同じ`mod_id + key + source_text`は1つのTranslationSourceとして共有
- `source_versions`でソースとバージョンをN:Mリンク
- `mod_versions.is_default=true`で現行バージョンを指定
- クエリ時はデフォルトバージョンにリンクされたソースのみを取得

### 翻訳ワークフロー（source_idベース）

**インポート動作:**
```
moddict import -jar [mod.jar]

1. ModVersionを作成し、is_default=trueに設定（既存のデフォルトはfalseに）
2. 各キーに対して:
   a. 既存ソース（mod_id + key + source_text）を検索
   b. あれば再利用、なければ新規作成
   c. source_versionsでバージョンとリンク
   d. 既存翻訳があれば再利用（新規作成しない）
```

**翻訳作業:**
```
# Mod単位で翻訳を管理（バージョン指定不要）
moddict translate -mod [mod_id] -status              # 進捗確認
moddict translate -mod [mod_id] -export pending.json # pendingをエクスポート
moddict translate -mod [mod_id] -json translated.json # 翻訳をインポート
moddict export -mod [mod_id]                          # 翻訳済みファイル出力
```

**バージョン更新時:**
- 新バージョンをインポート → 既存と同じsource_textなら同じSourceを再利用
- 翻訳が自動的に継承される（継承コマンド不要）
- source_textが変わったキーのみ新規ソースを作成（pending状態）

**メリット:**
- バージョン間で翻訳を自動共有（手動継承不要）
- Mod単位のシンプルなワークフロー
- ソーステキストの変更履歴を追跡可能

### 禁止事項

- サンプル翻訳だけ作成して「完了」とすること
- 翻訳率60%で次のModに進むこと
- 主要エントリだけ翻訳して残りを放置すること

---

## 翻訳エージェント構成

### サブエージェント戦略

翻訳作業では**1Modずつ確実に完了**させる。並列実行は**同一Mod内のカテゴリ分割**に限定する。

### サブエージェントの役割

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

## 対象Mod（ATM9 Tier 1）

| Mod | Mod ID | Version | MC Version | キー数 | ステータス |
|-----|--------|---------|------------|--------|----------|
| Blood Magic | `bloodmagic` | 3.0.0 | 1.16.3 | 1,378 | 100%完了 |
| Create | `create` | 0.5.1.f | 1.20.1 | 3,688 | 100%完了 |
| Botania | `botania` | 441 | 1.20.1 | 3,487 | 100%完了 |
| Mekanism | `mekanism` | 10.4.16 | 1.20.1 | 1,656 | 100%完了 |
| Tinkers' Construct | `tconstruct` | 3.10.2.92 | 1.20.1 | 2,883 | 100%完了 |
| Ars Nouveau | `ars_nouveau` | 4.10.0 | 1.20.1 | 1,545 | 100%完了 |

**合計: 14,637キー翻訳完了**

※ 新バージョンリリース時は `moddict import` で差分のみ追加翻訳

## 翻訳ガイドライン

### 用語の一貫性

```
global (低) → category (中) → mod (高)
```

### 必須ルール

1. **Minecraft公式訳に準拠** - バニラアイテム・ブロック名
2. **フォーマットコード保持** - `§`, `$(...)`, `%s`, `%d`
3. **改行・空白保持** - `\n`, `$(br)`, `$(br2)`

### Patchouliマクロ

- `$(br)` = 改行
- `$(br2)` = 空行
- `$(item)`, `$(thing)` = ハイライト
- `$(l:path)テキスト$(/)` = 内部リンク

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

## 注意事項

- `workspace/imports/` のJARファイルはGitにコミットしない
- 翻訳データはSQLite DB (`moddict.db`) で管理
- DBファイルのバックアップは `moddict.db.backup` として保持
- `data/` ディレクトリはパターン定義と用語辞書用
