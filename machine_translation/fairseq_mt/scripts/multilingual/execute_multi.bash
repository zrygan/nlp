# uv run --active python ./multi.py

# uv run --active fairseq-preprocess \
#   --task translation_multi_simple_epoch \
#   --source-lang eng \
#   --target-lang tgl \
#   --trainpref ../../data/spm_combined/train \
#   --validpref ../../data/spm_combined/valid \
#   --testpref  ../../data/spm_combined/test \
#   --joined-dictionary \
#   --destdir ../../data-bin/multilingual \
#   --workers 32
  
# fairseq-preprocess \
#   --source-lang en --target-lang fr,de \
#   --trainpref path/to/data/train.en-fr,path/to/data/train.en-de \
#   --validpref path/to/data/valid.en-fr,path/to/data/valid.en-de \
#   --destdir data-bin/multilingual \
#   --workers 10 \
#   --joined-dictionary \
#   --srcdict path/to/shared_vocab.txt

# uv run --active fairseq-preprocess \
#   --source-lang eng\
#   --target-lang bik,cbk,ceb,ilo,jil,ilo,jil,krj,msb,pag,pam,prf,rol,tao,tgl,tiu,tsg,war\
#   --trainpref ../../data/spm_combined/train \
#   --validpref ../../data/spm_combined/valid \
#   --testpref  ../../data/spm_combined/test \
#   --destdir ../../data-bin/multilingual\
#   --workers 8 \
#   --joined-dictionary 

# uv run --active fairseq-preprocess\
#   --source-lang eng\
#   --target-lang tgl\
#   --trainpref ../../data/spm/train \
#   --validpref ../../data/spm/valid \
#   --testpref  ../../data/spm/test \
#   --destdir data-bin/multilingual_data \
#   --workers 8 \
#   --joined-dictionary 

# uv run --active fairseq-preprocess \
#   --source-lang fil --target-lang en \
#   --trainpref $TEXT/train.spm.fil \
#   --validpref $TEXT/valid.spm.fil \
#   --testpref $TEXT/test.spm.fil \
#   --tgtdir $TEXT/train.spm.en \
#   --validdir $TEXT/valid.spm.en \
#   --testdir $TEXT/test.spm.en \
#   --destdir data-bin/multilingual_data \
#   --srcdict data-bin/multilingual_data/dict.en.txt \
#   --tgtdict data-bin/multilingual_data/dict.en.txt \
#   --workers 10

LANG1="eng"
LANG2="ceb"
SPM_DIR="../../data/spm/"

cat 

uv run --active fairseq-preprocess \
  --task translation_multi_simple_epoch \
  --trainpref ../../data/spm_combined/train \
  --validpref ../../data/spm_combined/valid \
  --testpref  ../../data/spm_combined/test \
  --source-lang eng \
  --target-lang ceb_tgl \
  --joined-dictionary \
  --bpe sentencepiece\
  --destdir ../../data-bin/multilingual \
  --workers 8
