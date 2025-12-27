#!/bin/bash
# LM翻訳品質チェック・修正スクリプト
# Usage: ./scripts/lm_quality_check.sh [options]
#   --fix       問題を修正（pendingに戻す）
#   --auto-fix  問題を検出したら即座にLLMで再翻訳
#   --dry-run   修正せずにレポートのみ

set -uo pipefail

DB_PATH="/home/iuif/dev/minecraft-mod-dictionary/moddict.db"
DRY_RUN=true
FIX=false
AUTO_FIX=false

# LLM設定
API_URL="${LM_API_URL:-http://192.168.11.8:1234}"
MODEL="${LM_MODEL:-openai/gpt-oss-20b}"
TEMPERATURE=0.2
BATCH_SIZE=10

# 一時ファイル
TMP_DIR=$(mktemp -d)
trap "rm -rf ${TMP_DIR}" EXIT

# 基本用語辞書（lm_translate.shから）
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

## その他
- Fluid: 液体
- Energy: エネルギー
- Mana: マナ
- Slot: スロット
- Upgrade: アップグレード
- Weapon: 武器
- Input: 入力
- Output: 出力
- Tank: タンク
- Storage: ストレージ'

# 翻訳用システムプロンプト
TRANSLATE_PROMPT="Minecraft Modの英語テキストを日本語に翻訳してJSON形式で出力してください。

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

# 引数解析
while [[ $# -gt 0 ]]; do
    case $1 in
        --fix) FIX=true; DRY_RUN=false; shift ;;
        --auto-fix) AUTO_FIX=true; FIX=false; DRY_RUN=false; shift ;;
        --dry-run) DRY_RUN=true; shift ;;
        *) echo "Unknown option: $1"; exit 1 ;;
    esac
done

log() {
    echo "[$(date '+%H:%M:%S')] $*"
}

# LLM API呼び出し
call_lm_api() {
    local user_content="$1"

    jq -n \
        --arg model "${MODEL}" \
        --arg system "${TRANSLATE_PROMPT}" \
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

# 問題のあるIDを収集する配列
declare -a PROBLEM_IDS=()

# 問題IDを記録する関数
record_problem_id() {
    local id=$1
    PROBLEM_IDS+=("$id")
}

# バッチで再翻訳する関数
retranslate_batch() {
    local ids=("$@")
    local count=${#ids[@]}

    if [[ $count -eq 0 ]]; then
        return 0
    fi

    log ""
    log "=== 再翻訳処理: ${count}件 ==="

    # バッチに分割して処理
    local offset=0
    local translated=0
    local errors=0

    while [[ $offset -lt $count ]]; do
        # 今回のバッチのID群
        local batch_end=$((offset + BATCH_SIZE))
        [[ $batch_end -gt $count ]] && batch_end=$count
        local batch_ids=("${ids[@]:$offset:$((batch_end - offset))}")
        local batch_count=${#batch_ids[@]}

        # IDリストをカンマ区切りに
        local id_list=$(IFS=,; echo "${batch_ids[*]}")

        # 原文を取得
        local batch_data=$(sqlite3 -separator $'\t' "${DB_PATH}" "
            SELECT t.id, ts.source_text
            FROM translations t
            JOIN translation_sources ts ON t.source_id = ts.id
            WHERE t.id IN (${id_list})
        ")

        if [[ -z "$batch_data" ]]; then
            log "  Warning: データ取得失敗 (ids: ${id_list})"
            ((offset += BATCH_SIZE))
            continue
        fi

        # 入力JSON構築
        declare -a batch_tid=()
        declare -a batch_source=()
        local input_obj="{"
        local first=true
        local idx=0

        while IFS=$'\t' read -r tid source_text; do
            [[ -z "$tid" ]] && continue
            batch_tid+=("$tid")
            batch_source+=("$source_text")

            local escaped=$(printf '%s' "${source_text}" | jq -Rs '.')
            if [[ "$first" == "true" ]]; then
                first=false
            else
                input_obj+=","
            fi
            input_obj+="\"k${idx}\":${escaped}"
            ((idx++))
        done <<< "$batch_data"

        input_obj+="}"

        local real_batch_count=${#batch_tid[@]}
        if [[ $real_batch_count -eq 0 ]]; then
            ((offset += BATCH_SIZE))
            continue
        fi

        log "  バッチ処理: ${real_batch_count}件"

        # LLM API呼び出し
        local response=$(call_lm_api "${input_obj}")
        local content=$(echo "${response}" | jq -r '.choices[0].message.content // empty')

        if [[ -z "$content" ]]; then
            log "  Error: 空のレスポンス"
            ((errors++))
            ((offset += BATCH_SIZE))
            batch_tid=()
            batch_source=()
            continue
        fi

        # JSON抽出
        local clean_content=$(echo "${content}" | sed 's/```json//g; s/```//g')
        if [[ "${clean_content}" == *"{"* ]]; then
            clean_content=$(echo "${clean_content}" | grep -o '{.*}' | tail -1)
        fi

        if ! echo "${clean_content}" | jq . > /dev/null 2>&1; then
            log "  Error: 無効なJSON応答"
            ((errors++))
            ((offset += BATCH_SIZE))
            batch_tid=()
            batch_source=()
            continue
        fi

        # 結果をDBに保存
        local batch_translated=0
        for i in $(seq 0 $((real_batch_count - 1))); do
            local idx_key="k${i}"
            local translation=$(echo "${clean_content}" | jq -r --arg k "${idx_key}" '.[$k] // empty')
            local tid="${batch_tid[$i]}"
            local source_text="${batch_source[$i]}"

            if [[ -n "${translation}" && "${translation}" != "null" ]]; then
                # 翻訳結果が原文と同じ場合はスキップ（翻訳不要と判断）
                if [[ "${translation}" == "${source_text}" ]]; then
                    log "    [SKIP] id=${tid}: 原文保持"
                else
                    local escaped_trans=$(printf '%s' "${translation}" | jq -Rs '.')
                    sqlite3 "${DB_PATH}" "UPDATE translations SET target_text=${escaped_trans}, status='translated', translator='lm:${MODEL}', updated_at=datetime('now') WHERE id=${tid};"
                    log "    [OK] id=${tid}: ${translation:0:40}"
                    ((batch_translated++))
                fi
            else
                log "    [FAIL] id=${tid}: 翻訳取得失敗"
            fi
        done

        ((translated += batch_translated))
        ((offset += BATCH_SIZE))

        # 配列クリア
        batch_tid=()
        batch_source=()

        # レート制限対策
        sleep 0.3
    done

    log "  再翻訳完了: ${translated}件成功, ${errors}件エラー"
}

log "=== LM翻訳品質チェック ==="
[[ "${DRY_RUN}" == "true" ]] && log "*** DRY RUN モード（--fix で修正実行） ***"
[[ "${AUTO_FIX}" == "true" ]] && log "*** AUTO-FIX モード（問題検出後に再翻訳） ***"

# AUTO_FIXモード時はLLM接続テスト
if [[ "${AUTO_FIX}" == "true" ]]; then
    log "LLM接続テスト: ${API_URL}"
    if ! curl -s "${API_URL}/v1/models" > /dev/null; then
        echo "Error: LLMサーバーに接続できません: ${API_URL}"
        exit 1
    fi
    log "接続OK (model: ${MODEL})"
fi

# 1. 未翻訳（原文と完全一致）
log ""
log "--- 1. 未翻訳（原文と完全一致） ---"
count_untranslated=$(sqlite3 "${DB_PATH}" "
SELECT COUNT(*)
FROM translations t
JOIN translation_sources ts ON t.source_id = ts.id
WHERE t.translator LIKE 'lm:%'
  AND t.status = 'translated'
  AND ts.source_text = t.target_text
  AND length(ts.source_text) >= 3
  AND ts.source_text NOT LIKE '%\%%'
  AND ts.source_text NOT LIKE '%@%'
  AND ts.source_text NOT LIKE '%§%';
")
log "検出数: ${count_untranslated}件"

if [[ "${FIX}" == "true" && "${count_untranslated}" -gt 0 ]]; then
    sqlite3 "${DB_PATH}" "
    UPDATE translations SET status='pending', target_text=NULL, translator=NULL
    WHERE id IN (
        SELECT t.id
        FROM translations t
        JOIN translation_sources ts ON t.source_id = ts.id
        WHERE t.translator LIKE 'lm:%'
          AND t.status = 'translated'
          AND ts.source_text = t.target_text
          AND length(ts.source_text) >= 3
          AND ts.source_text NOT LIKE '%\%%'
          AND ts.source_text NOT LIKE '%@%'
          AND ts.source_text NOT LIKE '%§%'
    );"
    log "→ ${count_untranslated}件をpendingに戻しました"
fi

# AUTO_FIXモード: 問題IDを収集
if [[ "${AUTO_FIX}" == "true" && "${count_untranslated}" -gt 0 ]]; then
    while read -r id; do
        [[ -n "$id" ]] && record_problem_id "$id"
    done < <(sqlite3 "${DB_PATH}" "
        SELECT t.id
        FROM translations t
        JOIN translation_sources ts ON t.source_id = ts.id
        WHERE t.translator LIKE 'lm:%'
          AND t.status = 'translated'
          AND ts.source_text = t.target_text
          AND length(ts.source_text) >= 3
          AND ts.source_text NOT LIKE '%\%%'
          AND ts.source_text NOT LIKE '%@%'
          AND ts.source_text NOT LIKE '%§%'
    ")
    log "→ ${count_untranslated}件を再翻訳キューに追加"
fi

# 2. 中途半端な翻訳（日本語と英字が混在、かつ英字が多い）
# ※Mod名・固有名詞が含まれる場合は除外
log ""
log "--- 2. 中途半端な翻訳（英字混在率50%以上） ---"
count_partial=$(sqlite3 "${DB_PATH}" "
SELECT COUNT(*)
FROM translations t
JOIN translation_sources ts ON t.source_id = ts.id
WHERE t.translator LIKE 'lm:%'
  AND t.status = 'translated'
  AND length(t.target_text) >= 5
  AND t.target_text GLOB '*[a-zA-Z]*'
  AND t.target_text GLOB '*[ぁ-んァ-ン一-龯]*'
  AND t.target_text NOT LIKE '%Config%'
  AND t.target_text NOT LIKE '%Settings%'
  AND t.target_text NOT LIKE '%Microsoft%'
  AND t.target_text NOT LIKE '%Voidstone%'
  AND t.target_text NOT LIKE '%Rune %'
  AND t.target_text NOT LIKE '%Tool%'
  AND t.target_text NOT LIKE '%ProjectE%'
  AND t.target_text NOT LIKE '%Overhaul%'
  AND t.target_text NOT LIKE '%apothic%'
  AND t.target_text NOT LIKE '%Glassential%'
  AND t.target_text NOT LIKE '%Bluprintz%'
  AND t.target_text NOT LIKE '%StateMerger%'
  AND t.target_text NOT LIKE '%Hohlraum%'
  AND t.target_text NOT LIKE '%Morph-o-%'
  AND t.target_text NOT LIKE '%Aeternalis%'
  AND t.target_text NOT LIKE '%Energised%'
  AND t.target_text NOT LIKE '%Phlogistate%'
  AND (length(t.target_text) - length(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(t.target_text,'a',''),'b',''),'c',''),'d',''),'e',''),'f',''),'g',''),'h',''),'i',''),'j',''),'k',''),'l',''),'m',''),'n',''),'o',''),'p',''),'q',''),'r',''),'s',''),'t',''),'u',''),'v',''),'w',''),'x',''),'y',''),'z',''))) * 100 / length(t.target_text) > 50;
")
log "検出数: ${count_partial}件"

# サンプル表示
if [[ "${count_partial}" -gt 0 ]]; then
    log "サンプル:"
    sqlite3 -separator ' → ' "${DB_PATH}" "
    SELECT ts.source_text, t.target_text
    FROM translations t
    JOIN translation_sources ts ON t.source_id = ts.id
    WHERE t.translator LIKE 'lm:%'
      AND t.status = 'translated'
      AND length(t.target_text) >= 5
      AND t.target_text GLOB '*[a-zA-Z]*'
      AND t.target_text GLOB '*[ぁ-んァ-ン一-龯]*'
      AND (length(t.target_text) - length(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(t.target_text,'a',''),'b',''),'c',''),'d',''),'e',''),'f',''),'g',''),'h',''),'i',''),'j',''),'k',''),'l',''),'m',''),'n',''),'o',''),'p',''),'q',''),'r',''),'s',''),'t',''),'u',''),'v',''),'w',''),'x',''),'y',''),'z',''))) * 100 / length(t.target_text) > 50
    LIMIT 10;" | while read line; do log "  $line"; done
fi

if [[ "${FIX}" == "true" && "${count_partial}" -gt 0 ]]; then
    sqlite3 "${DB_PATH}" "
    UPDATE translations SET status='needs_review'
    WHERE id IN (
        SELECT t.id
        FROM translations t
        JOIN translation_sources ts ON t.source_id = ts.id
        WHERE t.translator LIKE 'lm:%'
          AND t.status = 'translated'
          AND length(t.target_text) >= 5
          AND t.target_text GLOB '*[a-zA-Z]*'
          AND t.target_text GLOB '*[ぁ-んァ-ン一-龯]*'
          AND (length(t.target_text) - length(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(t.target_text,'a',''),'b',''),'c',''),'d',''),'e',''),'f',''),'g',''),'h',''),'i',''),'j',''),'k',''),'l',''),'m',''),'n',''),'o',''),'p',''),'q',''),'r',''),'s',''),'t',''),'u',''),'v',''),'w',''),'x',''),'y',''),'z',''))) * 100 / length(t.target_text) > 50
    );"
    log "→ ${count_partial}件をneeds_reviewに変更しました"
fi

# AUTO_FIXモード: 問題IDを収集
if [[ "${AUTO_FIX}" == "true" && "${count_partial}" -gt 0 ]]; then
    while read -r id; do
        [[ -n "$id" ]] && record_problem_id "$id"
    done < <(sqlite3 "${DB_PATH}" "
        SELECT t.id
        FROM translations t
        JOIN translation_sources ts ON t.source_id = ts.id
        WHERE t.translator LIKE 'lm:%'
          AND t.status = 'translated'
          AND length(t.target_text) >= 5
          AND t.target_text GLOB '*[a-zA-Z]*'
          AND t.target_text GLOB '*[ぁ-んァ-ン一-龯]*'
          AND t.target_text NOT LIKE '%Config%'
          AND t.target_text NOT LIKE '%Settings%'
          AND (length(t.target_text) - length(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(replace(t.target_text,'a',''),'b',''),'c',''),'d',''),'e',''),'f',''),'g',''),'h',''),'i',''),'j',''),'k',''),'l',''),'m',''),'n',''),'o',''),'p',''),'q',''),'r',''),'s',''),'t',''),'u',''),'v',''),'w',''),'x',''),'y',''),'z',''))) * 100 / length(t.target_text) > 50
    ")
    log "→ ${count_partial}件を再翻訳キューに追加"
fi

# 3. 小文字のみの翻訳（coming soon のような未翻訳）
log ""
log "--- 3. 小文字のみ（未翻訳の可能性） ---"
count_lowercase=$(sqlite3 "${DB_PATH}" "
SELECT COUNT(*)
FROM translations t
JOIN translation_sources ts ON t.source_id = ts.id
WHERE t.translator LIKE 'lm:%'
  AND t.status = 'translated'
  AND length(t.target_text) >= 5
  AND t.target_text = lower(t.target_text)
  AND t.target_text GLOB '*[a-z]*'
  AND t.target_text NOT GLOB '*[ぁ-んァ-ン一-龯]*';
")
log "検出数: ${count_lowercase}件"

if [[ "${count_lowercase}" -gt 0 ]]; then
    log "サンプル:"
    sqlite3 -separator ' → ' "${DB_PATH}" "
    SELECT ts.source_text, t.target_text
    FROM translations t
    JOIN translation_sources ts ON t.source_id = ts.id
    WHERE t.translator LIKE 'lm:%'
      AND t.status = 'translated'
      AND length(t.target_text) >= 5
      AND t.target_text = lower(t.target_text)
      AND t.target_text GLOB '*[a-z]*'
      AND t.target_text NOT GLOB '*[ぁ-んァ-ン一-龯]*'
    LIMIT 10;" | while read line; do log "  $line"; done
fi

if [[ "${FIX}" == "true" && "${count_lowercase}" -gt 0 ]]; then
    sqlite3 "${DB_PATH}" "
    UPDATE translations SET status='pending', target_text=NULL, translator=NULL
    WHERE id IN (
        SELECT t.id
        FROM translations t
        JOIN translation_sources ts ON t.source_id = ts.id
        WHERE t.translator LIKE 'lm:%'
          AND t.status = 'translated'
          AND length(t.target_text) >= 5
          AND t.target_text = lower(t.target_text)
          AND t.target_text GLOB '*[a-z]*'
          AND t.target_text NOT GLOB '*[ぁ-んァ-ン一-龯]*'
    );"
    log "→ ${count_lowercase}件をpendingに戻しました"
fi

# AUTO_FIXモード: 問題IDを収集
if [[ "${AUTO_FIX}" == "true" && "${count_lowercase}" -gt 0 ]]; then
    while read -r id; do
        [[ -n "$id" ]] && record_problem_id "$id"
    done < <(sqlite3 "${DB_PATH}" "
        SELECT t.id
        FROM translations t
        JOIN translation_sources ts ON t.source_id = ts.id
        WHERE t.translator LIKE 'lm:%'
          AND t.status = 'translated'
          AND length(t.target_text) >= 5
          AND t.target_text = lower(t.target_text)
          AND t.target_text GLOB '*[a-z]*'
          AND t.target_text NOT GLOB '*[ぁ-んァ-ン一-龯]*'
    ")
    log "→ ${count_lowercase}件を再翻訳キューに追加"
fi

# 4. 中国語混入（日本語では使わない繁体字・簡体字のみ検出）
log ""
log "--- 4. 中国語混入の可能性 ---"
count_chinese=$(sqlite3 "${DB_PATH}" "
SELECT COUNT(*)
FROM translations t
WHERE t.translator LIKE 'lm:%'
  AND t.status = 'translated'
  AND (t.target_text LIKE '%掉%' OR t.target_text LIKE '%絲%' OR t.target_text LIKE '%們%' OR t.target_text LIKE '%從%');
")
log "検出数: ${count_chinese}件"

if [[ "${count_chinese}" -gt 0 ]]; then
    log "サンプル:"
    sqlite3 -separator ' → ' "${DB_PATH}" "
    SELECT ts.source_text, t.target_text
    FROM translations t
    JOIN translation_sources ts ON t.source_id = ts.id
    WHERE t.translator LIKE 'lm:%'
      AND t.status = 'translated'
      AND (t.target_text LIKE '%掉%' OR t.target_text LIKE '%絲%' OR t.target_text LIKE '%們%' OR t.target_text LIKE '%從%')
    LIMIT 10;" | while read line; do log "  $line"; done
fi

if [[ "${FIX}" == "true" && "${count_chinese}" -gt 0 ]]; then
    sqlite3 "${DB_PATH}" "
    UPDATE translations SET status='needs_review'
    WHERE id IN (
        SELECT t.id
        FROM translations t
        WHERE t.translator LIKE 'lm:%'
          AND t.status = 'translated'
          AND (t.target_text LIKE '%掉%' OR t.target_text LIKE '%絲%' OR t.target_text LIKE '%們%' OR t.target_text LIKE '%從%')
    );"
    log "→ ${count_chinese}件をneeds_reviewに変更しました"
fi

# AUTO_FIXモード: 問題IDを収集
if [[ "${AUTO_FIX}" == "true" && "${count_chinese}" -gt 0 ]]; then
    while read -r id; do
        [[ -n "$id" ]] && record_problem_id "$id"
    done < <(sqlite3 "${DB_PATH}" "
        SELECT t.id
        FROM translations t
        WHERE t.translator LIKE 'lm:%'
          AND t.status = 'translated'
          AND (t.target_text LIKE '%掉%' OR t.target_text LIKE '%絲%' OR t.target_text LIKE '%們%' OR t.target_text LIKE '%從%')
    ")
    log "→ ${count_chinese}件を再翻訳キューに追加"
fi

# サマリー
log ""
log "=== サマリー ==="
total_issues=$((count_untranslated + count_partial + count_lowercase + count_chinese))
log "問題検出総数: ${total_issues}件"
log "  - 未翻訳（原文一致）: ${count_untranslated}件"
log "  - 中途半端: ${count_partial}件"
log "  - 小文字のみ: ${count_lowercase}件"
log "  - 中国語混入: ${count_chinese}件"

if [[ "${DRY_RUN}" == "true" ]]; then
    log ""
    log "修正を実行するには:"
    log "  --fix      問題をpending/needs_reviewに戻す"
    log "  --auto-fix 問題を検出して即座にLLMで再翻訳"
fi

# AUTO_FIXモード: 収集した問題IDを再翻訳
if [[ "${AUTO_FIX}" == "true" ]]; then
    if [[ ${#PROBLEM_IDS[@]} -gt 0 ]]; then
        retranslate_batch "${PROBLEM_IDS[@]}"
    else
        log ""
        log "再翻訳対象はありませんでした"
    fi
fi

exit 0
