#!/bin/bash
set -e

# ============================================
# Configuration
# ============================================

SRC_LANG="ceb"          # source language
PIVOT_LANG="tgl"        # pivot language (middle step)
TGT_LANG="eng"          # final output language

# Model checkpoints
MODEL_SRC_PIVOT="checkpoints/${SRC_LANG}-${PIVOT_LANG}/checkpoint_best.pt"
MODEL_PIVOT_TGT="checkpoints/${PIVOT_LANG}-${TGT_LANG}/checkpoint_best.pt"

# Data-bin directories for each pair
DATA_SRC_PIVOT="data-bin/unigram/${SRC_LANG}-${PIVOT_LANG}"
DATA_PIVOT_TGT="data-bin/unigram/${PIVOT_LANG}-${TGT_LANG}"

# Input and output
INPUT_FILE=$1
INTERMEDIATE_FILE="pivot_output.${PIVOT_LANG}"
FINAL_OUTPUT=$2

# ============================================
# Functions
# ============================================

translate() {
    local input=$1
    local data_dir=$2
    local checkpoint=$3
    local output=$4
    local source_lang=$5
    local target_lang=$6

    echo "Translating ${source_lang} → ${target_lang} ..."
    cat "${input}" | \
    uv run --active fairseq-interactive "${data_dir}" \
        --path "${checkpoint}" \
        --beam 5 \
        --source-lang "${source_lang}" \
        --target-lang "${target_lang}" \
        --buffer-size 64 --batch-size 32 \
        --remove-bpe \
    | grep "^H" | awk -F'\t' '{print $3}' > "${output}"

    echo "✓ Output saved to ${output}"
}

# ============================================
# Checks
# ============================================

if [ -z "$INPUT_FILE" ] || [ -z "$FINAL_OUTPUT" ]; then
    echo "Usage: ./pivot_translate.sh input.txt final_output.txt"
    exit 1
fi

if [ ! -f "${MODEL_SRC_PIVOT}" ]; then
    echo "ERROR: Missing model for ${SRC_LANG}→${PIVOT_LANG}"
    exit 1
fi

if [ ! -f "${MODEL_PIVOT_TGT}" ]; then
    echo "ERROR: Missing model for ${PIVOT_LANG}→${TGT_LANG}"
    exit 1
fi

# ============================================
# Pivot Translation Pipeline
# ============================================

echo "==============================================="
echo " Pivot MT: ${SRC_LANG} → ${PIVOT_LANG} → ${TGT_LANG}"
echo "==============================================="

# Step 1: SRC → PIVOT
translate \
    "${INPUT_FILE}" \
    "${DATA_SRC_PIVOT}" \
    "${MODEL_SRC_PIVOT}" \
    "${INTERMEDIATE_FILE}" \
    "${SRC_LANG}" \
    "${PIVOT_LANG}"

# Step 2: PIVOT → TGT
translate \
    "${INTERMEDIATE_FILE}" \
    "${DATA_PIVOT_TGT}" \
    "${MODEL_PIVOT_TGT}" \
    "${FINAL_OUTPUT}" \
    "${PIVOT_LANG}" \
    "${TGT_LANG}"

echo "==============================================="
echo "✓ Final Output Saved: ${FINAL_OUTPUT}"
echo "==============================================="
