import os
from pathlib import Path
import csv
import random

# --- CONFIG ---
TRAIN_PCT = 0.90
VALID_PCT = 0.05
TEST_PCT  = 0.05
assert abs(TRAIN_PCT + VALID_PCT + TEST_PCT - 1.0) < 1e-6
# -------------

root = Path("../../bible_cleaning/parallel_corpus/by_verses")

# Detect TSV files
tsv_files = list(root.glob("*.tsv"))
print(f"Found {len(tsv_files)} TSV files")

if len(tsv_files) == 0:
    print("ERROR: No TSV files found in", root)
    exit(1)

os.makedirs("data", exist_ok=True)

for tsv_file in tsv_files:
    # Parse filename: srclang_tgtlang.tsv
    filename = tsv_file.stem
    parts = filename.split('_')
    
    if len(parts) != 2:
        print(f"Skipping bad filename: {tsv_file.name} (expected format: srclang_tgtlang.tsv)")
        continue
    
    langA, langB = parts[0], parts[1]
    print(f"\n{'='*60}")
    print(f"Processing: {langA} → {langB}")
    print(f"{'='*60}")
    
    # Read TSV file
    aligned_pairs = []
    
    with open(tsv_file, 'r', encoding='utf-8') as f:
        reader = csv.DictReader(f, delimiter='\t')
        
        # Validate headers
        if not('source_text' in reader.fieldnames) or not ('target_text' in reader.fieldnames): #type: ignore
            print(f"ERROR: Missing 'source_text' or 'target_text' columns in {tsv_file.name}")
            print(f"Found columns: {reader.fieldnames}")
            continue
        
        for row in reader:
            source = row['source_text'].strip()
            target = row['target_text'].strip()
            
            # Skip empty pairs
            if source and target:
                aligned_pairs.append((source, target))
    
    print(f"Loaded {len(aligned_pairs)} aligned verse pairs")
    
    if len(aligned_pairs) == 0:
        print(f"WARNING: No valid pairs found in {tsv_file.name}, skipping...")
        continue
    
    # Shuffle before splitting
    random.seed(42)
    random.shuffle(aligned_pairs)
    
    # Split into train/valid/test
    n = len(aligned_pairs)
    n_train = int(n * TRAIN_PCT)
    n_valid = int(n * VALID_PCT)
    n_test  = n - n_train - n_valid
    
    train_data = aligned_pairs[:n_train]
    valid_data = aligned_pairs[n_train:n_train+n_valid]
    test_data  = aligned_pairs[n_train+n_valid:]
    
    print(f"Split: train={len(train_data)}, valid={len(valid_data)}, test={len(test_data)}")
    
    def write_split(split_name, data, langA, langB):
        """Write parallel data in fairseq format"""
        if len(data) == 0:
            print(f"  WARNING: {split_name} has no data, skipping...")
            return
        
        # Fairseq expects languages in alphabetical order
        if langA < langB:
            src, tgt = langA, langB
            pairs = data
        else:
            src, tgt = langB, langA
            # Swap source and target
            pairs = [(b, a) for (a, b) in data]
        
        src_path = f"data/{split_name}.{src}-{tgt}.{src}"
        tgt_path = f"data/{split_name}.{src}-{tgt}.{tgt}"
        
        with open(src_path, "w", encoding="utf-8") as fa, \
             open(tgt_path, "w", encoding="utf-8") as fb:
            for s, t in pairs:
                fa.write(s + "\n")
                fb.write(t + "\n")
        
        # Verify line counts
        with open(src_path, 'r', encoding='utf-8') as fa, \
             open(tgt_path, 'r', encoding='utf-8') as fb:
            src_lines = sum(1 for _ in fa)
            tgt_lines = sum(1 for _ in fb)
            
            if src_lines != tgt_lines:
                print(f"  ❌ ERROR: Line count mismatch! {src}={src_lines}, {tgt}={tgt_lines}")
            else:
                print(f"  ✓ {split_name}: {len(pairs)} lines → {src_path}, {tgt_path}")
    
    write_split("train", train_data, langA, langB)
    write_split("valid", valid_data, langA, langB)
    write_split("test",  test_data,  langA, langB)

print("\n" + "="*60)
print("✓ Preprocessing complete!")
print("="*60)
print("\nNext steps:")
print("1. ./fairseq_preprocess\n")
print("2. ./faisreq_train\n")
