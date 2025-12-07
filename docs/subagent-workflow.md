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

【厳守事項】
- 翻訳JSONに空文字("")を含めない
- 翻訳できないキーはJSONに含めない
- 各ステップでステータスを確認する

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
   （翻訳後 - 空文字を含めない）
   ./moddict translate -mod [mod_id] -json translated.json
   ./moddict translate -mod [mod_id] -status  # 毎回確認
   （繰り返し）

5. 100%になったら報告

翻訳ルール: フォーマットコード保持、Minecraft公式用語準拠
最終報告: Mod ID、総キー数、翻訳完了数、ステータス
```

**Haikuを使う理由**: 翻訳タスクは定型的で、コスト効率が良い。

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
