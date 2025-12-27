#!/usr/bin/env python3
"""Restore translations from CSV to moddict.db"""

import csv
import sqlite3
import sys

DB_PATH = "/home/iuif/dev/minecraft-mod-dictionary/moddict.db"
CSV_PATH = "/home/iuif/dev/minecraft-mod-dictionary/workspace/exports/all_mods_ja_jp.csv"

def main():
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()

    # Load CSV translations
    print("Loading CSV...")
    csv_translations = {}
    with open(CSV_PATH, 'r', encoding='utf-8-sig') as f:
        reader = csv.DictReader(f)
        for row in reader:
            key = row.get('key', '').strip()
            source = row.get('source_text', '').strip()
            target = row.get('target_text', '').strip()
            if key and target:
                # Use (key, source_text) as composite key for matching
                csv_translations[(key, source)] = target

    print(f"Loaded {len(csv_translations)} translations from CSV")

    # Get all sources with their translation IDs
    print("Finding matching sources...")
    cursor.execute("""
        SELECT ts.id, ts.key, ts.source_text, t.id as trans_id, t.target_text, t.status
        FROM translation_sources ts
        JOIN translations t ON t.source_id = ts.id
        WHERE t.target_lang = 'ja_jp'
    """)

    rows = cursor.fetchall()
    print(f"Found {len(rows)} translation entries in DB")

    # Match and update
    updated = 0
    skipped_official = 0
    skipped_same = 0
    not_found = 0

    for source_id, key, source_text, trans_id, current_target, status in rows:
        csv_key = (key, source_text)

        if csv_key in csv_translations:
            new_target = csv_translations[csv_key]

            # Skip if already has official translation
            if status == 'official':
                skipped_official += 1
                continue

            # Skip if same text
            if current_target == new_target:
                skipped_same += 1
                continue

            # Update translation
            cursor.execute("""
                UPDATE translations
                SET target_text = ?, status = 'translated'
                WHERE id = ?
            """, (new_target, trans_id))
            updated += 1

            if updated % 10000 == 0:
                print(f"  Updated {updated} translations...")
                conn.commit()
        else:
            not_found += 1

    conn.commit()
    conn.close()

    print(f"\n=== Restore Complete ===")
    print(f"Updated: {updated}")
    print(f"Skipped (official): {skipped_official}")
    print(f"Skipped (same text): {skipped_same}")
    print(f"Not in CSV: {not_found}")

if __name__ == "__main__":
    main()
