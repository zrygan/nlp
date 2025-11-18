cd parallel_corpus_preprocess

uv run --active python training_data_creation.py
ls

cd ..
cd unigram_preprocess

./unigram_conversion.bash

cd ..
cd train

./train_src_dest.bash ceb tgl

cd ..

./fairseq_evaluate.bash