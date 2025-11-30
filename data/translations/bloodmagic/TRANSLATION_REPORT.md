# Blood Magic 翻訳レポート

## プロジェクト概要

- **Mod名**: Blood Magic
- **対象バージョン**: 3.x.x
- **Minecraft バージョン**: 1.20.1
- **GitHubリポジトリ**: https://github.com/WayofTime/BloodMagic
- **ブランチ**: 1.20.1
- **作成日**: 2025-11-29

## 翻訳状況

### 完了した翻訳

| カテゴリ | 項目数 | 翻訳済み | 進捗率 |
|---------|--------|---------|--------|
| Book設定 | 1 | 1 | 100% |
| カテゴリ | 7 | 7 | 100% |
| エントリ（サンプル） | 5 | 5 | 100% |
| 用語辞書 | 約80項目 | 80 | 100% |

### 未完了の翻訳

| カテゴリ | 項目数 | 備考 |
|---------|--------|------|
| 残りのエントリファイル | 約109 | 全体で114エントリ、サンプル5を除く |
| 言語ファイル（en_us.json） | N/A | GitHubリポジトリに存在せず |

## ファイル構成

### 作成したファイル

```
data/translations/bloodmagic/
├── bloodmagic_book.yaml          # Book.json翻訳
├── bloodmagic_categories.yaml    # 全7カテゴリの翻訳
├── bloodmagic_entries_sample.yaml# 主要5エントリの翻訳サンプル
└── TRANSLATION_REPORT.md         # このレポート

data/terms/mod/
└── bloodmagic.yaml               # Blood Magic固有用語辞書
```

## 翻訳詳細

### 1. Book設定（bloodmagic_book.yaml）

Patchouliガイドブックの基本設定を翻訳。

**翻訳項目**:
- ブック名: "Sanguine Scientiem"
- サブタイトル: "あなたの生命、あなたの魔法"
- ランディングテキスト
- マクロ定義（色コード）の説明

### 2. カテゴリ（bloodmagic_categories.yaml）

全7カテゴリを完全翻訳。

**カテゴリ一覧**:
1. **Blood Altars（血の祭壇）** - sortnum: 0
   - 祭壇の構築ガイド

2. **Alchemy Table（錬金術テーブル）** - sortnum: 1
   - 錬金術システムの基礎

3. **Alchemy Arrays（錬金術陣）** - sortnum: 2
   - 地面に描く魔法陣

4. **Demon Will（悪魔の意志）** - sortnum: 3
   - 悪魔の意志システム

5. **Rituals（儀式）** - sortnum: 4
   - 儀式の構築と実行

6. **Dungeon Delving（ダンジョン探索）** - sortnum: 7
   - 悪魔の領域の探索

7. **Utility Blocks & Items（ユーティリティブロック＆アイテム）** - sortnum: 99
   - その他のアイテムと情報

### 3. エントリサンプル（bloodmagic_entries_sample.yaml）

重要な5つのエントリを翻訳。

**翻訳済みエントリ**:

1. **Tiers & Getting Started（ティアと始め方）**
   - カテゴリ: Utility
   - ページ数: 17ページ（うち5ページを翻訳）
   - 内容: Blood Magic 3の進行ガイド

2. **The Blood Altar（血の祭壇）**
   - カテゴリ: Altar
   - ページ数: 20+ページ
   - 内容: 祭壇の使い方とティアアップグレード

3. **Tiers of Slates（石板のティア）**
   - カテゴリ: Altar
   - 内容: 5段階の石板の作成

4. **Demon Will（悪魔の意志）**
   - カテゴリ: Demon Will
   - 内容: 悪魔の意志の収集方法

5. **Rituals - Getting Started（儀式 - 始め方）**
   - カテゴリ: Rituals
   - 内容: 儀式の基本と構築方法

### 4. 用語辞書（data/terms/mod/bloodmagic.yaml）

Blood Magic固有の約80項目の用語を定義。

**主要カテゴリ**:
- Core Concepts（核心概念）: LP, Life Essence等
- Structures & Blocks（構造物とブロック）: 祭壇、儀式石等
- Tools & Weapons（道具と武器）: ナイフ、剣等
- Items（アイテム）: オーブ、石板、シジル等
- Demon Will System（悪魔の意志システム）
- Alchemy（錬金術）
- Living Equipment（生きている装備）
- Rituals（儀式）
- Dimensions & Dungeons（次元とダンジョン）

## 重要な翻訳ルール

### Blood Magic固有の用語

| 英語 | 日本語 | 備考 |
|------|--------|------|
| Life Points (LP) | ライフポイント (LP) | 略称はそのまま |
| Blood Altar | 血の祭壇 | 中心的な構造物 |
| Sigil | シジル | 魔法印、カタカナ表記 |
| Demon Will | 悪魔の意志 | 重要リソース |
| Soul Network | ソウルネットワーク | システム名 |
| Ritual | 儀式 | 強力な効果 |
| Sentient Sword | 意識剣 | 意識を持つ剣 |
| Tartaric Gem | タルタロスの宝石 | ギリシャ神話由来 |
| Anointment | 膏薬 | 強化アイテム |

### Patchouliマクロの保持

翻訳時に必ず保持すべきマクロ：

- **改行**: `$(br)`, `$(br2)`
- **リンク**: `$(l:path)テキスト$(/l)`
- **強調**: `$(item)`, `$(thing)`
- **色コード**: `$(blood)`, `$(raw)`, `$(fire)`等
- **キーバインド**: `$(k:use)`, `$(k:sneak)`

### 翻訳の方針

1. **Minecraft公式訳に準拠**
   - バニラアイテム・ブロック名は公式日本語版に従う

2. **Blood Magicのテーマ性を保持**
   - 血と犠牲をテーマとした重厚で神秘的な雰囲気
   - ラテン語由来の用語は適切に意訳

3. **ゲームプレイの明確性**
   - システムやメカニクスの説明は分かりやすく
   - 手順は明確に、曖昧さを排除

4. **一貫性の維持**
   - 同じ用語は常に同じ訳語を使用
   - 用語辞書を参照して統一

## エントリファイル一覧

### 全114エントリの内訳

| カテゴリ | エントリ数 | サンプル翻訳済み |
|---------|-----------|----------------|
| Altar | 4 | 2 |
| Alchemy Table | 3 | 0 |
| Alchemy Array | 23 | 0 |
| - Functional Arrays | 5 | 0 |
| - Living Equipment | 4 | 0 |
| - Sigil | 14 | 0 |
| Demon Will | 23 | 1 |
| - Demonic Items | 4 | 0 |
| - Item Routing | 10 | 0 |
| - Will Manipulation | 9 | 1 |
| Dungeons | 11 | 0 |
| Rituals | 37 | 1 |
| - Ritual Basics | 5 | 1 |
| - Ritual List | 32 | 0 |
| Utility | 13 | 1 |
| **合計** | **114** | **5** |

### 未翻訳エントリ（優先度順）

#### 高優先度（基礎システム）

**Altar**:
- redstone_automation.json - 自動化
- soul_network.json - ソウルネットワーク

**Alchemy Table**:
- alchemy_table.json - テーブルの使い方
- anointments.json - 膏薬システム
- potions.json - ポーション作成

**Demon Will / Will Manipulation**:
- soul_snare.json - 最初の意志収集
- soul_gem.json - タルタロスの宝石
- soul_forge.json - 地獄の炉
- aspected_will.json - 属性付き意志

#### 中優先度（アイテムと機能）

**Alchemy Array / Functional Arrays**:
- arcane_ash.json - 秘術の灰
- crafting_array.json - クラフト陣
- movement_arrays.json - 移動陣
- time_arrays.json - 時間操作陣

**Alchemy Array / Sigil**（14ファイル）:
- divination.json - 占いのシジル（重要）
- water.json, lava.json, void.json等

**Rituals / Ritual List**（32ファイル）:
- ritual_well_of_suffering.json - LP自動化
- ritual_condor.json - クリエイティブ飛行
- ritual_simple_dungeon.json - ダンジョン入場
- 他29ファイル

#### 低優先度（発展システム）

**Living Equipment**:
- living_basics.json
- living_upgrades.json
- living_tomes.json
- training_bracelet.json

**Dungeons**（11ファイル）:
- 全てのダンジョン関連エントリ

## 技術的な課題と解決策

### 1. 言語ファイル（en_us.json）が存在しない

**問題**:
- GitHubリポジトリの1.20.1ブランチに言語ファイルが存在しない
- 通常のMinecraft Modには必須のファイル

**調査結果**:
- GitHub検索で確認: 言語ファイルなし
- リリースページ: リリースなし
- ブランチ: 1.20.1と1.21.1のみ

**推測される理由**:
1. 開発中のため、言語ファイルがまだコミットされていない
2. 自動生成されるため、リポジトリには含まれていない
3. ビルドプロセスで生成される

**代替アプローチ**:
- Patchouliのコンテンツのみを翻訳（現在のアプローチ）
- 将来的にリリース版JARファイルから抽出
- CurseForge/Modrinthからダウンロード

### 2. Patchouliエントリの膨大な数

**問題**:
- 合計114エントリファイル
- 各エントリに平均5-10ページ
- 総ページ数は500-1000ページ規模

**解決策**:
- 段階的な翻訳アプローチ
- 優先度を設定（基礎 → 発展 → 応用）
- サンプル翻訳で品質基準を確立

## 今後の作業計画

### Phase 1: 基礎システム翻訳（推定: 20-30エントリ）

1. Altarカテゴリの残り2エントリ
2. Alchemy Tableカテゴリ全3エントリ
3. Demon Will / Will Manipulation の残り8エントリ
4. 基本的なSigil（5-7個）
5. 基本的なRitual（5-7個）

### Phase 2: 機能拡張翻訳（推定: 40-50エントリ）

1. Alchemy Array / Functional Arrays 全5エントリ
2. Alchemy Array / Sigil 全14エントリ
3. Rituals / Ritual List の主要儀式（15-20個）
4. Demon Will / Demonic Items 全4エントリ

### Phase 3: 発展システム翻訳（推定: 40-50エントリ）

1. Living Equipment 全4エントリ
2. Item Routing 全10エントリ
3. Dungeons 全11エントリ
4. Utility 残り12エントリ
5. 残りのRitual List（10-15個）

### Phase 4: 言語ファイル統合

1. リリース版から言語ファイル抽出
2. アイテム名・ブロック名の翻訳
3. UI要素の翻訳
4. 統合テストと品質確認

## 注意が必要な項目

### 1. 複雑なマクロ構文

一部のエントリには複雑なPatchouliマクロが使用されている：

```
$(l:bloodmagic:path)$(item)テキスト$()$(/l)
```

翻訳時には：
- マクロの入れ子構造を保持
- 閉じタグの位置を確認
- リンクパスは変更しない

### 2. 儀式名の詩的表現

Blood Magicの儀式名は詩的で抽象的：
- "Edge of the Hidden Realm" → "隠された領域の境界"
- "Pathway to the Endless Realm" → "無限の領域への道"
- "Well of Suffering" → "苦しみの井戸"
- "Sound of the Cleansing Soul" → "浄化された魂の音"

これらは直訳ではなく、雰囲気を保った意訳が必要。

### 3. Living Armourのアップグレード名

Living Armourには多数のアップグレード/ダウングレードがあり、それぞれ固有の名前と効果を持つ。一貫性を保ちながら、効果が分かる名前にする必要がある。

例：
- Body Builder → ボディビルダー（筋力増強）
- Leadened Pick → 鉛の鶴嘴（採掘速度低下）
- Limp Leg → 不自由な脚（移動速度低下）

## 品質保証

### 翻訳品質チェックリスト

- [ ] Blood Magic用語辞書との整合性
- [ ] Patchouliマクロの完全保持
- [ ] リンクパスの正確性
- [ ] 日本語の自然さ
- [ ] ゲームプレイの理解しやすさ
- [ ] 既存翻訳（Botania, Ars Nouveau）との一貫性

### 推奨される検証方法

1. **マクロ検証**: 正規表現でマクロの開始/終了タグを検証
2. **用語統一**: 用語辞書に基づく用語の一貫性チェック
3. **実ゲームテスト**: 実際にMinecraftで表示を確認
4. **ピアレビュー**: 他の翻訳者による確認

## まとめ

### 達成したこと

✅ Blood Magic の基本構造を理解
✅ Patchouliガイドブックの全体像を把握（7カテゴリ、114エントリ）
✅ Book設定と全カテゴリの翻訳完了
✅ 重要な5エントリのサンプル翻訳完成
✅ 包括的な用語辞書の作成（約80項目）
✅ 翻訳ガイドラインとルールの確立

### 今後の課題

⚠️ 残り109エントリの翻訳（全体の95%）
⚠️ 言語ファイル（en_us.json）の取得と翻訳
⚠️ 実ゲームでの動作確認
⚠️ 翻訳の品質レビューと改善

### 推定作業量

- **Phase 1** (基礎): 約15-20時間
- **Phase 2** (機能): 約25-30時間
- **Phase 3** (発展): 約25-30時間
- **Phase 4** (統合): 約10-15時間
- **合計**: 約75-95時間

## ファイルパス

すべての翻訳ファイルは以下に格納：

```
/home/iuif/dev/minecraft-mod-dictionary/data/translations/bloodmagic/
```

用語辞書：

```
/home/iuif/dev/minecraft-mod-dictionary/data/terms/mod/bloodmagic.yaml
```

---

**レポート作成日**: 2025-11-29
**作成者**: Claude (AI Assistant)
**Blood Magic バージョン**: 3.x.x (1.20.1)
