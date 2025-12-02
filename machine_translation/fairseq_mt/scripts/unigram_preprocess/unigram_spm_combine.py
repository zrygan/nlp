import re
import os
from pathlib import Path
from collections import defaultdict

SRC_DIR = Path("../../data/spm")
DEST_DIR = Path("../../data/spm_combined")
DEST_DIR.mkdir(exist_ok=True, parents=True)

# Match filenames like:
#   train.eng-tgl.spm.eng
#   train.eng-bik.spm.bik
pattern = re.compile(r"(train|valid|test)\.[^.]+\.spm\.([a-z]+)$")

# Collect: split → lang → list of files
files_by_lang = defaultdict(lambda: defaultdict(list))

for f in SRC_DIR.iterdir():
    if not f.is_file():
        continue
    m = pattern.match(f.name)
    if not m:
        continue

    split, lang = m.group(1), m.group(2)
    files_by_lang[split][lang].append(f)

print("Combining SPM files (preserving parallel alignment)...")

for split in ["train", "valid", "test"]:
    for lang, files in files_by_lang[split].items():
        dest = DEST_DIR / f"{split}.{lang}"
        print(f" → {dest} ({len(files)} files)")

        with dest.open("w", encoding="utf-8") as out_f:
            for f in sorted(files):
                with f.open("r", encoding="utf-8") as in_f:
                    for line in in_f:
                        out_f.write(line)

print("✔ Done.")
