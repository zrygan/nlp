SRC_LANG="ceb"
TGT_LANG="tgl"
DATA_DIR="data-bin/${SRC_LANG}-${TGT_LANG}"
CHECKPOINT_DIR="checkpoints/${SRC_LANG}-${TGT_LANG}"

uv run --active fairseq-preprocess \
--source-lang ${SRC_LANG} \
--target-lang ${TGT_LANG} \
--trainpref data/train.${SRC_LANG}-${TGT_LANG} \
--destdir data-bin/${SRC_LANG}-${TGT_LANG} \
--validpref data/valid.${SRC_LANG}-${TGT_LANG} \
--testpref data/test.${SRC_LANG}-${TGT_LANG}