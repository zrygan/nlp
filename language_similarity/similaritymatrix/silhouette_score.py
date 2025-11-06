import pandas as pd
from sklearn.cluster import KMeans
from sklearn.metrics import silhouette_score

# === CONFIG ===
FILE_PATH = "orthographic_similarity_matrix.tsv"   # path to your TSV file

# === LOAD DATA ===
# Assumes numeric columns only (or already preprocessed)
data = pd.read_csv(FILE_PATH, sep="\t", index_col=0)

# If your TSV has headers and non-numeric columns, you can filter:
# data = data.select_dtypes(include='number')

# === FIT CLUSTERING ===
for k in range(2, 11):
    kmeans = KMeans(n_clusters=k, random_state=42)
    labels = kmeans.fit_predict(data)
    score = silhouette_score(data, labels)
    print(f"k={k}: Silhouette Score = {score:.4f}")