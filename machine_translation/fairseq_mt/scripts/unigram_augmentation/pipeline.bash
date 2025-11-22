#!/bin/bash

# Check if at least two arguments were provided
if [ "$#" -lt 2 ]; then
  echo "Usage: $0 <arg1> <arg2> [remaining args...]"
  exit 1
fi

# Store the first two arguments
LANG_SRC="$1"
LANG_DST="$2"

DATA_PATH="../../data/raw"
DATA_BIN_PATH="../../data/augmentation"
SPM_MODEL_PATH="../../spm_model/augmentation"

mkdir ../../data/augmentation

echo "### Required Arguments ###"
echo "SRC LANG (\$1): ${LANG_SRC}"
echo "DST LANG (\$2): ${LANG_DST}"

# Shift the argument list twice to discard $1 and $2
shift 2 
echo "### Remaining Arguments ###"

AUGMENTED=()
# The total number of remaining arguments is now $#
if [ "$#" -eq 0 ]; then
  echo "No remaining variables (n=0)."
fi

# Loop through the remaining arguments using "$@"
# Note: "$@" now only contains the arguments that were originally $3, $4, ..., $n
for remaining_arg in "$@"; do
  echo "Augmentation: ${remaining_arg}"
  AUGMENTED+=("$remaining_arg")
done

echo "Total Augmentation (n): $#"
TOTAL_AUG=$#
TOTAL_AUG=$(($TOTAL_AUG + 1))

for field in "train" "valid" "test"; do

  cp "$DATA_PATH/$field.$LANG_SRC-$LANG_DST.$LANG_SRC" "$field.raw.src"
  cp "$DATA_PATH/$field.$LANG_SRC-$LANG_DST.$LANG_DST" "$field.raw.tgt.$LANG_DST"

  
  
  for lang_arg in $AUGMENTED; do
      cat "$DATA_PATH/$field.$LANG_SRC-$lang_arg.$LANG_SRC" >> "$field.raw.src" 
      cp $DATA_PATH/$field.$LANG_SRC-$lang_arg.$lang_arg $field.raw.tgt.$lang_arg
      mv $field.raw.tgt.$lang_arg $DATA_BIN_PATH/
  done

  mv $field.raw.src $DATA_BIN_PATH/
  mv $field.raw.tgt.$LANG_DST $DATA_BIN_PATH/
done

# Print it out
for field in field valid test; do
  echo "$DATA_BIN_PATH/$field.raw.src" | wc -l
  echo "$DATA_BIN_PATH/$field.raw.tgt.$LANG_DST" | wc -l
  for lang_arg in $AUGMENTED; do
    echo "$DATA_BIN_PATH/$field.raw.src.$lang_arg" | wc -l
  done
done



cd $DATA_BIN_PATH
for field in "train" "valid" "test"; do
  echo "Concatenating $field.spm.raw.tgt"
  cp "$field.raw.src" "$field.spm.raw.src"
  cp "$field.raw.tgt.$LANG_DST" "$field.spm.raw.tgt"

  for lang_arg in $AUGMENTED; do
    echo "$AUGMENTED"
    cat "$field.raw.tgt.$lang_arg" >> "$field.spm.raw.tgt"
  done 
done

uv run --active python ../../scripts/unigram_augmentation/train_spm.py $1 $2 $3

mkdir ../../spm_model
mkdir ../../spm_model/augmentation

cd ../../data/augmentation

rm -rf $SPM_MODEL_PATH/spm.src.vocab
rm -rf $SPM_MODEL_PATH/spm.src.model
rm -rf $SPM_MODEL_PATH/spm.tgt.vocab
rm -rf $SPM_MODEL_PATH/spm.tgt.model

mv ./spm.src.vocab $SPM_MODEL_PATH/
mv ./spm.src.model $SPM_MODEL_PATH/
mv ./spm.tgt.vocab $SPM_MODEL_PATH/
mv ./spm.tgt.model $SPM_MODEL_PATH/

uv run --active python ../../scripts/unigram_augmentation/encode.py $LANG_SRC $LANG_DST $@


for field in "train" "valid" "test"; do
  cat $field.spm.src > $field.spm.final.$LANG_SRC
  cat $field.spm.tgt.$LANG_DST > $field.spm.final.$LANG_DST
  for lang_aug in $AUGMENTED; do
    cat $field.spm.tgt.$lang_aug >> $field.spm.final.$LANG_DST
  done

done

cd ../..
cd data-bin
mkdir augmentation
cd augmentation

uv run --active fairseq-preprocess \
    --source-lang $LANG_SRC \
    --target-lang $LANG_DST \
    --trainpref $DATA_BIN_PATH/train.spm.final \
    --validpref $DATA_BIN_PATH/valid.spm.final \
    --testpref $DATA_BIN_PATH/test.spm.final \
    --destdir ../../data-bin/augmentation/$LANG_SRC-$LANG_DST/bin \
    --joined-dictionary \
    --workers 8 
