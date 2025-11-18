import itertools
import os
import sentencepiece as spm

# -------------------------
# Configuration
# -------------------------

LANGUAGES = [
    "bik", "cbk", "ceb", "ilo", "jil", "krj", "msg",
    "pag", "pam", "prf", "rol", "tao", "tgl", "tiu",
    "tsg", "war"
]

SPLITS = ["train", "valid", "test"]
MODEL_PATH = "../../data/unigram/spm_unigram.model"
DATA_DIR = "../../data/raw/"
OUT_DIR = "../../data/spm"

# -------------------------
# Load SPM model
# -------------------------

sp = spm.SentencePieceProcessor()
sp.load(MODEL_PATH) # type: ignore

# -------------------------
# Generate unique sorted pairs (LANG1 < LANG2)
# -------------------------

unique_pairs = [
    (l1, l2)
    for l1, l2 in itertools.combinations(sorted(LANGUAGES), 2)
]

print(f"Generated {len(unique_pairs)} unique language pairs.")

# -------------------------
# Encode all files
# -------------------------
os.chdir("../../data/spm/")
for split in SPLITS:
    for lang1, lang2 in unique_pairs:
        pair = f"{lang1}-{lang2}"

        for lang in (lang1, lang2):
            input_file = os.path.join(DATA_DIR, f"{split}.{pair}.{lang}")
            output_file = os.path.join(OUT_DIR, f"{split}.{pair}.spm.{lang}")

            if not os.path.exists(input_file):
                # No parallel file → skip silently
                print(f"Skipping... {input_file}")
                continue

            print(f"Encoding {input_file} → {output_file}")

            with open(input_file, "r", encoding="utf-8") as f_in:
                lines = [line.strip() for line in f_in]

            encoded = [ " ".join(sp.encode(line, out_type=str)) for line in lines ] # type: ignore

            with open(output_file, "w", encoding="utf-8") as f_out:
                for line in encoded:
                    f_out.write(line + "\n")

print(":) Finished encoding all language pairs.")
