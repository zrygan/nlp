# # `zrygan/nlp/sentence_reconstruction`

This directory contains reconstructed sentences generated using **trigram** and **quadgram** language models.  

Each reconstructed sentence is defined as a string  

$$
\ell := p_1 \; p_2 \; \dots \; p_x
$$

where each $p_i$ is an n-gram that satisfies the following **boundary and overlap constraints**:

$$
p_1[0] = \textunderscore
$$
$$
p_k[-(n - 1):] = p_{k + 1}[:(n - 1)]
$$
$$
p_x[-1] = \textunderscore
$$

Here:
- $1 < k < k + 1 < x$
- `[:]` represents slice notation for word sequences.  
- The notation $p_i[-(n-1):]$ refers to the **last (n−1) words** of $p_i$, while $p_i[:(n-1)]$ refers to the **first (n−1) words**.
- Thus, each consecutive pair $(p_i, p_{i+1})$ shares **n−1 overlapping words**.  

## $n$-gram Models
The trigram and quadgram models may be accessed in the [`data/`](data/) directory.

## Reconstructed Sentences
The reconstructed sentences may be accessed in the [`output`](output/) directory.
