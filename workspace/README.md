# Workspace Directory

このディレクトリは翻訳作業用のファイル格納場所です。

## ディレクトリ構成

```
workspace/
├── imports/     # インポート対象のModファイル
├── exports/     # エクスポート出力先
└── temp/        # 一時ファイル（自動削除）
```

## imports/

翻訳辞書に取り込むModファイル（.jar/.zip）を配置します。

### 使用方法

1. CurseForge/ModrinthからModファイルをダウンロード
2. このディレクトリに配置
3. `moddict import ./workspace/imports/mod_name.jar` でインポート

### ファイル命名規則（推奨）

```
{mod_id}-{mc_version}-{mod_version}.jar

例:
create-1.20.1-0.5.1f.jar
mekanism-1.20.1-10.4.0.jar
```

### 注意事項

- 著作権に注意：Modファイル自体はGitにコミットしない
- `.gitignore` で `*.jar` は除外済み

## exports/

翻訳データのエクスポート先です。

### 出力形式

- `{mod_id}/lang/ja_jp.json` - Minecraft言語ファイル形式
- `{mod_id}/terms.yaml` - 用語辞書形式
- `{mod_id}/diff.json` - バージョン差分

## temp/

一時ファイル格納場所。jar展開時などに使用。

- 自動的にクリーンアップされる
- 手動で削除しても問題なし
