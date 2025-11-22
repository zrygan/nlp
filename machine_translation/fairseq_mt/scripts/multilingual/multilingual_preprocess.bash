#!/usr/bin/env bash
set -e

# Usage:
#   ./multilingual_prepare_raw.bash lang1 lang2 lang3 ...
#

LANGS=("$@")
RAW="../../data/spm" # regular test.src-dest.spm, valid.src-dest.spm, or train.src-dest.spm files

if [ ${#LANGS[@]} -lt 2 ]; then
    echo "Usage: $0 lang1 lang2 lang3 ..."
    exit 1
fi

FIRST_DICT=""

for SRC in "${LANGS[@]}"; do
  for TGT in "${LANGS[@]}"; do

    [ "$SRC" = "$TGT" ] && continue
    
    PAIR="$SRC-$TGT"

    # FIRST PAIR CREATES THE DICTIONARY
    if [ -z "$FIRST_DICT" ]; then
        fairseq-preprocess \
            --source-lang $SRC --target-lang $TGT \
            --trainpref $RAW/train.$PAIR.spm \
            --validpref $RAW/valid.$PAIR.spm \
            --testpref  $RAW/test.$PAIR.spm \
            --destdir $DESTDIR \
            --workers 8

        FIRST_DICT=$DESTDIR/dict.$TGT.txt

    else
        fairseq-preprocess \
            --source-lang $SRC --target-lang $TGT \
            --trainpref $RAW/train.$PAIR.spm \
            --validpref $RAW/valid.$PAIR.spm \
            --testpref  $RAW/test.$PAIR.spm \
            --tgtdict $FIRST_DICT \
            --srcdict $FIRST_DICT \
            --destdir $DESTDIR \
            --workers 8
    fi

  done
done