# 実装ロードマップ

## Phase 1: コア機能

### 1.1 SQLite Repository

```
internal/database/sqlite/
├── repository.go      # Repository実装
├── migrations.go      # マイグレーション
└── queries.go         # クエリヘルパー
```

**実装順序**:
1. GORM接続・初期化
2. AutoMigrate
3. Mod CRUD
4. Version CRUD
5. Term CRUD
6. Translation CRUD
7. Pattern CRUD
8. Diff CRUD

### 1.2 YAML読み込み

```
scripts/
├── build.go           # YAML→DB変換
└── validate.go        # データ検証
```

**実装順序**:
1. YAML構造体定義
2. ファイル読み込み
3. バリデーション
4. DB書き込み

## Phase 2: パーサー

### 2.1 Parser Registry

```go
type registry struct {
    parsers map[string]Parser
}

func NewRegistry() *registry {
    r := &registry{parsers: make(map[string]Parser)}
    r.Register("json_lang", &JSONLangParser{})
    r.Register("snbt", &SNBTParser{})
    r.Register("patchouli", &PatchouliParser{})
    return r
}
```

### 2.2 個別パーサー

| パーサー | 対象ファイル | 優先度 |
|----------|-------------|--------|
| json_lang | lang/*.json | 高 |
| snbt | *.snbt (FTB Quests) | 高 |
| patchouli | patchouli_books/**/*.json | 中 |
| lang_legacy | lang/*.lang (1.12-) | 低 |

## Phase 3: インポート/エクスポート

### 3.1 JAR展開

```go
func ExtractJar(jarPath, destDir string) error
func DetectModID(jarPath string) (string, error)
func ListTranslatableFiles(jarPath string, patterns []*FilePattern) ([]string, error)
```

### 3.2 エクスポート

```go
type ExportFormat string
const (
    FormatJSON ExportFormat = "json"
    FormatYAML ExportFormat = "yaml"
    FormatLang ExportFormat = "lang"  // Minecraft形式
    FormatCSV  ExportFormat = "csv"
)
```

## Phase 4: CLI

```bash
moddict build              # YAML→DB
moddict import <jar>       # Modインポート
moddict export <mod_id>    # 翻訳エクスポート
moddict info <mod_id>      # Mod情報表示
moddict diff <v1> <v2>     # バージョン差分
moddict validate           # データ検証
```

## Phase 5: MC Localizer統合

### 5.1 依存追加

```go
// mc_localizer/wails-app/go.mod
require github.com/iuif/minecraft-mod-dictionary v0.1.0
```

### 5.2 サービス実装

```go
// wails-app/internal/services/v2/translation_dictionary_service.go
type TranslationDictionaryService struct {
    client *dictionary.Client
}
```

### 5.3 翻訳フロー統合

1. 辞書から既存翻訳を検索
2. 用語辞書をLLMプロンプトに含める
3. 翻訳結果を辞書にフィードバック（オプション）

## テスト戦略

### ユニットテスト

- `pkg/models/*_test.go` - モデルテスト
- `internal/database/*_test.go` - リポジトリテスト
- `internal/parser/*_test.go` - パーサーテスト

### 統合テスト

- `tests/integration/` - E2Eテスト
- テスト用SQLite（in-memory）

### テストデータ

- `testdata/` - テスト用ファイル
  - `sample.jar` - テスト用Mod
  - `lang.json` - 言語ファイル
  - `terms.yaml` - 用語辞書
