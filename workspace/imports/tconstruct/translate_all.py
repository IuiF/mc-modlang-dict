#!/usr/bin/env python3
import json
import yaml
import sys
from pathlib import Path

# 翻訳辞書（Tinkers用語に基づく）
TRANSLATIONS = {
    # 素材名（金属）
    "Cobalt": "コバルト",
    "Manyullyn": "マニュリン",
    "Hepatizon": "ヘパティゾン",
    "Queens Slime": "女王スライム",
    "Pig Iron": "ピグアイアン",
    "Rose Gold": "ローズゴールド",
    "Knightslime": "ナイトスライム",
    
    # ツールパーツ
    "Pick Head": "ピッケルヘッド",
    "Pickaxe Head": "ピッケルヘッド",
    "Axe Head": "斧ヘッド",
    "Small Axe Head": "小さな斧ヘッド",
    "Broad Axe Head": "幅広斧ヘッド",
    "Shovel Head": "シャベルヘッド",
    "Sword Blade": "剣の刀身",
    "Small Blade": "小さな刀身",
    "Broad Blade": "幅広刀身",
    "Tool Binding": "ツール接合部",
    "Tool Handle": "ツール柄",
    "Tool Rod": "ツール棒",
    "Tough Handle": "頑丈な柄",
    "Large Plate": "大きなプレート",
    "Hammer Head": "ハンマーヘッド",
    "Bow Limb": "弓のアーム",
    "Bow Grip": "弓のグリップ",
    "Bowstring": "弓弦",
    "Arrow Head": "矢じり",
    "Arrow Shaft": "矢の棒",
    "Fletching": "矢羽根",
    "Repair Kit": "修理キット",
    "Plating": "プレート装甲",
    "Helmet Plating": "ヘルメット装甲",
    "Chestplate Plating": "チェストプレート装甲",
    "Chest Plating": "チェストプレート装甲",
    "Leggings Plating": "レギンス装甲",
    "Leg Plating": "レギンス装甲",
    "Boots Plating": "ブーツ装甲",
    "Boot Plating": "ブーツ装甲",
    "Shield Core": "盾の芯",
    
    # ツール名
    "Pickaxe": "ピッケル",
    "Sledge Hammer": "大ハンマー",
    "Vein Hammer": "鉱脈ハンマー",
    "Mattock": "マトック",
    "Excavator": "掘削機",
    "Broad Axe": "幅広斧",
    "Hand Axe": "手斧",
    "Kama": "鎌",
    "Scythe": "大鎌",
    "Dagger": "短剣",
    "Sword": "剣",
    "Cleaver": "包丁",
    
    # 型紙・鋳型
    "Pattern": "型紙",
    "Cast": "鋳型",
    
    # 材料アイテム
    "Ingot": "インゴット",
    "Nugget": "塊",
    "Block": "ブロック",
    "Dust": "粉",
    "Gem": "宝石",
    "Crystal": "結晶",
    "Shard": "欠片",
    "Rod": "棒",
    "Plate": "プレート",
    "Small Plate": "小さなプレート",
    "Gear": "ギア",
    "Coin": "コイン",
    "Wire": "ワイヤー",
    "Scale": "鱗",
    
    # Seared/Scorched
    "Seared": "焼成",
    "Scorched": "焦熱",
    
    # 液体
    "Molten": "溶融",
    "Liquid": "液体",
    
    # ブロック/建材
    "Stone": "石",
    "Cobblestone": "丸石",
    "Bricks": "レンガ",
    "Paver": "舗装",
    "Glass": "ガラス",
    "Slab": "ハーフブロック",
    "Stairs": "階段",
    "Wall": "塀",
    "Pane": "板ガラス",
    "Lantern": "ランタン",
    "Fence": "フェンス",
    "Fence Gate": "フェンスゲート",
    "Door": "ドア",
    "Trapdoor": "トラップドア",
    "Button": "ボタン",
    "Pressure Plate": "感圧板",
    "Planks": "板材",
    "Log": "原木",
    "Wood": "木",
    "Sign": "看板",
    
    # 形容詞
    "Cracked": "ひび割れた",
    "Fancy": "装飾",
    "Triangle": "三角",
    "Tinted": "色付き",
    "Soul": "ソウル",
    "Chiseled": "模様入り",
    "Smooth": "滑らかな",
    "Polished": "磨かれた",
    "Stripped": "樹皮を剥いだ",
    "Congealed": "凝固",
    "Tall": "高い",
    
    # コンポーネント
    "Controller": "コントローラー",
    "Drain": "排出口",
    "Duct": "ダクト",
    "Chute": "投入口",
    "Tank": "タンク",
    "Fuel Tank": "燃料タンク",
    "Faucet": "蛇口",
    "Channel": "水路",
    "Alloyer": "合金炉",
    "Casting Table": "鋳造台",
    "Casting Basin": "鋳造盤",
    
    # スライム
    "Slime": "スライム",
    "Slimeball": "スライムボール",
    "Sapling": "苗木",
    "Leaves": "葉",
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
    "Earth": "大地",
    
    # 本
    "Materials and You": "素材と君",
    "Puny Smelting": "小さな精錬",
    "Mighty Smelting": "偉大なる精錬",
    "Tinkers' Gadgetry": "ティンカーの道具術",
    "Fantastic Foundry": "素晴らしき鋳造所",
    "Encyclopedia": "大百科",
    
    # その他
    "Reinforcement": "補強",
    "Modifier": "モディファイア",
    "Book": "本",
    "Helmet": "ヘルメット",
    "Chestplate": "チェストプレート",
    "Leggings": "レギンス",
    "Boots": "ブーツ",
    "Shield": "盾",
    "Travelers": "旅人の",
    "Slotless": "スロットレス",
}

def translate_text(text):
    """英語テキストを日本語に翻訳（基本的な置換）"""
    result = text
    
    # 複合語を優先的に置換（長い順）
    for en, ja in sorted(TRANSLATIONS.items(), key=lambda x: len(x[0]), reverse=True):
        result = result.replace(en, ja)
    
    return result

def process_category(data, prefix, category_name):
    """特定のカテゴリの翻訳を処理"""
    entries = {}
    for key, value in data.items():
        if key.startswith(prefix):
            entries[key] = value
    
    translations = []
    for key, source in sorted(entries.items()):
        target = translate_text(source)
        translations.append({
            'key': key,
            'source': source,
            'target': target if target != source else f"[要翻訳] {source}"
        })
    
    return {
        'mod_id': 'tconstruct',
        'version': '3.8.x',
        'mc_version': '1.20.1',
        'category': category_name,
        'total_entries': len(translations),
        'translations': translations
    }

def main():
    input_file = '/home/iuif/dev/minecraft-mod-dictionary/workspace/imports/tconstruct/en_us.json'
    output_dir = Path('/home/iuif/dev/minecraft-mod-dictionary/workspace/imports/tconstruct')
    
    with open(input_file, 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    # カテゴリごとに処理
    categories = [
        ('block.tconstruct.', 'blocks'),
        ('item.tconstruct.', 'items'),
        ('material.tconstruct.', 'materials'),
        ('modifier.tconstruct.', 'modifiers'),
        ('fluid.tconstruct.', 'fluids'),
        ('book.tconstruct.', 'books'),
    ]
    
    for prefix, category in categories:
        result = process_category(data, prefix, category)
        output_file = output_dir / f'{category}_translation.yaml'
        
        with open(output_file, 'w', encoding='utf-8') as f:
            yaml.dump(result, f, allow_unicode=True, sort_keys=False, width=120, default_flow_style=False)
        
        print(f"✓ {category}: {result['total_entries']} entries → {output_file}")

if __name__ == '__main__':
    main()
