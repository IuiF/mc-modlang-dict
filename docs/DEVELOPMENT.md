# 開発ガイド

## 開発環境セットアップ

### 必要環境

- Go 1.21以上
- SQLite3

### 初期セットアップ

```bash
# リポジトリクローン
git clone https://github.com/iuif/minecraft-mod-dictionary.git
cd minecraft-mod-dictionary

# 依存関係取得
go mod tidy

# ビルド確認
go build ./...

# テスト実行
go test ./...
```

## プロジェクト構成

```
minecraft-mod-dictionary/
├── pkg/                    # 公開API（外部からimportされる）
│   ├── dictionary/         # メインクライアント
│   ├── models/             # データモデル
│   └── interfaces/         # インターフェース定義
├── internal/               # 内部実装（外部非公開）
│   ├── database/           # SQLite実装
│   ├── parser/             # ファイルパーサー
│   ├── diff/               # 差分計算
│   └── jar/                # jar展開
├── data/                   # 辞書データ（YAML）
├── workspace/              # 作業ディレクトリ
├── cmd/                    # CLIツール
└── scripts/                # ビルドスクリプト
```

## 実装ロードマップ

### Phase 1: コア機能

- [ ] SQLite Repository実装
- [ ] マイグレーション
- [ ] 基本CRUD操作

### Phase 2: パーサー

- [ ] Parser Registry
- [ ] JSON Lang Parser
- [ ] Patchouli Parser
- [ ] SNBT Parser
- [ ] Legacy Lang Parser

### Phase 3: データ読み込み

- [ ] YAML→DB変換スクリプト
- [ ] パターンマッチング
- [ ] 用語マージロジック

### Phase 4: インポート/エクスポート

- [ ] JAR展開
- [ ] mod_id自動検出
- [ ] 各種形式エクスポート

### Phase 5: 差分管理

- [ ] バージョン間差分計算
- [ ] 差分適用
- [ ] 継承ベース更新

### Phase 6: CLI

- [ ] moddict build
- [ ] moddict import
- [ ] moddict export
- [ ] moddict info

## コーディング規約

### 命名規則

- **変数・関数**: snake_case（プロジェクト標準）
- **公開関数/メソッド**: PascalCase（Go標準）
- **構造体**: PascalCase（Go標準）
- **ファイル名**: snake_case

### エラーハンドリング

```go
// エラーは具体的な型で返す
func (c *Client) GetMod(ctx context.Context, id string) (*models.Mod, error) {
    mod, err := c.repo.GetMod(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrModNotFound
        }
        return nil, fmt.Errorf("failed to get mod %s: %w", id, err)
    }
    return mod, nil
}
```

### テスト

- 各パッケージに `*_test.go` を作成
- Table-Driven Tests を使用
- モックは interface を活用

```go
func TestClient_GetTerms(t *testing.T) {
    tests := []struct {
        name    string
        query   TermQuery
        want    []*models.Term
        wantErr bool
    }{
        // test cases
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

## データ追加ワークフロー

### 共通用語の追加

1. `data/terms/global.yaml` を編集
2. `go run scripts/validate.go` で検証
3. `go run scripts/build.go` でDB生成
4. PRを作成

### 新しいカテゴリの追加

1. `data/terms/categories/{category}.yaml` を作成
2. 検証・ビルド
3. PRを作成

### 新しいパーサーの追加

1. `internal/parser/{name}.go` を実装
2. `Parser` インターフェースを満たす
3. `registry.go` に登録
4. テストを追加

## デバッグ

### ログ出力

```go
import "log/slog"

slog.Info("processing mod", "mod_id", modID, "version", version)
slog.Error("failed to parse", "error", err, "file", filePath)
```

### SQLiteデバッグ

```bash
# DBファイルを直接確認
sqlite3 dictionary.db

sqlite> .tables
sqlite> SELECT * FROM mods LIMIT 5;
sqlite> SELECT * FROM terms WHERE scope = 'global' LIMIT 10;
```

## リリース手順

1. バージョンタグ作成
2. GitHub Actions で自動ビルド
3. Releases に dictionary.db を添付
4. CHANGELOG.md を更新
