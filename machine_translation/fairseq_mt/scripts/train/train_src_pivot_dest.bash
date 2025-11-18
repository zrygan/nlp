#!/bin/bash
set -e

# ============================================
# Pivot MT Training:
# SRC → PIVOT and PIVOT → DST
#
# Example:
#   rol → ilo → tgl
# ============================================

SRC=$1
PIV=$2
DST=$3

if [ -z "$SRC" ] || [ -z "$PIV" ] || [ -z "$DST" ]; then
    echo "Usage: ./pivot_train.sh <SRC> <PIVOT> <DST>"
    echo "Example: ./pivot_train.sh rol ilo tgl"
    exit 1
fi

# ==== Paths ====
DIR_SRC_PIV="data-bin/unigram/${SRC}-${PIV}"
DIR_PIV_DST="data-bin/unigram/${PIV}-${DST}"

CKPT_SRC_PIV="checkpoints/${SRC}-${PIV}"
CKPT_PIV_DST="checkpoints/${PIV}-${DST}"

ARCH="transformer_iwslt_de_en"
MAX_TOKENS=4096
LR=0.0005
DROPOUT=0.3
MAX_EPOCH=50
WARMUP=4000

train_pair() {
    local SRC_LANG=$1
    local TGT_LANG=$2
    local DATA=$3
    local SAVE=$4

    echo "==========================================="
    echo " Training: ${SRC_LANG} → ${TGT_LANG}"
    echo " Data: $DATA"
    echo " Save: $SAVE"
    echo "==========================================="

    mkdir -p "$SAVE"

    uv run --active fairseq-train "$DATA" \
      --arch $ARCH \
      --encoder-normalize-before --decoder-normalize-before \
      --share-all-embeddings \
      --optimizer adam --adam-betas '(0.9,0.98)' \
      --clip-norm 1.0 \
      --lr $LR \
      --lr-scheduler inverse_sqrt \
      --warmup-updates $WARMUP \
      --dropout $DROPOUT \
      --criterion label_smoothed_cross_entropy \
      --label-smoothing 0.1 \
      --max-tokens $MAX_TOKENS \
      --max-epoch $MAX_EPOCH \
      --seed 42 \
      --fp16 \
      --save-dir "$SAVE" \
      --keep-best-checkpoints 3 \
      --no-last-checkpoints \
      --eval-bleu \
      --eval-bleu-args '{"beam":5,"max_len_a":1.2,"max_len_b":10}' \
      --eval-bleu-detok space \
      --eval-bleu-remove-bpe \
      --best-checkpoint-metric bleu \
      --maximize-best-checkpoint-metric

    echo "✓ Done: ${SRC_LANG} → ${TGT_LANG}"
}


# ======= TRAIN STAGE 1: SRC → PIVOT =======
if [ ! -d "$DIR_SRC_PIV" ]; then
    echo "❌ Missing data directory: $DIR_SRC_PIV"
    exit 1
fi

train_pair "$SRC" "$PIV" "$DIR_SRC_PIV" "$CKPT_SRC_PIV"


# ======= TRAIN STAGE 2: PIVOT → DST =======
if [ ! -d "$DIR_PIV_DST" ]; then
    echo "❌ Missing data directory: $DIR_PIV_DST"
    exit 1
fi

train_pair "$PIV" "$DST" "$DIR_PIV_DST" "$CKPT_PIV_DST"


echo "==========================================="
echo "✓ Pivot MT Training Complete!"
echo " Models:"
echo "   SRC → PIV  : $CKPT_SRC_PIV/checkpoint_best.pt"
echo "   PIV → DST  : $CKPT_PIV_DST/checkpoint_best.pt"
echo "==========================================="
echo ""
echo "To translate:"
echo "  1. ./pivot_translate.sh input.txt tmp_pivot.txt $SRC $PIV"
echo "  2. ./pivot_translate.sh tmp_pivot.txt output.txt $PIV $DST"
