#!/usr/bin/env python3
import os
import subprocess
from pathlib import Path

#############################################
# USER CONFIG
#############################################

LANGS = [
    "bik","cbk","ceb","eng","ilo","jil","krj",
    "msb","pag","pam","prf","rol","tao","tgl",
    "tiu","tsg","war"
]

CANONICAL_SRC = "eng"
CANONICAL_TGT = "tgl"
CANONICAL_PAIR = f"{CANONICAL_SRC}-{CANONICAL_TGT}"

DATA_SPM = "../../data/spm"
DESTDIR = "../../data-bin/multilingual"

#############################################
# SETUP
#############################################

os.makedirs(DESTDIR, exist_ok=True)

def file_exists(prefix, src, tgt):
    return (
        Path(f"{prefix}.{src}").exists()
        and Path(f"{prefix}.{tgt}").exists()
    )

#############################################
# 1. Process canonical pair FIRST
#############################################

print(f"\n=== Processing canonical pair {CANONICAL_PAIR} ===\n")

trainpref = f"{DATA_SPM}/train.{CANONICAL_PAIR}.spm"
validpref = f"{DATA_SPM}/valid.{CANONICAL_PAIR}.spm"
testpref  = f"{DATA_SPM}/test.{CANONICAL_PAIR}.spm"

cmd = [
    "fairseq-preprocess",
    "--source-lang", CANONICAL_SRC,
    "--target-lang", CANONICAL_TGT,
    "--trainpref", trainpref,
    "--validpref", validpref,
    "--testpref", testpref,
    "--destdir", DESTDIR,
    "--joined-dictionary",   # <-- FIXED
    "--workers", "8",
]

print(" ".join(cmd))
subprocess.run(cmd, check=True)

canonical_srcdict = f"{DESTDIR}/dict.{CANONICAL_SRC}.txt"
canonical_tgtdict = f"{DESTDIR}/dict.{CANONICAL_TGT}.txt"

print(f"Canonical dict src = {canonical_srcdict}")
print(f"Canonical dict tgt = {canonical_tgtdict}")

#############################################
# 2. Process ALL other pairs reusing ONLY canonical dicts
#############################################

print("\n=== Processing all other language pairs ===\n")

for src in LANGS:
    for tgt in LANGS:
        if src == tgt:
            continue
        if f"{src}-{tgt}" == CANONICAL_PAIR:
            continue

        pair = f"{src}-{tgt}"
        trainpref = f"{DATA_SPM}/train.{pair}.spm"
        validpref = f"{DATA_SPM}/valid.{pair}.spm"
        testpref  = f"{DATA_SPM}/test.{pair}.spm"

        if not file_exists(trainpref, src, tgt):
            print(f"[Skipping] Missing: {pair}")
            continue

        print(f"\n--- Processing {pair} ---")

        cmd = [
            "fairseq-preprocess",
            "--task", "translation",
            "--joined-dictionary",
            "--source-lang", src,
            "--target-lang", tgt,
            "--srcdict", canonical_srcdict,
            "--trainpref", trainpref,
            "--validpref", validpref,
            "--testpref", testpref,
            "--destdir", DESTDIR,
            "--dict-only",
            "--thresholdsrc", "0",
            "--thresholddst", "0",
            "--workers", "8",
        ]


        print(" ".join(cmd))
        subprocess.run(cmd, check=True)

print("\n=== ALL DONE ===\n")
