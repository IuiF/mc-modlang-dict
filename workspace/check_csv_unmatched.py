#!/usr/bin/env python3
"""Check CSV entries that don't match any DB source"""

import csv
import sqlite3

DB_PATH = "/home/iuif/dev/minecraft-mod-dictionary/moddict.db"
CSV_PATH = "/home/iuif/dev/minecraft-mod-dictionary/workspace/exports/all_mods_ja_jp.csv"

def main():
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()

    # Load all DB sources as set of (key, source_text)
    print("Loading DB sources...")
    cursor.execute("SELECT key, source_text FROM translation_sources")
    db_sources = set()
    for key, source_text in cursor.fetchall():
        db_sources.add((key, source_text))
    print(f"Loaded {len(db_sources)} sources from DB")

    # Check CSV
    print("Checking CSV...")
    csv_total = 0
    matched = 0
    unmatched = []

    with open(CSV_PATH, 'r', encoding='utf-8-sig') as f:
        reader = csv.DictReader(f)
        for row in reader:
            csv_total += 1
            key = row.get('key', '').strip()
            source = row.get('source_text', '').strip()
            target = row.get('target_text', '').strip()

            if (key, source) in db_sources:
                matched += 1
            else:
                unmatched.append({
                    'key': key,
                    'source': source[:50] + '...' if len(source) > 50 else source,
                    'target': target[:50] + '...' if len(target) > 50 else target
                })

    print(f"\n=== Results ===")
    print(f"CSV total: {csv_total}")
    print(f"Matched with DB: {matched}")
    print(f"Unmatched: {len(unmatched)}")

    if unmatched:
        print(f"\n=== Sample unmatched entries (first 30) ===")
        # Group by key prefix to see patterns
        prefixes = {}
        for entry in unmatched:
            prefix = entry['key'].split('.')[0] if '.' in entry['key'] else entry['key']
            prefixes[prefix] = prefixes.get(prefix, 0) + 1

        print("\nKey prefix distribution:")
        for prefix, count in sorted(prefixes.items(), key=lambda x: -x[1])[:20]:
            print(f"  {prefix}: {count}")

        print("\nSample entries:")
        for entry in unmatched[:30]:
            print(f"  Key: {entry['key']}")
            print(f"  Source: {entry['source']}")
            print(f"  Target: {entry['target']}")
            print()

    # Save full unmatched list
    if unmatched:
        with open('/home/iuif/dev/minecraft-mod-dictionary/workspace/unmatched_csv.txt', 'w') as f:
            for entry in unmatched:
                f.write(f"{entry['key']}\t{entry['source']}\t{entry['target']}\n")
        print(f"\nFull list saved to workspace/unmatched_csv.txt")

    conn.close()

if __name__ == "__main__":
    main()
