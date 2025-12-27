#!/bin/bash
cd /home/iuif/dev/minecraft-mod-dictionary

moved=0
while IFS= read -r filename; do
    src="workspace/all_mods/$filename"
    if [ -f "$src" ]; then
        mv "$src" workspace/no_lang_mods/
        ((moved++))
    fi
done < /tmp/zero_names.txt

echo "Moved $moved additional files"
