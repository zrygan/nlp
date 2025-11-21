#!/bin/bash
set -e

# Configuration
SRC_LANG="ceb"
TGT_LANG="tgl"
TRAIN_REVERSE=false  # Set to true to also train reverse direction
MT_MODEL="mt_base"

# Training hyperparameters
MAX_TOKENS=4096
LR=0.0005
DROPOUT=0.3
MAX_EPOCH=50
WARMUP_UPDATES=000
CLIP_NORM=1.0

# Function to train a model
train_model() {
    local src=$1
    local tgt=$2
    local data_dir="data-bin/unigram/${src}-${tgt}"
    local checkpoint_dir="checkpoints/${src}-${tgt}"
    local mt_model="${MT_MODEL}"

    echo "========================================"
    echo "Training ${src} → ${tgt}"
    echo "========================================"
    
    # Check if data exists
    if [ ! -d "${data_dir}" ]; then
        echo "ERROR: Data directory not found: ${data_dir}"
        echo "Please run preprocessing first!"
        exit 1
    fi
    
    # Create checkpoint directory
    mkdir -p ${checkpoint_dir}
    
    # Train
    
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
    
    echo ""
    echo "✓ Training complete: ${src} → ${tgt}"
    echo "Best model saved in: ${checkpoint_dir}/checkpoint_best.pt"
    echo ""
}


# Train reverse direction if enabled (TGL → CEB)
if [ "$TRAIN_REVERSE" = true ]; then
cat << EOF
╔══════════════════════════════════════╗
║ Training in Reverse Direction ...    ║
╚══════════════════════════════════════╝
EOF
    train_model ${TGT_LANG} ${SRC_LANG}
else 
    # Train forward direction (CEB → TGL)
    train_model ${SRC_LANG} ${TGT_LANG}
fi

cat << EOF
╔══════════════════════════════════════╗
║ All Training Complete                ║
╚══════════════════════════════════════╝

To evaluate your model, execute: ./evalute.sh

To translate text for your model: ./translate.sh input.txt output.txt ${SRC_LANG}

EOF