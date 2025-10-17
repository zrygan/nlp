//! # N-gram Sentence Reconstruction
//!
//! Reconstructs sentences using 3-gram or 4-gram language models.
//! Each reconstructed sentence is made of three overlapping n-grams where
//! the intersection between successive n-grams equals `n – 1` tokens.
//!
//! ## Notes
//! * Sentences start and end with `_`.
//! * Casing and punctuation are preserved exactly.

use rayon::prelude::*;
use std::{
    collections::{HashMap, HashSet},
    fs::File,
    io::{BufRead, BufReader, BufWriter, Write},
    sync::{Arc, Mutex},
};

/// Represents a single n-gram entry consisting of:
/// - the phrase itself (as a space-separated string)
/// - its frequency count
#[derive(Clone)]
struct NGram {
    phrase: String,
    freq: usize,
}

/// Reads an n-gram dataset and parses each line into an `NGram` struct.
/// 
/// Each line is expected to end with a frequency value.
/// For instance: `_ Let's all vote 23`
///
/// # Arguments
/// * `file_name` – Base name of the file (e.g., "3gram" or "4gram").
///
/// # Returns
/// A vector of parsed `NGram` objects.
fn data_file(file_name: &str) -> Vec<NGram> {
    let n = if file_name.contains('3') { 3 } else { 4 };
    let input_path = format!("data/{file_name}.txt");

    // Ensure the output directory exists before writing
    std::fs::create_dir_all("output").unwrap();

    // Open and read the input file line-by-line
    let file = File::open(&input_path).expect("Failed to read input file");
    let reader = BufReader::new(file);

    let mut ngrams = Vec::new();

    for line in reader.lines().map_while(Result::ok) {
        let parts: Vec<&str> = line.split_whitespace().collect();

        // Skip malformed lines (must have n tokens + 1 frequency)
        if parts.len() < n + 1 {
            continue;
        }

        // Parse the last token as frequency
        let freq_str = parts.last().unwrap();
        if let Ok(freq) = freq_str.parse::<usize>() {
            // Join the rest as the n-gram phrase
            let phrase = parts[..parts.len() - 1].join(" ");
            ngrams.push(NGram { phrase, freq });
        }
    }

    ngrams
}

/// Constructs a lookup table (LUT) mapping the prefix of each n-gram (of length n–1)
/// to all n-grams that share that prefix.
///
/// This enables constant-time retrieval of potential continuations
/// when expanding a sentence.
///
/// # Arguments
/// * `ngrams` – A slice of parsed `NGram` entries.
/// * `n` – The size of each n-gram (3 or 4).
fn build_lut(ngrams: &[NGram], n: usize) -> HashMap<String, Vec<NGram>> {
    let mut lut: HashMap<String, Vec<NGram>> = HashMap::new();

    for ng in ngrams {
        let words: Vec<&str> = ng.phrase.split_whitespace().collect();

        // Only include complete n-grams
        if words.len() != n {
            continue;
        }

        // Prefix = first (n–1) tokens
        let prefix = words[..n - 1].join(" ");
        lut.entry(prefix).or_default().push(ng.clone());
    }

    lut
}

/// Expands a starting n-gram into a full reconstructed sentence by
/// chaining overlapping n-grams based on prefix–suffix matches.
/// 
/// Expansion stops when:
/// 1. No valid continuation exists, or
/// 2. The next n-gram ends with `_`.
///
/// Internal underscores are strictly forbidden to ensure only sentence
/// boundaries have `_`.
///
/// # Arguments
/// * `start` – The starting `_` n-gram.
/// * `lut` – Lookup table mapping (n–1)-word prefixes to continuation n-grams.
/// * `n` – The size of n-gram (3 or 4).
///
/// # Returns
/// A reconstructed sentence and its cumulative frequency.
fn expand_sentence(
    start: &NGram,
    lut: &HashMap<String, Vec<NGram>>,
    n: usize,
) -> Option<(String, usize, Vec<NGram>)> {
    // i will just fucking spam clonsed here, BITE ME!
    let mut current = start.phrase.clone();
    let mut total_freq = start.freq;
    let mut used = vec![start.clone()];
    let mut seen = HashSet::new();
    seen.insert(current.clone());

    // i know, dangerous BITE ME
    loop {
        let words: Vec<&str> = current.split_whitespace().collect();
        if words.len() < n - 1 {
            break;
        }

        let suffix = words[words.len() - (n - 1)..].join(" ");
        let candidates = lut.get(&suffix)?;

        // find valid continuation (refer to the constraints in the README)
        if let Some(next_ngram) = candidates.iter().find(|ng| {
            let p = &ng.phrase;
            (!p.contains('_')) || (p.ends_with('_') && !p[..p.len() - 1].contains('_'))
        }) {
            if seen.contains(&next_ngram.phrase) {
                break;
            }
            seen.insert(next_ngram.phrase.clone());

            let next_words: Vec<&str> = next_ngram.phrase.split_whitespace().collect();
            let new_token = next_words.last().unwrap();
            current.push(' ');
            current.push_str(new_token);
            total_freq += next_ngram.freq;
            used.push(next_ngram.clone());

            if next_ngram.phrase.ends_with('_') {
                break;
            }
        } else {
            break;
        }
    }

    Some((current, total_freq, used))
}

/// Reconstructs full sentences from an n-gram dataset.
/// 
/// Each sentence starts from an n-gram beginning with `_`
/// and is expanded forward until an n-gram ending with `_`
/// or no continuation remains.
///
/// Uses Rayon for parallel expansion — each starting `_`
/// n-gram runs independently in parallel.
///
/// # Arguments
/// * `file_name` – Base name of the input file (e.g., "3gram" or "4gram").
fn reconstruct(file_name: &str) {
    let n = if file_name.contains('3') { 3 } else { 4 };
    let ngrams = data_file(file_name);
    let lut = build_lut(&ngrams, n);

    let output_path = format!("output/{file_name}_output.txt");
    let out_file = File::create(&output_path).expect("Failed to create output file");
    let writer = Arc::new(Mutex::new(BufWriter::new(out_file)));

    ngrams
        .par_iter()
        .filter(|ng| ng.phrase.starts_with('_') && !ng.phrase[1..].contains('_'))
        .for_each(|start_ngram| {
            if let Some((sentence, total_freq, used_ngrams)) =
                expand_sentence(start_ngram, &lut, n)
            {
                let mut lock = writer.lock().unwrap();

                // Write final sentence
                writeln!(lock, "{sentence} => {total_freq}").unwrap();

                // Write construction tree (or, what n-grams were used to construct this sentence?)
                // this is soo beautiful btw
                for (i, ng) in used_ngrams.iter().enumerate() {
                    writeln!(lock, "    [{}] {} ({})", i + 1, ng.phrase, ng.freq).unwrap();
                }

                writeln!(lock).unwrap(); // blank line between sentences
            }
        });

    println!("{file_name} > {output_path}");
}

fn main() {
    reconstruct("3gram");
    reconstruct("4gram");
}
