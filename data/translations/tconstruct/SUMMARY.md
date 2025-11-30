# Tinkers' Construct 翻訳データ作成 - 完了報告

## 作成日時
2025-11-29

## 作業概要

Tinkers' Construct Mod (1.20.1版) の日本語翻訳データを作成しました。

## 成果物

### 1. 用語辞書
**ファイル**: `/home/iuif/dev/minecraft-mod-dictionary/data/terms/mod/tconstruct.yaml`

Tinkers' Construct固有の用語を定義した辞書ファイル。主要な用語を優先度付きで管理。

**主要用語例**:
- Smeltery → 精錬炉
- Foundry → 鋳造所
- Seared → 焼成
- Scorched → 焦熱
- Cast → 鋳型
- Pattern → 型紙
- Modifier → モディファイア

### 2. 翻訳データファイル

#### メインファイル
**ファイル**: `/home/iuif/dev/minecraft-mod-dictionary/data/translations/tconstruct/ja_jp.yaml`

**統計情報**:
- 総翻訳キー数: 3,257
- 検証済み翻訳: 107エントリ
- 重点カテゴリ:
  - Smeltery/Foundry関連: 206エントリ
  - ツール/武器: 45エントリ
  - 素材: 100エントリ
  - モディファイア: 100エントリ

#### カテゴリ別ファイル（作業用）
**ディレクトリ**: `/home/iuif/dev/minecraft-mod-dictionary/workspace/imports/tconstruct/`

- `blocks_translation.yaml` - 543エントリ (65KB)
- `items_translation.yaml` - 416エントリ (50KB)
- `materials_translation.yaml` - 404エントリ (69KB)
- `modifiers_translation.yaml` - 1,001エントリ (131KB)
- `fluids_translation.yaml` - 205エントリ (30KB)
- `books_translation.yaml` - 24エントリ (3.6KB)

### 3. ドキュメント

- `README.md` - 翻訳データの概要と使用方法
- `SUMMARY.md` - この完了報告書

## 翻訳品質

### 検証済み高品質翻訳（107エントリ）

以下のカテゴリは手動で翻訳を確認し、高品質な訳語を適用：

1. **Smeltery/Foundryシステム** - 精錬炉と鋳造所の主要コンポーネント
   - 精錬炉コントローラー、鋳造所コントローラー
   - 焼成/焦熱ブロック（石、丸石、レンガ、ガラスなど）
   - 排出口、ダクト、投入口、タンク、蛇口、水路
   - 鋳造台、鋳造盤、合金炉

2. **クラフト系ブロック**
   - 作業場、改造ステーション、パーツ作成台
   - モディファイア作業台、ティンカーの金床

3. **書籍**
   - 素材と君 (Materials and You)
   - 小さな精錬 (Puny Smelting)
   - 偉大なる精錬 (Mighty Smelting)
   - ティンカーの道具術 (Tinkers' Gadgetry)
   - 素晴らしき鋳造所 (Fantastic Foundry)
   - ティンカリング大百科 (Encyclopedia)

4. **型紙・鋳型**
   - 各種パーツ用の型紙と鋳型（ピッケルヘッド、斧ヘッド、刀身など）

5. **主要素材**
   - バニラ素材（木、石、鉄、金、ダイヤモンド、ネザライト）
   - Tinkers素材（コバルト、マニュリン、ヘパティゾン、女王スライム、ピグアイアン、ローズゴールド）

6. **主要モディファイア**
   - 基本強化（ダイヤモンド、エメラルド、ネザライト、補強、不壊）
   - 採掘系（採掘速度上昇、幸運、シルクタッチ、自動精錬、範囲拡大）
   - 戦闘系（鋭さ、特効系、火属性、ドロップ増加、ノックバック）
   - 防具系（各種ダメージ軽減、水中系）

### 下書き翻訳（残り3,150エントリ）

機械的な置換により基本的な翻訳を生成。以下の状態：
- 主要用語は適切に翻訳済み
- 複雑な文章や説明文は `[要翻訳]` マーク付き
- 今後の手動翻訳・レビューが必要

## 翻訳方針

### 準拠基準

1. **Minecraft公式日本語訳** - バニラアイテム・ブロック名は公式訳に従う
2. **フォーマットコード保持** - `%s`, `%d`, `$(...)` などは変更しない
3. **改行・空白保持** - `\n`, `$(br)`, `$(br2)` は保持
4. **用語の一貫性** - 同じ英語用語は常に同じ日本語訳を使用

### Tinkers固有用語の統一

| カテゴリ | 英語 | 日本語 |
|---------|------|--------|
| システム | Smeltery | 精錬炉 |
| システム | Foundry | 鋳造所 |
| 建材 | Seared | 焼成 |
| 建材 | Scorched | 焦熱 |
| クラフト | Cast | 鋳型 |
| クラフト | Pattern | 型紙 |
| クラフト | Casting | 鋳造 |
| 液体 | Molten | 溶融 |
| ツール | Tool Part | ツールパーツ |
| 強化 | Modifier | モディファイア |
| 強化 | Trait | 特性 |
| 強化 | Ability | 能力 |
| ステータス | Durability | 耐久値 |
| ステータス | Mining Speed | 採掘速度 |
| ステータス | Attack Damage | 攻撃力 |
| ステータス | Harvest Tier | 採掘層 |

## 注意が必要な項目

### 1. Mantle Book（ゲーム内ガイドブック）

現在のバージョンでは、Mantle Bookのコンテンツは言語ファイル内に24エントリのみ存在。
大規模なブックシステムは別途JSONファイルで管理されている可能性がありますが、
1.20.1ブランチのGitHubリポジトリでは該当ディレクトリが見つかりませんでした。

**今後の対応**:
- 実際のゲームファイル（JARファイル）を確認
- Mantleライブラリ側のリポジトリも調査
- ブックコマンドでエクスポートされたデータを入手

### 2. 複雑な説明文

以下のタイプの翻訳は今後の手動作業が必要：
- モディファイアの詳細な効果説明
- 素材の特性（Trait）説明
- フレーバーテキスト
- 進捗（Advancement）の説明文
- ツールチップの詳細情報

### 3. 固有名詞

以下は原語のまま、またはカタカナ表記を検討：
- Manyullyn（マニュリン）
- Hepatizon（ヘパティゾン）
- Nahuatl（ナワトル）
- Knightslime（ナイトスライム）

## ファイル配置

```
minecraft-mod-dictionary/
├── data/
│   ├── terms/
│   │   └── mod/
│   │       └── tconstruct.yaml          # Tinkers固有用語辞書
│   └── translations/
│       └── tconstruct/
│           ├── ja_jp.yaml               # メイン翻訳ファイル
│           ├── README.md                # 使用方法ドキュメント
│           └── SUMMARY.md               # この完了報告書
└── workspace/
    └── imports/
        └── tconstruct/
            ├── en_us.json               # 元の英語ファイル
            ├── blocks_translation.yaml   # カテゴリ別翻訳（作業用）
            ├── items_translation.yaml
            ├── materials_translation.yaml
            ├── modifiers_translation.yaml
            ├── fluids_translation.yaml
            ├── books_translation.yaml
            └── tconstruct_translation_final.yaml  # 統合版
```

## 今後の推奨作業

### 優先度：高

1. **Mantle Bookコンテンツの翻訳**
   - ゲーム内ガイドブックの全ページ
   - 精錬炉・鋳造所の構築方法
   - ツール作成チュートリアル

2. **主要モディファイアの説明文翻訳**
   - 上位100モディファイアの詳細説明
   - 効果の数値情報を含む正確な翻訳

3. **素材特性の翻訳**
   - 全素材のTraitとAbilityの説明
   - フレーバーテキストの自然な日本語化

### 優先度：中

4. **進捗（Advancement）の翻訳**
   - タイトルと説明文
   - ゲーム進行のガイドとして重要

5. **ツールチップの翻訳**
   - アイテム詳細情報
   - 使用方法のヒント

6. **GUI・コマンド関連**
   - ユーザーインターフェース
   - コマンドメッセージ

### 優先度：低

7. **効果音字幕**
   - subtitle系のキー（アクセシビリティ向上）

8. **デバッグ・開発者向けメッセージ**
   - コメント、エラーメッセージ

## まとめ

Tinkers' Construct Modの基本的な翻訳データ構造を確立し、主要な107エントリについて
高品質な日本語訳を適用しました。残り3,150エントリは下書き状態であり、今後の
継続的な翻訳作業により品質を向上させていく必要があります。

特に、精錬炉・鋳造所システムとツール関連の用語は統一された訳語で管理されており、
一貫性のある翻訳が可能になっています。
