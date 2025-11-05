import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

# Load your TSV or CSV file
# For TSV: sep='\t'
df = pd.read_csv('orthographic_similarity_matrix.tsv', sep='\t', index_col=0)

# Create a heatmap
plt.figure(figsize=(8, 6))
sns.heatmap(df, cmap='viridis', annot=False)

plt.title('Orthographic Similarity Matrix')
plt.xlabel('Language A')
plt.ylabel('Language B')

# Save as PNG or PDF for LaTeX
plt.tight_layout()
plt.savefig('orthographic_heat_map', dpi=300)
plt.show()
