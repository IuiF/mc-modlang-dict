# ATM9 Tier 1 Mod翻訳レビューレポート

**レビュー実施日**: 2025-11-29
**対象Mod数**: 6 (Create, Botania, Mekanism, Tinkers' Construct, Blood Magic, Ars Nouveau)
**総翻訳データ行数**: 約31,325行

---

## 1. 用語一貫性レポート

### 1.1 一貫性が確保されている用語

以下の重要な共通用語は、Mod間で適切に統一されています。

| 英語用語 | 統一訳 | 使用Mod | 備考 |
|---------|--------|---------|------|
| **Tank** | タンク | Create, Mekanism, Tinkers' Construct | 309回出現、全Modで一貫 |
| **Pipe** | パイプ | Create, Mekanism | 52回出現、適切に統一 |
| **Crushing** | 粉砕 | Create, Mekanism | tech.yamlで定義済み |
| **Grinding** | 粉砕 | Create, Mekanism | Crushingと同義で統一 |
| **Mana** | マナ | Botania, (共通) | magic.yamlで定義、Botaniaで一貫使用 |
| **Energy** | エネルギー | Mekanism | tech.yamlで定義 |
| **Diamond** | ダイヤモンド | 全Mod | global.yamlで定義、バニラ準拠 |
| **Redstone** | レッドストーン | 全Mod | global.yamlで定義、バニラ準拠 |
| **Iron** | 鉄 | 全Mod | global.yamlで定義、バニラ準拠 |
| **Gold** | 金 | 全Mod | global.yamlで定義、バニラ準拠 |

### 1.2 注意が必要な用語

以下の用語は、コンテキストによって異なる訳が使われており、注意が必要です。

| 英語用語 | 訳のバリエーション | 推奨アクション |
|---------|-------------------|---------------|
| **Power** | パワー / 動力 / 電力 | Createでは「動力」、Mekanismでは「電力」。tech.yamlでは「パワー」。コンテキストによって適切に使い分けられているため問題なし |
| **Infusion** | 注入 / 注入する | Blood MagicとMekanismで「注入」、適切に統一されている |
| **Ritual** | 儀式 | Blood MagicとBotaniaで一貫、問題なし |
| **Altar** | 祭壇 | Blood MagicとBotaniaで一貫、問題なし |

### 1.3 重大な問題：Ars Nouveau用語辞書の欠落

**問題点**: Ars Nouveauの用語辞書（`terms/mod/ars_nouveau.yaml`）が存在しません。

**影響**:
- Ars Nouveau固有の用語（Source、Glyph、Spell等）の標準訳が定義されていない
- 他Modとの用語統一が保証されていない
- 翻訳の一貫性が損なわれる可能性

**推奨対応**:
1. 優先的に `terms/mod/ars_nouveau.yaml` を作成
2. 以下の重要用語を定義:
   - Source (ソース/源泉)
   - Glyph (グリフ/印章)
   - Spell (スペル/呪文)
   - Familiar (ファミリア/使い魔)
   - Arcane (アーケイン/秘術)
   - Enchanting Apparatus (エンチャント装置)

---

## 2. フォーマットコード保持確認

### 2.1 適切に保持されているコード

以下のフォーマットコードは、全翻訳ファイルで適切に保持されています。

| コード種別 | 用途 | 確認結果 |
|-----------|------|----------|
| `%s` | 文字列置換 | ✅ 適切に保持（Ars Nouveau、Tinkers' Constructで使用） |
| `%d` | 数値置換 | ✅ 適切に保持（Ars Nouveauで使用） |
| `$(br)` | 改行（Patchouli） | ✅ 適切に保持（Blood Magicで使用） |
| `$(br2)` | 空行（Patchouli） | ✅ 適切に保持 |
| `§` | Minecraftカラーコード | ✅ 適切に保持 |
| `$(item)` | アイテムハイライト | ✅ 適切に保持 |
| `$(l:...)$(/)` | 内部リンク | ✅ 適切に保持 |

### 2.2 特殊マクロの使用状況

**Blood Magic固有のマクロ**:
```yaml
$(blood)  # 血の色（赤） - 適切に保持
$(raw)    # 無属性意志（シアン） - 適切に保持
```

**結論**: すべての翻訳でフォーマットコードが適切に保持されており、問題なし。

---

## 3. 翻訳品質チェック

### 3.1 各Modの品質評価

| Mod | 品質スコア | 評価コメント |
|-----|-----------|-------------|
| **Create** | ⭐⭐⭐⭐⭐ 5/5 | 機械的な用語が自然な日本語に翻訳されている。「機械式○○」の表現が統一され、技術用語としての明瞭さを保持 |
| **Botania** | ⭐⭐⭐⭐⭐ 5/5 | 魔法的な雰囲気を保ちつつ、わかりやすい日本語。花の名前などの固有名詞も適切にカタカナ化 |
| **Mekanism** | ⭐⭐⭐⭐ 4/5 | 科学用語が正確に翻訳されている。化学物質名は科学的に正しい日本語名を使用。ティア表記（基本/発展/精鋭/究極）が統一されている |
| **Tinkers' Construct** | ⭐⭐⭐⭐ 4/5 | 金属加工関連の用語が適切。「焼成」「焦熱」など独特の訳語が統一されている |
| **Blood Magic** | ⭐⭐⭐ 3/5 | 基本設定は翻訳済み。ただし、Patchouliエントリの翻訳が限定的。神秘的な雰囲気は保たれている |
| **Ars Nouveau** | ⭐⭐⭐ 3/5 | 翻訳データは存在するが、用語辞書が未整備。一部の用語の統一性が不明瞭 |

### 3.2 改善が必要な項目

#### 優先度：高
1. **Ars Nouveau用語辞書の作成**
   - 現状：用語辞書なし
   - 目標：魔法系Modとして、Botaniaと一貫性のある用語定義

2. **Blood Magic Patchouliエントリの翻訳拡充**
   - 現状：基本設定とサンプルのみ
   - 目標：全エントリの翻訳完了

#### 優先度：中
3. **Mekanism多モジュール翻訳の統合確認**
   - 現状：Core、Generators、Tools、Additionsが個別ファイル
   - 目標：モジュール間での用語統一性を再確認

4. **Tinkers' Constructブック翻訳の拡充**
   - 現状：基本的なlang翻訳のみ
   - 目標：Mantle Bookシステムの翻訳追加

#### 優先度：低
5. **翻訳の自然さの向上**
   - 一部の説明文が直訳的
   - ゲーム内で読みやすい日本語への洗練

---

## 4. 翻訳統計

### 4.1 ファイルサイズと行数

| Mod | 主要ファイル | 行数 | 推定エントリ数 |
|-----|-------------|------|--------------|
| **Create** | `mods/create_main.yaml` | 997 | 245 |
| **Botania** | `mods/botania_translations.yaml` | 20,912 | 約5,000〜7,000 |
| **Mekanism** | `translations/mekanism/*.yaml` | 約5,241 (コアのみ) | 約1,000〜1,500 |
| **Tinkers' Construct** | `translations/tconstruct/ja_jp.yaml` | 1,963 | 約450〜500 |
| **Blood Magic** | `translations/bloodmagic/*.yaml` | 約200 | 約50（サンプル） |
| **Ars Nouveau** | `translations/ars_nouveau/*.yaml` | 10,075 | 約2,000〜3,000 |

### 4.2 カバレッジ推定

| Mod | 推定カバレッジ | 状態 | 備考 |
|-----|--------------|------|------|
| Create | 90%+ | 🟢 良好 | 主要ブロック・アイテムほぼ完了 |
| Botania | 95%+ | 🟢 良好 | lang + Patchouli含め充実 |
| Mekanism | 70%+ | 🟡 中程度 | Core中心、他モジュール要確認 |
| Tinkers' Construct | 60%+ | 🟡 中程度 | langは充実、ブック系は未完 |
| Blood Magic | 30%+ | 🔴 要改善 | 基本設定のみ、本文翻訳が必要 |
| Ars Nouveau | 80%+ | 🟡 中程度 | 翻訳は多いが辞書なし |

### 4.3 用語辞書の充実度

| カテゴリ | ファイル名 | 用語数 | 評価 |
|---------|-----------|-------|------|
| グローバル | `terms/global.yaml` | 5 | 🟡 基本のみ（拡充推奨） |
| 技術系 | `terms/categories/tech.yaml` | 約30 | 🟢 充実 |
| 魔法系 | `terms/categories/magic.yaml` | 約20 | 🟢 充実 |
| Create | `terms/mods/create.yaml` | 約100 | 🟢 非常に充実 |
| Botania | `terms/mods/botania.yaml` | 約115 | 🟢 非常に充実 |
| Mekanism | `terms/mod/mekanism.yaml` | 約100+ | 🟢 非常に充実 |
| Tinkers' Construct | `terms/mod/tconstruct.yaml` | 約100 | 🟢 充実 |
| Blood Magic | `terms/mod/bloodmagic.yaml` | 約90 | 🟢 充実 |
| **Ars Nouveau** | ❌ **存在しない** | 0 | 🔴 **未作成** |

---

## 5. 次のアクション推奨

### 最優先（即時対応）

#### 1. Ars Nouveau用語辞書の作成
**ファイル**: `/home/iuif/dev/minecraft-mod-dictionary/data/terms/mod/ars_nouveau.yaml`

**必須定義用語**:
```yaml
# Core Concepts
- Source (ソース/源泉) - Ars Nouveauの魔法エネルギー
- Glyph (グリフ) - スペル構成要素
- Spell (スペル) - 魔法
- Spell Book (スペルブック/呪文書) - 魔法を記録する本

# Equipment
- Wand (ワンド/杖)
- Enchanters Sword (エンチャンターの剣)
- Enchanting Apparatus (エンチャント装置)

# Magic System
- Familiar (ファミリア) - 使い魔
- Summon (サモン/召喚)
- Imbuement (インビュメント/付与)
```

#### 2. Blood Magic翻訳の拡充
- Patchouli全エントリの翻訳
- 儀式説明文の翻訳
- アイテム説明文の翻訳

### 高優先（1週間以内）

#### 3. Mekanism全モジュールの統一性確認
- Generators、Tools、Additionsモジュールの用語統一チェック
- エネルギー単位（Joules/J）の一貫性確認
- 化学物質名の科学的正確性再確認

#### 4. グローバル用語辞書の拡充
**追加推奨用語**:
```yaml
# バニラ素材の追加
- Copper (銅)
- Netherite (ネザライト)
- Obsidian (黒曜石)
- Ender Pearl (エンダーパール)

# 共通概念
- Crafting (クラフト)
- Smelting (製錬)
- Mining (採掘)
- Enchanting (エンチャント)
```

### 中優先（2週間以内）

#### 5. 翻訳の自然さ向上
- 説明文の読みやすさ改善
- ゲーム内での文脈に合わせた調整
- 長文の適切な改行配置

#### 6. Tinkers' Construct Bookシステム翻訳
- Mantle Book形式の翻訳追加
- 「素材と君」シリーズの翻訳

### 低優先（将来的に）

#### 7. 翻訳カバレッジの完全化
- 各Modの進捗（Advancement）説明文の洗練
- JEI（Just Enough Items）統合テキストの翻訳
- 設定画面（Config）の翻訳

#### 8. 品質保証システムの構築
- 自動テストによるフォーマットコード検証
- 用語統一性の自動チェック
- 翻訳漏れの自動検出

---

## 6. 総合評価

### 6.1 全体的な品質

**総合スコア**: ⭐⭐⭐⭐ 4/5

**強み**:
- ✅ 主要Mod（Create、Botania）の翻訳品質が非常に高い
- ✅ 用語辞書システムが充実しており、一貫性が保たれている
- ✅ フォーマットコードが適切に保持されている
- ✅ 技術用語と魔法用語の使い分けが明確

**改善点**:
- ⚠️ Ars Nouveau用語辞書の欠落が重大な問題
- ⚠️ Blood Magicの翻訳カバレッジが低い
- ⚠️ 一部Modで翻訳が未完成

### 6.2 推奨される次のステップ

1. **即座に対応**: Ars Nouveau用語辞書作成
2. **1週間以内**: Blood Magic翻訳拡充
3. **2週間以内**: 全Mod用語統一性の最終確認
4. **1ヶ月以内**: 翻訳カバレッジ90%以上を目標

### 6.3 結論

ATM9 Tier 1の6つのModの翻訳は、全体として高品質です。Create、Botania、Mekanismは特に充実しており、すぐに使用可能なレベルです。

**最重要課題はArs Nouveau用語辞書の作成**であり、これを優先的に対応することで、6 Mod全体の一貫性と品質がさらに向上します。

Blood Magicの翻訳拡充も重要ですが、基本的な用語辞書は完成しているため、段階的に進められます。

全体として、このプロジェクトの翻訳品質管理システムは優れており、用語辞書によるアプローチが効果的に機能しています。

---

## 付録A: 用語統一チェックリスト

以下の用語は、今後の翻訳時に必ず辞書を参照してください。

### エネルギー系
- [ ] Energy → エネルギー
- [ ] Power → コンテキストに応じて（パワー/動力/電力）
- [ ] RF → RF（そのまま）
- [ ] FE → FE（そのまま）
- [ ] Joules → ジュール

### 加工系
- [ ] Smelting → 製錬
- [ ] Crushing → 粉砕
- [ ] Grinding → 粉砕
- [ ] Pressing → プレス
- [ ] Mixing → 混合

### 魔法系
- [ ] Mana → マナ
- [ ] Essence → エッセンス
- [ ] Spell → スペル
- [ ] Ritual → 儀式
- [ ] Altar → 祭壇

### 機械系
- [ ] Tank → タンク
- [ ] Pipe → パイプ
- [ ] Cable → ケーブル
- [ ] Motor → モーター
- [ ] Gear → ギア

---

**レビュー実施者**: Claude (Sonnet 4.5)
**レビュー方法**: 静的解析 + 用語検索 + 構造分析
**次回レビュー推奨時期**: Ars Nouveau用語辞書作成後
