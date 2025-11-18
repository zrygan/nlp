#!/bin/bash
set -e

# ============================================
# Direct Training: SRC → DST
# Example: tgl → ceb
# ============================================

SRC=$1
DST=$2

DATA_DIR="../../data-bin/unigram/${SRC}-${DST}"
SAVE_DIR="../../checkpoints/${SRC}-${DST}"

# Training hyperparameters (optimized for low-resource)
ARCH=""
LR=0.0005
DROPOUT=0.3
MAX_TOKENS=4096
MAX_EPOCH=50
WARMUP=1000

if [ -z "$SRC" ] || [ -z "$DST" ]; then
    echo "Usage: ./direct_train.sh tgl ceb"
    exit 1
fi

if [ ! -d "$DATA_DIR" ]; then
    echo "❌ ERROR: Missing data directory: $DATA_DIR"
    echo "Run preprocessing first."
    exit 1
fi

mkdir -p "$SAVE_DIR"

echo "============================================"
echo " Training Direct MT: ${SRC} → ${DST}"
echo " Data: $DATA_DIR"
echo " Save: $SAVE_DIR"
echo "============================================"

uv run --active fairseq-train "$DATA_DIR" \
  --arch $ARCH \
  --encoder-normalize-before --decoder-normalize-before \
  --share-all-embeddings \
  --optimizer adam --adam-betas '(0.9,0.98)' \
  --clip-norm 1.0 \
  --lr $LR \
  --lr-scheduler inverse_sqrt \
  --warmup-updates $WARMUP \
  --warmup-init-lr 1e-07 \
  --dropout $DROPOUT \
  --weight-decay 0.0001 \
  --criterion label_smoothed_cross_entropy \
  --label-smoothing 0.1 \
  --max-tokens $MAX_TOKENS \
  --max-epoch $MAX_EPOCH \
  --seed 42 \
  --fp16 \
  --save-dir "$SAVE_DIR" \
  --keep-best-checkpoints 3 \
  --no-epoch-checkpoints --no-last-checkpoints \
  --eval-bleu \
  --eval-bleu-args '{"beam": 5, "max_len_a": 1.2, "max_len_b": 10}' \
  --eval-bleu-detok space \
  --eval-bleu-remove-bpe \
  --best-checkpoint-metric bleu \
  --maximize-best-checkpoint-metric

echo "✓ Training complete!"
echo "Best model saved at: ${SAVE_DIR}/checkpoint_best.pt"
