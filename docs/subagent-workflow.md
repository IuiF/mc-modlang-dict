# サブエージェント運用ガイド

## エージェント構成

| エージェント | 役割 | 並列実行 |
|-------------|------|----------|
| 翻訳Agent | Mod単位・ファイル単位の翻訳実行 | 可 |
| 用語抽出Agent | 新規用語の抽出・辞書提案 | 可 |
| レビュアーAgent | 翻訳品質・一貫性の検証 | 最終段階 |
| 調査Agent | Mod構造・既存翻訳の調査 | 可 |

## 並列実行の原則

1. **Mod単位では順次実行**
   - 1つのModを100%完了してから次へ
   - 中途半端な状態で複数Modを並列しない

2. **同一Mod内のカテゴリは並列可**
   - block / item / advancement 等のカテゴリ分割
   - lang / patchouli の同時処理

3. **レビューは各Mod完了時に実施**
   - 翻訳率100%を確認
   - 用語の一貫性、フォーマット保持を検証

## コンテキスト効率化

**重要: サブエージェントにDB操作を委任する**

```
❌ NG: メインでJSON出力 → 内容確認 → サブエージェントに渡す
       （メインコンテキストにJSONが展開されてしまう）

✅ OK: サブエージェント起動 → エージェント内でDB操作 → 完了報告のみ返す
       （メインコンテキストを消費しない）
```

## サブエージェント指示テンプレート

```
[Mod名]の翻訳を行ってください。

【用語辞書 - 必ず従うこと】
色名: White→白, Orange→橙, Magenta→赤紫, Light Blue→空色,
      Yellow→黄, Lime→黄緑, Pink→桃色, Gray→灰色,
      Light Gray→薄灰色, Cyan→青緑, Purple→紫, Blue→青,
      Brown→茶, Green→緑, Red→赤, Black→黒

接尾辞: Block→ブロック, Stairs→階段, Slab→ハーフブロック,
        Wall→塀, Fence→フェンス, Pillar→柱, Bricks→レンガ

接頭辞: Chiseled→彫られた, Polished→磨かれた, Smooth→滑らかな,
        Cracked→ひび割れた, Mossy→苔むした

【Mod固有用語】（必要に応じて追加）
例: Marble→大理石, Limestone→石灰岩, Basalt→玄武岩

【厳守事項】
- 翻訳JSONに空文字("")を含めない
- 翻訳できないキーはJSONに含めない
- 各ステップでステータスを確認する
- 用語辞書に従って一貫性を保つ

1. JARをインポート:
   ./moddict import -jar workspace/imports/mods/[mod.jar]
   ./moddict repair

2. 公式翻訳確認・インポート:
   unzip -l [mod.jar] | grep ja_jp
   （あれば）./moddict translate -mod [mod_id] -official [ja_jp.json]

3. ステータス確認:
   ./moddict translate -mod [mod_id] -status

4. pendingがあれば20件ずつ翻訳:
   ./moddict translate -mod [mod_id] -export pending.json -limit 20
   （翻訳後 - 用語辞書に従い、空文字を含めない）
   ./moddict translate -mod [mod_id] -json translated.json
   ./moddict translate -mod [mod_id] -status  # 毎回確認
   （繰り返し）

5. 100%になったら報告

翻訳ルール: 用語辞書準拠、フォーマットコード保持、Minecraft公式用語準拠
最終報告: Mod ID、総キー数、翻訳完了数、ステータス
```

**Haikuを使う理由**: 翻訳タスクは定型的で、コスト効率が良い。

**用語辞書の詳細**: `docs/translation-consistency.md` を参照。

## 翻訳JSON生成ルール（厳守）

```json
// ✅ 正しい翻訳JSON:
{
  "item.mod.example": "サンプルアイテム",
  "block.mod.example": "サンプルブロック"
}

// ❌ 禁止パターン1: 空文字
{
  "item.mod.example": "",      // ← データ破損の原因
  "block.mod.example": "サンプルブロック"
}

// ❌ 禁止パターン2: 翻訳できないキーを含める
{
  "item.mod.example": "item.mod.example",  // ← 原文のままは無意味
  "block.mod.example": "サンプルブロック"
}

// ✅ 翻訳できないキーは含めない:
{
  "block.mod.example": "サンプルブロック"
}
```

## 翻訳ガイドライン

### 用語の優先順位

```
global (低) → category (中) → mod (高)
```

### 必須ルール

1. **Minecraft公式訳に準拠** - バニラアイテム・ブロック名
2. **フォーマットコード保持** - `§`, `$(...)`, `%s`, `%d`
3. **改行・空白保持** - `\n`, `$(br)`, `$(br2)`

### Patchouliマクロ

| マクロ | 説明 |
|-------|------|
| `$(br)` | 改行 |
| `$(br2)` | 空行 |
| `$(item)`, `$(thing)` | ハイライト |
| `$(l:path)テキスト$(/)` | 内部リンク |

## 一貫性の仕組み

### 自動伝播（インポート時）

同じ`source_text`を持つ複数のキーがある場合、1つに翻訳をインポートすると他のpendingキーにも自動適用される。

```
例: "Stained Glass" が17キーで使用されている場合
  - key1に "色付きガラス" をインポート
  - key2〜key17のpendingにも自動で "色付きガラス" が適用
```

### 用語辞書（翻訳時）

部品レベル（色、素材、接辞）の一貫性は、用語辞書をサブエージェントに渡すことで担保。

**重要**: 用語辞書はコンパクトに保つ（50〜100項目以内）。全翻訳を参照として渡すとコンテキストを圧迫する。

詳細: `docs/translation-consistency.md`
