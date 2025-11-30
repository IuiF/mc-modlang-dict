#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Convert translated JSON data to clean YAML format
"""

import json
import yaml
from pathlib import Path

def main():
    # Load translated data
    json_path = Path("/home/iuif/dev/minecraft-mod-dictionary/workspace/bloodmagic_alchemy_table/translated_data.json")

    with open(json_path, 'r', encoding='utf-8') as f:
        all_entries = json.load(f)

    # Prepare YAML structure
    yaml_data = {
        "mod_id": "bloodmagic",
        "mc_version": "1.16.3",
        "source_language": "en_us",
        "target_language": "ja_jp",
        "category": "alchemy_table",
        "translator": "Claude Code",
        "translation_date": "2025-11-29",
        "total_entries": len(all_entries),
        "entries": []
    }

    for entry in all_entries:
        entry_obj = {
            "id": entry["file"].replace(".json", ""),
            "file": entry["file"],
            "name": entry["name"],
            "icon": entry["icon"],
            "category": entry["category"],
            "page_count": len(entry["pages"]),
            "pages": entry["pages"]
        }

        yaml_data["entries"].append(entry_obj)

    # Save as YAML
    output_path = Path("/home/iuif/dev/minecraft-mod-dictionary/data/translations/bloodmagic/patchouli/alchemy_table_complete.yaml")

    # Create directory if not exists
    output_path.parent.mkdir(parents=True, exist_ok=True)

    with open(output_path, 'w', encoding='utf-8') as f:
        # Write header manually for better formatting
        f.write("# Blood Magic - Patchouli Entries Translation: alchemy_table Category\\n")
        f.write("# Blood Magic - Patchouli エントリ翻訳: alchemy_table カテゴリ\\n")
        f.write("# Complete translation of all 37 entries\\n")
        f.write("\\n")

        # Write YAML (excluding entries for now, will format manually)
        yaml.dump(yaml_data, f, allow_unicode=True, default_flow_style=False, sort_keys=False)

    print(f"YAML file created: {output_path}")
    print(f"Total entries: {len(all_entries)}")

if __name__ == "__main__":
    main()
