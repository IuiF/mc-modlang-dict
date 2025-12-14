# 翻訳一貫性ガイド

## 概要

大規模Mod（数千〜数万キー）の翻訳において、一貫性を保つための仕組みとプロセスを定義する。

## 一貫性の2つのレベル

### レベル1: 完全一致（同一source_text）

同じ英語テキストが複数のキーで使用されているケース。

```
例: Chiselの「Stained Glass」
  - tile.chisel.stained_glass_white.name: "Stained Glass"
  - tile.chisel.stained_glass_orange.name: "Stained Glass"
  - ... (計17キー)

→ すべて同じ「色付きガラス」に翻訳すべき
```

**解決策: インポート時自動伝播**
- 1つのキーに翻訳をインポートすると、同じsource_textを持つ他のpendingキーにも自動適用
- 翻訳者は意識する必要なし

### レベル2: 部品一致（共通要素）

異なる英語テキストだが、共通の要素（色、素材、接辞）を含むケース。

```
例: 色名「White」の揺れ
  ❌ "White Wool" → "白い羊毛"
  ❌ "White Glass" → "ホワイトガラス"  ← 揺れ！
  ❌ "White Concrete" → "白色コンクリート"  ← 揺れ！

  ✅ すべて「白い〜」または「白〜」で統一すべき
```

**解決策: 用語辞書（Terms）**
- 色名、素材名、接頭辞・接尾辞を事前定義
- 翻訳時に参照して一貫性を保つ

## 用語辞書の構造

### スコープ階層

```
global（全Mod共通）
  └─ Minecraft公式用語: 色名、バニラアイテム、共通接辞

category:{name}（カテゴリ共通）
  └─ 例: category:building → 建材系接辞

mod:{mod_id}（Mod固有）
  └─ 例: mod:chisel → Chisel特有の用語
```

### 優先順位

```
global (低) → category (中) → mod (高)
```

Mod固有の定義がある場合、それが優先される。

### 標準グローバル用語

#### 色名（16色）

| 英語 | 日本語 |
|------|--------|
| White | 白 |
| Orange | 橙 |
| Magenta | 赤紫 |
| Light Blue | 空色 |
| Yellow | 黄 |
| Lime | 黄緑 |
| Pink | 桃色 |
| Gray | 灰色 |
| Light Gray | 薄灰色 |
| Cyan | 青緑 |
| Purple | 紫 |
| Blue | 青 |
| Brown | 茶 |
| Green | 緑 |
| Red | 赤 |
| Black | 黒 |

#### 共通接尾辞

| 英語 | 日本語 |
|------|--------|
| Block | ブロック |
| Stairs | 階段 |
| Slab | ハーフブロック |
| Wall | 塀 |
| Fence | フェンス |
| Door | ドア |
| Trapdoor | トラップドア |
| Button | ボタン |
| Pressure Plate | 感圧板 |
| Pillar | 柱 |
| Bricks | レンガ |
| Tiles | タイル |

#### 共通接頭辞

| 英語 | 日本語 |
|------|--------|
| Chiseled | 彫られた / 模様入り |
| Polished | 磨かれた |
| Smooth | 滑らかな |
| Cracked | ひび割れた |
| Mossy | 苔むした |
| Weathered | 風化した |

## 翻訳プロセス

### フェーズ1: 準備（Mod固有用語の定義）

```bash
# 1. source_textを分析して頻出単語を確認
sqlite3 moddict.db "
SELECT
  CASE
    WHEN source_text LIKE '%Marble%' THEN 'Marble'
    WHEN source_text LIKE '%Limestone%' THEN 'Limestone'
    WHEN source_text LIKE '%Basalt%' THEN 'Basalt'
    -- 必要に応じて追加
  END as term,
  COUNT(*) as cnt
FROM translation_sources
WHERE mod_id = 'chisel'
GROUP BY term
HAVING term IS NOT NULL
ORDER BY cnt DESC;
"

# 2. Mod固有用語を決定（チームで合意）
# 例: Marble → 大理石, Limestone → 石灰岩, Basalt → 玄武岩
```

### フェーズ2: 用語辞書の作成

翻訳開始前に、用語辞書をJSONで作成：

```json
{
  "mod_id": "chisel",
  "terms": {
    "materials": {
      "Marble": "大理石",
      "Limestone": "石灰岩",
      "Basalt": "玄武岩",
      "Andesite": "安山岩",
      "Diorite": "閃緑岩",
      "Granite": "花崗岩"
    },
    "mod_specific": {
      "Antiblock": "アンチブロック",
      "Factory Block": "工場ブロック",
      "Laboratory Block": "研究所ブロック"
    }
  }
}
```

### フェーズ3: 翻訳実行

サブエージェントへの指示に用語辞書を含める：

```
[Mod名]の翻訳を行ってください。

【用語辞書 - 必ず従うこと】
色名: White→白, Orange→橙, Magenta→赤紫, Light Blue→空色,
      Yellow→黄, Lime→黄緑, Pink→桃色, Gray→灰色,
      Light Gray→薄灰色, Cyan→青緑, Purple→紫, Blue→青,
      Brown→茶, Green→緑, Red→赤, Black→黒

素材: Marble→大理石, Limestone→石灰岩, Basalt→玄武岩

接尾辞: Block→ブロック, Stairs→階段, Slab→ハーフブロック,
        Wall→塀, Pillar→柱, Bricks→レンガ

接頭辞: Chiseled→彫られた, Polished→磨かれた, Smooth→滑らかな

【翻訳ルール】
1. 上記用語辞書に従って一貫性を保つ
2. フォーマットコード保持: §, $(...)、%s, %d, \n
3. 空文字("")は禁止
4. 翻訳できないキーはJSONに含めない

【手順】
1. ./moddict translate -mod [mod_id] -status
2. ./moddict translate -mod [mod_id] -export /tmp/pending.json -limit 20
3. 用語辞書に従って翻訳
4. ./moddict translate -mod [mod_id] -json /tmp/translated.json
5. 繰り返し
```

### フェーズ4: インポートと自動伝播

翻訳をインポートすると：
1. 指定されたキーに翻訳が適用される
2. **同じsource_textを持つ他のpendingキーにも自動適用**（自動伝播）

```bash
# 例: "Stained Glass" を翻訳
# pending.json に key1, key2, key3 があり、すべて "Stained Glass"
# key1 だけ翻訳して "色付きガラス" をインポート
# → key2, key3 にも自動的に "色付きガラス" が適用される
```

### フェーズ5: 一貫性チェック

翻訳完了後、一貫性をチェック：

```bash
# 同じsource_textに異なる翻訳がないか確認
sqlite3 moddict.db "
SELECT ts.source_text, t.target_text, COUNT(*) as cnt
FROM translation_sources ts
JOIN translations t ON ts.id = t.source_id
WHERE ts.mod_id = 'chisel'
  AND t.status IN ('translated', 'verified')
  AND t.target_text IS NOT NULL
GROUP BY ts.source_text, t.target_text
HAVING (SELECT COUNT(DISTINCT t2.target_text)
        FROM translation_sources ts2
        JOIN translations t2 ON ts2.id = t2.source_id
        WHERE ts2.source_text = ts.source_text
          AND ts2.mod_id = ts.mod_id
          AND t2.status IN ('translated', 'verified')
          AND t2.target_text IS NOT NULL) > 1
ORDER BY ts.source_text;
"
```

## スケーラビリティ

### 大規模Mod対応

| Mod規模 | 用語辞書サイズ | 注意点 |
|---------|---------------|--------|
| 小（〜500キー） | 10〜20項目 | global + 少数のmod固有 |
| 中（500〜2000キー） | 20〜50項目 | カテゴリ別に整理 |
| 大（2000〜5000キー） | 50〜100項目 | 素材・パターンを網羅 |
| 特大（5000キー〜） | 100〜200項目 | 事前分析が重要 |

### 用語辞書のサイズを抑える理由

1. **LLMコンテキストの節約**: 毎バッチ添付するため、小さいほど良い
2. **翻訳者の認知負荷**: 多すぎると覚えられない
3. **メンテナンス性**: 少ないほど更新・管理が容易

### 用語辞書に含めるもの

✅ 含める:
- 複数回出現する共通要素（色、素材、接辞）
- Mod特有の固有名詞
- 揺れやすい表現

❌ 含めない:
- 1回しか出現しない単語
- 文脈で意味が変わる単語
- 一般的な動詞・形容詞

## まとめ

```
翻訳一貫性 = 用語辞書（事前定義）+ 自動伝播（インポート時）

用語辞書: 部品レベルの一貫性を担保（色、素材、接辞）
自動伝播: 完全一致の一貫性を担保（同一source_text）
```
