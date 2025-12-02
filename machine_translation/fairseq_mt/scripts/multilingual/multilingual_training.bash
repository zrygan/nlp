#!/usr/bin/env bash
set -e

#############################################
# 1. Define Languages
#############################################

LANGUAGES=("bik" "cbk" "ceb" "eng" "ilo" "jil" "krj" "msb" "pag" "pam" "prf" "rol" "tao" "tgl" "tiu" "tsg" "war")

echo "Training Multilingual MT Model for:"
printf '%s ' "${LANGUAGES[@]}"
echo ""

#############################################
# 2. Generate unique language pairs (LANG1 < LANG2)
#############################################

ALL_PAIRS=()
for ((i=0; i<${#LANGUAGES[@]}; i++)); do
  for ((j=i+1; j<${#LANGUAGES[@]}; j++)); do
    LANG1=${LANGUAGES[$i]}
    LANG2=${LANGUAGES[$j]}
    ALL_PAIRS+=("${LANG1}-${LANG2}")
  done
done

echo "Generating unique ${#ALL_PAIRS[@]} language pairs..."
echo "Example pairs: ${ALL_PAIRS[@]:0:10}"

JOINED=$(IFS=,; echo "${ALL_PAIRS[*]}")

#############################################
# 3. Paths
#############################################

DATA_BIN="../../data-bin/multilingual"
SAVE_DIR="../../checkpoints/multilingual"
mkdir -p "$SAVE_DIR"

#############################################
# 4. Multilingual Fairseq Training
#############################################

fairseq-train $DATA_BIN \
    --task multilingual_translation \
    --lang-pairs "$JOINED" \
    --share-all-embeddings \
    --arch transformer_iwslt_de_en \
    --share-decoder-input-output-embed \
    --optimizer adam --adam-betas "(0.9, 0.98)" \
    --lr 0.0007 \
    --lr-scheduler inverse_sqrt \
    --warmup-updates 4000 \
    --dropout 0.3 \
    --criterion label_smoothed_cross_entropy \
    --label-smoothing 0.1 \
    --max-tokens 4096 \
    --update-freq 8 \
    --max-epoch 40 \
    --fp16 \
    --save-dir "$SAVE_DIR"

echo "✔ Training done!"
