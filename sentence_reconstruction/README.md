# README

This directory contains a **subset** of reconstructed sentences generated using **trigram** and **quadgram** language models.  

Each reconstructed sentence is defined as a string  

\[
\ell := p_1 \; \mathbf{|} \; p_2 \; \mathbf{|} \; p_3
\]

where each \( p_i \) is an n-gram that satisfies the following **boundary and overlap constraints**:

\[
\begin{aligned}
& p_1[0] = \_ \\
& p_1[-(n-1):] = p_2[:(n-1)] \\
& p_2[-(n-1):] = p_3[:(n-1)] \\
& p_3[-1] = \_
\end{aligned}
\]

Here:
- `[:]` represents slice notation (Python-style) for word sequences.  
- The notation \( p_i[-(n-1):] \) refers to the **last (n−1) words** of \( p_i \), while \( p_i[:(n-1)] \) refers to the **first (n−1) words**.  
- Thus, each consecutive pair \( (p_i, p_{i+1}) \) shares **n−1 overlapping words**.  

The vertical bar (`|`) separates each \( p_i \) in the output for clarity.  
String equality is **case-sensitive**, and underscores (`_`) denote **sentence boundaries**.  

The output files retain `_` and `|` annotations; you may remove these when presenting the cleaned sentences.