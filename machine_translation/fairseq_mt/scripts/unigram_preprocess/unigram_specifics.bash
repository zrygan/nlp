cat <<EOF
Converting sentence data to Unigram...
EOF

LANG1="ENG"
LANG2=""

mkdir ../../data/unigram

cat << EOF
  1) Concatanating all training data to a single file...
EOF

cat ../../data/raw/train.* > ../../data/unigram/all_train.txt

cat << EOF
  2) Converting all_train.txt to unigram main vocabulary and model
EOF

uv run --active python ./unigram_convert_vocab.py

mkdir ../../data/spm

cat << EOF
  3) Encoding raw data to unigram format...
EOF

uv run --active python ./unigram_encode.py

uv run --active python ./unigram_spm_combine.py

cat << EOF
  :) Unigram conversion is complete.
EOF 

# uv run --active python ./unigram_fairseq_preprocess.py


## CEB --> TGL

# SRC_LANG="ceb"
# TGT_LANG="tgl"
# DATA_DIR="../../data-bin/unigram/${SRC_LANG}-${TGT_LANG}"
# DATA_SPM_DIR="../../data/spm/"

# uv run --active fairseq-preprocess \
# --source-lang ${SRC_LANG} \
# --target-lang ${TGT_LANG} \
# --trainpref ${DATA_SPM_DIR}train.${SRC_LANG}-${TGT_LANG}.spm \
# --validpref ${DATA_SPM_DIR}valid.${SRC_LANG}-${TGT_LANG}.spm \
# --testpref ${DATA_SPM_DIR}test.${SRC_LANG}-${TGT_LANG}.spm \
# --destdir ${DATA_DIR} \
# --joined-dictionary \
# --workers 8

# uv run --active


# CEB -> ENG
SRC_LANG="ceb"
TGT_LANG="eng"
DATA_DIR="../../data-bin/unigram/${SRC_LANG}-${TGT_LANG}"
DATA_SPM_DIR="../../data/spm/"

uv run --active fairseq-preprocess \
--source-lang ${SRC_LANG} \
--target-lang ${TGT_LANG} \
--trainpref ${DATA_SPM_DIR}train.${SRC_LANG}-${TGT_LANG}.spm \
--validpref ${DATA_SPM_DIR}valid.${SRC_LANG}-${TGT_LANG}.spm \
--testpref ${DATA_SPM_DIR}test.${SRC_LANG}-${TGT_LANG}.spm \
--destdir ${DATA_DIR} \
--joined-dictionary \
--workers 8

