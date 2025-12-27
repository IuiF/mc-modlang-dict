#!/bin/bash
cd /home/iuif/dev/minecraft-mod-dictionary

count=0
success=0
total=$(ls -1 translations/*_ja_jp.csv 2>/dev/null | grep -v all_mods | wc -l)
echo "Restoring translations from $total CSV files..."

for csv in translations/*_ja_jp.csv; do
    # Skip all_mods (it's a combined file)
    if [[ "$csv" == *"all_mods"* ]]; then
        continue
    fi

    ((count++))

    # Extract mod_id from filename (e.g., translations/mekanism_ja_jp.csv -> mekanism)
    filename=$(basename "$csv")
    mod_id="${filename%_ja_jp.csv}"

    # Check if mod exists in DB
    mod_exists=$(sqlite3 moddict.db "SELECT COUNT(*) FROM mods WHERE id='$mod_id';" 2>/dev/null)
    if [ "$mod_exists" != "1" ]; then
        echo "[$count/$total] SKIP: $mod_id (mod not in DB)"
        continue
    fi

    # Convert CSV to JSON
    json_file="/tmp/${mod_id}_restore.json"

    # Use Python to convert CSV to JSON (handles special characters better)
    python3 -c "
import csv
import json
import sys

translations = {}
with open('$csv', 'r', encoding='utf-8-sig') as f:
    reader = csv.DictReader(f)
    for row in reader:
        key = row.get('key', '')
        target = row.get('target_text', '')
        if key and target:
            translations[key] = target

with open('$json_file', 'w', encoding='utf-8') as f:
    json.dump(translations, f, ensure_ascii=False, indent=2)
print(len(translations))
" 2>/dev/null

    if [ -f "$json_file" ]; then
        result=$(./moddict translate -mod "$mod_id" -json "$json_file" 2>&1)
        if echo "$result" | grep -q "Applied\|Updated\|imported"; then
            ((success++))
            applied=$(echo "$result" | grep -oE '[0-9]+ translation' | head -1 | grep -oE '[0-9]+')
            echo "[$count/$total] OK: $mod_id ($applied translations)"
        else
            echo "[$count/$total] PARTIAL: $mod_id"
        fi
        rm -f "$json_file"
    else
        echo "[$count/$total] ERROR: $mod_id (CSV conversion failed)"
    fi
done

echo "=============================="
echo "Completed: $success / $total mods restored"
