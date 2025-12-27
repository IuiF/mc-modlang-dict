#!/bin/bash
# LMStudio バッチ翻訳スクリプト
# Usage: ./scripts/lm_translate.sh [options]
#   -n NUM    翻訳する件数 (default: 100)
#   -b NUM    バッチサイズ (default: 15)
#   -m MOD    特定のModのみ翻訳
#   -l NUM    最大文字数 (default: 40)
#   --dry-run 実際にDBを更新しない

set -uo pipefail

# 設定
API_URL="${LM_API_URL:-http://192.168.11.8:1234}"
MODEL="${LM_MODEL:-openai/gpt-oss-20b}"
DB_PATH="/home/iuif/dev/minecraft-mod-dictionary/moddict.db"
BATCH_SIZE=15
TOTAL_COUNT=100
MAX_LENGTH=100
MIN_LENGTH=3
MOD_FILTER=""
DRY_RUN=false
TEMPERATURE=0.2

# 一時ファイル
TMP_DIR=$(mktemp -d)
trap "rm -rf ${TMP_DIR}" EXIT

# 基本用語辞書
TERMS='## Minecraft公式用語
- Advancements: 進捗
- Achievement: 実績
- Armor: 防具
- Block: ブロック
- Chest: チェスト
- Crafting: クラフト
- Damage: ダメージ
- Dimension: ディメンション
- Drops/Drop: ドロップ
- Durability: 耐久値
- Effect: 効果
- Enchantment: エンチャント
- Entity: エンティティ
- Experience/XP: 経験値
- Fuel: 燃料
- Health: 体力
- Inventory: インベントリ
- Item: アイテム
- Level: レベル
- Ore: 鉱石
- Potion: ポーション
- Recipe: レシピ
- Redstone: レッドストーン
- Smelting: 製錬
- Stack: スタック
- Tool: ツール
- Villager: 村人

## 鉱石・金属
- Tin: スズ
- Lead: 鉛
- Iron: 鉄
- Gold: 金
- Copper: 銅
- Steel: 鋼鉄
- Bronze: 青銅
- Silver: 銀
- Zinc: 亜鉛
- Nickel: ニッケル
- Aluminum/Aluminium: アルミニウム
- Platinum: プラチナ
- Titanium: チタン
- Uranium: ウラン
- Coal: 石炭
- Diamond: ダイヤモンド
- Emerald: エメラルド
- Ruby: ルビー
- Sapphire: サファイア
- Quartz: クォーツ
- Obsidian: 黒曜石

## 岩石・ブロック
- Basalt: 玄武岩
- Granite: 花崗岩
- Diorite: 閃緑岩
- Andesite: 安山岩
- Deepslate: 深層岩
- Tuff: 凝灰岩
- Calcite: 方解石
- Sandstone: 砂岩
- Cobblestone: 丸石
- Gravel: 砂利
- Clay: 粘土

## ツール
- Chisel: ノミ
- Hammer: ハンマー
- Wrench: レンチ
- Saw: ノコギリ
- Pickaxe: ツルハシ
- Axe: 斧
- Shovel: シャベル
- Hoe: クワ
- Sword: 剣

## 木材
- Oak: オーク
- Birch: シラカバ
- Spruce: トウヒ
- Fir: トウヒ
- Pine: マツ
- Jungle: ジャングル
- Acacia: アカシア
- Dark Oak: ダークオーク
- Mangrove: マングローブ
- Cherry: サクラ
- Log: 原木
- Plank: 板材
- Sap: 樹液

## 動物・Mob
- Owl: フクロウ
- Bee: ミツバチ
- Wolf: オオカミ
- Cat: ネコ
- Pig: ブタ
- Cow: ウシ
- Sheep: ヒツジ
- Chicken: ニワトリ
- Horse: ウマ
- Creeper: クリーパー
- Zombie: ゾンビ
- Skeleton: スケルトン
- Spider: クモ
- Enderman: エンダーマン

## その他
- Fluid: 液体
- Energy: エネルギー
- Mana: マナ
- Modifier: 修飾子
- Slot: スロット
- Upgrade: アップグレード
- Weapon: 武器
- Input: 入力
- Output: 出力
- Tank: タンク
- Storage: ストレージ'

# 引数解析
while [[ $# -gt 0 ]]; do
    case $1 in
        -n) TOTAL_COUNT="$2"; shift 2 ;;
        -b) BATCH_SIZE="$2"; shift 2 ;;
        -m) MOD_FILTER="$2"; shift 2 ;;
        -l) MAX_LENGTH="$2"; shift 2 ;;
        --dry-run) DRY_RUN=true; shift ;;
        *) echo "Unknown option: $1"; exit 1 ;;
    esac
done

# システムプロンプト
SYSTEM_PROMPT="Minecraft Modの英語テキストを日本語に翻訳してJSON形式で出力してください。

入力例:
{\"k0\": \"Iron Pickaxe\", \"k1\": \"Crafting Table\", \"k2\": \"Energy: %d RF\"}

出力例:
{\"k0\": \"鉄のツルハシ\", \"k1\": \"作業台\", \"k2\": \"エネルギー: %d RF\"}

ルール:
- JSONのみ出力（説明不要）
- フォーマットコード(%s, %d, %1\$s, §, \\n等)は保持
- 数字のみ・記号のみの場合はそのまま返す
- 翻訳不要なら原文をそのまま返す

用語:
${TERMS}"

# 進捗表示
log() {
    echo "[$(date '+%H:%M:%S')] $*"
}

# API接続テスト
log "API接続テスト: ${API_URL}"
if ! curl -s "${API_URL}/v1/models" > /dev/null; then
    echo "Error: LMStudioに接続できません: ${API_URL}"
    exit 1
fi
log "接続OK (model: ${MODEL})"

# 未翻訳テキスト取得（TSV形式）
get_pending_batch() {
    local offset=$1
    local limit=$2
    local where_mod=""
    if [[ -n "${MOD_FILTER}" ]]; then
        where_mod="AND ts.mod_id = '${MOD_FILTER}'"
    fi

    sqlite3 -separator $'\t' "${DB_PATH}" "
SELECT t.id, ts.key, ts.source_text
FROM translations t
JOIN translation_sources ts ON t.source_id = ts.id
WHERE t.status = 'pending'
  AND length(ts.source_text) BETWEEN ${MIN_LENGTH} AND ${MAX_LENGTH}
  ${where_mod}
ORDER BY length(ts.source_text) ASC
LIMIT ${limit} OFFSET ${offset}
"
}

# LMStudio API呼び出し
call_lm_api() {
    local user_content="$1"

    jq -n \
        --arg model "${MODEL}" \
        --arg system "${SYSTEM_PROMPT}" \
        --arg user "${user_content}" \
        --argjson temp "${TEMPERATURE}" \
        '{
            model: $model,
            messages: [
                {role: "system", content: $system},
                {role: "user", content: $user}
            ],
            temperature: $temp,
            max_tokens: 2048,
            stream: false
        }' > "${TMP_DIR}/request.json"

    curl -s "${API_URL}/v1/chat/completions" \
        -H "Content-Type: application/json" \
        -d @"${TMP_DIR}/request.json"
}

# メイン処理
translated_count=0
error_count=0
offset=0

log "翻訳開始: 総数=${TOTAL_COUNT}, バッチ=${BATCH_SIZE}, 文字数=${MIN_LENGTH}-${MAX_LENGTH}"
[[ -n "${MOD_FILTER}" ]] && log "Modフィルター: ${MOD_FILTER}"
[[ "${DRY_RUN}" == "true" ]] && log "*** DRY RUN モード ***"

while [[ ${translated_count} -lt ${TOTAL_COUNT} ]]; do
    # バッチデータ取得
    batch_tsv=$(get_pending_batch ${offset} ${BATCH_SIZE})

    if [[ -z "${batch_tsv}" ]]; then
        log "翻訳対象がなくなりました"
        break
    fi

    # TSVをJSONに変換して保存
    declare -a batch_ids=()
    declare -a batch_keys=()
    declare -a batch_sources=()
    input_obj="{"
    first=true
    idx=0

    while IFS=$'\t' read -r id key source_text; do
        [[ -z "${id}" ]] && continue
        batch_ids+=("${id}")
        batch_keys+=("${key}")
        batch_sources+=("${source_text}")

        # JSONエスケープ
        escaped=$(printf '%s' "${source_text}" | jq -Rs '.')

        if [[ "${first}" == "true" ]]; then
            first=false
        else
            input_obj+=","
        fi
        input_obj+="\"k${idx}\":${escaped}"
        ((idx++))
    done <<< "${batch_tsv}"

    input_obj+="}"

    batch_count=${#batch_ids[@]}
    if [[ ${batch_count} -eq 0 ]]; then
        log "翻訳対象がなくなりました"
        break
    fi

    log "バッチ処理: ${batch_count}件 (offset=${offset})"

    # API呼び出し
    response=$(call_lm_api "${input_obj}")

    # レスポンス解析
    content=$(echo "${response}" | jq -r '.choices[0].message.content // empty')

    if [[ -z "${content}" ]]; then
        log "Error: 空のレスポンス"
        echo "${response}" | jq . 2>/dev/null || echo "${response}"
        ((error_count++))
        ((offset += BATCH_SIZE))
        batch_ids=()
        batch_keys=()
        batch_sources=()
        continue
    fi

    # JSON抽出（マークダウンコードブロックを除去し、JSONオブジェクトを抽出）
    # 1. マークダウンコードブロックを除去
    clean_content=$(echo "${content}" | sed 's/```json//g; s/```//g')

    # 2. 最初の { から最後の } までを抽出
    if [[ "${clean_content}" == *"{"* ]]; then
        clean_content=$(echo "${clean_content}" | grep -o '{.*}' | tail -1)
    fi

    # JSONとしてパース可能か確認
    if ! echo "${clean_content}" | jq . > /dev/null 2>&1; then
        log "Error: 無効なJSON応答"
        log "  Content: ${content:0:200}"
        ((error_count++))
        ((offset += BATCH_SIZE))
        batch_ids=()
        batch_keys=()
        batch_sources=()
        continue
    fi

    # 翻訳結果をDBに保存
    batch_translated=0
    for i in $(seq 0 $((batch_count - 1))); do
        idx_key="k${i}"
        translation=$(echo "${clean_content}" | jq -r --arg k "${idx_key}" '.[$k] // empty')
        source_text="${batch_sources[$i]:-}"

        if [[ -n "${translation}" && "${translation}" != "null" ]]; then
            tid="${batch_ids[$i]}"
            orig_key="${batch_keys[$i]}"

            # 品質チェック
            quality_status="translated"
            skip_reason=""
            use_source_as_target=false

            # 1. 原文と完全一致 → 原文をそのまま保存
            if [[ "${translation}" == "${source_text}" ]]; then
                use_source_as_target=true
            fi

            # 2. 小文字のみで日本語なし → 原文をそのまま保存
            if [[ "${#translation}" -ge 5 ]] && \
               [[ "${translation}" == "${translation,,}" ]] && \
               [[ "${translation}" =~ ^[a-z\ ]+$ ]]; then
                use_source_as_target=true
            fi

            # 3. 英字混在率が高い（50%以上）→ needs_review
            if [[ "${#translation}" -ge 5 ]]; then
                # 英字をカウント
                alpha_only=$(echo "${translation}" | tr -cd 'a-zA-Z')
                alpha_len=${#alpha_only}
                total_len=${#translation}
                if [[ ${total_len} -gt 0 ]]; then
                    ratio=$((alpha_len * 100 / total_len))
                    if [[ ${ratio} -gt 50 ]] && [[ "${translation}" =~ [ぁ-んァ-ン一-龯] ]]; then
                        quality_status="needs_review"
                    fi
                fi
            fi

            # 保存するテキストを決定
            final_translation="${translation}"
            if [[ "${use_source_as_target}" == "true" ]]; then
                final_translation="${source_text}"
            fi

            if [[ "${DRY_RUN}" == "true" ]]; then
                if [[ "${use_source_as_target}" == "true" ]]; then
                    log "  [DRY][SAME] ${orig_key}: ${final_translation:0:40}"
                elif [[ "${quality_status}" == "needs_review" ]]; then
                    log "  [DRY][REVIEW] ${orig_key}: ${final_translation:0:40}"
                else
                    log "  [DRY] ${orig_key}: ${final_translation:0:40}"
                fi
            else
                # DBを更新
                escaped_trans=$(printf '%s' "${final_translation}" | jq -Rs '.')
                sqlite3 "${DB_PATH}" "UPDATE translations SET target_text=${escaped_trans}, status='${quality_status}', translator='lm:${MODEL}', updated_at=datetime('now') WHERE id=${tid};"
            fi

            ((batch_translated++))
            ((translated_count++))
        else
            log "Warning: ${idx_key} の翻訳が取得できませんでした"
        fi

        [[ ${translated_count} -ge ${TOTAL_COUNT} ]] && break
    done

    log "  → ${batch_translated}件翻訳完了"
    ((offset += BATCH_SIZE))

    # 配列クリア
    batch_ids=()
    batch_keys=()
    batch_sources=()

    # 進捗表示
    if [[ $((translated_count % 50)) -eq 0 && ${translated_count} -gt 0 ]]; then
        log "========== 進捗: ${translated_count}/${TOTAL_COUNT} (${error_count}エラー) =========="
    fi

    # レート制限対策
    sleep 0.3
done

log "=========================================="
log "完了: ${translated_count}件翻訳, ${error_count}件エラー"
[[ "${DRY_RUN}" == "true" ]] && log "*** DRY RUN: DBは更新されていません ***"

exit 0
