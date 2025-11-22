import sentencepiece as spm
import os
from pathlib import Path 

def load_user_defined_symbols(path="./special_tokens.txt"):
    with open(path, "r", encoding="utf-8") as f:
        return [line.strip() for line in f if line.strip()]

def train_spm():
    user_symbols = load_user_defined_symbols("./special_tokens.txt")

    os.chdir("../../data/unigram/")

    spm.SentencePieceTrainer.train( # type: ignore
        input="./all_train.txt",
        model_prefix="spm_unigram",
        vocab_size=32000,
        character_coverage=1.0,
        model_type="unigram",
        treat_whitespace_as_suffix=False,
        normalization_rule_name="nmt_nfkc",
        input_sentence_size=8000000,
        user_defined_symbols=user_symbols,
        shuffle_input_sentence=True
    )
    
    print(":) SentencePiece training complete!")
    print("Generated: spm_unigram.model, spm_unigram.vocab")

if __name__ == "__main__":
    train_spm()
