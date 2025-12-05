#!/bin/bash

mkdir -p translations

for mod in $(sqlite3 moddict.db "SELECT id FROM mods ORDER BY id"); do
  echo "Exporting: $mod"
  ./moddict export -mod "$mod" -format csv -out translations 2>&1 | tail -1
done

echo ""
echo "=== Export Complete ==="
ls -la translations/*.csv | wc -l
echo "CSV files created in translations/"
