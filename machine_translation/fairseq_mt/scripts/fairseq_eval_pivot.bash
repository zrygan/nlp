#!/bin/bash

# Configuration
SRC=ceb
PIVOT=eng
TGT=tgl

MODEL1="../../../evaluation/unigram-label-pivot-ceb-eng/model/checkpoint_best.pt"
MODEL2="../../../evaluation/unigram-label-pivot-eng-tgl/model/checkpoint_best.pt"

DATA1="../../../evaluation/unigram-label-pivot-ceb-eng/bin"
DATA2="../../../evaluation/unigram-label-pivot-eng-tgl/bin"

RESULTS_DIR="results/augmentation/{$SRC}-{$PIVOT}-{$TGT}"
REF="../../../evaluation/unigram-label-pivot-eng-tgl/bin/test.eng-tgl.tgl.bin"
# Create results directory
mkdir -p ${RESULTS_DIR}

cat << EOF
==========================================
Evaluating ${SRC_LANG} → ${TGT_LANG} Model
==========================================
Model: ${MODEL1}
Model: ${MODEL2}
Data: ${DATA1}
Data: ${DATA2}
EOF

# Check if model exists
# if [ ! -f "${MODEL1}" ] || [ ! -f "${MODEL2}" ] ; then
# cat << EOF
# :( ERROR: Model not found for MODEL1: ${MODEL1} or MODEL2: ${MODEL2}
# :( Please train the model first or check the path.
# EOF
#     exit 1
# fi
cd results
mkdir pivot
cd pivot
mkdir "$SRC-$PIVOT-$TGT"
# Step 1: SRC → PIVOT
# uv run --active fairseq-generate $DATA1 --gen-subset test --path $MODEL1 --beam 5 --remove-bpe > "${SRC}_2_${PIVOT}.out"
grep '^H-' "${SRC}_2_${PIVOT}.out" | cut -f3 > "pivot.txt"

# Step 2: PIVOT → TGT
cat pivot.txt  | uv run --active fairseq-interactive $DATA2 --path $MODEL2 --beam 5 --remove-bpe | cat "${PIVOT}_2_${TGT}.out"
grep '^H-' "${PIVOT}_2_${TGT}.out" | cut -f3 > $PIVOT.final.txt

uv run --active fairseq-generate $DATA2 \
    --path /dev/null \
    --gen-subset test \
    --beam 1 --remove-bpe \
    > tmp.out

grep '^T-' tmp.out | cut -f2 > test.tgl.recovered


# Step 3: Evaluate BLEU
echo "Pivot BLEU:"
uv run --active sacrebleu test.tgl.recovered < $PIVOT.final.txt