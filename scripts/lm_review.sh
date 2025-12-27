#!/bin/bash
# LLM翻訳品質レビュースクリプト
# Usage: ./scripts/lm_review.sh [options]
#   -n NUM        レビュー件数 (default: 100)
#   -b NUM        バッチサイズ (default: 10)
#   -m MOD        特定のModのみレビュー
#   --dry-run     DB更新なし（レポートのみ）
#   --auto-fix    問題検出時に即座にLLMで再翻訳
#   --report FILE レポート出力先
#   --format json|text  レポート形式 (default: text)

set -uo pipefail

# 設定
API_URL="${LM_API_URL:-http://192.168.11.8:1234}"
MODEL="${LM_MODEL:-openai/gpt-oss-20b}"
DB_PATH="/home/iuif/dev/minecraft-mod-dictionary/moddict.db"
TERMS_PATH="/home/iuif/dev/minecraft-mod-dictionary/data/terms/global.yaml"
BATCH_SIZE=30
TOTAL_COUNT=100
MOD_FILTER=""
DRY_RUN=false
AUTO_FIX=false
REPORT_FILE=""
REPORT_FORMAT="text"
TEMPERATURE=0.1

# 一時ファイル
TMP_DIR=$(mktemp -d)
trap "rm -rf ${TMP_DIR}" EXIT

# 問題追跡用
declare -a ISSUES=()
declare -A ISSUE_COUNTS=(
    [mistranslation]=0
    [unnatural]=0
    [terminology]=0
    [format_lost]=0
    [partial]=0
)

# 引数解析
while [[ $# -gt 0 ]]; do
    case $1 in
        -n) TOTAL_COUNT="$2"; shift 2 ;;
        -b) BATCH_SIZE="$2"; shift 2 ;;
        -m) MOD_FILTER="$2"; shift 2 ;;
        --dry-run) DRY_RUN=true; shift ;;
        --auto-fix) AUTO_FIX=true; shift ;;
        --report) REPORT_FILE="$2"; shift 2 ;;
        --format) REPORT_FORMAT="$2"; shift 2 ;;
        *) echo "Unknown option: $1"; exit 1 ;;
    esac
done

# 進捗表示
log() {
    echo "[$(date '+%H:%M:%S')] $*"
}

# 用語辞書読み込み（YAMLからシンプルなリストに変換）
load_terms() {
    if [[ -f "${TERMS_PATH}" ]]; then
        # YAMLからsource → target形式でリスト化
        grep -E "^\s*-?\s*source:|^\s*target:" "${TERMS_PATH}" | \
        paste - - | \
        sed 's/.*source:\s*"\?\([^"]*\)"\?.*/\1 → /' | \
        sed 's/.*target:\s*"\?\([^"]*\)"\?$/\1/' | \
        paste -d '' - - | \
        head -50  # 最大50件（プロンプトサイズ制限）
    fi
}

# 用語辞書を取得
TERMS_LIST=$(load_terms)

# システムプロンプト
SYSTEM_PROMPT="あなたはMinecraft Modの翻訳品質をチェックする専門家です。

## 重要な前提
- Minecraft Modのテキストは主に「アイテム名」「ブロック名」「説明文」「UI文字列」です
- 英単語は原則として**名詞**として解釈してください（例: Perch = パーチ（魚）、止まるではない）
- Minecraft公式の日本語訳に準拠しているかを重視してください

## 入力形式
原文（英語）と翻訳（日本語）のペアをJSON形式で受け取ります

## 出力形式
各キーについて品質問題があればJSONで報告、問題なければ\"ok\"のみ出力

## チェック項目（重大な問題のみ報告）
1. mistranslation: 意味が原文と明らかに異なる誤訳（軽微なニュアンス差は除く）
2. unnatural: 明らかな文法ミス、意味不明な日本語
3. terminology: Minecraft公式用語と明らかに異なる翻訳（下記リスト参照）
4. format_lost: フォーマットコード（§, %s, %d, %1\$s, \$(...)、\\n等）の欠落・変更
5. partial: 一般的な英単語が翻訳されずに残っている（固有名詞・Mod名は除く）

## Minecraft公式用語（必ず準拠）
- Blackstone → ブラックストーン（黒曜石ではない！黒曜石=Obsidian）
- Stone Brick(s) → 石レンガ（石ブロックではない）
- Nether Brick(s) → ネザーレンガ
- End Stone → エンドストーン
- Purpur → プルプァ（紫水晶ではない）
- Log → 原木（例: Oak Log = オークの原木）
- Stem → 幹（ネザーの木材、例: Crimson Stem = 真紅の幹）
- Fence Gate → フェンスゲート
- Deepslate → 深層岩
- Calcite → 方解石
- Tuff → 凝灰岩

## 問題としない（必ず\"ok\"を返す）
- カタカナ表記の揺れ（セーター/スウェーター、ハット/帽子、等）
- 「の」の有無（鉄ツルハシ vs 鉄のツルハシ、粗い土スラブ vs 粗い土のスラブ）
- 送り仮名の違い（取り付ける vs 取付ける）
- 漢字/ひらがなの違い（出来る vs できる）
- 固有名詞・Mod名・技術用語が英語のまま
- 意味は正しいが表現が異なる程度の差
- Slab/スラブ と ハーフブロック は同義（両方OK）
- 「苔の〇〇」と「苔むした〇〇」は同義（両方OK）
- 翻訳として正しければ細かい表現差は全てOK

## 入力例
{\"k0\": {\"s\": \"Iron Pickaxe\", \"t\": \"鉄のツルハシ\"},
 \"k1\": {\"s\": \"Energy: %d RF\", \"t\": \"エネルギー: RF\"},
 \"k2\": {\"s\": \"Stone Brick Fence\", \"t\": \"石ブロックフェンス\"}}

## 出力例（JSONのみ、説明不要）
{\"k0\": \"ok\", \"k1\": {\"issue\": \"format_lost\", \"detail\": \"%dが欠落\"}, \"k2\": {\"issue\": \"terminology\", \"detail\": \"Stone Brick→石レンガ（公式）\"}}

用語辞書（参考）:
${TERMS_LIST}"

# 翻訳用システムプロンプト（再翻訳用）
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
${TERMS_LIST}"

# 翻訳API呼び出し（再翻訳用）
call_translate_api() {
    local user_content="$1"

    jq -n \
        --arg model "${MODEL}" \
        --arg system "${TRANSLATE_PROMPT}" \
        --arg user "${user_content}" \
        --argjson temp 0.2 \
        '{
            model: $model,
            messages: [
                {role: "system", content: $system},
                {role: "user", content: $user}
            ],
            temperature: $temp,
            max_tokens: 2048,
            stream: false
        }' > "${TMP_DIR}/translate_request.json"

    curl -s --max-time 60 "${API_URL}/v1/chat/completions" \
        -H "Content-Type: application/json" \
        -d @"${TMP_DIR}/translate_request.json"
}

# 単一エントリを再翻訳する関数
retranslate_entry() {
    local tid="$1"
    local source_text="$2"

    # 入力JSON構築
    local escaped=$(printf '%s' "${source_text}" | jq -Rs '.')
    local input_obj="{\"k0\":${escaped}}"

    # 翻訳API呼び出し
    local response=$(call_translate_api "${input_obj}")
    local content=$(echo "${response}" | jq -r '.choices[0].message.content // empty')

    if [[ -z "$content" ]]; then
        log "    [FAIL] id=${tid}: 翻訳API応答なし"
        return 1
    fi

    # JSON抽出
    local clean_content=$(echo "${content}" | sed 's/```json//g; s/```//g')
    if [[ "${clean_content}" == *"{"* ]]; then
        clean_content=$(echo "${clean_content}" | grep -o '{.*}' | tail -1)
    fi

    if ! echo "${clean_content}" | jq . > /dev/null 2>&1; then
        log "    [FAIL] id=${tid}: 無効なJSON応答"
        return 1
    fi

    local translation=$(echo "${clean_content}" | jq -r '.k0 // empty')

    if [[ -n "${translation}" && "${translation}" != "null" ]]; then
        if [[ "${translation}" == "${source_text}" ]]; then
            log "    [SKIP] id=${tid}: 原文保持"
            return 0
        else
            local escaped_trans=$(printf '%s' "${translation}" | jq -Rs '.')
            sqlite3 "${DB_PATH}" "UPDATE translations SET target_text=${escaped_trans}, status='translated', translator='lm:${MODEL}', notes=NULL, updated_at=datetime('now') WHERE id=${tid};"
            log "    [FIX] id=${tid}: ${translation:0:50}"
            return 0
        fi
    else
        log "    [FAIL] id=${tid}: 翻訳取得失敗"
        return 1
    fi
}

# API接続テスト
log "API接続テスト: ${API_URL}"
if ! curl -s "${API_URL}/v1/models" > /dev/null; then
    echo "Error: LLMサーバーに接続できません: ${API_URL}"
    exit 1
fi
log "接続OK (model: ${MODEL})"

# 翻訳済みテキスト取得
get_translated_batch() {
    local offset=$1
    local limit=$2
    local where_mod=""
    if [[ -n "${MOD_FILTER}" ]]; then
        where_mod="AND ts.mod_id = '${MOD_FILTER}'"
    fi

    sqlite3 -separator $'\t' "${DB_PATH}" "
SELECT t.id, ts.mod_id, ts.key, ts.source_text, t.target_text
FROM translations t
JOIN translation_sources ts ON t.source_id = ts.id
WHERE t.status = 'translated'
  AND t.target_text IS NOT NULL
  AND length(ts.source_text) >= 3
  ${where_mod}
ORDER BY ts.mod_id, ts.key
LIMIT ${limit} OFFSET ${offset}
"
}

# LLM API呼び出し
call_review_api() {
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

    curl -s --max-time 120 "${API_URL}/v1/chat/completions" \
        -H "Content-Type: application/json" \
        -d @"${TMP_DIR}/request.json"
}

# False Positiveフィルタ（許容すべき表現差をスキップ）
is_false_positive() {
    local detail="$1"
    local target="$2"

    # 日本語は小文字変換されないので、元のdetailも使用
    local detail_lower="${detail,,}"

    # ハット/帽子の揺れ（日本語なので元のdetailで比較）
    if [[ "$detail" == *"帽子"* && "$detail" == *"ハット"* ]]; then
        return 0  # False Positive
    fi

    # スラブ/ハーフブロックの揺れ（日本語で比較）
    if [[ "$detail" == *"スラブ"* && "$detail" == *"ハーフブロック"* ]]; then
        return 0
    fi
    if [[ "$detail_lower" == *"slab"* && "$detail" == *"ハーフブロック"* ]]; then
        return 0
    fi

    # 苔の/苔むしたの揺れ
    if [[ "$detail" == *"苔"* ]] && [[ "$detail" == *"苔むし"* || "$detail" == *"苔の"* ]]; then
        return 0
    fi

    # カタカナ表記の揺れ（サクラ/桜、マップル/メープル、クォーツ/クオーツ）
    if [[ "$detail" == *"サクラ"* && "$detail" == *"桜"* ]]; then
        return 0
    fi
    if [[ "$detail" == *"マップル"* && "$detail" == *"メープル"* ]]; then
        return 0
    fi
    if [[ "$detail" == *"クォーツ"* && "$detail" == *"クオーツ"* ]]; then
        return 0
    fi

    # 「の」の有無だけの問題
    if [[ "$detail" == *"のが欠落"* ]] || [[ "$detail" == *"「の」"* ]]; then
        return 0
    fi
    # 「の」の有無（桑のクッキー vs 桑クッキー等）
    if [[ "$detail" == *"桑クッキー"* && "$target" == *"桑のクッキー"* ]]; then
        return 0
    fi

    # Spruce/トウヒ - 両方正しい公式訳
    if [[ "$detail_lower" == *"spruce"* && "$detail" == *"トウヒ"* ]]; then
        return 0
    fi
    if [[ "$detail_lower" == *"spruce"* && "$detail" == *"スプルース"* ]]; then
        return 0
    fi

    # Quartz/クォーツ/石英 - 全て許容
    if [[ "$detail_lower" == *"quartz"* ]] && [[ "$detail" == *"クォーツ"* || "$detail" == *"石英"* || "$detail" == *"クオーツ"* ]]; then
        return 0
    fi

    # Birch/シラカバ - 公式訳
    if [[ "$detail_lower" == *"birch"* && "$detail" == *"シラカバ"* ]]; then
        # ログ/原木の問題でなければOK
        if [[ "$detail" != *"ログ"* && "$detail" != *"原木"* ]]; then
            return 0
        fi
    fi

    # Oak/オーク - 公式訳
    if [[ "$detail_lower" == *"oak"* && "$detail" == *"オーク"* ]]; then
        if [[ "$detail" != *"ログ"* && "$detail" != *"原木"* ]]; then
            return 0
        fi
    fi

    # Skull/スカル/頭骨 - 揺れを許容
    if [[ "$detail_lower" == *"skull"* ]] && [[ "$detail" == *"スカル"* || "$detail" == *"頭骨"* || "$detail" == *"頭部"* ]]; then
        return 0
    fi

    # 「レンガの」vs「レンガ」- 表現の揺れ
    if [[ "$detail" == *"レンガの"* && "$detail" == *"レンガフェンス"* ]]; then
        return 0
    fi

    # 翻訳が同じなのに問題と報告された場合
    if [[ "$detail" == *"$target"* && "$detail" == *"→"* ]]; then
        # 変換前後が同じかチェック
        local before after
        before=$(echo "$detail" | sed -n 's/.*→ *\([^（]*\).*/\1/p' | tr -d ' ')
        if [[ "$target" == *"$before"* ]]; then
            return 0
        fi
    fi

    return 1  # 正当な問題
}

# 問題を記録
record_issue() {
    local id="$1"
    local mod_id="$2"
    local key="$3"
    local source="$4"
    local target="$5"
    local issue_type="$6"
    local detail="$7"

    # False Positiveフィルタ
    if is_false_positive "$detail" "$target"; then
        return 0  # スキップ
    fi

    # JSON形式で追加
    local issue_json
    issue_json=$(jq -n \
        --arg id "$id" \
        --arg mod_id "$mod_id" \
        --arg key "$key" \
        --arg source "$source" \
        --arg target "$target" \
        --arg issue "$issue_type" \
        --arg detail "$detail" \
        '{id: $id, mod_id: $mod_id, key: $key, source: $source, target: $target, issue: $issue, detail: $detail}')

    ISSUES+=("$issue_json")
    ((ISSUE_COUNTS[$issue_type]++)) || true
}

# レポート生成（JSON形式）
generate_json_report() {
    local total_reviewed=$1
    local timestamp
    timestamp=$(date -Iseconds)

    # 問題配列をJSONに変換
    local issues_json="[]"
    if [[ ${#ISSUES[@]} -gt 0 ]]; then
        issues_json=$(printf '%s\n' "${ISSUES[@]}" | jq -s '.')
    fi

    jq -n \
        --arg timestamp "$timestamp" \
        --arg mod_id "${MOD_FILTER:-all}" \
        --argjson total_reviewed "$total_reviewed" \
        --argjson issues_found "${#ISSUES[@]}" \
        --argjson issues "$issues_json" \
        --argjson mistranslation "${ISSUE_COUNTS[mistranslation]}" \
        --argjson unnatural "${ISSUE_COUNTS[unnatural]}" \
        --argjson terminology "${ISSUE_COUNTS[terminology]}" \
        --argjson format_lost "${ISSUE_COUNTS[format_lost]}" \
        --argjson partial "${ISSUE_COUNTS[partial]}" \
        '{
            timestamp: $timestamp,
            mod_id: $mod_id,
            total_reviewed: $total_reviewed,
            issues_found: $issues_found,
            issues: $issues,
            summary: {
                mistranslation: $mistranslation,
                unnatural: $unnatural,
                terminology: $terminology,
                format_lost: $format_lost,
                partial: $partial
            }
        }'
}

# レポート生成（テキスト形式）
generate_text_report() {
    local total_reviewed=$1
    local timestamp
    timestamp=$(date '+%Y-%m-%d %H:%M:%S')

    echo "=== LLM翻訳品質レビュー結果 ==="
    echo "日時: ${timestamp}"
    echo "対象Mod: ${MOD_FILTER:-全Mod}"
    echo "レビュー件数: ${total_reviewed}件"
    echo "問題検出: ${#ISSUES[@]}件 ($(awk "BEGIN {printf \"%.1f\", ${#ISSUES[@]}*100/${total_reviewed}}")%)"
    echo ""

    if [[ ${#ISSUES[@]} -gt 0 ]]; then
        echo "--- 問題一覧 ---"
        local idx=1
        for issue_json in "${ISSUES[@]}"; do
            local id mod_id key source target issue detail
            id=$(echo "$issue_json" | jq -r '.id')
            mod_id=$(echo "$issue_json" | jq -r '.mod_id')
            key=$(echo "$issue_json" | jq -r '.key')
            source=$(echo "$issue_json" | jq -r '.source')
            target=$(echo "$issue_json" | jq -r '.target')
            issue=$(echo "$issue_json" | jq -r '.issue')
            detail=$(echo "$issue_json" | jq -r '.detail')

            echo "[${idx}] id=${id} mod=${mod_id} key=${key}"
            echo "    原文: ${source:0:60}"
            echo "    翻訳: ${target:0:60}"
            echo "    問題: ${issue} - ${detail}"
            echo ""
            ((idx++))
        done
    fi

    echo "--- サマリー ---"
    echo "誤訳(mistranslation):      ${ISSUE_COUNTS[mistranslation]}件"
    echo "不自然(unnatural):         ${ISSUE_COUNTS[unnatural]}件"
    echo "用語不統一(terminology):   ${ISSUE_COUNTS[terminology]}件"
    echo "フォーマット破損(format_lost): ${ISSUE_COUNTS[format_lost]}件"
    echo "部分未翻訳(partial):       ${ISSUE_COUNTS[partial]}件"

    if [[ "${DRY_RUN}" == "true" ]]; then
        echo ""
        echo "*** DRY RUN: DBは更新されていません ***"
    else
        echo ""
        echo "--- 次のアクション ---"
        echo "- 問題のある翻訳は status='needs_review' に変更されました"
        echo "- 修正後は ./moddict translate -mod [mod_id] -json で再インポートしてください"
    fi
}

# メイン処理
reviewed_count=0
error_count=0
offset=0

log "品質レビュー開始: 総数=${TOTAL_COUNT}, バッチ=${BATCH_SIZE}"
[[ -n "${MOD_FILTER}" ]] && log "Modフィルター: ${MOD_FILTER}"
[[ "${DRY_RUN}" == "true" ]] && log "*** DRY RUN モード ***"
[[ "${AUTO_FIX}" == "true" ]] && log "*** AUTO-FIX モード（問題検出時に即座に再翻訳） ***"

# 修正カウンター
fixed_count=0

while [[ ${reviewed_count} -lt ${TOTAL_COUNT} ]]; do
    # バッチデータ取得
    batch_tsv=$(get_translated_batch ${offset} ${BATCH_SIZE})

    if [[ -z "${batch_tsv}" ]]; then
        log "レビュー対象がなくなりました"
        break
    fi

    # TSVをJSONに変換
    declare -a batch_ids=()
    declare -a batch_mod_ids=()
    declare -a batch_keys=()
    declare -a batch_sources=()
    declare -a batch_targets=()
    input_obj="{"
    first=true
    idx=0

    while IFS=$'\t' read -r id mod_id key source_text target_text; do
        [[ -z "${id}" ]] && continue
        batch_ids+=("${id}")
        batch_mod_ids+=("${mod_id}")
        batch_keys+=("${key}")
        batch_sources+=("${source_text}")
        batch_targets+=("${target_text}")

        # JSONエスケープ
        escaped_source=$(printf '%s' "${source_text}" | jq -Rs '.')
        escaped_target=$(printf '%s' "${target_text}" | jq -Rs '.')

        if [[ "${first}" == "true" ]]; then
            first=false
        else
            input_obj+=","
        fi
        input_obj+="\"k${idx}\":{\"s\":${escaped_source},\"t\":${escaped_target}}"
        ((idx++))
    done <<< "${batch_tsv}"

    input_obj+="}"

    batch_count=${#batch_ids[@]}
    if [[ ${batch_count} -eq 0 ]]; then
        log "レビュー対象がなくなりました"
        break
    fi

    log "バッチ処理: ${batch_count}件 (offset=${offset})"

    # API呼び出し
    response=$(call_review_api "${input_obj}")

    # レスポンス解析
    content=$(echo "${response}" | jq -r '.choices[0].message.content // empty')

    if [[ -z "${content}" ]]; then
        log "Error: 空のレスポンス"
        echo "${response}" | jq . 2>/dev/null || echo "${response}"
        ((error_count++))
        ((offset += BATCH_SIZE))
        # 配列クリア
        batch_ids=()
        batch_mod_ids=()
        batch_keys=()
        batch_sources=()
        batch_targets=()
        continue
    fi

    # JSON抽出（マークダウンコードブロックを除去）
    clean_content=$(echo "${content}" | sed 's/```json//g; s/```//g')
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
        batch_mod_ids=()
        batch_keys=()
        batch_sources=()
        batch_targets=()
        continue
    fi

    # 結果解析
    batch_issues=0
    for i in $(seq 0 $((batch_count - 1))); do
        idx_key="k${i}"
        result=$(echo "${clean_content}" | jq -r --arg k "${idx_key}" '.[$k] // "ok"')

        if [[ "$result" != "ok" && "$result" != "null" ]]; then
            # 問題検出
            issue_type=$(echo "$result" | jq -r '.issue // "unknown"')
            detail=$(echo "$result" | jq -r '.detail // "詳細なし"')

            tid="${batch_ids[$i]}"
            mod_id="${batch_mod_ids[$i]}"
            orig_key="${batch_keys[$i]}"
            source_text="${batch_sources[$i]}"
            target_text="${batch_targets[$i]}"

            # 有効な問題タイプか確認
            if [[ "$issue_type" =~ ^(mistranslation|unnatural|terminology|format_lost|partial)$ ]]; then
                record_issue "$tid" "$mod_id" "$orig_key" "$source_text" "$target_text" "$issue_type" "$detail"

                log "  [ISSUE] ${orig_key}: ${issue_type} - ${detail:0:40}"

                # AUTO_FIXモード: 即座に再翻訳
                if [[ "${AUTO_FIX}" == "true" ]]; then
                    if retranslate_entry "$tid" "$source_text"; then
                        ((fixed_count++))
                    fi
                # 通常モード: needs_reviewに変更（dry-runでなければ）
                elif [[ "${DRY_RUN}" != "true" ]]; then
                    # SQLite用にシングルクォートをエスケープ
                    escaped_detail="${detail//\'/\'\'}"
                    sqlite3 "${DB_PATH}" "UPDATE translations SET status='needs_review', notes='${escaped_detail}', updated_at=datetime('now') WHERE id=${tid};"
                fi

                ((batch_issues++))
            fi
        fi

        ((reviewed_count++))
        [[ ${reviewed_count} -ge ${TOTAL_COUNT} ]] && break
    done

    log "  → ${batch_issues}件の問題検出"
    ((offset += BATCH_SIZE))

    # 配列クリア
    batch_ids=()
    batch_mod_ids=()
    batch_keys=()
    batch_sources=()
    batch_targets=()

    # 進捗表示
    if [[ $((reviewed_count % 50)) -eq 0 && ${reviewed_count} -gt 0 ]]; then
        log "========== 進捗: ${reviewed_count}/${TOTAL_COUNT} (${#ISSUES[@]}件問題検出, ${error_count}エラー) =========="
    fi

done

log "=========================================="
if [[ "${AUTO_FIX}" == "true" ]]; then
    log "完了: ${reviewed_count}件レビュー, ${#ISSUES[@]}件問題検出, ${fixed_count}件修正, ${error_count}件エラー"
else
    log "完了: ${reviewed_count}件レビュー, ${#ISSUES[@]}件問題検出, ${error_count}件エラー"
fi

# レポート出力
if [[ "${REPORT_FORMAT}" == "json" ]]; then
    report=$(generate_json_report ${reviewed_count})
else
    report=$(generate_text_report ${reviewed_count})
fi

if [[ -n "${REPORT_FILE}" ]]; then
    echo "${report}" > "${REPORT_FILE}"
    log "レポート出力: ${REPORT_FILE}"
else
    echo ""
    echo "${report}"
fi

exit 0
