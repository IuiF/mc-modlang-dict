# Mekanism バージョン間差分分析レポート

生成日: 2025年11月29日

## 概要

Mekanism Modの過去バージョン（1.18.x, 1.19.x）と最新版（1.20.x）の翻訳データを比較分析し、バージョン間の変更点を特定しました。

## データソース

| バージョン | Mekanismバージョン | ブランチ | 取得元 |
|-----------|------------------|---------|--------|
| 1.18.2 | v10.2.x | 1.18.x | https://github.com/mekanism/Mekanism/tree/1.18.x |
| 1.19.2 | v10.3.x | 1.19.x | https://github.com/mekanism/Mekanism/tree/1.19.x |
| 1.20.x | v10.4.x | 1.20.x | https://github.com/mekanism/Mekanism/tree/1.20.x |

## エントリ数の推移

| モジュール | 1.18.x | 1.19.x | 1.20.x | 変化 (1.18→1.20) |
|-----------|--------|--------|--------|-----------------|
| Core | 1,446 | 1,647 | 1,655 | +209 (+14.5%) |
| Generators | 197 | 197 | 202 | +5 (+2.5%) |
| Additions | 309 | 319 | 321 | +12 (+3.9%) |
| Tools | 169 | 181 | 183 | +14 (+8.3%) |
| **合計** | **2,121** | **2,344** | **2,361** | **+240 (+11.3%)** |

## 主要な変更点

### 1. 流体システムの大規模変更（1.18.x → 1.19.x）

#### 削除された流体キー: 44個

**1.19.xでの変更内容:**
- 流体の言語キーシステムが完全に廃止
- 従来の `fluid.*` および `fluid.*.flowing_*` キーがすべて削除
- 新しい流体表示システムに移行

**影響を受けたモジュール:**

**Core (36キー削除):**
- 基本流体: brine, ethene, steam, heavy_water
- 酸: sulfuric_acid, hydrofluoric_acid, hydrogen_chloride
- ガス類: oxygen, hydrogen, chlorine
- 金属: sodium, lithium, superheated_sodium
- 核燃料: uranium_oxide, uranium_hexafluoride
- その他: nutritional_paste, sulfur_dioxide, sulfur_trioxide

各流体には `flowing_*` バリアントも存在（計36キー）

**Generators (8キー削除):**
- bioethanol, deuterium, tritium, fusion_fuel
- 各流体の `flowing_*` バリアント

### 2. 進捗システムの大幅追加（1.19.x）

**追加された進捗数:**
- Core: 186個
- Generators: 8個
- Additions: 10個
- Tools: 12個
- **合計: 216個**

**主要な進捗カテゴリ（Core）:**
- 機械製作関連
- マルチブロック建造
- アップグレード・モジュール
- 装備（MekaSuit, Meka-Tool）
- 核分裂・核融合
- QIOシステム

**Generators進捗例:**
- 「はじめての発電機」（Heat Generator）
- 「太陽の力」（Solar Generator）
- 「回る赤ちゃん回る」（Wind Generator）
- 「ガスを燃やす」（Gas-Burning Generator）

**Additions進捗例:**
- 「空に手を伸ばせ」（Balloon）
- 「暗がりに光る」（Glow Panel）
- 「ポップ・ポップ」（風船を飛ばす）

**Tools進捗例:**
- 「より多くの種類の防具」
- 「Better Than Netherite」（精製黒曜石防具）
- 「ピグリンからの愛」（精製グロウストーン防具）
- 「マルチツール」（Paxel）

### 3. UI/機能の変更

#### 組立機（Assemblicator）
**1.18.x-1.19.x:**
- `assemblicator.mekanism.fill_empty` → アイテムの配置/除去

**1.20.x:**
- `assemblicator.mekanism.fill` → グリッドを埋める
- `assemblicator.mekanism.empty` → グリッドを空にする

機能が分離され、より明確なUIに改善されました。

#### 分解機（Atomic Disassembler）モード
**1.18.xのみに存在:**
- `disassembler.mekanism.off` → オフ
- `disassembler.mekanism.slow` → 遅い
- `disassembler.mekanism.normal` → 通常
- `disassembler.mekanism.fast` → 早い
- `disassembler.mekanism.vein` → 鉱脈採掘

**1.19.x以降:**
モード切り替えシステムが変更され、別のUI方式に移行。

#### フィルターシステム
**1.18.x-1.19.x:**
- `button.mekanism.material_filter` → マテリアル
- `filter.mekanism.material` → マテリアルフィルター
- `filter.mekanism.material.details` → 使用素材:

**1.20.x:**
マテリアルフィルターが廃止され、より汎用的なフィルターシステムに統合。

#### デジタルマイナー
**1.18.x-1.19.x:**
- `miner.mekanism.visuals` → 範囲の可視化: %1$s
- `miner.mekanism.visuals.too_big` → 範囲を可視化するには半径が大きすぎます

**1.20.x:**
可視化システムが改善され、別の実装に変更。

### 4. 1.20.xでの微細な改善

**追加された機能（全モジュール共通）:**
- `constants.{module}.mod_name` → モジュール名定数
- `constants.{module}.pack_description` → リソースパック説明

**Core での追加:**
- ラジアルメニュー関連（爆破力設定）
- アップグレード関連の新メッセージ
- 伝送装置の新表示形式

**Generators での追加:**
- タービンの蒸気破棄機能
  - `turbine.mekanismgenerators.tooltip.steamdump` → 蒸気を破棄する
  - `turbine.mekanismgenerators.tooltip.steamdump.warning` → 水は再利用されません
  - `turbine.mekanismgenerators.tooltip.steamdump_excess` → 余分な蒸気を破棄する

## 過去バージョン限定エントリ

### 統計

| モジュール | 1.18.xのみ | 1.19.xのみ | 両方 | 合計 |
|-----------|----------|----------|------|------|
| Core | 43 | 0 | 5 | 48 |
| Generators | 8 | 0 | 0 | 8 |
| Additions | 0 | 0 | 0 | 0 |
| Tools | 0 | 0 | 0 | 0 |
| **合計** | **51** | **0** | **5** | **56** |

### カテゴリ別内訳（Core）

| カテゴリ | エントリ数 | 主な内容 |
|---------|----------|---------|
| fluid | 36 | 流体の旧表記 |
| disassembler | 5 | 分解機モード |
| button | 1 | マテリアルフィルターボタン |
| filter | 2 | マテリアルフィルター |
| miner | 2 | デジタルマイナー可視化 |
| assemblicator | 1 | 組立機UI |
| generic | 1 | 汎用フォーマット |

### カテゴリ別内訳（Generators）

| カテゴリ | エントリ数 | 主な内容 |
|---------|----------|---------|
| fluid | 8 | 発電機用流体（bioethanol, deuterium, tritium, fusion_fuel） |

## 推奨事項

### 1. 翻訳リソースパックの対応

**1.18.x対応リソースパック:**
- Core: 1,446エントリ（内48エントリは1.18.x限定）
- Generators: 197エントリ（内8エントリは1.18.x限定）

**1.19.x対応リソースパック:**
- Core: 1,647エントリ（内5エントリは1.19.x限定）
- Generators: 197エントリ

**1.20.x対応リソースパック:**
- 現行の翻訳データをそのまま使用

### 2. 後方互換性の考慮

過去バージョン向けの翻訳リソースパックを作成する場合:
1. 最新版の翻訳をベースにする
2. 各バージョンの legacy ファイルから削除されたエントリを復元
3. バージョン固有の変更点を考慮してマージ

### 3. ドキュメント化

以下の情報をプロジェクトに含める:
- ✅ バージョン間の変更点リスト（本レポート）
- ✅ 過去バージョン限定エントリのYAMLファイル
- ✅ README.mdへの過去バージョン情報追加

## 成果物

### 生成されたファイル

1. **mekanism_core_legacy.yaml**
   - 48エントリの過去バージョン限定翻訳
   - 1.18.xのみ: 43エントリ
   - 1.18.x-1.19.x共通: 5エントリ

2. **mekanism_generators_legacy.yaml**
   - 8エントリの過去バージョン限定翻訳
   - すべて1.18.xのみ

3. **更新されたREADME.md**
   - 過去バージョン対応の情報追加
   - バージョン間の変更点の概要
   - 注意事項の追加

## 結論

Mekanismは1.18.x → 1.19.x → 1.20.xの間で以下の主要な変更がありました:

1. **流体システムの刷新**（1.19.x）
   - 44キーの削除
   - 新しい流体表示メカニズムへの移行

2. **進捗システムの追加**（1.19.x）
   - 216個の新しい進捗
   - ゲームプレイガイドの強化

3. **UI/UXの改善**（1.20.x）
   - より明確な機能分離
   - 一貫性のある用語使用

過去バージョン対応として56エントリの翻訳を保存し、後方互換性を確保しました。
