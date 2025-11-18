!/bin/bash
set -e

# Configuration
SRC_LANG="ceb"
TGT_LANG="tgl"
TRAIN_REVERSE=false  # Set to true to also train reverse direction


# Training hyperparameters
MAX_TOKENS=4096
LR=0.0005
DROPOUT=0.3
MAX_EPOCH=50
WARMUP_UPDATES=4000
CLIP_NORM=1.0

uv run --active fairseq-train ${data_dir} \
  --arch transformer_iwslt_de_en \  
  --encoder-normalize-before --decoder-normalize-before \
  --share-all-embeddings
  --optimizer adam \
  --adam-betas '(0.9, 0.98)' \
  --clip-norm ${CLIP_NORM} \
  --lr ${LR} \
  --lr-scheduler inverse_sqrt \
  --warmup-updates ${WARMUP_UPDATES} \
  --warmup-init-lr 1e-07 \
  --dropout ${DROPOUT} \
  --weight-decay 0.0001 \
  --criterion label_smoothed_cross_entropy \
  --label-smoothing 0.1 \
  --max-tokens ${MAX_TOKENS} \
  --max-epoch ${MAX_EPOCH} \
  --save-dir ${checkpoint_dir} \
  --keep-best-checkpoints 3 \
  --no-epoch-checkpoints \
  --no-last-checkpoints \
  --no-save-optimizer-state \
  --log-interval 50 \
  --seed 42 \
  --fp16 \
  --eval-bleu \
  --eval-bleu-args '{"beam": 5, "max_len_a": 1.2, "max_len_b": 10}' \
  --eval-bleu-detok space \
  --eval-bleu-remove-bpe \
  --eval-bleu-print-samples \
  --best-checkpoint-metric bleu \
  --maximize-best-checkpoint-metric