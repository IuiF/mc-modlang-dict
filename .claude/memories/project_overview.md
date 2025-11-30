# Minecraft Mod Dictionary プロジェクト概要

## 目的

Minecraft Modの翻訳辞書とファイル構造パターンを管理するGoライブラリ。
MC Localizerと連携して、高品質で一貫性のある翻訳を実現する。

## 技術スタック

- **言語**: Go 1.21+
- **データベース**: SQLite (GORM)
- **データ形式**: YAML (辞書定義), JSON (エクスポート)
- **提供形態**: Go Library + CLI Tool

## 主要コンポーネント

### 公開API (pkg/)

- `dictionary/` - メインクライアント
  - `client.go` - Client構造体とメソッド
  - `options.go` - 設定オプション
  - `errors.go` - エラー定義

- `models/` - データモデル
  - `mod.go` - Mod情報
  - `version.go` - バージョン情報
  - `term.go` - 用語辞書
  - `translation.go` - 翻訳データ
  - `pattern.go` - ファイルパターン
  - `diff.go` - バージョン差分

- `interfaces/` - インターフェース
  - `repository.go` - データアクセス
  - `parser.go` - ファイルパーサー

### 内部実装 (internal/) - 未実装

- `database/sqlite/` - SQLite Repository
- `parser/` - 各種パーサー
- `diff/` - 差分計算
- `jar/` - JAR展開

### データ (data/)

- `patterns/global.yaml` - 標準ファイルパターン
- `terms/global.yaml` - 共通用語
- `terms/categories/` - カテゴリ別用語
  - `tech.yaml` - テクノロジー系
  - `magic.yaml` - 魔法系
  - `adventure.yaml` - 冒険系
  - `utility.yaml` - ユーティリティ系

### 作業ディレクトリ (workspace/)

- `imports/` - Modファイル格納（.gitignore対象）
- `exports/` - 出力先
- `temp/` - 一時ファイル

## 実装状況

### 完了

- [x] プロジェクト構造
- [x] Go module初期化
- [x] データモデル定義
- [x] インターフェース定義
- [x] 基本クライアントAPI
- [x] 初期用語辞書（4カテゴリ）
- [x] ドキュメント

### 未実装

- [ ] SQLite Repository
- [ ] マイグレーション
- [ ] YAML読み込み
- [ ] パーサー実装
- [ ] CLIツール
- [ ] MC Localizer統合

## 関連プロジェクト

- MC Localizer: `/home/iuif/dev/mc_localizer/`
