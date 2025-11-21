#!/bin/bash
set -e

echo "Installing fairseq with uv..."
uv pip uninstall fairseq 
uv pip install fairseq  # Install latest (same version you trained with)

FAIRSEQ_PATH=$(uv run python -c "import fairseq, os; print(os.path.dirname(fairseq.__file__))")
echo "Fairseq installed at: $FAIRSEQ_PATH"

# Patch 1: checkpoint_utils.py
echo "Applying checkpoint patch..."
uv run python << EOF
with open("$FAIRSEQ_PATH/checkpoint_utils.py", 'r') as f:
    content = f.read()
content = content.replace(
    'state = torch.load(f, map_location=torch.device("cpu"))',
    'state = torch.load(f, map_location=torch.device("cpu"), weights_only=False)'
)
with open("$FAIRSEQ_PATH/checkpoint_utils.py", 'w') as f:
    f.write(content)
print("✓ Checkpoint patch applied")
EOF

# Patch 2: transformer_layer.py - disable fast path
echo "Applying transformer patch..."
uv run python << EOF
with open("$FAIRSEQ_PATH/modules/transformer_layer.py", 'r') as f:
    lines = f.readlines()

# Find and disable the fast path by adding False at the start
new_lines = []
for i, line in enumerate(lines):
    if 'if (' in line and i+1 < len(lines) and 'self.BT_version' in lines[i+1]:
        new_lines.append(line)
        new_lines.append(lines[i+1].replace('self.BT_version', 'False and self.BT_version'))
        i += 1
    else:
        new_lines.append(line)

with open("$FAIRSEQ_PATH/modules/transformer_layer.py", 'w') as f:
    f.writelines(new_lines)
print("✓ Transformer patch applied")
EOF

echo "Testing..."
uv run fairseq-generate data-bin/ceb-tgl \
    --path checkpoints/ceb-tgl/checkpoint_best.pt \
    --batch-size 32 \
    --beam 5 \
    --gen-subset test \
    | head -50