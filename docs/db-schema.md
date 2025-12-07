# DBスキーマ

## データ管理の原則

- YAML/JSONファイルではなくSQLite DBで一元管理
- 正確なカウント、クエリ、検証が可能

## テーブル構造

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

## テーブル説明

| テーブル | 説明 |
|---------|------|
| `mods` | Modメタデータ |
| `mod_versions` | バージョン情報（MC版、ローダー含む） |
| `translation_sources` | ソーステキスト管理（mod_id + key + source_textでユニーク） |
| `source_versions` | ソースとバージョンのN:Mリンク |
| `translations` | 翻訳データ（source_idで紐付け） |

## スキーマの特徴

- 翻訳はMod単位で管理（バージョン非依存）
- 同じ`mod_id + key + source_text`は1つのTranslationSourceとして共有
- `source_versions`でソースとバージョンをN:Mリンク
- `mod_versions.is_default=true`で現行バージョンを指定
- クエリ時はデフォルトバージョンにリンクされたソースのみを取得

## バージョン更新時の動作

- 新バージョンをインポート → 既存と同じsource_textなら同じSourceを再利用
- 翻訳が自動的に継承される（継承コマンド不要）
- source_textが変わったキーのみ新規ソースを作成（pending状態）

### 継承の仕組み

```
旧バージョン (v1.0)              新バージョン (v2.0)
┌────────────────────┐          ┌────────────────────┐
│ key: item.sword    │          │ key: item.sword    │
│ source: "Sword"    │ ──同じ→  │ source: "Sword"    │
│ 翻訳: "剣" ✓       │          │ 翻訳: "剣" ✓ (継承)│
└────────────────────┘          └────────────────────┘

┌────────────────────┐          ┌────────────────────┐
│ key: item.axe      │          │ key: item.axe      │
│ source: "Axe"      │ ──変更→  │ source: "Battle Axe"│
│ 翻訳: "斧" ✓       │          │ 翻訳: pending ✗    │
└────────────────────┘          └────────────────────┘
```

### 実際の作業量

| キーの状況 | 作業量 |
|-----------|--------|
| 95%同じ | 5%だけ翻訳 → すぐ終わる |
| 80%同じ | 20%翻訳 → 少し作業 |
| 50%同じ | 50%翻訳 → それなりに作業 |
