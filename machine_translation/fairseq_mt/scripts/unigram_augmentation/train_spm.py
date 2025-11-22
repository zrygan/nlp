import os
import re
import subprocess
from pathlib import Path
import sentencepiece as spm


def load_user_defined_symbols(path="./special_tokens.txt"):
    with open(path, "r", encoding="utf-8") as f:
        return [line.strip() for line in f if line.strip()]

def train_spm():
  
    user_symbols = load_user_defined_symbols("../../scripts/unigram_augmentation/special_tokens.txt")

    os.chdir("../../data/augmentation/")
    print(os.listdir())
    spm.SentencePieceTrainer.train( # type: ignore
        input="./train.spm.raw.src",
        model_prefix="spm.src",
        vocab_size=8000,
        character_coverage=1.0,
        model_type="unigram",
        normalization_rule_name="nmt_nfkc",
        user_defined_symbols=user_symbols,
        treat_whitespace_as_suffix=False,
        shuffle_input_sentence=True,
        input_sentence_size=10000000,
        num_threads=8
    )
    
    spm.SentencePieceTrainer.train( # type: ignore
        input="./train.spm.raw.tgt",
        model_prefix="spm.tgt",
        vocab_size=8000,
        character_coverage=1.0,
        model_type="unigram",
        normalization_rule_name="nmt_nfkc",
        user_defined_symbols=user_symbols,
        treat_whitespace_as_suffix=False,
        shuffle_input_sentence=True,
        input_sentence_size=10000000,
        num_threads=8
    )
    
    print(":) SentencePiece training complete!")
    print("Generated: spm_unigram.model, spm_unigram.vocab")
if __name__ == "__main__":
  train_spm()