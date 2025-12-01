# バッチ翻訳エージェント

複数のModを効率的に翻訳するためのバッチ処理エージェント。
**各Modを100%完了してから次に進む。並列処理は同一Mod内のみ許可。**

---

## 重要: 並列処理のルール

### 許可される並列処理

- 同一Mod内の異なるカテゴリ（block/item/advancement等）
- 同一Modの翻訳ループ（100件以内のバッチ）

### 禁止される並列処理

- 複数Modの同時翻訳
- 60%完了で次のModに移行

---

## 一時ファイル管理（厳守）

### 並列処理時のファイル命名

同一Mod内で並列処理する場合、一時ファイルにサフィックスを付ける：

| エージェント | pending | translated |
|-------------|---------|------------|
| Agent 1 | `/tmp/moddict_pending_1.json` | `/tmp/moddict_translated_1.json` |
| Agent 2 | `/tmp/moddict_pending_2.json` | `/tmp/moddict_translated_2.json` |
| Agent 3 | `/tmp/moddict_pending_3.json` | `/tmp/moddict_translated_3.json` |

### 禁止

- `workspace/` 以下に一時ファイルを作成
- プロジェクトルートに一時ファイルを作成
- 完了後に一時ファイルを残す

### クリーンアップ

```bash
# 作業開始時・終了時
rm -f /tmp/moddict_*.json
```

---

## バッチ処理フロー

### フェーズ1: 対象Modのリストアップ

```bash
# 全Modのステータス確認
./moddict translate -status-all

# または個別確認
./moddict translate -mod [mod_id] -status
```

### フェーズ2: Mod単位で順次処理

```
for each mod in mod_list:
    1. インポート・公式翻訳適用
    2. 翻訳ループ（100%まで）
    3. 検証・クリーンアップ
    4. 完了報告
    5. 次のModへ
```

### フェーズ3: 最終報告

```
## バッチ翻訳完了報告

| Mod ID | 総キー数 | 公式 | 新規 | ステータス |
|--------|----------|------|------|------------|
| mod_a  | 500      | 300  | 200  | 100%完了   |
| mod_b  | 300      | 150  | 150  | 100%完了   |
| mod_c  | 200      | 0    | 200  | 100%完了   |

総計: 1000キー翻訳完了
```

---

## 並列翻訳ループ（同一Mod内）

大規模Mod（1000キー以上）の場合、同一Mod内で並列処理可能。

### 手順

1. **pendingを分割エクスポート**

```bash
# エージェント1用
./moddict translate -mod [mod_id] -export /tmp/moddict_pending_1.json -limit 50 -offset 0

# エージェント2用
./moddict translate -mod [mod_id] -export /tmp/moddict_pending_2.json -limit 50 -offset 50

# エージェント3用
./moddict translate -mod [mod_id] -export /tmp/moddict_pending_3.json -limit 50 -offset 100
```

2. **各エージェントで翻訳実行**

3. **順次インポート**（DBロック回避のため順次）

```bash
./moddict translate -mod [mod_id] -json /tmp/moddict_translated_1.json
./moddict translate -mod [mod_id] -json /tmp/moddict_translated_2.json
./moddict translate -mod [mod_id] -json /tmp/moddict_translated_3.json
```

4. **クリーンアップ**

```bash
rm -f /tmp/moddict_pending_*.json /tmp/moddict_translated_*.json
```

5. **ステータス確認**

```bash
./moddict translate -mod [mod_id] -status
```

---

## コンテキスト効率化

### 100件ごとにエージェントリセット

長時間の翻訳作業ではコンテキストが肥大化する。
100件の翻訳ごとに新しいエージェントを起動し、コンテキストをリセット。

### 最小限の報告

エージェントからメインへの報告は以下のみ：

```
Mod: [mod_id]
処理済み: [n]件
残り: [m]件
ステータス: [x]%
```

JSON内容や翻訳テキストは報告に含めない。

---

## エラー対応

### DBロックエラー

```bash
# 少し待ってリトライ
sleep 2
./moddict translate -mod [mod_id] -json /tmp/moddict_translated.json
```

### インポートエラー

```bash
./moddict repair
```

### 部分的な失敗

1. 成功した翻訳は保持される
2. 失敗したキーを再エクスポート
3. 再翻訳・再インポート

---

## 禁止事項

- 複数Modを同時に翻訳開始すること
- 60%で次のModに移行すること
- workspace/以下に一時ファイルを作成すること
- 連番なしで並列処理すること（ファイル衝突）
- JSON内容をメインに報告すること
- 一時ファイルを削除せずに終了すること
