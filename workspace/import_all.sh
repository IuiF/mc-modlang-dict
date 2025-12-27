#!/bin/bash
cd /home/iuif/dev/minecraft-mod-dictionary

count=0
total=$(ls -1 workspace/imports/mods/*.jar 2>/dev/null | wc -l)
echo "Importing $total jar files..."

for jar in workspace/imports/mods/*.jar; do
    if [ -f "$jar" ]; then
        result=$(./moddict import -jar "$jar" 2>&1)
        if echo "$result" | grep -q "Detected\|Already"; then
            ((count++))
            echo "[$count/$total] OK: $(basename "$jar")"
        else
            echo "[$count/$total] SKIP: $(basename "$jar")"
        fi
    fi
done

echo "============================="
echo "Completed: $count / $total mods imported"
