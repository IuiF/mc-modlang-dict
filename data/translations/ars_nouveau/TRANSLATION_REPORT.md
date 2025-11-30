# Ars Nouveau 日本語翻訳データ レポート

生成日時: 2025-11-29 02:04:57
GitHubリポジトリ: https://github.com/baileyholl/Ars-Nouveau
ブランチ: main
Minecraftバージョン: 1.21.1

## 概要

Ars Nouveau Modの言語ファイルを分析し、既存の日本語翻訳と未翻訳項目を整理しました。

### 統計サマリー

- **英語版キー総数**: 2,342
- **日本語版キー総数**: 1,480
- **翻訳済みキー数**: 1,480
- **未翻訳キー数**: 862
- **翻訳カバレッジ**: 63.2%

## カテゴリ別翻訳状況

### 完全翻訳済み（100%）のカテゴリ
- **advancements**: 80/80
- **armor**: 3/3
- **arsnouveau**: 1/1
- **attribute**: 2/2
- **biome**: 1/1
- **config**: 1/1
- **curios**: 1/1
- **enchantments**: 6/6
- **familiars**: 19/19
- **itemGroup**: 2/2
- **narrator**: 1/1
- **overworld**: 1/1
- **the_end**: 1/1
- **the_nether**: 1/1
- **whirlisprig**: 13/13

### 高翻訳率（90%以上）のカテゴリ
- **book**: 18/19 (94.7%)
- **entities**: 46/51 (90.2%)
- **glyphs**: 158/175 (90.3%)
- **key**: 17/18 (94.4%)
- **rituals**: 27/28 (96.4%)

### 中翻訳率（50-90%）のカテゴリ
- **alert**: 1/2 (50.0%)
- **blocks**: 169/195 (86.7%)
- **death**: 2/4 (50.0%)
- **items**: 220/246 (89.4%)
- **mob_jar**: 13/18 (72.2%)
- **tooltips**: 48/54 (88.9%)

### 低翻訳率（50%未満）のカテゴリ【優先翻訳推奨】
- **effects**: 16/38 (42.1%) - 未翻訳: 22件
- **emi**: 0/9 (0.0%) - 未翻訳: 9件
- **general**: 561/1349 (41.6%) - 未翻訳: 788件
- **jukebox_song**: 0/3 (0.0%) - 未翻訳: 3件

## 主要カテゴリの詳細

### グリフ（Glyphs）
Ars Nouveauの魔法システムの中核を成すグリフの翻訳状況です。

- 総数: 522
- 翻訳済み: 169
- 未翻訳: 353
- カバレッジ: 32.4%

**未翻訳グリフ例（最初の10件）:**
- `ars_nouveau.augment_desc.glyph_animate_block_glyph_duration_down`: Reduces the duration of the summon.
- `ars_nouveau.augment_desc.glyph_animate_block_glyph_extend_time`: Extends the duration of the summon.
- `ars_nouveau.augment_desc.glyph_blink_glyph_amplify`: Increases the distance of the teleport.
- `ars_nouveau.augment_desc.glyph_blink_glyph_dampen`: Decreases the distance of the teleport.
- `ars_nouveau.augment_desc.glyph_bounce_glyph_amplify`: Increases the level of the effect.
- `ars_nouveau.augment_desc.glyph_bounce_glyph_duration_down`: Reduces the duration of the effect.
- `ars_nouveau.augment_desc.glyph_bounce_glyph_extend_time`: Extends the duration of the effect.
- `ars_nouveau.augment_desc.glyph_break_glyph_amplify`: Increases the harvest level.
- `ars_nouveau.augment_desc.glyph_break_glyph_aoe`: Increases the radius of targeted blocks.
- `ars_nouveau.augment_desc.glyph_break_glyph_dampen`: Decreases the harvest level.

### Ars Nouveau固有用語の翻訳ルール

以下の用語は一貫した翻訳を使用してください：

| 英語 | 日本語 | 備考 |
|------|--------|------|
| Source | ソース | 魔力を意味する |
| Glyph | グリフ | スペルの構成要素 |
| Spell | スペル | 魔法 |
| Familiar | ファミリア | 使い魔 |
| Starbuncle | スターバンクル | 魔法生物の名前 |
| Drygmy | ドライグミー | 魔法生物の名前 |
| Whirlisprig | ワーリスプリグ | 魔法生物の名前 |
| Wixie | ウィクシー | 魔法生物の名前 |
| Archwood | アーチウッド | 魔法の木材 |
| Magebloom | メイジブルーム | 魔法の花 |
| Scribe's Table | スクライブの作業台 | |
| Enchanting Apparatus | エンチャント装置 | |
| Imbuement Chamber | 付与室 | |
| Worn Notebook | 使い古されたノート | Patchouliブック名 |

## Patchouliブック（Worn Notebook）について

Ars Nouveauには「Worn Notebook」というPatchouliブックがありますが、
翻訳は言語ファイル（en_us.json）に統合されています。

`ars_nouveau.book.*` のキーがブック関連の翻訳です。

## 出力ファイル

以下のファイルが生成されました：

1. **ars_nouveau_translation.yaml** - 完全な翻訳データベース（全カテゴリ）
2. **ars_nouveau_categories.yaml** - 主要カテゴリ別詳細データ
3. **category_summary.json** - カテゴリ統計サマリー
4. **ars_nouveau_en_us.json** - 英語言語ファイル（元データ）
5. **ars_nouveau_ja_jp.json** - 日本語言語ファイル（既存翻訳）

## 注意が必要な項目

### 1. グリフの説明文（Augment Descriptions）
グリフの強化効果の説明が大量に未翻訳です（約350件）。
ゲームプレイに重要な情報のため、優先的な翻訳を推奨します。

### 2. エフェクトの説明（Effect Descriptions）
ポーション効果の説明が約22件未翻訳です。
プレイヤーが効果を理解するために重要です。

### 3. 新規追加アイテム
Alakarkinos（カニの魔法生物）関連のアイテムが未翻訳です。
最近追加された機能と思われます。

### 4. EMI統合
EMI（Extended Mod Interface）関連の9件が未翻訳です。
レシピ表示等のUI統合に関する項目です。

## 推奨作業フロー

1. **Phase 1**: 低翻訳率カテゴリの翻訳（general, effects）
2. **Phase 2**: グリフ説明文の翻訳（augment_desc.*）
3. **Phase 3**: 新規アイテム・ブロックの翻訳
4. **Phase 4**: EMI統合等の周辺機能の翻訳

---

生成場所: /home/iuif/dev/minecraft-mod-dictionary/data/translations/ars_nouveau/
