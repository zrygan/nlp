mkdir ../../data/unigram

cat ../../data/raw/train.* > ../../data/unigram/all_raw_languages_training_data.txt

cat << EOF
Converting parallel corpus training data to unigram files....
EOF

uv run --active python ./unigram_convert_vocab.py

mkdir ../../data/spm

uv run --active python ./unigram_encode.py

# uv run --active python ./unigram_fairseq_preprocess.py


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

SRC_LANG="ceb"
TGT_LANG="tgl"
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

uv run --active