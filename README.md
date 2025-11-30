# Minecraft Mod Dictionary

Minecraft Modの高品質な翻訳辞書とファイル構造パターンを管理するGoライブラリ。

## 概要

このプロジェクトは以下の問題を解決します：

- **翻訳の一貫性**: Mod間・バージョン間で統一された用語翻訳を提供
- **パターンの多様性**: 様々なModのファイル構造に対応
- **効率的な更新**: バージョン差分管理で重複作業を削減
- **再利用可能**: ライブラリとして他のツールから利用可能

## 特徴

- **汎用設計**: 任意のModに対応可能なスキーマ
- **階層型用語辞書**: 共通 → カテゴリ → Mod固有の優先度管理
- **プラグイン可能パーサー**: JSON、SNBT、Patchouli等に対応
- **差分管理**: バージョン間の効率的な差分追跡
- **LLM統合対応**: 翻訳ツールへのプロンプト生成機能

## インストール

```bash
go get github.com/iuif/minecraft-mod-dictionary
```

## 使用方法

### ライブラリとして

```go
import (
    "github.com/iuif/minecraft-mod-dictionary/pkg/dictionary"
)

func main() {
    // クライアント初期化
    client, err := dictionary.NewFromFile("./dictionary.db")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 用語取得
    terms, err := client.GetTerms(dictionary.TermQuery{
        ModID:      stringPtr("create"),
        TargetLang: "ja_jp",
        Categories: []string{"tech"},
    })

    // LLMプロンプト用フォーマット
    prompt, err := client.FormatTermsForLLM(dictionary.TermQuery{
        TargetLang: "ja_jp",
    })
}
```

### CLIツール

```bash
# データベースビルド
moddict build

# Mod情報表示
moddict info <mod_id>

# jarからインポート
moddict import ./path/to/mod.jar

# 翻訳エクスポート
moddict export <mod_id> --format json
```

## プロジェクト構造

```
minecraft-mod-dictionary/
├── pkg/                    # 公開API
│   ├── dictionary/         # メインクライアント
│   ├── models/             # データモデル
│   └── interfaces/         # インターフェース定義
├── internal/               # 内部実装
│   ├── database/           # SQLite実装
│   ├── parser/             # ファイルパーサー
│   ├── diff/               # 差分計算
│   └── jar/                # jar展開
├── data/                   # 辞書データ（YAML）
│   ├── patterns/           # ファイルパターン定義
│   ├── terms/              # 用語辞書
│   └── mods/               # Mod登録情報
├── workspace/              # 作業ディレクトリ
│   ├── imports/            # インポート対象Mod
│   ├── exports/            # エクスポート出力
│   └── temp/               # 一時ファイル
├── cmd/moddict/            # CLIツール
├── scripts/                # ビルドスクリプト
└── docs/                   # ドキュメント
```

## 開発ガイド

### 必要環境

- Go 1.21以上
- SQLite3

### ビルド

```bash
# 依存関係取得
go mod tidy

# ビルド
go build ./...

# テスト
go test ./...

# 辞書DBビルド
go run scripts/build.go
```

### 新しいModの追加

1. `workspace/imports/` にjarファイルを配置
2. `moddict import` でインポート
3. 必要に応じて `data/terms/mods/{mod_id}.yaml` を作成
4. PRを作成

詳細は [CONTRIBUTING.md](CONTRIBUTING.md) を参照。

## データ形式

### 用語辞書 (YAML)

```yaml
scope: global  # global, category:{name}, mod:{mod_id}
terms:
  - source: "Redstone"
    target: "レッドストーン"
    tags: [item, material]
    priority: 100
```

### ファイルパターン (YAML)

```yaml
patterns:
  - pattern: "assets/{mod_id}/lang/{lang}.json"
    type: lang
    parser: json_lang
    priority: 100
```

## ライセンス

MIT License

## 関連プロジェクト

- [MC Localizer](https://github.com/iuif/mc_localizer) - Minecraft翻訳ツール
