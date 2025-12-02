import itertools
import os
import sentencepiece as spm

# -------------------------
# Configuration
# -------------------------

LANGUAGES = [
    "bik", "cbk", "ceb", "eng", "ilo", "jil", "krj", "msg", 
    "pag", "pam", "prf", "rol", "tao", "tgl", "tiu",
    "tsg", "war"
]

SPLITS = ["train", "valid", "test"]
MODEL_PATH = "../unigram/spm_unigram.model"
DATA_DIR = "../../data/raw/"
OUT_DIR = "../../data/spm"
os.chdir("../../data/spm/")

# -------------------------
# Load SPM model
# -------------------------

sp = spm.SentencePieceProcessor()
sp.load(MODEL_PATH) # type: ignore

# -------------------------
# Generate unique sorted pairs (LANG1 < LANG2)
# -------------------------

# -------------------------
# Encode all files
# -------------------------

for split in SPLITS:
    for lang1 in LANGUAGES:
        for lang2 in LANGUAGES:
            if lang1 == lang2:
                continue
            
            pair_str = (lang1, lang2)
            
            pair = f"{lang1}-{lang2}"

            for lang in pair_str:
                input_file = os.path.join(DATA_DIR, f"{split}.{pair}.{lang}")
                output_file = os.path.join(OUT_DIR, f"{split}.{pair}.spm.{lang}")

                if not os.path.exists(input_file):
                    # No parallel file → skip silently
                    print(f"Skipping... {input_file}")
                    continue

                print(f"Encoding {input_file} → {output_file}")

                with open(input_file, "r", encoding="utf-8") as f_in:
                    lines = [line.strip() for line in f_in]

                encoded = [ " ".join(sp.encode(line, out_type=str )[1:]) for line in lines ] # type: ignore

                with open(output_file, "w", encoding="utf-8") as f_out:
                    for line in encoded:
                        f_out.write(line + "\n")

print(":) Finished encoding all language pairs.")
