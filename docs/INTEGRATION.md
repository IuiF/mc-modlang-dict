# MC Localizer統合ガイド

## 概要

このドキュメントはminecraft-mod-dictionaryをMC Localizerに統合する方法を説明します。

## インストール

```bash
# MC Localizerプロジェクトで
go get github.com/iuif/minecraft-mod-dictionary
```

## 基本的な使用方法

### クライアント初期化

```go
import (
    "github.com/iuif/minecraft-mod-dictionary/pkg/dictionary"
    "github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
    // 実際のrepository実装をimport
)

func NewDictionaryClient(dbPath string) (*dictionary.Client, error) {
    // repository実装を初期化
    repo, err := sqlite.NewRepository(dbPath)
    if err != nil {
        return nil, err
    }

    return dictionary.New(repo,
        dictionary.WithTargetLang("ja_jp"),
        dictionary.WithGlobalTerms(true),
    )
}
```

### サービス統合例

```go
// wails-app/internal/services/v2/translation_dictionary_service.go

package v2

import (
    "context"

    "github.com/iuif/minecraft-mod-dictionary/pkg/dictionary"
    "github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

type TranslationDictionaryService struct {
    common.ServiceBase
    client *dictionary.Client
}

func NewTranslationDictionaryService(
    dictClient *dictionary.Client,
    errorHandler *common.ErrorHandler,
    logger *common.Logger,
) *TranslationDictionaryService {
    return &TranslationDictionaryService{
        ServiceBase: common.NewServiceBase(
            "TranslationDictionaryService",
            errorHandler,
            logger,
        ),
        client: dictClient,
    }
}

// GetCachedTranslation は辞書から既存翻訳を取得
func (s *TranslationDictionaryService) GetCachedTranslation(
    ctx context.Context,
    modID, mcVersion, key string,
) (*models.Translation, error) {
    // バージョンを取得
    version, err := s.client.GetVersionBySpec(ctx, modID, "", mcVersion)
    if err != nil {
        return nil, err
    }

    return s.client.GetTranslation(ctx, version.ID, key)
}

// GetTermsForLLM はLLMプロンプト用の用語辞書を取得
func (s *TranslationDictionaryService) GetTermsForLLM(
    ctx context.Context,
    modID string,
    categories []string,
) (string, error) {
    return s.client.FormatTermsForLLM(ctx, dictionary.TermQuery{
        ModID:         &modID,
        Categories:    categories,
        TargetLang:    "ja_jp",
        IncludeGlobal: true,
    })
}

// GetFilePatterns は翻訳対象ファイルパターンを取得
func (s *TranslationDictionaryService) GetFilePatterns(
    ctx context.Context,
    modID string,
) ([]*models.FilePattern, error) {
    return s.client.GetPatterns(ctx, modID)
}
```

### 翻訳フローへの統合

```go
// wails-app/internal/services/v2/translation_engine_service.go

func (s *TranslationEngineService) TranslateEntry(
    ctx context.Context,
    entry *models.TranslationEntry,
    modID string,
) (*models.TranslationEntry, error) {

    // 1. 辞書から既存翻訳を検索
    if s.dictService != nil {
        cached, err := s.dictService.GetCachedTranslation(
            ctx, modID, s.mcVersion, entry.Key,
        )
        if err == nil && cached != nil && cached.TargetText != nil {
            entry.TranslatedText = *cached.TargetText
            entry.TranslatedBy = "dictionary:" + *cached.Translator
            entry.Status = "translated"
            s.logger.Info("translation from dictionary",
                "key", entry.Key,
                "translator", *cached.Translator,
            )
            return entry, nil
        }
    }

    // 2. 用語辞書を取得
    var terms string
    if s.dictService != nil {
        terms, _ = s.dictService.GetTermsForLLM(ctx, modID, nil)
    }

    // 3. LLMプロンプトに用語辞書を含める
    prompt := s.buildPrompt(entry.OriginalText, terms)

    // 4. LLMで翻訳
    translated, err := s.llmProvider.Translate(ctx, prompt)
    if err != nil {
        return nil, err
    }

    entry.TranslatedText = translated
    entry.TranslatedBy = s.llmProvider.Name()
    entry.Status = "translated"

    return entry, nil
}

func (s *TranslationEngineService) buildPrompt(text, terms string) string {
    if terms == "" {
        return fmt.Sprintf(`
Translate the following Minecraft mod text to Japanese.
Keep formatting codes (§a, §l, etc.) and placeholders ({0}, %%s) unchanged.

Text: %s
`, text)
    }

    return fmt.Sprintf(`
Translate the following Minecraft mod text to Japanese.

## Terminology Dictionary (MUST follow these translations):
%s

## Text to translate:
%s

## Rules:
- Use the exact translations from the terminology dictionary
- Maintain formatting codes (§a, §l, etc.)
- Keep placeholders like {0}, %%s unchanged
`, terms, text)
}
```

## DIコンテナへの登録

```go
// wails-app/internal/container/container.go

func (c *Container) RegisterDictionaryServices() error {
    // 辞書DBパス
    dictDBPath := filepath.Join(c.dataDir, "dictionary.db")

    // Repository初期化
    repo, err := sqlite.NewRepository(dictDBPath)
    if err != nil {
        return err
    }

    // Client初期化
    dictClient, err := dictionary.New(repo,
        dictionary.WithTargetLang("ja_jp"),
    )
    if err != nil {
        return err
    }

    // サービス登録
    c.dictClient = dictClient
    c.dictService = v2.NewTranslationDictionaryService(
        dictClient,
        c.errorHandler,
        c.logger,
    )

    return nil
}
```

## 辞書DBの配布

### オプション1: 埋め込み

```go
import "embed"

//go:embed data/dictionary.db
var embeddedDB embed.FS

func NewEmbeddedClient() (*dictionary.Client, error) {
    // 埋め込みDBを一時ファイルに展開して使用
    // ...
}
```

### オプション2: 外部ダウンロード

```go
func EnsureDictionaryDB(dataDir string) (string, error) {
    dbPath := filepath.Join(dataDir, "dictionary.db")

    // 既存チェック
    if _, err := os.Stat(dbPath); err == nil {
        return dbPath, nil
    }

    // GitHubリリースからダウンロード
    url := "https://github.com/iuif/minecraft-mod-dictionary/releases/latest/download/dictionary.db"
    // ダウンロード処理...

    return dbPath, nil
}
```

### オプション3: 起動時同期

```go
func SyncDictionaryDB(dataDir string) error {
    // リモートのバージョンをチェック
    // 新しいバージョンがあればダウンロード
    // ...
}
```

## Mod検出との連携

```go
// jarファイルからmod_idを検出してパターンを取得
func (s *ModDetectionService) GetTranslationPatterns(
    ctx context.Context,
    jarPath string,
) ([]*models.FilePattern, error) {
    // 1. jarからmod_idを検出
    modID, err := s.detectModID(jarPath)
    if err != nil {
        // 不明なModでもグローバルパターンは使用可能
        modID = ""
    }

    // 2. パターン取得
    patterns, err := s.dictService.GetFilePatterns(ctx, modID)
    if err != nil {
        return nil, err
    }

    return patterns, nil
}
```

## エラーハンドリング

```go
import "github.com/iuif/minecraft-mod-dictionary/pkg/dictionary"

func (s *TranslationService) handleDictionaryError(err error) {
    switch {
    case errors.Is(err, dictionary.ErrModNotFound):
        // Mod未登録 - グローバル辞書のみ使用
        s.logger.Warn("mod not in dictionary, using global terms only")
    case errors.Is(err, dictionary.ErrVersionNotFound):
        // バージョン未登録 - 最新バージョンにフォールバック
        s.logger.Warn("version not in dictionary, falling back to latest")
    default:
        s.logger.Error("dictionary error", "error", err)
    }
}
```
