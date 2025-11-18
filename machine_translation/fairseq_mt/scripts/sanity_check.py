#!/usr/bin/env python3
import csv
import re
import statistics
from pathlib import Path
from collections import Counter

import sentencepiece as spm
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity
from sentence_transformers import SentenceTransformer, util
from sacrebleu import CHRF

# -----------------------------
# CONFIG
# -----------------------------
TSV_FILE = "../../../bible_cleaning/parallel_corpus/by_verses/ceb_tgl.tsv"
SPM_MODEL = "../data/unigram/spm_unigram.model"

LEN_RATIO_THRESHOLD = 3.0

# -----------------------------
# HELPERS
# -----------------------------

def normalize(text):
    return text.strip().lower()

def whitespace_tokenize(text):
    return text.split()

def compute_ttr(tokens):
    if len(tokens) == 0:
        return 0.0
    return len(set(tokens)) / len(tokens)

def safe_len_ratio(a, b):
    if min(len(a), len(b)) == 0:
        return 999
    return max(len(a), len(b)) / min(len(a), len(b))

def cosine(a, b):
    return util.cos_sim([a], [b])[0][0]



# LOAD DATA
print("Loading TSV...")

pairs = []  # (src, tgt)

with open(TSV_FILE, newline="", encoding="utf-8") as f:
    reader = csv.DictReader(f, delimiter="\t")
    for row in reader:
        if "source_text" not in row or "target_text" not in row:
            raise ValueError("TSV must have source_text and target_text columns.")
        src = normalize(row["source_text"])
        tgt = normalize(row["target_text"])
        pairs.append((src, tgt))

print(f"Loaded {len(pairs)} parallel pairs.\n")



# SETUP MODELS
# print("Loading SentencePiece...")
# sp = spm.SentencePieceProcessor()
# sp.load(SPM_MODEL) # type: ignore shush language parser-
model = SentenceTransformer("sentence-transformers/LaBSE")

# print("Loading LaBSE embeddings...")
# embedder = SentenceTransformer("sentence-transformers/LaBSE")

chrF = CHRF()



# GLOBAL TOKEN COUNTS
src_tokens_all = []
tgt_tokens_all = []

for src, tgt in pairs:
    src_tokens_all.extend(whitespace_tokenize(src))
    tgt_tokens_all.extend(whitespace_tokenize(tgt))

src_vocab = Counter(src_tokens_all)
tgt_vocab = Counter(tgt_tokens_all)


# -----------------------------
# METRIC STORAGE
# -----------------------------
length_diffs = []
length_ratios = []
rare_counts_src = []
rare_counts_tgt = []
chrf_scores = []
embedding_distances = []
fertilities_src = []
fertilities_tgt = []
flagged_len_outliers = []
flagged_embed_outliers = []


# -----------------------------
# MAIN LOOP
# -----------------------------
print("Computing metrics...\n")

for idx, (src, tgt) in enumerate(pairs):

    # ----- Tokenization -----
    src_tok = whitespace_tokenize(src)
    tgt_tok = whitespace_tokenize(tgt)

    # ----- Sentence Length Difference -----
    length_diffs.append(abs(len(src_tok) - len(tgt_tok)))

    # ----- Length Ratio -----
    ratio = safe_len_ratio(src_tok, tgt_tok)
    length_ratios.append(ratio)
    if ratio > LEN_RATIO_THRESHOLD:
        flagged_len_outliers.append((idx, ratio))

    # ----- Rare Token Frequency -----
    rare_counts_src.append(sum(1 for t in src_tok if src_vocab[t] == 1))
    rare_counts_tgt.append(sum(1 for t in tgt_tok if tgt_vocab[t] == 1))

    # ----- LASER/LaBSE Embedding Distance -----
    emb_src = model.encode(src, convert_to_tensor=True)
    emb_tgt = model.encode(tgt, convert_to_tensor=True)
    dist = cosine(emb_src, emb_tgt)
    embedding_distances.append(dist)
    if dist > 0.4:  # heuristic threshold, this sentence is an outlier
        flagged_embed_outliers.append((idx, dist))

    # ----- ChrF++ Similarity -----
    chrF_score = chrF.sentence_score(tgt, [src]).score
    chrf_scores.append(chrF_score)

# -----------------------------
# SUMMARY OUTPUT
# -----------------------------

print("====================================")
print("TYPE–TOKEN RATIO (TTR)")
print("====================================")
print(f"Source TTR: {compute_ttr(src_tokens_all):.4f}")
print(f"Target TTR: {compute_ttr(tgt_tokens_all):.4f}")
print(f"Combined TTR: {compute_ttr(src_tokens_all + tgt_tokens_all):.4f}\n")

print("====================================")
print("AVERAGE SENTENCE LENGTH")
print("====================================")
print(f"Mean sentence length difference: {statistics.mean(length_diffs):.3f}")
print(f"Mean length ratio: {statistics.mean(length_ratios):.3f}")
print(f"Length ratio outliers (> {LEN_RATIO_THRESHOLD}): {len(flagged_len_outliers)}\n")

print("====================================")
print("RARE TOKEN FREQUENCY")
print("====================================")
print(f"Avg rare tokens per source sentence: {statistics.mean(rare_counts_src):.3f}")
print(f"Avg rare tokens per target sentence: {statistics.mean(rare_counts_tgt):.3f}\n")

print("====================================")
print("SEMANTIC SIMILARITY (LaBSE)")
print("====================================")
print(f"Mean embedding distance: {statistics.mean(embedding_distances):.3f}")
print(f"Embedding outliers (>0.4): {len(flagged_embed_outliers)}\n")

print("====================================")
print("ChrF++ SIMILARITY")
print("====================================")
print(f"Mean ChrF++ score: {statistics.mean(chrf_scores):.3f}\n")

print("====================================")
print("SUBWORD FERTILITY (SentencePiece)")
print("====================================")
print(f"Avg fertility (src): {statistics.mean(fertilities_src):.3f}")
print(f"Avg fertility (tgt): {statistics.mean(fertilities_tgt):.3f}\n")

print("====================================")
print("DONE – FULL DATASET QUALITY REPORT")
print("====================================")
