import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

# Load your TSV or CSV file
# For TSV: sep='\t'
df = pd.read_csv('orthographic_similarity_matrix.tsv', sep='\t', index_col=0)

# The matrix is symmetric and its diagonal is 1.0 by construction, so only the
# strict lower triangle carries information. Dropping the first row and last
# column removes the two bands that the diagonal mask would leave empty, giving a
# 15x15 staircase. Excluding the diagonal also frees the colour scale for the
# values that vary rather than pinning it at the self-similarity of 1.0.
df = df.iloc[1:, :-1]
mask = np.triu(np.ones_like(df, dtype=bool), k=1)

# Values are read from the colour bar rather than printed in the cells.
plt.figure(figsize=(7.0, 6.4))
sns.heatmap(df, mask=mask, cmap='viridis', annot=False,
            linewidths=0.4, linecolor='white', square=True,
            cbar_kws={'shrink': 0.6, 'pad': 0.02, 'aspect': 40})

plt.xlabel('Language A')
plt.ylabel('Language B')
plt.yticks(rotation=0)

# Save as PNG or PDF for LaTeX
plt.tight_layout()
plt.savefig('orthographic_heat_map', dpi=300, bbox_inches='tight')
plt.show()
