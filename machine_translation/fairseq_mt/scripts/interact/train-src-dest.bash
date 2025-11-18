#!/bin/bash
set -e

# ============================================
# Direct Machine Translation
# TGL → CEB
# ============================================

SRC_LANG="tgl"
TGT_LANG="ceb"

DATA_DIR="../../data-bin/unigram/${SRC_LANG}-${TGT_LANG}"
CHECKPOINT="../checkpoints/${SRC_LANG}-${TGT_LANG}/checkpoint_best.pt"

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
echo " Translating ${SRC_LANG} → ${TGT_LANG}"
echo " Model: ${CHECKPOINT}"
echo "==============================================="

cat "${INPUT_FILE}" | \
uv run --active fairseq-interactive "${DATA_DIR}" \
    --path "${CHECKPOINT}" \
    --beam 5 \
    --buffer-size 64 \
    --batch-size 32 \
    --source-lang "${SRC_LANG}" \
    --target-lang "${TGT_LANG}" \
    --remove-bpe \
| grep "^H" | awk -F'\t' '{print $3}' > "${OUTPUT_FILE}"

echo "==============================================="
echo "✓ Translation complete!"
echo "✓ Output saved to: ${OUTPUT_FILE}"
echo "==============================================="
