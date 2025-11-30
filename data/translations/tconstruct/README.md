# Tinkers' Construct 翻訳データ

## 概要

- **Mod ID**: tconstruct
- **Mod名**: Tinkers' Construct
- **バージョン**: 3.8.x
- **Minecraft バージョン**: 1.20.1
- **翻訳日**: 2025-11-29

## 統計情報

- **総翻訳キー数**: 3,257
- **検証済み翻訳**: 107エントリ
- **カテゴリ別エントリ数**:
  - ブロック: 543
  - アイテム: 416
  - 素材: 404
  - モディファイア: 1,001
  - 液体: 205
  - 書籍: 24
  - その他: 664

## 重点翻訳カテゴリ

### 1. Smeltery/Foundry関連 (206エントリ)
精錬炉と鋳造所システムに関する翻訳。
- 精錬炉コントローラー、鋳造所コントローラー
- 焼成/焦熱ブロック各種
- 鋳造台、鋳造盤
- 排出口、ダクト、投入口、タンクなど

### 2. ツール/武器 (45エントリ)
Tinkers' Constructの主要ツールと武器。
- ピッケル、ハンマー、掘削機
- 斧、マトック、大鎌
- 剣、包丁、短剣

### 3. 素材 (100エントリ)
ツール素材とその特性。
- コバルト、マニュリン、ヘパティゾン
- ローズゴールド、ピグアイアン
- 女王スライム、ナイトスライム

### 4. モディファイア (100エントリ)
ツール強化用のモディファイア。
- ダイヤモンド、エメラルド、ネザライト
- 補強、不壊、採掘速度上昇
- 鋭さ、幸運、自動精錬

## ファイル構成

```
/data/translations/tconstruct/
├── ja_jp.yaml              # 最終統合翻訳ファイル
└── README.md               # このファイル

/workspace/imports/tconstruct/
├── en_us.json              # 元の英語言語ファイル
├── blocks_translation.yaml # ブロック翻訳（543エントリ）
├── items_translation.yaml  # アイテム翻訳（416エントリ）
├── materials_translation.yaml # 素材翻訳（404エントリ）
├── modifiers_translation.yaml # モディファイア翻訳（1001エントリ）
├── fluids_translation.yaml # 液体翻訳（205エントリ）
└── books_translation.yaml  # 書籍翻訳（24エントリ）
```

## 翻訳品質ステータス

- **verified**: 手動で検証済みの高品質翻訳（107エントリ）
- **draft**: 機械的に生成された下書き翻訳（残りのエントリ）
- **[要翻訳]**: 未翻訳または翻訳が必要なエントリ

## Tinkers固有用語

| 英語 | 日本語 | 備考 |
|------|--------|------|
| Smeltery | 精錬炉 | メインのマルチブロック構造 |
| Foundry | 鋳造所 | 高度な精錬炉 |
| Seared | 焼成 | 精錬炉用建材 |
| Scorched | 焦熱 | 鋳造所用建材 |
| Cast | 鋳型 | 金属製の再利用可能な型 |
| Pattern | 型紙 | 紙や木製の使い捨て型 |
| Casting | 鋳造 | 液体金属を型に流し込む工程 |
| Molten | 溶融 | 溶けた状態の金属 |
| Tool Part | ツールパーツ | ツールの構成部品 |
| Modifier | モディファイア | ツール強化要素 |
| Trait | 特性 | 素材固有の性質 |
| Durability | 耐久値 | ツールの耐久性 |
| Mining Speed | 採掘速度 | 採掘の速さ |
| Attack Damage | 攻撃力 | 武器の攻撃力 |
| Harvest Tier | 採掘層 | 採掘可能なブロックのレベル |

## 注意事項

1. **フォーマットコード保持**: `%s`, `%d`, `$(...)` などのコードは変更しないでください
2. **改行保持**: `\n` や `$(br)` などの改行コードは保持してください
3. **Minecraft公式訳準拠**: バニラアイテム・ブロック名は公式訳に従ってください
4. **一貫性**: 同じ用語は常に同じ訳語を使用してください

## 今後の作業

- [ ] draft状態のエントリを手動で翻訳・検証
- [ ] Mantle Book（ゲーム内ガイドブック）のコンテンツ翻訳
- [ ] ツールチップやフレーバーテキストの丁寧な翻訳
- [ ] モディファイア説明文の詳細翻訳
- [ ] 実際のゲームプレイでの動作確認

## 参考リンク

- **GitHubリポジトリ**: https://github.com/SlimeKnights/TinkersConstruct
- **CurseForge**: https://www.curseforge.com/minecraft/mc-mods/tinkers-construct
- **公式Wiki**: https://tinkers-construct.fandom.com/wiki/Tinkers'_Construct_Wiki
