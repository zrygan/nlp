#!/bin/bash
set -e

# ============================================
# Direct Machine Translation
# TGL → CEB
# ============================================

SRC_LANG="eng"
TGT_LANG="tgl"
REVERSE=false
DIR_REVERSE=false

TEST_DIR="../../tests/"
if [ "$DIR_REVERSE" = true ] ; then
    DATA_DIR="../../evaluation/unigram-label-pivot-eng-tgl/bin"
    CHECKPOINT="../../evaluation/unigram-label-pivot-eng-tgl/model/checkpoint_best.pt"
else
    DATA_DIR="../../evaluation/unigram-label-pivot-eng-tgl/bin"
    CHECKPOINT="../../evaluation/unigram-label-pivot-eng-tgl/model/checkpoint_best.pt"
fi
INPUT_FILE=$1
OUTPUT_FILE=$2

# ============================================
# Usage Check
# ============================================

if [ -z "$INPUT_FILE" ] || [ -z "$OUTPUT_FILE" ]; then
    echo "Usage: ./train-src-dest.sh input_tgl.txt output_ceb.txt"
    exit 1
fi

if [ ! -f "$CHECKPOINT" ]; then
    echo "ERROR: Missing checkpoint: $CHECKPOINT"
    exit 1
fi

if [ ! -d "$DATA_DIR" ]; then
    echo "ERROR: Missing data directory: $DATA_DIR"
    exit 1
fi

# ============================================
# Perform Translation
# ============================================

echo "==============================================="
echo " Translating ${SRC_LANG} → ${TGT_LANG}" AT ${TEST_DIR} DIRECTORY
echo " Model: ${CHECKPOINT}"
echo "==============================================="

if [ "$REVERSE" = true ] ; then
    echo "${TGT_LANG} --> ${SRC_LANG}"

    cat "${TEST_DIR}${INPUT_FILE}" | \
    uv run --active fairseq-interactive "${DATA_DIR}" \
        --path "${CHECKPOINT}" \
        --beam 5 \
        --buffer-size 64 \
        --batch-size 32 \
        --source-lang "${TGT_LANG}" \
        --target-lang "${SRC_LANG}" \
        --remove-bpe \
    | grep "^H" | awk -F'\t' '{print $3}' > "${TEST_DIR}${OUTPUT_FILE}"
else
    echo "${SRC_LANG} --> ${TGT_LANG}"

    cat "${TEST_DIR}${INPUT_FILE}" | \
    uv run --active fairseq-interactive "${DATA_DIR}" \
        --path "${CHECKPOINT}" \
        --beam 5 \
        --buffer-size 64 \
        --batch-size 32 \
        --source-lang "${SRC_LANG}" \
        --target-lang "${TGT_LANG}" \
        --remove-bpe \
    | grep "^H" | awk -F'\t' '{print $3}' > "${TEST_DIR}${OUTPUT_FILE}"
fi


echo "==============================================="
echo "✓ Translation complete!"
echo "✓ Output saved to: ${TEST_DIR}${OUTPUT_FILE}"
echo "==============================================="
