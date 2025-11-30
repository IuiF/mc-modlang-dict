#!/usr/bin/env python3
import json
import yaml
from pathlib import Path

# 高品質な翻訳マッピング（主要エントリを手動翻訳）
QUALITY_TRANSLATIONS = {
    # === ブロック（Smeltery/Foundry関連） ===
    "block.tconstruct.smeltery_controller": "精錬炉コントローラー",
    "block.tconstruct.foundry_controller": "鋳造所コントローラー",
    "block.tconstruct.seared_heater": "焼成ヒーター",
    "block.tconstruct.scorched_alloyer": "焦熱合金炉",
    "block.tconstruct.seared_melter": "焼成溶解炉",
    "block.tconstruct.seared_drain": "焼成排出口",
    "block.tconstruct.scorched_drain": "焦熱排出口",
    "block.tconstruct.seared_duct": "焼成ダクト",
    "block.tconstruct.scorched_duct": "焦熱ダクト",
    "block.tconstruct.seared_chute": "焼成投入口",
    "block.tconstruct.scorched_chute": "焦熱投入口",
    "block.tconstruct.seared_tank": "焼成タンク",
    "block.tconstruct.scorched_tank": "焦熱タンク",
    "block.tconstruct.seared_fuel_tank": "焼成燃料タンク",
    "block.tconstruct.scorched_fuel_tank": "焦熱燃料タンク",
    "block.tconstruct.seared_lantern": "焼成ランタン",
    "block.tconstruct.scorched_lantern": "焦熱ランタン",
    "block.tconstruct.seared_faucet": "焼成蛇口",
    "block.tconstruct.scorched_faucet": "焦熱蛇口",
    "block.tconstruct.seared_channel": "焼成水路",
    "block.tconstruct.scorched_channel": "焦熱水路",
    "block.tconstruct.seared_table": "焼成作業台",
    "block.tconstruct.scorched_basin": "焦熱鋳造盤",
    "block.tconstruct.casting_table": "鋳造台",
    "block.tconstruct.casting_basin": "鋳造盤",
    
    # === 素材ブロック ===
    "block.tconstruct.seared_stone": "焼成石",
    "block.tconstruct.seared_cobble": "焼成丸石",
    "block.tconstruct.seared_bricks": "焼成レンガ",
    "block.tconstruct.seared_paver": "焼成舗装ブロック",
    "block.tconstruct.seared_glass": "焼成ガラス",
    "block.tconstruct.scorched_stone": "焦熱石",
    "block.tconstruct.scorched_bricks": "焦熱レンガ",
    "block.tconstruct.scorched_road": "焦熱道路",
    
    # === クラフト系ブロック ===
    "block.tconstruct.crafting_station": "作業場",
    "block.tconstruct.tinker_station": "改造ステーション",
    "block.tconstruct.part_builder": "パーツ作成台",
    "block.tconstruct.part_chest": "パーツ収納箱",
    "block.tconstruct.tinkers_chest": "ティンカー収納箱",
    "block.tconstruct.modifier_worktable": "モディファイア作業台",
    "block.tconstruct.scorched_anvil": "ティンカーの金床",
    
    # === アイテム（本） ===
    "item.tconstruct.materials_and_you": "素材と君",
    "item.tconstruct.puny_smelting": "小さな精錬",
    "item.tconstruct.mighty_smelting": "偉大なる精錬",
    "item.tconstruct.tinkers_gadgetry": "ティンカーの道具術",
    "item.tconstruct.fantastic_foundry": "素晴らしき鋳造所",
    "item.tconstruct.encyclopedia": "ティンカリング大百科",
    
    # === アイテム（型紙・鋳型） ===
    "item.tconstruct.pattern": "空白の型紙",
    "item.tconstruct.ingot_cast": "インゴット鋳型",
    "item.tconstruct.nugget_cast": "塊鋳型",
    "item.tconstruct.gem_cast": "宝石鋳型",
    "item.tconstruct.rod_cast": "棒鋳型",
    "item.tconstruct.repair_kit_cast": "修理キット鋳型",
    "item.tconstruct.pick_head_cast": "ピッケルヘッド鋳型",
    "item.tconstruct.small_axe_head_cast": "小さな斧ヘッド鋳型",
    "item.tconstruct.small_blade_cast": "小さな刀身鋳型",
    
    # === アイテム（ツール素材） ===
    "item.tconstruct.seared_brick": "焼成レンガ",
    "item.tconstruct.scorched_brick": "焦熱レンガ",
    "item.tconstruct.grout": "砂利土",
    "item.tconstruct.netherite_nugget": "ネザライト塊",
    "item.tconstruct.debris_nugget": "古代の残骸の塊",
    
    # === 素材名 ===
    "material.tconstruct.wood": "木",
    "material.tconstruct.stone": "石",
    "material.tconstruct.flint": "火打石",
    "material.tconstruct.iron": "鉄",
    "material.tconstruct.copper": "銅",
    "material.tconstruct.gold": "金",
    "material.tconstruct.diamond": "ダイヤモンド",
    "material.tconstruct.netherite": "ネザライト",
    "material.tconstruct.cobalt": "コバルト",
    "material.tconstruct.manyullyn": "マニュリン",
    "material.tconstruct.hepatizon": "ヘパティゾン",
    "material.tconstruct.queens_slime": "女王スライム",
    "material.tconstruct.pig_iron": "ピグアイアン",
    "material.tconstruct.rose_gold": "ローズゴールド",
    "material.tconstruct.slimesteel": "スライム鋼",
    "material.tconstruct.tinkers_bronze": "ティンカーブロンズ",
    "material.tconstruct.nahuatl": "ナワトル",
    
    # === 主要モディファイア ===
    "modifier.tconstruct.diamond": "ダイヤモンド",
    "modifier.tconstruct.emerald": "エメラルド",
    "modifier.tconstruct.netherite": "ネザライト",
    "modifier.tconstruct.reinforced": "補強",
    "modifier.tconstruct.unbreakable": "不壊",
    "modifier.tconstruct.haste": "採掘速度上昇",
    "modifier.tconstruct.sharpness": "鋭さ",
    "modifier.tconstruct.knockback": "ノックバック",
    "modifier.tconstruct.experienced": "経験値獲得",
    "modifier.tconstruct.magnetic": "磁力",
    "modifier.tconstruct.silky": "シルクタッチ",
    "modifier.tconstruct.fortune": "幸運",
    "modifier.tconstruct.luck": "幸運",
    "modifier.tconstruct.autosmelt": "自動精錬",
    "modifier.tconstruct.expanded": "範囲拡大",
    "modifier.tconstruct.reach": "リーチ",
    "modifier.tconstruct.smite": "アンデッド特効",
    "modifier.tconstruct.bane_of_arthropods": "虫特効",
    "modifier.tconstruct.antiaquatic": "水棲特効",
    "modifier.tconstruct.sharpness": "鋭さ",
    "modifier.tconstruct.sweeping_edge": "スイープ攻撃",
    "modifier.tconstruct.fire_aspect": "火属性",
    "modifier.tconstruct.looting": "ドロップ増加",
    "modifier.tconstruct.protection": "ダメージ軽減",
    "modifier.tconstruct.blast_protection": "爆発耐性",
    "modifier.tconstruct.projectile_protection": "飛び道具耐性",
    "modifier.tconstruct.fire_protection": "火炎耐性",
    "modifier.tconstruct.feather_falling": "落下耐性",
    "modifier.tconstruct.respiration": "水中呼吸",
    "modifier.tconstruct.aqua_affinity": "水中採掘",
    "modifier.tconstruct.thorns": "棘の鎧",
}

def load_json():
    """元のJSONファイルを読み込み"""
    input_file = '/home/iuif/dev/minecraft-mod-dictionary/workspace/imports/tconstruct/en_us.json'
    with open(input_file, 'r', encoding='utf-8') as f:
        return json.load(f)

def create_translation_entry(key, source):
    """翻訳エントリを作成"""
    # 高品質翻訳があればそれを使用
    if key in QUALITY_TRANSLATIONS:
        return {
            'key': key,
            'source': source,
            'target': QUALITY_TRANSLATIONS[key],
            'status': 'verified'
        }
    else:
        return {
            'key': key,
            'source': source,
            'target': f"[要翻訳] {source}",
            'status': 'draft'
        }

def main():
    data = load_json()
    output_dir = Path('/home/iuif/dev/minecraft-mod-dictionary/workspace/imports/tconstruct')
    
    # 重要なカテゴリを抽出
    important_categories = {
        'smeltery_foundry': [],
        'tools_weapons': [],
        'materials': [],
        'modifiers_top': [],
    }
    
    # カテゴリ分類
    for key, value in sorted(data.items()):
        entry = create_translation_entry(key, value)
        
        # Smeltery/Foundry関連
        if any(x in key for x in ['smeltery', 'foundry', 'seared', 'scorched', 'casting', 'melter', 'alloyer']):
            important_categories['smeltery_foundry'].append(entry)
        
        # ツール/武器
        elif key.startswith('item.tconstruct.') and any(x in key for x in ['pickaxe', 'hammer', 'axe', 'sword', 'cleaver', 'dagger', 'scythe', 'bow', 'crossbow']):
            important_categories['tools_weapons'].append(entry)
        
        # 素材（上位100エントリのみ）
        elif key.startswith('material.tconstruct.') and len(important_categories['materials']) < 100:
            important_categories['materials'].append(entry)
        
        # モディファイア（上位100エントリのみ）
        elif key.startswith('modifier.tconstruct.') and len(important_categories['modifiers_top']) < 100:
            important_categories['modifiers_top'].append(entry)
    
    # 統計情報
    stats = {
        'total_keys': len(data),
        'verified_translations': len(QUALITY_TRANSLATIONS),
        'categories': {k: len(v) for k, v in important_categories.items()}
    }
    
    # 最終出力YAML
    output = {
        'mod_id': 'tconstruct',
        'mod_name': "Tinkers' Construct",
        'version': '3.8.x',
        'mc_version': '1.20.1',
        'translation_date': '2025-11-29',
        'stats': stats,
        'important_translations': important_categories
    }
    
    output_file = output_dir / 'tconstruct_translation_final.yaml'
    with open(output_file, 'w', encoding='utf-8') as f:
        yaml.dump(output, f, allow_unicode=True, sort_keys=False, width=120, default_flow_style=False)
    
    print(f"✓ 最終翻訳ファイル作成完了: {output_file}")
    print(f"  - 総キー数: {stats['total_keys']}")
    print(f"  - 検証済み翻訳: {stats['verified_translations']}")
    print(f"  - Smeltery/Foundry: {stats['categories']['smeltery_foundry']}エントリ")
    print(f"  - ツール/武器: {stats['categories']['tools_weapons']}エントリ")
    print(f"  - 素材: {stats['categories']['materials']}エントリ")
    print(f"  - モディファイア: {stats['categories']['modifiers_top']}エントリ")

if __name__ == '__main__':
    main()
