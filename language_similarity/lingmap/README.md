# `zrygan/nlp/language_similarity/lingmap`

Computing the similarity of languages using geographic and historical data.

This project is written in Python using uv, geopandas, matplotlib, and Jupyter Notebook.

##

The main source code is `main.ipynb` containing the class definitions and plot styles for
the maps.

### Shape Files

This project uses the **Philippines administrative level 0-4 boundaries (COD-AB) dataset**.
In particular, the administrative boundaries (level 1) by the Philippine Statistics Authority
and National Mapping and Resource Information Authority. Version 20231106, released on
November 6, 2023.

For brevity, the source code and the repository refers to this as `psa_namria_sf`.