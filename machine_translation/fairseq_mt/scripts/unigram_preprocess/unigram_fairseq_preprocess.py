import os
import re
import subprocess
from pathlib import Path

DATA_DIR = Path("../../data/spm/")
DESTDIR = "../../data-bin/"

PAIR_REGEX = re.compile(r"(train|valid|test)\.([a-z]+)-([a-z]+)\.spm\.[a-z]+")

def find_unique_pairs():
    """Scan data/ folder for all valid parallel corpus pairs."""
    pairs = set()
    
    for fname in os.listdir(DATA_DIR):
        match = PAIR_REGEX.match(fname)
        if match:
            _, src, tgt = match.groups()
            pair = tuple(sorted((src, tgt)))
            pairs.add(pair)
    
    return sorted(list(pairs))


def build_argument_list(pairs):
    """Build the full argument list in the bash printf command."""
    args = []
    
    for src, tgt in pairs:
        prefix = f"{src}-{tgt}"
        args.extend([
            f"--trainpref train.{prefix}.spm",
            f"--validpref valid.{prefix}.spm",
            f"--testpref test.{prefix}.spm",
        ])
    
    return args


def run_fairseq_preprocess(args):
    """Run fairseq-preprocess with all generated pairs."""
    base_cmd = [
        "fairseq-preprocess",
        "--joined-dictionary",
        "--workers", "8",
        "--destdir", DESTDIR
    ]

    full_cmd = base_cmd + args

    print("\nRunning command:\n")
    print(" ".join(full_cmd), "\n")

    subprocess.run(full_cmd, check=True)
    print("\n✓ Preprocessing complete!")


def main():
    pairs = find_unique_pairs()

    if not pairs:
        print(":( No valid language pairs found in data/ folder!")
        return

    print("Detected language pairs:")
    for src, tgt in pairs:
        print(f" - {src}-{tgt}")
    os.chdir(DATA_DIR)
    args = build_argument_list(pairs)
    run_fairseq_preprocess(args)


if __name__ == "__main__":
    main()
