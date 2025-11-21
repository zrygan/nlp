#!/bin/bash

# Configuration
SRC_LANG="ceb"
TGT_LANG="tgl"
MODEL_PATH="../checkpoints/${SRC_LANG}-${TGT_LANG}/checkpoint_best.pt"
DATA_BIN="../data-bin/unigram/${SRC_LANG}-${TGT_LANG}"
RESULTS_DIR="results/unigram/${SRC_LANG}-${TGT_LANG}"

# Create results directory
mkdir -p ${RESULTS_DIR}

echo "=========================================="
echo "Evaluating ${SRC_LANG} → ${TGT_LANG} Model"
echo "=========================================="
echo ""
echo "Model: ${MODEL_PATH}"
echo "Data: ${DATA_BIN}"
echo ""

# Check if model exists
if [ ! -f "${MODEL_PATH}" ]; then
    echo "ERROR: Model not found at ${MODEL_PATH}"
    echo "Please train the model first or check the path."
    exit 1
fi

# Evaluate on test set
echo "Running evaluation on test set..."
uv run --active fairseq-generate ${DATA_BIN} \
    --path ${MODEL_PATH} \
    --batch-size 128 \
    --beam 5 \
    --remove-bpe \
    --gen-subset test \
    --results-path ${RESULTS_DIR} \
    2>&1 | tee ${RESULTS_DIR}/evaluation.log

# The actual translations are in generate-test.txt
GENERATE_FILE="${RESULTS_DIR}/generate-test.txt"

echo ""
echo "=========================================="
echo "Results Summary"
echo "=========================================="

# Extract BLEU score from evaluation.log (console output)
echo ""
echo "BLEU Score:"
BLEU=$(grep "BLEU4 = " ${RESULTS_DIR}/generate-test.txt | grep -oP "BLEU4 = \K[0-9.]+")
if [ -z "$BLEU" ]; then
    BLEU=$(grep "BLEU = " ${RESULTS_DIR}/generate-test.txt | grep -oP "BLEU = \K[0-9.]+")
fi
if [ -n "$BLEU" ]; then
    echo "  BLEU: $BLEU"
else
    echo "  (BLEU score not found - check evaluation.log)"
fi

echo ""
echo "Full evaluation log: ${RESULTS_DIR}/evaluation.log"
echo "Generated translations: ${GENERATE_FILE}"

# Check if generate-test.txt exists
if [ ! -f "${GENERATE_FILE}" ]; then
    echo ""
    echo "ERROR: ${GENERATE_FILE} not found!"
    exit 1
fi

# Extract sample translations from generate-test.txt
echo ""
echo "Sample Translations (first 10):"
grep "^S-\|^T-\|^H-\|^D-" ${GENERATE_FILE} | head -40

echo ""
echo "=========================================="
echo "Additional Analysis"
echo "=========================================="

# Count translations from generate-test.txt
NUM_TRANSLATIONS=$(grep "^H-" ${GENERATE_FILE} | wc -l)
echo "Total translations: ${NUM_TRANSLATIONS}"

if [ ${NUM_TRANSLATIONS} -gt 0 ]; then
    # Extract all hypotheses (3rd field)
    grep "^H-" ${GENERATE_FILE} | cut -f3 > ${RESULTS_DIR}/hypotheses.txt
    echo "All translations saved to: ${RESULTS_DIR}/hypotheses.txt"
    
    # Extract all references (2nd field)
    grep "^T-" ${GENERATE_FILE} | cut -f2 > ${RESULTS_DIR}/references.txt
    echo "All references saved to: ${RESULTS_DIR}/references.txt"
    
    # Extract all sources (2nd field)
    grep "^S-" ${GENERATE_FILE} | cut -f2 > ${RESULTS_DIR}/sources.txt
    echo "All sources saved to: ${RESULTS_DIR}/sources.txt"
    
    # Show statistics
    echo ""
    echo "Translation Statistics:"
    echo "  Average source length: $(awk '{print NF}' ${RESULTS_DIR}/sources.txt | awk '{s+=$1} END {printf "%.1f", s/NR}') words"
    echo "  Average translation length: $(awk '{print NF}' ${RESULTS_DIR}/hypotheses.txt | awk '{s+=$1} END {printf "%.1f", s/NR}') words"
    
    # Show a few example translations side-by-side
    echo ""
    echo "Example Translations:"
    echo "---"
    paste ${RESULTS_DIR}/sources.txt ${RESULTS_DIR}/hypotheses.txt ${RESULTS_DIR}/references.txt | head -5 | awk -F'\t' '{printf "SRC: %s\nHYP: %s\nREF: %s\n---\n", $1, $2, $3}'
else
    echo "WARNING: No translations found. Check ${GENERATE_FILE} for errors."
fi

echo ""
echo "Evaluation complete!"