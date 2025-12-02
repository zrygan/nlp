
cd unigram_augmentation

./pipeline.bash eng tgl ceb

cd ..
cd train

./train_src_dest.bash

cd ..

./fairseq_evaluate.bash

Augmentation
