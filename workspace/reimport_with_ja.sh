#!/bin/bash
cd /home/iuif/dev/minecraft-mod-dictionary

count=0
total=$(ls -1 workspace/all_mods/*.jar 2>/dev/null | wc -l)
echo "Re-importing $total jar files with official Japanese support..."

# Clear previous logs
> workspace/import_results.log
> workspace/zero_keys.log
> workspace/official_ja.log

total_official=0

for jar in workspace/all_mods/*.jar; do
    if [ -f "$jar" ]; then
        ((count++))
        result=$(./moddict import -jar "$jar" 2>&1)

        # Extract key count from result
        keys=$(echo "$result" | grep "Total keys:" | sed 's/.*Total keys: //')

        # Extract official Japanese count
        official=$(echo "$result" | grep "Official Japanese:" | sed 's/.*Official Japanese: //' | sed 's/ keys//')
        if [ -n "$official" ] && [ "$official" != "0" ]; then
            echo "$(basename "$jar")|$official" >> workspace/official_ja.log
            total_official=$((total_official + official))
        fi

        # Log result
        echo "$(basename "$jar")|$keys" >> workspace/import_results.log

        # Track zero-key mods
        if [ "$keys" = "0" ]; then
            echo "$jar" >> workspace/zero_keys.log
        fi

        if [ $((count % 200)) -eq 0 ]; then
            zero_count=$(wc -l < workspace/zero_keys.log 2>/dev/null || echo 0)
            echo "[$count/$total] Zero-key: $zero_count, Official JA total: $total_official"
        fi
    fi
done

echo "=============================="
echo "Completed: $count mods processed"
echo "Zero-key mods: $(wc -l < workspace/zero_keys.log)"
echo "Total official Japanese imported: $total_official"
