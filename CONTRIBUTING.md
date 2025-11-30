# 貢献ガイド

Minecraft Mod Dictionaryへの貢献をありがとうございます！

## 貢献の種類

### 1. 翻訳の追加・修正

最も歓迎される貢献です。プログラミング知識は不要です。

### 2. 新しいModの対応

ファイルパターンや固有用語の追加。

### 3. コードの改善

バグ修正、新機能、パフォーマンス改善。

---

## 翻訳の追加・修正

### 共通用語の追加

`data/terms/global.yaml` を編集：

```yaml
terms:
  - source: "English Term"
    target: "日本語訳"
    tags: [item, block, ui]  # 適切なタグ
    priority: 100            # 優先度（高いほど優先）
    context: "使用される文脈の説明"  # 任意
```

### カテゴリ用語の追加

`data/terms/categories/{category}.yaml` を編集：

利用可能なカテゴリ：
- `tech.yaml` - テクノロジー系Mod
- `magic.yaml` - 魔法系Mod
- `adventure.yaml` - 冒険・RPG系Mod
- `utility.yaml` - ユーティリティ系Mod

### Mod固有用語の追加

`data/terms/mods/{mod_id}.yaml` を新規作成または編集：

```yaml
scope: "mod:{mod_id}"
terms:
  - source: "Mod Specific Term"
    target: "Mod固有の訳"
    tags: [block]
```

### 翻訳のルール

1. **公式翻訳を尊重**: Minecraft本体の公式翻訳に従う
2. **一貫性**: 同じ用語は同じ訳を使う
3. **文脈を考慮**: 同じ英語でも文脈で訳が変わる場合はcontextを記載
4. **タグを適切に**: 検索・フィルタリングに使用される

---

## 新しいModの対応

### 1. Modファイルの準備

```bash
# workspace/imports/ にjarを配置
cp /path/to/mod.jar workspace/imports/
```

### 2. ファイル構造の調査

```bash
# jarの中身を確認
unzip -l workspace/imports/mod.jar | grep -E "\.(json|lang|snbt)$"
```

### 3. パターンの追加（必要な場合）

標準パターンで対応できない場合のみ `data/patterns/overrides/{mod_id}.yaml` を作成：

```yaml
mod_id: "example_mod"
patterns:
  - pattern: "custom/path/{lang}.json"
    type: lang
    parser: json_lang
    priority: 100
    description: "This mod uses non-standard path"
```

### 4. Mod情報の登録（任意）

`data/mods/{mod_id}.yaml` を作成：

```yaml
id: "example_mod"
display_name: "Example Mod"
author: "Author Name"
tags: [tech, utility]
metadata:
  curseforge_id: "123456"
  modrinth_id: "abcdef"
```

---

## コード貢献

### 開発環境セットアップ

```bash
# リポジトリクローン
git clone https://github.com/iuif/minecraft-mod-dictionary.git
cd minecraft-mod-dictionary

# 依存関係
go mod tidy

# テスト実行
go test ./...
```

### ブランチ戦略

```
main          # 安定版
├── develop   # 開発版
└── feature/* # 機能ブランチ
```

### コミットメッセージ

```
feat: 新機能追加
fix: バグ修正
docs: ドキュメント
data: 辞書データ追加・修正
refactor: リファクタリング
test: テスト追加・修正
```

### プルリクエスト

1. `develop` から機能ブランチを作成
2. 変更を実装
3. テストを追加・実行
4. PRを作成

---

## ファイル構造リファレンス

```
data/
├── patterns/
│   ├── global.yaml           # 全Mod共通パターン
│   └── overrides/
│       └── {mod_id}.yaml     # Mod固有パターン上書き
│
├── terms/
│   ├── global.yaml           # 共通用語
│   ├── categories/
│   │   ├── tech.yaml         # テクノロジー系
│   │   ├── magic.yaml        # 魔法系
│   │   ├── adventure.yaml    # 冒険系
│   │   └── utility.yaml      # ユーティリティ系
│   └── mods/
│       └── {mod_id}.yaml     # Mod固有用語
│
└── mods/
    └── {mod_id}.yaml         # Mod情報（任意）
```

---

## 用語のタグ一覧

| タグ | 説明 |
|------|------|
| `item` | アイテム名 |
| `block` | ブロック名 |
| `entity` | エンティティ名 |
| `ui` | UI要素 |
| `tooltip` | ツールチップ |
| `advancement` | 進捗 |
| `mechanic` | ゲームメカニクス用語 |
| `material` | 素材名 |
| `dimension` | ディメンション名 |
| `biome` | バイオーム名 |
| `effect` | 効果・ステータス |
| `enchantment` | エンチャント |
| `official` | Minecraft公式翻訳 |

---

## 質問・サポート

- GitHub Issues: バグ報告、機能要望
- Discussions: 質問、議論
