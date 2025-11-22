import sentencepiece as spm
from pathlib import Path
import sys

# Example usage:
# python encode.py eng tgl ceb ilo ...
# LANG_SRC = eng
# LANG_DST = tgl
# AUG_LANGS = [ceb, ilo]

LANG_SRC = sys.argv[1]
LANG_DST = sys.argv[2]
AUG_LANGS = sys.argv[3:]  # augmented target languages

DATA_PATH = "../../data/augmentation"
SPM_MODEL_PATH = "../../spm_model/augmentation"


# ----------------------------------
# helpers
# ----------------------------------

def get_tag(lang):
    return f"<2{lang}>"


def encode_file(sp, infile, outfile, tag):
    with open(infile, "r", encoding="utf-8") as fin, \
         open(outfile, "w", encoding="utf-8") as fout:

        for line in fin:
            line = line.strip()
            if not line:
                fout.write("\n")
                continue
            # 1. SPM segmentation
            pieces = sp.encode(line, out_type=str)

            # 2. Add language tag for this file
            fout.write(tag + " " + " ".join(pieces) + "\n")


# ----------------------------------
# main
# ----------------------------------

def main():
    # Load SPM models
    sp_src = spm.SentencePieceProcessor()
    sp_src.load(f"{SPM_MODEL_PATH}/spm.src.model")

    sp_tgt = spm.SentencePieceProcessor()
    sp_tgt.load(f"{SPM_MODEL_PATH}/spm.tgt.model")

    for split in ["train", "valid", "test"]:
        print(split)
        # ----------------------------------------
        # ENCODE SOURCE ONCE
        # ----------------------------------------
        raw_src = Path(f"{DATA_PATH}/{split}.raw.src")
        out_src = Path(f"{DATA_PATH}/{split}.spm.src")

        encode_file(
            sp_src,
            raw_src,
            out_src,
            get_tag(LANG_SRC)  # source tag
        )

        # ----------------------------------------
        # ENCODE TARGET FOR DST (main direction)
        # ----------------------------------------
        raw_tgt = Path(f"{DATA_PATH}/{split}.raw.tgt.{LANG_DST}")
        out_tgt = Path(f"{DATA_PATH}/{split}.spm.tgt.{LANG_DST}")

        encode_file(
            sp_tgt,
            raw_tgt,
            out_tgt,
            get_tag(LANG_DST)
        )

        # ----------------------------------------
        # ENCODE EACH AUGMENTED TARGET LANGUAGE
        # ----------------------------------------
        for lang in AUG_LANGS:
            raw_aug = Path(f"{DATA_PATH}/{split}.raw.tgt.{lang}")
            out_aug = Path(f"{DATA_PATH}/{split}.spm.tgt.{lang}")

            encode_file(
                sp_tgt,
                raw_aug,
                out_aug,
                get_tag(lang)
            )


if __name__ == "__main__":
    main()
