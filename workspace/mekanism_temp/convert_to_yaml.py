#!/usr/bin/env python3
"""
Mekanismの言語ファイル（JSON）をYAML形式の翻訳データに変換するスクリプト
"""
import json
import yaml
from pathlib import Path
from collections import defaultdict

def categorize_key(key: str) -> str:
    """キーから適切なカテゴリを抽出"""
    parts = key.split('.')
    if len(parts) >= 2:
        return parts[0]  # 最初のセグメントをカテゴリとする
    return "other"

def load_json(filepath: Path) -> dict:
    """JSONファイルを読み込む"""
    with open(filepath, 'r', encoding='utf-8') as f:
        return json.load(f)

def create_translation_yaml(en_file: Path, ja_file: Path, module_name: str, version: str = "10.4.x") -> dict:
    """英語と日本語のJSONから翻訳YAMLデータを作成"""
    en_data = load_json(en_file)
    ja_data = load_json(ja_file)

    # カテゴリ別に分類
    categories = defaultdict(list)

    for key, en_value in sorted(en_data.items()):
        ja_value = ja_data.get(key, "")
        category = categorize_key(key)

        translation_entry = {
            'key': key,
            'en': en_value,
            'ja': ja_value
        }

        categories[category].append(translation_entry)

    # 統計情報を収集
    total_entries = len(en_data)
    translated_entries = sum(1 for key in en_data.keys() if ja_data.get(key))
    translation_rate = (translated_entries / total_entries * 100) if total_entries > 0 else 0

    # YAMLデータ構造を作成
    result = {
        'module': module_name,
        'version': version,
        'statistics': {
            'total_entries': total_entries,
            'translated_entries': translated_entries,
            'translation_rate': f"{translation_rate:.1f}%"
        },
        'categories': {}
    }

    # カテゴリ別にデータを追加
    for category, entries in sorted(categories.items()):
        result['categories'][category] = {
            'count': len(entries),
            'translations': entries
        }

    return result

def save_yaml(data: dict, output_file: Path):
    """YAMLファイルとして保存"""
    with open(output_file, 'w', encoding='utf-8') as f:
        yaml.dump(data, f, allow_unicode=True, sort_keys=False, default_flow_style=False, width=120)

def main():
    """メイン処理"""
    temp_dir = Path(__file__).parent
    output_dir = temp_dir / 'output'
    output_dir.mkdir(exist_ok=True)

    # モジュール定義
    modules = [
        ('core', 'core_en_us.json', 'core_ja_jp.json', 'Mekanism Core'),
        ('generators', 'generators_en_us.json', 'generators_ja_jp.json', 'Mekanism Generators'),
        ('tools', 'tools_en_us.json', 'tools_ja_jp.json', 'Mekanism Tools'),
        ('additions', 'additions_en_us.json', 'additions_ja_jp.json', 'Mekanism Additions'),
    ]

    all_modules = []

    for module_id, en_filename, ja_filename, module_name in modules:
        print(f"Processing {module_name}...")

        en_file = temp_dir / en_filename
        ja_file = temp_dir / ja_filename

        if not en_file.exists():
            print(f"  Warning: {en_filename} not found, skipping...")
            continue

        if not ja_file.exists():
            print(f"  Warning: {ja_filename} not found, skipping...")
            continue

        # YAML変換
        yaml_data = create_translation_yaml(en_file, ja_file, module_name)

        # 個別モジュールとして保存
        output_file = output_dir / f'mekanism_{module_id}_translations.yaml'
        save_yaml(yaml_data, output_file)
        print(f"  Saved to {output_file}")
        print(f"  Total entries: {yaml_data['statistics']['total_entries']}")
        print(f"  Translated: {yaml_data['statistics']['translated_entries']} ({yaml_data['statistics']['translation_rate']})")

        # 統合用データに追加
        all_modules.append({
            'id': module_id,
            'name': module_name,
            'version': yaml_data['version'],
            'statistics': yaml_data['statistics'],
            'categories': yaml_data['categories']
        })

    # 全モジュールを統合したYAMLを作成
    integrated_data = {
        'mod_id': 'mekanism',
        'mod_name': 'Mekanism',
        'minecraft_version': '1.20.x',
        'modules': all_modules
    }

    integrated_file = output_dir / 'mekanism_all_translations.yaml'
    save_yaml(integrated_data, integrated_file)
    print(f"\nIntegrated file saved to {integrated_file}")

    # サマリー
    total_all = sum(m['statistics']['total_entries'] for m in all_modules)
    translated_all = sum(m['statistics']['translated_entries'] for m in all_modules)
    print(f"\n=== Summary ===")
    print(f"Total modules: {len(all_modules)}")
    print(f"Total entries: {total_all}")
    print(f"Total translated: {translated_all} ({translated_all/total_all*100:.1f}%)")

if __name__ == '__main__':
    main()
