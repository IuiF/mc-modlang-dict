#!/bin/bash

MOD="herbsandharvest"
ROUND=4

while true; do
    # ステータス確認
    STATUS=$(./moddict translate -mod $MOD -status)
    PENDING=$(echo "$STATUS" | grep "Pending:" | awk '{print $NF}')
    PROGRESS=$(echo "$STATUS" | grep "Progress:" | awk '{print $NF}')
    
    echo "=== ラウンド$ROUND: Pending=$PENDING, Progress=$PROGRESS ==="
    
    # pendingがないなら終了
    if [ "$PENDING" = "0" ]; then
        echo "完了！"
        ./moddict translate -mod $MOD -status
        break
    fi
    
    # 20件エクスポート
    ./moddict translate -mod $MOD -export /tmp/pending_$ROUND.json -limit 20
    
    # JSONから翻訳を自動生成（簡易版）
    python3 << 'EOFPYTHON'
import json
import sys

with open('/tmp/pending_' + str($ROUND) + '.json', 'r') as f:
    data = json.load(f)

# 翻訳パターン辞書
patterns = {
    "mustard_placeable_herb": "カラシナ",
    "onion_basket": "玉ねぎのバスケット",
    "onion_plant": "玉ねぎの種",
    "oregano_herb": "オレガノの種",
    "oregano_placeable_herb": "オレガノ",
    "parmesan_leaf_herb": "パセリの種",
    "parmesan_leaf_placeable_herb": "パセリ",
    "pea_basket": "豌豆のバスケット",
    "pea_plant": "豌豆の種",
    "pepper_basket": "唐辛子のバスケット",
    "pepper_plant": "唐辛子の種",
    "pineapple_basket": "パイナップルのバスケット",
    "pineapple_fruit_leaves": "パイナップルの実のある葉",
    "pineapple_fruit_sapling": "パイナップルの苗木",
    "pomegranate_basket": "ザクロのバスケット",
    "pomegranate_fruit_leaves": "ザクロの実のある葉",
    "pomegranate_fruit_sapling": "ザクロの苗木",
    "potato_basket": "じゃがいものバスケット",
    "radish_basket": "大根のバスケット",
    "radish_plant": "大根の種",
    "raspberry_basket": "ラズベリーのバスケット",
    "raspberry_plant": "ラズベリーの種",
    "red_pepper_basket": "赤唐辛子のバスケット",
    "red_pepper_plant": "赤唐辛子の種",
    "rice_basket": "米のバスケット",
    "rice_plant": "米の種",
    "rosemary_herb": "ローズマリーの種",
    "rosemary_placeable_herb": "ローズマリー",
    "soybean_basket": "大豆のバスケット",
    "soybean_plant": "大豆の種",
    "spinach_basket": "ホウレン草のバスケット",
    "spinach_plant": "ホウレン草の種",
    "strawberry_basket": "いちごのバスケット",
    "strawberry_plant": "いちごの種",
    "sunflower_plant": "ひまわりの種",
    "thyme_herb": "タイムの種",
    "thyme_placeable_herb": "タイム",
    "tomato_basket": "トマトのバスケット",
    "tomato_plant": "トマトの種",
    "turnip_basket": "かぶのバスケット",
    "turnip_plant": "かぶの種",
    "watermelon_basket": "スイカのバスケット",
    "watermelon_plant": "スイカの種",
    "item.herbsandharvest": "ハーブと収穫",
}

translated = {}
for key, en_text in data.items():
    # キー名から該当パターンを探す
    for pattern, ja_text in patterns.items():
        if pattern in key:
            translated[key] = ja_text
            break

# 既にくあるキーはスキップ
if translated:
    with open('/tmp/translated_' + str($ROUND) + '.json', 'w', encoding='utf-8') as f:
        json.dump(translated, f, ensure_ascii=False, indent=2)
    print(f"翻訳完了: {len(translated)}件")
EOFPYTHON
    
    # 翻訳結果をインポート
    if [ -f /tmp/translated_$ROUND.json ]; then
        ./moddict translate -mod $MOD -json /tmp/translated_$ROUND.json
    fi
    
    ROUND=$((ROUND + 1))
done
