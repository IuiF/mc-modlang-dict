package database

import (
	"context"
	"testing"

	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

func TestNewRepository(t *testing.T) {
	repo, err := NewRepository(":memory:")
	if err != nil {
		t.Fatalf("NewRepository() error = %v", err)
	}
	defer repo.Close()

	if repo == nil {
		t.Fatal("NewRepository() returned nil")
	}
}

func TestRepository_Migrate(t *testing.T) {
	repo, err := NewRepository(":memory:")
	if err != nil {
		t.Fatalf("NewRepository() error = %v", err)
	}
	defer repo.Close()

	if err := repo.Migrate(); err != nil {
		t.Errorf("Migrate() error = %v", err)
	}
}

func TestRepository_Mod_CRUD(t *testing.T) {
	repo := setupTestRepository(t)

	ctx := context.Background()

	t.Run("SaveMod", func(t *testing.T) {
		mod := &models.Mod{
			ID:          "create",
			DisplayName: "Create",
			Author:      "simibubi",
			Description: "A mod about building",
			Tags:        []string{"tech", "automation"},
		}

		if err := repo.SaveMod(ctx, mod); err != nil {
			t.Errorf("SaveMod() error = %v", err)
		}
	})

	t.Run("GetMod", func(t *testing.T) {
		got, err := repo.GetMod(ctx, "create")
		if err != nil {
			t.Errorf("GetMod() error = %v", err)
			return
		}

		if got.ID != "create" {
			t.Errorf("GetMod() ID = %v, want %v", got.ID, "create")
		}
		if got.DisplayName != "Create" {
			t.Errorf("GetMod() DisplayName = %v, want %v", got.DisplayName, "Create")
		}
		if got.Author != "simibubi" {
			t.Errorf("GetMod() Author = %v, want %v", got.Author, "simibubi")
		}
	})

	t.Run("GetMod_NotFound", func(t *testing.T) {
		_, err := repo.GetMod(ctx, "nonexistent")
		if err == nil {
			t.Error("GetMod() expected error for nonexistent mod")
		}
	})

	t.Run("ListMods", func(t *testing.T) {
		// Add another mod
		mod2 := &models.Mod{
			ID:          "botania",
			DisplayName: "Botania",
			Author:      "Vazkii",
			Tags:        []string{"magic"},
		}
		if err := repo.SaveMod(ctx, mod2); err != nil {
			t.Fatalf("SaveMod() error = %v", err)
		}

		mods, err := repo.ListMods(ctx, interfaces.ModFilter{})
		if err != nil {
			t.Errorf("ListMods() error = %v", err)
			return
		}

		if len(mods) != 2 {
			t.Errorf("ListMods() got %d mods, want 2", len(mods))
		}
	})

	t.Run("UpdateMod", func(t *testing.T) {
		mod := &models.Mod{
			ID:          "create",
			DisplayName: "Create Mod",
			Author:      "simibubi",
			Description: "Updated description",
		}

		if err := repo.SaveMod(ctx, mod); err != nil {
			t.Errorf("SaveMod() update error = %v", err)
		}

		got, err := repo.GetMod(ctx, "create")
		if err != nil {
			t.Errorf("GetMod() error = %v", err)
			return
		}

		if got.DisplayName != "Create Mod" {
			t.Errorf("GetMod() DisplayName = %v, want %v", got.DisplayName, "Create Mod")
		}
	})

	t.Run("DeleteMod", func(t *testing.T) {
		if err := repo.DeleteMod(ctx, "botania"); err != nil {
			t.Errorf("DeleteMod() error = %v", err)
		}

		_, err := repo.GetMod(ctx, "botania")
		if err == nil {
			t.Error("GetMod() expected error after delete")
		}
	})
}

func TestRepository_Version_CRUD(t *testing.T) {
	repo := setupTestRepository(t)
	ctx := context.Background()

	// Create a mod first (foreign key)
	mod := &models.Mod{ID: "create", DisplayName: "Create"}
	if err := repo.SaveMod(ctx, mod); err != nil {
		t.Fatalf("SaveMod() error = %v", err)
	}

	var savedVersionID int64

	t.Run("SaveVersion", func(t *testing.T) {
		version := &models.ModVersion{
			ModID:     "create",
			Version:   "0.5.1",
			MCVersion: "1.20.1",
			Loader:    "forge",
			Stats: models.VersionStats{
				TotalKeys:      100,
				TranslatedKeys: 50,
			},
		}

		if err := repo.SaveVersion(ctx, version); err != nil {
			t.Errorf("SaveVersion() error = %v", err)
		}

		if version.ID == 0 {
			t.Error("SaveVersion() did not set ID")
		}
		savedVersionID = version.ID
	})

	t.Run("GetVersion", func(t *testing.T) {
		got, err := repo.GetVersion(ctx, savedVersionID)
		if err != nil {
			t.Errorf("GetVersion() error = %v", err)
			return
		}

		if got.ModID != "create" {
			t.Errorf("GetVersion() ModID = %v, want %v", got.ModID, "create")
		}
		if got.Version != "0.5.1" {
			t.Errorf("GetVersion() Version = %v, want %v", got.Version, "0.5.1")
		}
		if got.Stats.TotalKeys != 100 {
			t.Errorf("GetVersion() TotalKeys = %v, want %v", got.Stats.TotalKeys, 100)
		}
	})

	t.Run("GetVersionBySpec", func(t *testing.T) {
		got, err := repo.GetVersionBySpec(ctx, "create", "0.5.1", "1.20.1")
		if err != nil {
			t.Errorf("GetVersionBySpec() error = %v", err)
			return
		}

		if got.ID != savedVersionID {
			t.Errorf("GetVersionBySpec() ID = %v, want %v", got.ID, savedVersionID)
		}
	})

	t.Run("ListVersions", func(t *testing.T) {
		// Add another version
		version2 := &models.ModVersion{
			ModID:     "create",
			Version:   "0.5.0",
			MCVersion: "1.19.2",
			Loader:    "forge",
		}
		if err := repo.SaveVersion(ctx, version2); err != nil {
			t.Fatalf("SaveVersion() error = %v", err)
		}

		versions, err := repo.ListVersions(ctx, "create", interfaces.VersionFilter{})
		if err != nil {
			t.Errorf("ListVersions() error = %v", err)
			return
		}

		if len(versions) != 2 {
			t.Errorf("ListVersions() got %d versions, want 2", len(versions))
		}
	})

	t.Run("ListVersions_ByMCVersion", func(t *testing.T) {
		versions, err := repo.ListVersions(ctx, "create", interfaces.VersionFilter{
			MCVersion: "1.20.1",
		})
		if err != nil {
			t.Errorf("ListVersions() error = %v", err)
			return
		}

		if len(versions) != 1 {
			t.Errorf("ListVersions(mc=1.20.1) got %d versions, want 1", len(versions))
		}
	})

	t.Run("DeleteVersion", func(t *testing.T) {
		if err := repo.DeleteVersion(ctx, savedVersionID); err != nil {
			t.Errorf("DeleteVersion() error = %v", err)
		}

		_, err := repo.GetVersion(ctx, savedVersionID)
		if err == nil {
			t.Error("GetVersion() expected error after delete")
		}
	})
}

func TestRepository_Term_CRUD(t *testing.T) {
	repo := setupTestRepository(t)
	ctx := context.Background()

	t.Run("SaveTerm", func(t *testing.T) {
		term := &models.Term{
			Scope:      "global",
			SourceText: "Redstone",
			TargetText: "レッドストーン",
			SourceLang: "en_us",
			TargetLang: "ja_jp",
			Tags:       []string{"block", "material"},
			Priority:   100,
			Source:     "official",
		}

		if err := repo.SaveTerm(ctx, term); err != nil {
			t.Errorf("SaveTerm() error = %v", err)
		}

		if term.ID == 0 {
			t.Error("SaveTerm() did not set ID")
		}
	})

	t.Run("GetTerm", func(t *testing.T) {
		got, err := repo.GetTerm(ctx, 1)
		if err != nil {
			t.Errorf("GetTerm() error = %v", err)
			return
		}

		if got.SourceText != "Redstone" {
			t.Errorf("GetTerm() SourceText = %v, want %v", got.SourceText, "Redstone")
		}
		if got.TargetText != "レッドストーン" {
			t.Errorf("GetTerm() TargetText = %v, want %v", got.TargetText, "レッドストーン")
		}
	})

	t.Run("ListTerms_ByScope", func(t *testing.T) {
		// Add mod-specific term
		modTerm := &models.Term{
			Scope:      "mod:create",
			SourceText: "Mechanical Bearing",
			TargetText: "メカニカルベアリング",
			TargetLang: "ja_jp",
			Priority:   200,
		}
		if err := repo.SaveTerm(ctx, modTerm); err != nil {
			t.Fatalf("SaveTerm() error = %v", err)
		}

		// List global terms only
		terms, err := repo.ListTerms(ctx, interfaces.TermFilter{Scope: "global"})
		if err != nil {
			t.Errorf("ListTerms() error = %v", err)
			return
		}

		if len(terms) != 1 {
			t.Errorf("ListTerms(global) got %d terms, want 1", len(terms))
		}
	})

	t.Run("ListTerms_MultipleScopes", func(t *testing.T) {
		terms, err := repo.ListTerms(ctx, interfaces.TermFilter{
			Scopes: []string{"global", "mod:create"},
		})
		if err != nil {
			t.Errorf("ListTerms() error = %v", err)
			return
		}

		if len(terms) != 2 {
			t.Errorf("ListTerms(multiple) got %d terms, want 2", len(terms))
		}
	})

	t.Run("BulkSaveTerms", func(t *testing.T) {
		bulkTerms := []*models.Term{
			{Scope: "global", SourceText: "Diamond", TargetText: "ダイヤモンド", TargetLang: "ja_jp"},
			{Scope: "global", SourceText: "Emerald", TargetText: "エメラルド", TargetLang: "ja_jp"},
			{Scope: "global", SourceText: "Gold", TargetText: "金", TargetLang: "ja_jp"},
		}

		if err := repo.BulkSaveTerms(ctx, bulkTerms); err != nil {
			t.Errorf("BulkSaveTerms() error = %v", err)
		}

		terms, err := repo.ListTerms(ctx, interfaces.TermFilter{Scope: "global"})
		if err != nil {
			t.Errorf("ListTerms() error = %v", err)
			return
		}

		if len(terms) != 4 {
			t.Errorf("ListTerms() after bulk got %d terms, want 4", len(terms))
		}
	})

	t.Run("DeleteTerm", func(t *testing.T) {
		if err := repo.DeleteTerm(ctx, 1); err != nil {
			t.Errorf("DeleteTerm() error = %v", err)
		}

		_, err := repo.GetTerm(ctx, 1)
		if err == nil {
			t.Error("GetTerm() expected error after delete")
		}
	})
}

func TestRepository_Translation_CRUD(t *testing.T) {
	repo := setupTestRepository(t)
	ctx := context.Background()

	// Setup: create mod and version
	mod := &models.Mod{ID: "create", DisplayName: "Create"}
	repo.SaveMod(ctx, mod)
	version := &models.ModVersion{ModID: "create", Version: "0.5.1", MCVersion: "1.20.1"}
	repo.SaveVersion(ctx, version)

	var savedID int64

	t.Run("SaveTranslation", func(t *testing.T) {
		targetText := "メカニカルベアリング"
		trans := &models.Translation{
			ModVersionID: version.ID,
			Key:          "item.create.mechanical_bearing",
			SourceText:   "Mechanical Bearing",
			TargetText:   &targetText,
			SourceLang:   "en_us",
			TargetLang:   "ja_jp",
			Status:       models.StatusTranslated,
		}

		if err := repo.SaveTranslation(ctx, trans); err != nil {
			t.Errorf("SaveTranslation() error = %v", err)
		}
		savedID = trans.ID
	})

	t.Run("GetTranslation", func(t *testing.T) {
		got, err := repo.GetTranslation(ctx, version.ID, "item.create.mechanical_bearing")
		if err != nil {
			t.Errorf("GetTranslation() error = %v", err)
			return
		}

		if got.SourceText != "Mechanical Bearing" {
			t.Errorf("GetTranslation() SourceText = %v, want %v", got.SourceText, "Mechanical Bearing")
		}
	})

	t.Run("ListTranslations", func(t *testing.T) {
		// Add more translations
		texts := []string{"Gearbox", "Cogwheel"}
		for _, text := range texts {
			repo.SaveTranslation(ctx, &models.Translation{
				ModVersionID: version.ID,
				Key:          "item.create." + text,
				SourceText:   text,
				TargetLang:   "ja_jp",
				Status:       models.StatusPending,
			})
		}

		translations, err := repo.ListTranslations(ctx, version.ID, interfaces.TranslationFilter{})
		if err != nil {
			t.Errorf("ListTranslations() error = %v", err)
			return
		}

		if len(translations) != 3 {
			t.Errorf("ListTranslations() got %d, want 3", len(translations))
		}
	})

	t.Run("ListTranslations_ByStatus", func(t *testing.T) {
		translations, err := repo.ListTranslations(ctx, version.ID, interfaces.TranslationFilter{
			Status: models.StatusPending,
		})
		if err != nil {
			t.Errorf("ListTranslations() error = %v", err)
			return
		}

		if len(translations) != 2 {
			t.Errorf("ListTranslations(pending) got %d, want 2", len(translations))
		}
	})

	t.Run("BulkSaveTranslations", func(t *testing.T) {
		bulk := []*models.Translation{
			{ModVersionID: version.ID, Key: "block.create.gearbox", SourceText: "Gearbox", TargetLang: "ja_jp"},
			{ModVersionID: version.ID, Key: "block.create.shaft", SourceText: "Shaft", TargetLang: "ja_jp"},
		}

		if err := repo.BulkSaveTranslations(ctx, bulk); err != nil {
			t.Errorf("BulkSaveTranslations() error = %v", err)
		}
	})

	t.Run("DeleteTranslation", func(t *testing.T) {
		if err := repo.DeleteTranslation(ctx, savedID); err != nil {
			t.Errorf("DeleteTranslation() error = %v", err)
		}
	})
}

func TestRepository_Pattern_CRUD(t *testing.T) {
	repo := setupTestRepository(t)
	ctx := context.Background()

	var savedID int64

	t.Run("SavePattern", func(t *testing.T) {
		desc := "Standard Minecraft lang file"
		pattern := &models.FilePattern{
			Scope:       "global",
			Pattern:     "assets/{mod_id}/lang/{lang}.json",
			Type:        models.PatternTypeLang,
			Parser:      models.ParserJSONLang,
			Priority:    100,
			Required:    true,
			Description: &desc,
		}

		if err := repo.SavePattern(ctx, pattern); err != nil {
			t.Errorf("SavePattern() error = %v", err)
		}
		savedID = pattern.ID
	})

	t.Run("GetPattern", func(t *testing.T) {
		got, err := repo.GetPattern(ctx, savedID)
		if err != nil {
			t.Errorf("GetPattern() error = %v", err)
			return
		}

		if got.Pattern != "assets/{mod_id}/lang/{lang}.json" {
			t.Errorf("GetPattern() Pattern = %v", got.Pattern)
		}
	})

	t.Run("ListPatterns", func(t *testing.T) {
		// Add mod-specific pattern
		repo.SavePattern(ctx, &models.FilePattern{
			Scope:   "mod:patchouli",
			Pattern: "data/{mod_id}/patchouli_books/**/*.json",
			Type:    models.PatternTypeBook,
			Parser:  models.ParserPatchouli,
		})

		patterns, err := repo.ListPatterns(ctx, "global")
		if err != nil {
			t.Errorf("ListPatterns() error = %v", err)
			return
		}

		if len(patterns) != 1 {
			t.Errorf("ListPatterns(global) got %d, want 1", len(patterns))
		}
	})

	t.Run("DeletePattern", func(t *testing.T) {
		if err := repo.DeletePattern(ctx, savedID); err != nil {
			t.Errorf("DeletePattern() error = %v", err)
		}
	})
}

func TestRepository_Diff_CRUD(t *testing.T) {
	repo := setupTestRepository(t)
	ctx := context.Background()

	// Setup versions
	mod := &models.Mod{ID: "create", DisplayName: "Create"}
	repo.SaveMod(ctx, mod)
	v1 := &models.ModVersion{ModID: "create", Version: "0.5.0", MCVersion: "1.20.1"}
	v2 := &models.ModVersion{ModID: "create", Version: "0.5.1", MCVersion: "1.20.1"}
	repo.SaveVersion(ctx, v1)
	repo.SaveVersion(ctx, v2)

	t.Run("SaveDiff", func(t *testing.T) {
		newText := "Updated Text"
		diff := &models.VersionDiff{
			FromVersionID: v1.ID,
			ToVersionID:   v2.ID,
			Type:          models.DiffTypeChanged,
			Key:           "item.create.name",
			OldText:       nil,
			NewText:       &newText,
		}

		if err := repo.SaveDiff(ctx, diff); err != nil {
			t.Errorf("SaveDiff() error = %v", err)
		}
	})

	t.Run("BulkSaveDiffs", func(t *testing.T) {
		added := "New Item"
		diffs := []*models.VersionDiff{
			{FromVersionID: v1.ID, ToVersionID: v2.ID, Type: models.DiffTypeAdded, Key: "item.create.new", NewText: &added},
			{FromVersionID: v1.ID, ToVersionID: v2.ID, Type: models.DiffTypeRemoved, Key: "item.create.old"},
		}

		if err := repo.BulkSaveDiffs(ctx, diffs); err != nil {
			t.Errorf("BulkSaveDiffs() error = %v", err)
		}
	})

	t.Run("ListDiffs", func(t *testing.T) {
		diffs, err := repo.ListDiffs(ctx, v1.ID, v2.ID)
		if err != nil {
			t.Errorf("ListDiffs() error = %v", err)
			return
		}

		if len(diffs) != 3 {
			t.Errorf("ListDiffs() got %d, want 3", len(diffs))
		}
	})
}

// setupTestRepository creates a new in-memory repository for testing.
func setupTestRepository(t *testing.T) *Repository {
	t.Helper()

	repo, err := NewRepository(":memory:")
	if err != nil {
		t.Fatalf("NewRepository() error = %v", err)
	}

	if err := repo.Migrate(); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	t.Cleanup(func() {
		repo.Close()
	})

	return repo
}
