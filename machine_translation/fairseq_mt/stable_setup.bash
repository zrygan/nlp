#!/bin/bash
# Complete stable fairseq setup for training

echo "Setting up stable fairseq environment..."

# Use uv to install stable versions
uv pip uninstall fairseq torch torchvision torchaudio 

# Install PyTorch 2.0.1 (most stable with fairseq)
uv pip install torch==2.0.1 torchvision torchaudio

# Install fairseq 0.12.2
uv pip install fairseq==0.12.2

echo "✓ Installed PyTorch 2.0.1 + fairseq 0.12.2"
echo ""
echo "This combination requires NO patches and is fully stable!"
echo ""
echo "Now you can:"
echo "1. Delete old incompatible checkpoints: rm -rf checkpoints/ceb-tgl/*"
echo "2. Retrain: ./fairseq_train.bash"
echo "3. Evaluate: ./evaluate.sh"