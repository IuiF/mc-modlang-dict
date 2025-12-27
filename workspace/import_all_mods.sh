#!/bin/bash
cd /home/iuif/dev/minecraft-mod-dictionary

count=0
success=0
total=$(ls -1 workspace/all_mods/*.jar 2>/dev/null | wc -l)
echo "Importing $total jar files..."

for jar in workspace/all_mods/*.jar; do
    if [ -f "$jar" ]; then
        ((count++))
        result=$(./moddict import -jar "$jar" 2>&1)
        if echo "$result" | grep -q "Detected\|Already\|Imported"; then
            ((success++))
            if [ $((count % 100)) -eq 0 ]; then
                echo "[$count/$total] Progress: $success imported"
            fi
        fi
    fi
done

echo "=============================="
echo "Completed: $success / $total mods imported"
