# Create Mod 日本語翻訳データ

## 概要

Create Mod (v0.5.1 / Minecraft 1.20.1) の日本語翻訳データです。

## ファイル構成

### `create_main.yaml`
主要なブロック、アイテム、実績の翻訳データを含む統合ファイル。

**含まれる翻訳:**
- **ブロック**: 89エントリー（歯車、シャフト、機械、ケーシング等）
- **アイテム**: 46エントリー（素材、ツール、装備、食料等）
- **実績**: 110エントリー（ゲーム内の達成項目と説明文）
- **合計**: 245エントリー

### `create.yaml` (用語辞書)
Create Mod固有の用語辞書。翻訳の一貫性を保つために使用します。

## 主要な翻訳ガイドライン

### Create固有用語の統一

| 英語 | 日本語 | 備考 |
|------|-------|------|
| Kinetic | 回転動力 | Create Modの中核概念 |
| Stress | ストレス | 機械の負荷を表す |
| RPM | RPM | 回転数（そのまま使用） |
| Contraption | 装置 | 可動構造物の総称 |
| Schematic | 設計図 | - |
| Cogwheel | 歯車 | - |
| Large Cogwheel | 大きな歯車 | - |
| Shaft | シャフト | - |
| Encased | 覆われた | 接頭辞として使用 |
| Casing | ケーシング | - |

### 素材・合金

| 英語 | 日本語 |
|------|-------|
| Andesite Alloy | 安山岩合金 |
| Brass | 真鍮 |
| Zinc | 亜鉛 |
| Rose Quartz | ローズクォーツ |
| Sturdy Sheet | 頑丈なシート |

### 機械名の翻訳パターン

- **Mechanical ~**: 機械式〜
  - Mechanical Press → 機械式プレス
  - Mechanical Mixer → 機械式ミキサー
  - Mechanical Pump → 機械式ポンプ

- **Encased ~**: 覆われた〜
  - Encased Fan → 覆われたファン
  - Encased Shaft → 覆われたシャフト

## データ構造

```yaml
mod_id: create
mod_name: Create
version: 0.5.1
mc_version: 1.20.1
total_entries: 245

categories:
  blocks:
    count: 89
    translations:
      - key: block.create.cogwheel
        source: Cogwheel
        target: 歯車
        category: block

  items:
    count: 46
    translations:
      - key: item.create.andesite_alloy
        source: Andesite Alloy
        target: 安山岩合金
        category: item

  advancements:
    count: 110
    translations:
      - key: advancement.create.root
        target: Createへようこそ
        category: advancement
        is_description: false
```

## カバー範囲

### 翻訳済みカテゴリ

✅ **ブロック** (89/701)
- 動力伝達系（歯車、シャフト、ベルト）
- 加工機械（プレス、ミキサー、粉砕ホイール）
- 液体処理（ポンプ、タンク、パイプ）
- 物流（ファンネル、シュート、アーム）
- 鉄道（線路、駅、信号機）
- ケーシング各種

✅ **アイテム** (46/204)
- 素材・合金
- ツール・装備
- 食料
- 設計図関連

✅ **実績** (110/204)
- 主要な進行実績
- 隠し実績を含む

### 未翻訳カテゴリ

以下のカテゴリは今後追加予定：

- `create.ponder.*` (1052エントリー) - チュートリアルテキスト
- `create.gui.*` (268エントリー) - UI関連
- `create.tooltip.*` (37エントリー) - ツールチップ
- `create.subtitle.*` (60エントリー) - 効果音字幕
- その他（tag、schedule等）

## 翻訳の追加方法

### 1. 新しいブロック/アイテムの追加

```yaml
- key: block.create.new_machine
  source: New Machine
  target: 新しい機械
  category: block
```

### 2. 実績の追加

```yaml
- key: advancement.create.new_achievement
  target: 新しい達成項目
  category: advancement
  is_description: false

- key: advancement.create.new_achievement.desc
  target: 達成項目の説明文
  category: advancement
  is_description: true
```

## 翻訳品質基準

1. **Minecraft公式訳に準拠**
   - バニラアイテム・ブロック名は公式日本語訳を使用

2. **フォーマットコード保持**
   - `§7(隠し実績)` などの色コードは保持
   - `%s`, `%d`, `%1$s` などのプレースホルダーは保持

3. **自然な日本語**
   - 技術用語は適切にカタカナ化
   - 文脈に応じた自然な訳文

4. **一貫性**
   - 用語辞書（`create.yaml`）に従う
   - 同じ用語は常に同じ訳語を使用

## 使用例

このYAMLデータは、Minecraft Mod翻訳システムで以下のように利用できます：

1. **リソースパック生成**: YAMLからja_jp.jsonを生成
2. **用語検索**: 特定のブロック・アイテム名の翻訳を検索
3. **一貫性チェック**: 翻訳の統一性を検証

## ライセンス・クレジット

- **Create Mod**: Simibubi氏他による開発
- **公式翻訳プロジェクト**: [Crowdin - Create Mod](https://crowdin.com/project/createmod)
- **このデータ**: プロジェクト固有の翻訳データ

## 今後の予定

- [ ] Ponderテキストの翻訳追加
- [ ] GUIテキストの翻訳追加
- [ ] ツールチップの翻訳追加
- [ ] 全ブロック・アイテムの完全カバー
- [ ] Create Add-onの翻訳対応

## 参考リンク

- [Create Mod GitHub](https://github.com/Creators-of-Create/Create)
- [Create Mod Crowdin](https://crowdin.com/project/createmod)
- [Create Mod Wiki (日本語)](https://seesaawiki.jp/minecraft_create_mod/)
