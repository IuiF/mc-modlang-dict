#!/bin/bash
cd /home/iuif/dev/minecraft-mod-dictionary

count=0
moved=0
total=$(ls -1 workspace/all_mods/*.jar 2>/dev/null | wc -l)
echo "Checking $total jar files..."

for jar in workspace/all_mods/*.jar; do
    if [ -f "$jar" ]; then
        ((count++))
        result=$(./moddict import -jar "$jar" 2>&1)

        # インポート成功またはAlready importedならスキップ
        if echo "$result" | grep -q "Detected\|Already\|Imported"; then
            : # 何もしない
        else
            # 言語ファイルがないModを移動
            mv "$jar" workspace/no_lang_mods/
            ((moved++))
        fi

        if [ $((count % 200)) -eq 0 ]; then
            echo "[$count/$total] Moved: $moved"
        fi
    fi
done

echo "=============================="
echo "Completed: $moved mods moved to workspace/no_lang_mods/"
