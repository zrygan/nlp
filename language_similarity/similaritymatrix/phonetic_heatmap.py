import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

# Load your TSV or CSV file
# For TSV: sep='\t'
df = pd.read_csv('phonetic_similarity_matrix.tsv', sep='\t', index_col=0)

# The matrix is symmetric and its diagonal is 1.0 by construction, so only the
# strict lower triangle carries information. Dropping the first row and last
# column removes the two bands that the diagonal mask would leave empty, giving a
# 15x15 staircase. Excluding the diagonal also frees the colour scale for the
# values that vary, and lets each cell be labelled without a leading zero, which
# is what keeps the annotations legible once the figure is scaled to a column.
df = df.iloc[1:, :-1]
mask = np.triu(np.ones_like(df, dtype=bool), k=1)

# ".7622" rather than "0.7622": every off-diagonal value is < 1, and dropping the
# leading zero is what lets four decimal places fit the cell.
labels = np.vectorize(lambda v: f'{v:.4f}'[1:])(df.values)

plt.figure(figsize=(7.0, 6.4))
sns.heatmap(df, mask=mask, cmap='viridis', annot=labels, fmt='',
            annot_kws={'size': 9}, linewidths=0.4, linecolor='white',
            square=True, cbar=False)

plt.xlabel('Language A')
plt.ylabel('Language B')
plt.yticks(rotation=0)

# Save as PNG or PDF for LaTeX
plt.tight_layout()
plt.savefig('phonetic_heat_map', dpi=300, bbox_inches='tight')
plt.show()
