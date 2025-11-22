# SRC_LANG="ceb"
# TGT_LANG="tgl"
# DATA_DIR="data-bin/${SRC_LANG}-${TGT_LANG}"
# CHECKPOINT_DIR="checkpoints/${SRC_LANG}-${TGT_LANG}"

# uv run --active fairseq-preprocess \
# --source-lang ${SRC_LANG} \
# --target-lang ${TGT_LANG} \
# --trainpref data/raw/train.${SRC_LANG}-${TGT_LANG} \
# --destdir data-bin/${SRC_LANG}-${TGT_LANG} \
# --validpref data/raw/valid.${SRC_LANG}-${TGT_LANG} \
# --testpref data/raw/test.${SRC_LANG}-${TGT_LANG} \
# --joined-dictionary \
# --workers 8
cd parallel_corpus_preprocess
uv run --active python training_data_creation.py

cd ..
cd unigram_preprocess

./unigram_conversion.bash

cd ..
cd multilingual

uv run --active python multi.py