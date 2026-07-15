"""Cross-metric comparison of the orthographic and phonetic similarity matrices.

Produces the summary statistics quoted in the paper's final Results and
Discussion: the per-metric off-diagonal ranges, and the rank agreement between
the two metrics over all 120 language pairs.
"""

import numpy as np
import pandas as pd
from scipy.stats import spearmanr

ortho = pd.read_csv('orthographic_similarity_matrix.tsv', sep='\t', index_col=0)
phon = pd.read_csv('phonetic_similarity_matrix.tsv', sep='\t', index_col=0)

# Align phonetic to orthographic ordering; the two files are not stored in the
# same row order.
phon = phon.loc[ortho.index, ortho.columns]

for name, df in [('Orthographic', ortho), ('Phonetic', phon)]:
    off = df.values[~np.eye(len(df), dtype=bool)]
    print(f'{name}: min={off.min():.4f} max={off.max():.4f} '
          f'mean={off.mean():.4f} spread={off.max() - off.min():.4f}')

# Rank agreement over the 120 distinct pairs (upper triangle, diagonal excluded).
iu = np.triu_indices(len(ortho), k=1)
rho, p = spearmanr(ortho.values[iu], phon.values[iu])
print(f'\nPairwise rank agreement (n={len(iu[0])}): Spearman rho={rho:.3f}, p={p:.2e}')

# Per-language mean similarity to all others, which gives the centre-periphery
# ordering discussed in the text.
means = pd.DataFrame({
    'ortho': {lang: ortho.loc[lang].drop(lang).mean() for lang in ortho.index},
    'phon': {lang: phon.loc[lang].drop(lang).mean() for lang in phon.index},
})
rho_m, p_m = spearmanr(means['ortho'], means['phon'])
print(f'Per-language mean similarity: Spearman rho={rho_m:.3f}, p={p_m:.2e}\n')
print(means.sort_values('ortho').round(4).to_string())

# --- Corpus-size confound check -------------------------------------------
# Jaccard over trigram sets is in principle sensitive to disparities in set
# size, so the wider orthographic range could be an artefact of the eight
# New-Testament-only corpora rather than a property of the languages.
CHAPTERS = {'tgl': 977, 'ceb': 977, 'ilo': 977, 'jil': 1104, 'bik': 1104,
            'war': 1104, 'pam': 1104, 'pag': 1104, 'tiu': 209, 'cbk': 209,
            'prf': 209, 'tsg': 209, 'rol': 209, 'msb': 209, 'krj': 209,
            'tao': 209}

langs = list(ortho.index)
pairs = [(langs[i], langs[j]) for i, j in zip(*iu)]
ov, pv = ortho.values[iu], phon.values[iu]

print('\nCorpus-size confound (Spearman rho against each measure):')
for label, x in [
    ('|size difference|', [abs(CHAPTERS[a] - CHAPTERS[b]) for a, b in pairs]),
    ('min(size) of pair', [min(CHAPTERS[a], CHAPTERS[b]) for a, b in pairs]),
]:
    r_o, _ = spearmanr(x, ov)
    r_p, _ = spearmanr(x, pv)
    print(f'  {label:18s} ortho={r_o:+.3f}  phon={r_p:+.3f}')


def stratum(a, b):
    big_a, big_b = CHAPTERS[a] > 900, CHAPTERS[b] > 900
    if big_a and big_b:
        return 'Both complete'
    return 'Both New Testament' if not (big_a or big_b) else 'Mixed'


# The scale difference should survive within strata of equal completeness if it
# is a property of the measures rather than of corpus size. This reproduces the
# stratified table in the paper's Results and Discussion.
cats = np.array([stratum(a, b) for a, b in pairs])
print(f'\n{"Pairs":<20}{"n":>4}{"o.mean":>9}{"o.range":>9}{"p.mean":>9}{"p.range":>9}')
print(f'{"All":<20}{len(ov):>4}{ov.mean():>9.4f}{np.ptp(ov):>9.4f}'
      f'{pv.mean():>9.4f}{np.ptp(pv):>9.4f}')
for c in ['Both complete', 'Both New Testament', 'Mixed']:
    m = cats == c
    print(f'{c:<20}{m.sum():>4}{ov[m].mean():>9.4f}{np.ptp(ov[m]):>9.4f}'
          f'{pv[m].mean():>9.4f}{np.ptp(pv[m]):>9.4f}')
