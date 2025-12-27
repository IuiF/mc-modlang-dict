#!/usr/bin/env python3
"""
mmorpg (Mine and Slash) の日本語翻訳テーブル
キー: Minecraft ID -> 値: 日本語翻訳
"""

MMORPG_TRANSLATIONS = {
    # チャット関連
    "mmorpg.chat.error_type_and_reason": "%1$s %2$s",
    "mmorpg.chat.req_lvl": "必要レベル: %1$s",
    "mmorpg.chat.reset_points": "あなたの %1$s ポイントがリセットされました。",
    "mmorpg.chat.this_item_cant_be_used_more_than_x_times": "このアイテムはすでに最大使用回数に達しています: (%1$s)",
    "mmorpg.chat.tool_add_stat": "%1$s が %2$s にアップグレードされました!",
    "mmorpg.chat.tool_level_up": "%1$s がレベル %2$s に到達しました!",
    "mmorpg.chat.total_chests": "マップ内のチェスト: %1$s/%2$s",
    "mmorpg.chat.total_mobs": "マップ内のモブ: %1$s/%2$s",
    "mmorpg.chat.tp_to_dungeon_mapname": "[%1$s] ダンジョンにテレポートしました",

    # チェストコンテンツ
    "mmorpg.chest_content.aura": "オーラ",
    "mmorpg.chest_content.currency": "通貨",
    "mmorpg.chest_content.gear": "装備",
    "mmorpg.chest_content.gem": "宝石",
    "mmorpg.chest_content.harvest_blue_chest": "ハーベストから取得した装備",
    "mmorpg.chest_content.harvest_green_chest": "ハーベストから取得した通貨",
    "mmorpg.chest_content.harvest_purple_chest": "ハーベストから取得した宝石",
    "mmorpg.chest_content.map": "マップ",
    "mmorpg.chest_content.rune": "ルーン",
    "mmorpg.chest_content.support_gem": "サポート宝石",

    # チェストタイプ
    "mmorpg.chest_type.aura": "オーラ",
    "mmorpg.chest_type.currency": "通貨",
    "mmorpg.chest_type.gear": "装備",
    "mmorpg.chest_type.gem": "宝石",
    "mmorpg.chest_type.harvest_blue_chest": "ハーベスト青チェスト",
    "mmorpg.chest_type.harvest_green_chest": "ハーベスト緑チェスト",
    "mmorpg.chest_type.harvest_purple_chest": "ハーベスト紫チェスト",
    "mmorpg.chest_type.map": "マップ",
    "mmorpg.chest_type.rune": "ルーン",
    "mmorpg.chest_type.support_gem": "サポート宝石",

    # コマンド関連
    "mmorpg.command.auto_pve_teammate": " (自動PVE)",
    "mmorpg.command.been_invited": "%1$s があなたをチームに招待しました。",
    "mmorpg.command.been_kicked": "チームからキックされました。",
    "mmorpg.command.cant_invite_yourself": "自分自身を招待することはできません!",
    "mmorpg.command.cant_make_leader": "すでにチームリーダーです。",
    "mmorpg.command.click_to_create": "ここをクリックしてチームを作成できます!",
    "mmorpg.command.click_to_join": "ここをクリックして参加することもできます!",
    "mmorpg.command.click_to_leave_and_create": "ここをクリックして離脱し、新しいチームを作成します。",
    "mmorpg.command.create_when_in_a_team": "すでにチームに属しています。それでも新しいチームを作成しますか? ",
    "mmorpg.command.invite_info": "%1$s をチームに招待しました。",
    "mmorpg.command.join_team": "%1$s のチームに参加しました。",
    "mmorpg.command.join_tip": "受け入れるには /mine_and_slash teams join %1$s を実行してください",
    "mmorpg.command.kick_info": "%1$s をチームからキックしました。",
    "mmorpg.command.leave_team": "チームを離脱しました。",
    "mmorpg.command.not_in_a_team_currently": "プレイヤーはチームに属していません",
    "mmorpg.command.not_in_same_team_due_to_distance": " (距離が遠すぎます)",
    "mmorpg.command.not_in_team": "チームに属していません。",
    "mmorpg.command.not_invited": "彼らのチームに招待されていません",
    "mmorpg.command.player_not_in_team": "彼らはあなたのチームに属していません",
    "mmorpg.command.somebody_joined_your_team": "%1$s があなたのチームに参加しました!",
    "mmorpg.command.team_created": "チームを作成しました。他のプレイヤーを招待できるようになります。",

    # RPG用語辞書：能力値・ステータス関連
    "mmorpg.stat.strength": "力",
    "mmorpg.stat.dexterity": "敏捷性",
    "mmorpg.stat.intelligence": "知力",
    "mmorpg.stat.vitality": "生命力",
    "mmorpg.stat.health": "ヘルスポイント",
    "mmorpg.stat.mana": "マナ",
    "mmorpg.stat.experience": "経験値",
    "mmorpg.stat.level": "レベル",

    # 装備・アイテム関連
    "mmorpg.item.weapon": "武器",
    "mmorpg.item.armor": "防具",
    "mmorpg.item.amulet": "アミュレット",
    "mmorpg.item.ring": "リング",
    "mmorpg.item.rare": "レア",
    "mmorpg.item.legendary": "伝説",
    "mmorpg.item.unique": "ユニーク",

    # ダンジョン関連
    "mmorpg.dungeon.difficulty": "難易度",
    "mmorpg.dungeon.easy": "イージー",
    "mmorpg.dungeon.normal": "ノーマル",
    "mmorpg.dungeon.hard": "ハード",
    "mmorpg.dungeon.boss": "ボス",

    # スキル・能力関連
    "mmorpg.skill.fireball": "ファイアボール",
    "mmorpg.skill.frostbolt": "フロストボルト",
    "mmorpg.skill.heal": "ヒール",
    "mmorpg.skill.buff": "バフ",
    "mmorpg.skill.debuff": "デバフ",

    # UI・メニュー関連
    "mmorpg.ui.inventory": "インベントリ",
    "mmorpg.ui.equipment": "装備中",
    "mmorpg.ui.quest": "クエスト",
    "mmorpg.ui.skills": "スキル",
    "mmorpg.ui.map": "マップ",
    "mmorpg.ui.character": "キャラクター",
    "mmorpg.ui.stats": "ステータス",
}

if __name__ == "__main__":
    import json
    print(json.dumps(MMORPG_TRANSLATIONS, ensure_ascii=False, indent=2))
