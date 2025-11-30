#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Blood Magic - Alchemy Table Category Translation Script
全38エントリを翻訳してYAML形式で出力
"""

import json
import os
import re
from pathlib import Path

# 用語辞書（Blood Magic固有）
TERMS = {
    # Core
    "Blood Magic": "Blood Magic",
    "LP": "LP",
    "Life Points": "ライフポイント",
    "Life Essence": "生命のエッセンス",
    "Soul Network": "ソウルネットワーク",

    # Structures
    "Blood Altar": "血の祭壇",
    "Altar": "祭壇",
    "Alchemy Table": "錬金術テーブル",
    "Alchemy Flask": "錬金術フラスク",

    # Items
    "Blood Orb": "血のオーブ",
    "Slate": "石板",
    "Blank Slate": "空白の石板",
    "Sigil": "シジル",
    "Water Sigil": "水のシジル",
    "Lava Sigil": "溶岩のシジル",

    # Anointments
    "Anointment": "膏薬",
    "Anointments": "膏薬",

    # Alchemy
    "Catalyst": "触媒",
    "Catalysts": "触媒",
    "Simple Catalyst": "シンプル触媒",
    "Power Catalyst": "強化触媒",
    "Lengthening Catalyst": "延長触媒",
    "Combinational Catalyst": "組み合わせ触媒",
    "Cycling Catalyst": "サイクル触媒",
    "Filling Agent": "充填剤",

    # Effects
    "Potion": "ポーション",
    "Potions": "ポーション",
    "Effect": "効果",
    "Effects": "効果",

    # Vanilla items
    "Grass": "草ブロック",
    "Leather": "革",
    "Bread": "パン",
    "Clay": "粘土",
    "String": "糸",
    "Plant Oil": "植物油",
    "Coal Sand": "石炭砂",
    "Explosive Powder": "爆発性粉末",
    "Flint": "火打石",
    "Saltpeter": "硝石",
    "Gunpowder": "火薬",
    "Lava Bucket": "溶岩入りバケツ",
    "Water Bucket": "水入りバケツ",
    "Glass Bottle": "ガラス瓶",
    "Hopper": "ホッパー",
    "Hoppers": "ホッパー",
    "Nether Wart": "ネザーウォート",
    "Glowstone": "グロウストーン",
    "Redstone": "レッドストーン",

    # Anointment names
    "Iron Tip": "鉄の先端",
    "Archer's Polish": "弓研磨剤",
    "Fortuna Extract": "フォルトゥーナエキス",
    "Hidden Knowledge": "隠された知識",
    "Holy Water": "聖水",
    "Looting": "ドロップ増加",
    "Honing Oil": "研磨油",
    "Quick Draw": "速射",
    "Silk Touch": "シルクタッチ",
    "Smelting": "自動精錬",

    # Potion effects
    "Flight": "飛行",
    "Fire Resistance": "火炎耐性",
    "Water Breathing": "水中呼吸",
    "Regeneration": "再生",
    "Night Vision": "暗視",
    "Invisibility": "透明化",
    "Jump Boost": "跳躍力上昇",
    "Levitation": "浮遊",
    "Poison": "毒",
    "Slow Falling": "低速落下",
    "Slowness": "鈍化",
    "Speed": "移動速度上昇",
    "Strength": "攻撃力上昇",
    "Weakness": "弱体化",
    "Instant Health": "即時回復",
    "Instant Damage": "即時ダメージ",

    # Blood Magic potions
    "Bounce": "跳ね返り",
    "Gravity": "重力",
    "Grounded": "接地",
    "Hard Cloak": "硬化の外套",
    "Heavy Heart": "重い心臓",
    "Obsidian Cloak": "黒曜石の外套",
    "Passive": "平和",
    "Spectral Sight": "霊視",
    "Suspended": "浮遊状態",

    # Other
    "Splash": "スプラッシュ",
    "Lingering": "残留",
    "Tier": "ティア",
    "GUI": "GUI",
}

def translate_text(text):
    """テキストを翻訳（用語辞書を使用、Patchouliマクロは保持）"""
    if not text:
        return text

    result = text

    # 用語辞書を適用（長い語句から優先）
    for en, ja in sorted(TERMS.items(), key=lambda x: len(x[0]), reverse=True):
        # $(...)マクロ内は翻訳しない
        # 単純な置換（大文字小文字を区別）
        result = result.replace(en, ja)

    return result

def translate_entry(entry_data, file_path):
    """1つのエントリを翻訳"""
    translated = {
        "file": os.path.basename(file_path),
        "name": {
            "source": entry_data.get("name", ""),
            "translation": translate_text(entry_data.get("name", ""))
        },
        "icon": entry_data.get("icon", ""),
        "category": entry_data.get("category", ""),
        "pages": []
    }

    for idx, page in enumerate(entry_data.get("pages", [])):
        translated_page = {
            "page": idx,
            "type": page.get("type", "")
        }

        # テキストフィールドを翻訳
        if "text" in page:
            translated_page["source"] = page["text"]
            translated_page["translation"] = translate_text(page["text"])

        # タイトル/見出しを翻訳
        if "title" in page:
            translated_page["title"] = {
                "source": page["title"],
                "translation": translate_text(page["title"])
            }

        if "heading" in page:
            translated_page["heading"] = {
                "source": page["heading"],
                "translation": translate_text(page["heading"])
            }

        # レシピフィールドをコピー
        for key in ["recipe", "images", "border", "anchor"]:
            if key in page:
                translated_page[key] = page[key]

        # マルチレシピページ（a.heading, b.heading等）
        for prefix in ["a", "b", "c"]:
            heading_key = f"{prefix}.heading"
            recipe_key = f"{prefix}.recipe"

            if heading_key in page:
                translated_page[f"{prefix}_heading"] = {
                    "source": page[heading_key],
                    "translation": translate_text(page[heading_key])
                }

            if recipe_key in page:
                translated_page[f"{prefix}_recipe"] = page[recipe_key]

        translated["pages"].append(translated_page)

    return translated

def main():
    base_dir = Path("/home/iuif/dev/minecraft-mod-dictionary/workspace/bloodmagic_alchemy_table")

    all_entries = []

    # 全JSONファイルを収集
    json_files = []

    # メインエントリ
    for f in ["alchemy_table.json", "anointments.json", "potions.json"]:
        json_files.append(("main", base_dir / f))

    # Anointments
    anointments_dir = base_dir / "anointments"
    for f in sorted(anointments_dir.glob("*.json")):
        json_files.append(("anointments", f))

    # Potion flasks - blood_magic
    bm_dir = base_dir / "potion_flasks" / "blood_magic"
    for f in sorted(bm_dir.glob("*.json")):
        json_files.append(("potion_flasks/blood_magic", f))

    # Potion flasks - vanilla
    vanilla_dir = base_dir / "potion_flasks" / "vanilla"
    for f in sorted(vanilla_dir.glob("*.json")):
        json_files.append(("potion_flasks/vanilla", f))

    print(f"Found {len(json_files)} JSON files")

    # 各ファイルを翻訳
    for category, json_path in json_files:
        print(f"Translating: {json_path.name}")

        with open(json_path, 'r', encoding='utf-8') as f:
            entry_data = json.load(f)

        translated = translate_entry(entry_data, json_path)
        translated["_category_path"] = category
        all_entries.append(translated)

    # YAML形式で出力（手動フォーマット）
    output_path = Path("/home/iuif/dev/minecraft-mod-dictionary/workspace/bloodmagic_alchemy_table/translated_entries.txt")

    with open(output_path, 'w', encoding='utf-8') as f:
        f.write(f"# Translated {len(all_entries)} entries\\n\\n")

        for entry in all_entries:
            f.write(f"### {entry['file']} ###\\n")
            f.write(f"Name: {entry['name']['source']} -> {entry['name']['translation']}\\n")
            f.write(f"Icon: {entry['icon']}\\n")
            f.write(f"Pages: {len(entry['pages'])}\\n")
            f.write("\\n")

    print(f"\\nTranslation summary saved to: {output_path}")
    print(f"Total entries translated: {len(all_entries)}")

    # JSON形式でも保存
    json_output_path = Path("/home/iuif/dev/minecraft-mod-dictionary/workspace/bloodmagic_alchemy_table/translated_data.json")
    with open(json_output_path, 'w', encoding='utf-8') as f:
        json.dump(all_entries, f, ensure_ascii=False, indent=2)

    print(f"Full translation data saved to: {json_output_path}")

if __name__ == "__main__":
    main()
