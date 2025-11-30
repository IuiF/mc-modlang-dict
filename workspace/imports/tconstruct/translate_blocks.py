#!/usr/bin/env python3
import json
import yaml
import sys

# 翻訳辞書（Tinkers用語に基づく）
TRANSLATIONS = {
    # Seared (焼成) 関連
    "Seared": "焼成",
    "Smeltery": "精錬炉",
    
    # Scorched (焦熱) 関連
    "Scorched": "焦熱",
    "Foundry": "鋳造所",
    
    # ブロックタイプ
    "Stone": "石",
    "Cobblestone": "丸石",
    "Bricks": "レンガ",
    "Paver": "舗装",
    "Glass": "ガラス",
    "Slab": "ハーフブロック",
    "Stairs": "階段",
    "Wall": "塀",
    "Pane": "板ガラス",
    
    # 形容詞
    "Cracked": "ひび割れた",
    "Fancy": "装飾",
    "Triangle": "三角",
    "Tinted": "色付き",
    "Soul": "ソウル",
    "Chiseled": "模様入り",
    "Smooth": "滑らかな",
    "Polished": "磨かれた",
    
    # Smeltery/Foundry コンポーネント
    "Controller": "コントローラー",
    "Drain": "排出口",
    "Duct": "ダクト",
    "Chute": "投入口",
    "Tank": "タンク",
    "Fuel Tank": "燃料タンク",
    "Lantern": "ランタン",
    "Faucet": "蛇口",
    "Channel": "水路",
    "Alloyer": "合金炉",
    
    # Casting関連
    "Casting Table": "鋳造台",
    "Casting Basin": "鋳造盤",
    
    # 木材関連
    "Planks": "板材",
    "Log": "原木",
    "Wood": "木",
    "Stripped": "樹皮を剥いだ",
    "Fence": "フェンス",
    "Fence Gate": "フェンスゲート",
    "Door": "ドア",
    "Trapdoor": "トラップドア",
    "Button": "ボタン",
    "Pressure Plate": "感圧板",
    "Sign": "看板",
    
    # スライム関連
    "Congealed Slime": "凝固スライム",
    "Slime": "スライム",
    "Sapling": "苗木",
    "Leaves": "葉",
    "Tall": "高い",
    "Grass": "草",
    "Fern": "シダ",
    "Dirt": "土",
    
    # 色
    "Orange": "オレンジ色",
    "Green": "緑色",
    "Blue": "青色",
    "Purple": "紫色",
    "Magenta": "マゼンタ色",
    "Yellow": "黄色",
    "Pink": "ピンク色",
    "Red": "赤色",
    "Sky": "空色",
    "Ender": "エンダー",
    "Ichor": "エイカー",
    "Blood": "血",
    
    # その他
    "Mud": "泥",
    "Cake": "ケーキ",
    "Magma": "マグマ",
    "Graveyard": "墓地",
    "Earth": "地球",
    "Sky": "空",
    "Ender": "エンダー",
}

def translate_text(text):
    """英語テキストを日本語に翻訳"""
    result = text
    
    # 複合語を優先的に置換
    for en, ja in sorted(TRANSLATIONS.items(), key=lambda x: len(x[0]), reverse=True):
        result = result.replace(en, ja)
    
    return result

def translate_blocks(input_file, output_file):
    """ブロック翻訳を処理"""
    with open(input_file, 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    # ブロック関連のキーのみ抽出
    blocks = {}
    for key, value in data.items():
        if key.startswith('block.tconstruct.'):
            blocks[key] = value
    
    # 翻訳データ作成
    translations = []
    for key, source in sorted(blocks.items()):
        target = translate_text(source)
        translations.append({
            'key': key,
            'source': source,
            'target': target
        })
    
    # YAML出力用のデータ構造
    output = {
        'mod_id': 'tconstruct',
        'version': '3.8.x',
        'mc_version': '1.20.1',
        'category': 'blocks',
        'total_entries': len(translations),
        'translations': translations[:100]  # 最初の100エントリのみ（サンプル）
    }
    
    with open(output_file, 'w', encoding='utf-8') as f:
        yaml.dump(output, f, allow_unicode=True, sort_keys=False, width=120)
    
    print(f"Processed {len(translations)} block entries")
    print(f"Saved first 100 entries to {output_file}")

if __name__ == '__main__':
    translate_blocks(
        '/home/iuif/dev/minecraft-mod-dictionary/workspace/imports/tconstruct/en_us.json',
        '/home/iuif/dev/minecraft-mod-dictionary/workspace/imports/tconstruct/blocks_translation.yaml'
    )
