import sentencepiece as spm
import os
from pathlib import Path 

def load_user_defined_symbols(path="./special_tokens.txt"):
    """Load user-defined symbols as a comma-separated string."""
    p = Path(path)
    with open(p, "r", encoding="utf-8") as f:
        symbols = [line.strip() for line in f if line.strip()]
    return ",".join(symbols)

def train_spm():
    user_symbols = load_user_defined_symbols("./special_tokens.txt")

    os.chdir("../../data/unigram/")

    spm.SentencePieceTrainer.train( # type: ignore
        input="all_raw_languages_training_data.txt",
        model_prefix="spm_unigram",
        vocab_size=4000,
        character_coverage=1.0,
        model_type="unigram",
        user_defined_symbols=user_symbols
    )
    
    print(":) SentencePiece training complete!")
    print("Generated: spm_unigram.model, spm_unigram.vocab")

if __name__ == "__main__":
    train_spm()
