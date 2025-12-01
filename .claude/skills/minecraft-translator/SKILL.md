---
name: minecraft-translator
description: Minecraft Modの日本語翻訳を実行。SQLite DBでデータ管理、moddict CLIでDB操作。翻訳作業やMod翻訳の話題で自動的に使用。
allowed-tools: Read, Grep, Bash, Glob, Write
---

# Minecraft Mod 翻訳スキル

Minecraft Modの英語テキストを日本語に翻訳するためのスキル。

## DB中心のワークフロー

翻訳データは**SQLite DB (moddict.db)** で一元管理する。

```
DB → export → JSON(pending) → 翻訳 → JSON(translated) → import → DB → cleanup
```

JSONはCLIとの受け渡し用中間フォーマット。最終的なデータはすべてDBに格納。

---

## 一時ファイル管理（厳守）

### 使用する一時ファイルパス

| ファイル | パス | 用途 |
|----------|------|------|
| pending | `/tmp/moddict_pending.json` | DBからエクスポート |
| translated | `/tmp/moddict_translated.json` | 翻訳結果 |
| official | `/tmp/moddict_official.json` | 公式翻訳抽出 |

### 禁止

- `workspace/` 以下に一時JSONを作成
- ファイル名に連番を付ける（`pending_1.json`等）
- プロジェクトルートに一時ファイルを作成
- 作業完了後に一時ファイルを残す

### クリーンアップコマンド

```bash
# 作業開始時・終了時に実行
rm -f /tmp/moddict_*.json
```

---

## CLIコマンド（DB操作）

```bash
# DBにソースをインポート
./moddict import -jar [mod.jar]

# DB整合性修復
./moddict repair

# DBステータス確認
./moddict translate -mod [mod_id] -status

# DBからpendingをエクスポート（一時ファイルへ）
./moddict translate -mod [mod_id] -export /tmp/moddict_pending.json -limit 20

# 翻訳をDBにインポート（一時ファイルから）
./moddict translate -mod [mod_id] -json /tmp/moddict_translated.json

# 公式翻訳をDBにインポート（verified状態）
./moddict translate -mod [mod_id] -official /tmp/moddict_official.json

# DBから翻訳済みファイルを出力
./moddict export -mod [mod_id]
```

---

## 翻訳ループテンプレート

```bash
# 1. クリーンアップ
rm -f /tmp/moddict_*.json

# 2. エクスポート
./moddict translate -mod [mod_id] -export /tmp/moddict_pending.json -limit 20

# 3. 翻訳（JSONを読み、翻訳を/tmp/moddict_translated.jsonに書き出す）

# 4. インポート
./moddict translate -mod [mod_id] -json /tmp/moddict_translated.json

# 5. 一時ファイル削除
rm -f /tmp/moddict_pending.json /tmp/moddict_translated.json

# 6. ステータス確認
./moddict translate -mod [mod_id] -status

# 7. Pending: 0 になるまで 2-6 を繰り返し
```

---

## 翻訳JSON生成ルール（厳守）

### 正しいパターン

```json
{
  "item.mod.example": "サンプルアイテム",
  "block.mod.example": "サンプルブロック"
}
```

### 禁止パターン

```json
// 空文字は絶対禁止（DB破損の原因）
{ "item.mod.example": "" }

// 原文のままは無意味
{ "item.mod.example": "item.mod.example" }
```

**翻訳できないキーはJSONに含めない。**

---

## フォーマットコード保持（必須）

| コード | 説明 |
|--------|------|
| `%s`, `%d` | プレースホルダー |
| `%1$s`, `%2$d` | 順序付きプレースホルダー |
| `$(...)` | Patchouliマクロ |
| `§` | Minecraftカラーコード |
| `\n` | 改行 |

---

## Minecraft公式用語

### 色名

| 英語 | 日本語 |
|------|--------|
| White | 白色 |
| Orange | 橙色 |
| Magenta | 赤紫色 |
| Light Blue | 空色 |
| Yellow | 黄色 |
| Lime | 黄緑色 |
| Pink | 桃色 |
| Gray | 灰色 |
| Light Gray | 薄灰色 |
| Cyan | 青緑色 |
| Purple | 紫色 |
| Blue | 青色 |
| Brown | 茶色 |
| Green | 緑色 |
| Red | 赤色 |
| Black | 黒色 |

### 基本用語

| 英語 | 日本語 |
|------|--------|
| Redstone | レッドストーン |
| Nether | ネザー |
| End | エンド |
| Overworld | オーバーワールド |
| Experience | 経験値 |
| Durability | 耐久値 |
| Enchantment | エンチャント |

---

## Patchouliマクロ

| マクロ | 意味 |
|--------|------|
| `$(br)` | 改行 |
| `$(br2)` | 空行 |
| `$(item)` | アイテム名ハイライト |
| `$(thing)` | 重要事項ハイライト |
| `$(l:path)text$(/)` | 内部リンク |
