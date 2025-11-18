#!/bin/bash
# Complete stable fairseq setup for training
echo "Setting up stable fairseq environment..."

# Use uv to install stable versions
uv pip uninstall fairseq torch torchvision torchaudio 

# Install PyTorch 2.0.1 (most stable with fairseq)
uv pip install torch==2.0.1 torchvision torchaudio

# Install fairseq 0.12.2
uv pip install fairseq==0.12.2

cat << EOF
:) Installed PyTorch 2.0.1 + fairseq 0.12.2
Now, you can perform the following:
  1. Training data creation:
    uv run --active python ./training_data creation
  2. Execute the three training scripts:
    ./fairseq_train.bash
EOF
