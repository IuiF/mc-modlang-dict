# Mekanism 翻訳データ

Mekanism Modの日本語翻訳データをYAML形式で管理しています。

## ファイル一覧

### 個別モジュールファイル
- `mekanism_core_translations.yaml` - Mekanism Core（メインモジュール）
- `mekanism_generators_translations.yaml` - Mekanism Generators（発電機モジュール）
- `mekanism_tools_translations.yaml` - Mekanism Tools（ツールモジュール）
- `mekanism_additions_translations.yaml` - Mekanism Additions（追加要素モジュール）

### 統合ファイル
- `mekanism_all_translations.yaml` - 全モジュールを統合した翻訳データ

### 過去バージョン対応ファイル
- `mekanism_core_legacy.yaml` - 過去バージョン（1.18.x/1.19.x）のみに存在するCore翻訳
- `mekanism_generators_legacy.yaml` - 過去バージョン（1.18.x/1.19.x）のみに存在するGenerators翻訳

## 翻訳統計

| モジュール | エントリ数 | 翻訳済み | 翻訳率 |
|-----------|----------|---------|--------|
| Core | 1,656 | 1,655 | 99.9% |
| Generators | 202 | 202 | 100.0% |
| Tools | 183 | 183 | 100.0% |
| Additions | 321 | 321 | 100.0% |
| **合計** | **2,362** | **2,361** | **99.9%** |

## データソース

- **GitHubリポジトリ**: https://github.com/mekanism/Mekanism
- **最新版ブランチ**: 1.20.x
- **バージョン**: 10.4.x系
- **取得日**: 2025年11月29日

### 過去バージョン対応

以下のバージョンの翻訳データも取得済み:

| Minecraftバージョン | Mekanismバージョン | ブランチ | 対応ファイル |
|-------------------|------------------|---------|-------------|
| 1.18.2 | v10.2.x | 1.18.x | mekanism_*_legacy.yaml |
| 1.19.2 | v10.3.x | 1.19.x | mekanism_*_legacy.yaml |

過去バージョンのみに存在する翻訳エントリ:
- **Core**: 48エントリ（主に流体、UI、分解機モード）
- **Generators**: 8エントリ（流体関連）
- **Additions**: 0エントリ
- **Tools**: 0エントリ

### 英語ソース
```
src/datagen/generated/{module}/assets/{module_id}/lang/en_us.json
```

### 日本語翻訳
```
src/{module}/resources/assets/{module_id}/lang/ja_jp.json
```

## カテゴリ構成

### Core モジュール (84カテゴリ)
主要カテゴリ:
- **advancements** (186) - 進捗
- **block** (199) - ブロック
- **item** (168) - アイテム
- **description** (114) - 説明文
- **container** (98) - コンテナUI
- **gui** (89) - GUI要素
- **module** (71) - モジュール（MekaSuit用）
- **gas** (25) - ガス
- **pigment** (19) - 色素
- **slurry** (15) - スラリー
- **infuse_type** (9) - 注入タイプ

その他: sound_event, tooltip, qio, boiler, transmitter, miner, robit, configurator, security, etc.

### Generators モジュール (15カテゴリ)
- **advancements** (12) - 進捗
- **block** (54) - ブロック（発電機、反応炉部品など）
- **item** (39) - アイテム
- **container** (30) - コンテナUI
- **gui** (14) - GUI要素
- **reactor** (17) - 反応炉関連
- **turbine** (14) - タービン関連
- **gas** (4) - ガス（重水素、三重水素など）

### Tools モジュール (4カテゴリ)
- **advancements** (6) - 進捗
- **item** (173) - アイテム（各種ツール、防具、盾）
- **tooltip** (3) - ツールチップ
- **constants** (1) - 定数

### Additions モジュール (9カテゴリ)
- **advancements** (6) - 進捗
- **block** (304) - ブロック（プラスチックブロック各種）
- **item** (16) - アイテム（風船）
- **entity** (7) - エンティティ（ベビーモブ）
- **description** (3) - 説明文

## 主要翻訳用語

詳細は `/data/terms/mod/mekanism.yaml` を参照してください。

### エネルギー
- Joules (J) → ジュール (J)
- Heat → 熱

### 物質タイプ
- Gas → ガス
- Infusion → 注入
- Pigment → 色素
- Slurry → スラリー

### ティア
- Basic → 基本
- Advanced → 発展
- Elite → 精鋭
- Ultimate → 究極

### 伝送装置
- Universal Cable → ユニバーサルケーブル
- Pressurized Tube → 加圧チューブ
- Logistical Transporter → ロジスティカルトランスポーター
- Mechanical Pipe → メカニカルパイプ
- Thermodynamic Conductor → 熱力学導管

### 主要機械
- Enrichment Chamber → 濃縮チャンバー
- Metallurgic Infuser → 冶金注入機
- Purification Chamber → 精製チャンバー
- Chemical Injection Chamber → 化学注入チャンバー
- Digital Miner → デジタルマイナー
- Energized Smelter → 電動かまど
- Formulaic Assemblicator → 定型組立機

### マルチブロック
- Fission Reactor → 核分裂炉
- Fusion Reactor → 核融合炉
- Boiler → ボイラー
- Turbine → タービン
- Evaporation Plant → 蒸発プラント
- Induction Matrix → 誘導マトリクス
- SPS (Supercritical Phase Shifter) → 超臨界相シフター

### 装備
- MekaSuit → メカスーツ
- Meka-Tool → メカツール
- Atomic Disassembler → 原子分解機
- Jetpack → ジェットパック

## 注意が必要な項目

### 未翻訳エントリ (1件)

Core モジュールに1件の未翻訳エントリがあります:

```yaml
key: constants.mekanism.recipe_warning
en: "Broken tags in Mekanism recipes detected, please check server logs for details. You will be missing some recipes and machines may not accept expected inputs."
```

これは開発者向けの警告メッセージであり、通常プレイヤーが目にすることは少ないため、翻訳の優先度は低いと考えられます。

### 過去バージョン特有の項目

#### 1.18.x での主要な変更点

**流体システムの旧表記**（1.19.xで廃止）:
- `fluid.mekanism.*` および `fluid.mekanismgenerators.*` のキーが存在
- 各流体に `flowing_*` バリアントが存在
- 例: `fluid.mekanism.chlorine` → 液体塩素

**分解機（Atomic Disassembler）のモード**（1.19.xで廃止）:
- `disassembler.mekanism.off` → オフ
- `disassembler.mekanism.slow` → 遅い
- `disassembler.mekanism.normal` → 通常
- `disassembler.mekanism.fast` → 早い
- `disassembler.mekanism.vein` → 鉱脈採掘

**デジタルマイナーの可視化**（1.19.xで変更）:
- `miner.mekanism.visuals` → 範囲の可視化: %1$s

#### 1.19.x での主要な追加機能

**進捗システムの大幅追加**:
- Core: 186個の進捗
- Generators: 8個の進捗（発電機関連）
- Additions: 10個の進捗（バルーン、ベビーモブ関連）
- Tools: 12個の進捗（防具、ツール関連）

### フォーマットコード

多くのエントリにフォーマットコードが含まれています:
- `%s` - 文字列置換
- `%d` - 数値置換
- `$(br)` - 改行（Patchouliマクロ）
- `$(br2)` - 空行（Patchouliマクロ）
- `§` - Minecraftカラーコード

これらは必ず原文と同じ位置・形式で保持する必要があります。

### 長文の説明

一部の進捗説明やGUIテキストは複数行にわたる長文です。翻訳時は自然な日本語表現を心がけつつ、原文の意図を正確に伝える必要があります。

## 使用方法

### YAMLデータの読み込み例（Python）

```python
import yaml

# 個別モジュールの読み込み
with open('mekanism_core_translations.yaml', 'r', encoding='utf-8') as f:
    core_data = yaml.safe_load(f)

# 統合ファイルの読み込み
with open('mekanism_all_translations.yaml', 'r', encoding='utf-8') as f:
    all_data = yaml.safe_load(f)

# カテゴリ別のアクセス
for translation in core_data['categories']['block']['translations']:
    print(f"{translation['key']}: {translation['en']} → {translation['ja']}")
```

### 翻訳の検索

特定のキーの翻訳を探す:

```bash
grep -A 2 "key: block.mekanism.osmium_ore" mekanism_core_translations.yaml
```

カテゴリ内のすべての翻訳を抽出:

```bash
python3 -c "
import yaml
data = yaml.safe_load(open('mekanism_core_translations.yaml'))
for t in data['categories']['block']['translations']:
    print(f\"{t['en']} → {t['ja']}\")
"
```

## 更新履歴

- **2025-11-29 (2回目)**: 過去バージョン対応追加
  - 1.18.x (v10.2.x) および 1.19.x (v10.3.x) の翻訳データを取得
  - 過去バージョン限定の翻訳エントリを抽出（Core: 48件、Generators: 8件）
  - `mekanism_core_legacy.yaml` および `mekanism_generators_legacy.yaml` を作成
  - バージョン間の主要な変更点を分析:
    - 1.19.xで流体の言語キーシステムが変更（44キー削除）
    - 1.19.xで進捗システムが大幅追加（全モジュールで216エントリ追加）
    - 1.20.xで定数・パック情報の整理
  - モジュール別エントリ数の推移を記録

- **2025-11-29**: 初版作成
  - Mekanism 1.20.x (10.4.x系) の翻訳データを取得
  - 4モジュール（Core, Generators, Tools, Additions）をYAML化
  - 合計2,362エントリ、翻訳率99.9%

## ライセンス・クレジット

- **Mekanism Mod**: MIT License
- **翻訳データ**: Mekanism公式リポジトリより取得
- **翻訳者**: Mekanismコミュニティの翻訳貢献者の皆様

翻訳データは元のMIT Licenseに従います。
