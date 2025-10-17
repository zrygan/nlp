//! # N-gram Sentence Reconstruction
//!
//! Reconstructs sentences using 3-gram or 4-gram language models.
//! Each reconstructed sentence is made of three overlapping n-grams where
//! the intersection between successive n-grams equals `n – 1` tokens.
//!
//! ## Notes
//! * Sentences start and end with `_`.
//! * Casing and punctuation are preserved exactly.
//! * Reconstruction runs in parallel using Rayon.

use rayon::prelude::*;
use std::{
    collections::HashMap,
    fs::File,
    io::{BufRead, BufReader, BufWriter, Write},
    sync::{Arc, Mutex},
};

/// Represents one entry in the n-gram dataset: the phrase and its frequency.
#[derive(Clone)]
struct NGram {
    phrase: String,
    freq: usize,
}

/// Reads an n-gram file and parses each line into an `NGram`.
///
/// The last token of each line is treated as the frequency count.
///
/// # Arguments
/// * `file_name` – Base name of the file (e.g. `"3gram"` or `"4gram"`).
///
/// # Returns
/// A vector of `(phrase, frequency)` pairs.
fn data_file(file_name: &str) -> Vec<NGram> {
    let n = if file_name.contains('3') { 3 } else { 4 };
    let input_path = format!("data/{file_name}.txt");

    // Ensure output directory exists
    std::fs::create_dir_all("output").expect("Failed to create output directory");

    // Read input files
    let file = File::open(&input_path).expect("Failed to read input file");
    let reader = BufReader::new(file);

    let mut ngrams: Vec<NGram> = Vec::new();
    for line in reader.lines().map_while(Result::ok) {
        let parts: Vec<&str> = line.split_whitespace().collect();
        if parts.len() < n + 1 {
            continue;
        }

        let freq_str = parts.last().unwrap();
        if let Ok(freq) = freq_str.parse::<usize>() {
            let phrase = parts[..parts.len() - 1].join(" ");
            ngrams.push(NGram { phrase, freq });
        }
    }

    ngrams
}

/// Builds a lookup table mapping the prefix of length `n-1` words
/// to all n-grams that share that prefix.
///
/// # Arguments
/// * `ngrams` – Slice of `NGram` entries.
/// * `n` – N-gram size (3 or 4).
fn build_lut(ngrams: &[NGram], n: usize) -> HashMap<String, Vec<NGram>> {
    let mut lut: HashMap<String, Vec<NGram>> = HashMap::new();

    for ng in ngrams {
        let words: Vec<&str> = ng.phrase.split_whitespace().collect();
        if words.len() != n {
            continue;
        }

        let prefix = words[..n - 1].join(" ");
        lut.entry(prefix).or_default().push(ng.clone());
    }

    lut
}

/// Reads an n-gram file from the `input/` directory, reconstructs sentences based on n-gram overlaps,
/// and writes the results to the `output/` directory.
///
/// # Arguments
///
/// * `file_name` - The name of the file (without directory path or extension, e.g. `"3gram"`).
///
/// # Behavior
/// - Determines n (3 or 4) based on the filename.
/// - Reads all n-grams and their frequencies.
/// - Reconstructs sentences where the overlap between consecutive n-grams is `n - 1`.
/// - Writes all reconstructed sentences and their total frequency counts.
fn reconstruct(file_name: &str) {
    let n = if file_name.contains('3') { 3 } else { 4 };

    let output_path = format!("output/{file_name}_output.txt");
    let ngrams = data_file(file_name);
    let lut = build_lut(&ngrams, n);

    // Output file
    let out_file = File::create(&output_path).expect("Failed to create output file");
    let writer = Arc::new(Mutex::new(BufWriter::new(out_file)));

    // Parallel exhaustive search

    ngrams.par_iter().for_each(|ng1| {
        let p1 = &ng1.phrase;
        let f1 = ng1.freq;

        // the starting phrase MUST start with an _
        if !p1.starts_with('_') {
            return;
        }

        let p1_words: Vec<&str> = p1.split_whitespace().collect();
        if p1_words.len() != n {
            return;
        }

        let p1_suffix = p1_words[1..].join(" ");
        if let Some(p2_candidates) = lut.get(&p1_suffix) {
            for ng2 in p2_candidates {
                let p2 = &ng2.phrase;
                let f2 = ng2.freq;

                // this is the middle phrase, it CANNOT have an _ anywhere
                if p2.ends_with('_') {
                    continue;
                }

                let p2_words: Vec<&str> = p2.split_whitespace().collect();
                let p2_suffix = p2_words[1..].join(" ");

                if let Some(p3_candidates) = lut.get(&p2_suffix) {
                    for ng3 in p3_candidates {
                        let p3 = &ng3.phrase;
                        let f3 = ng3.freq;

                        // the ending phrase MUST end with an _
                        if p3.ends_with('_') {
                            // Compute total frequency count
                            let total_freq = f1 + f2 + f3;
                            let result = format!("{p1} | {p2} | {p3} => {total_freq}\n");

                            let mut lock = writer.lock().unwrap();
                            lock.write_all(result.as_bytes()).unwrap();
                        }
                    }
                }
            }
        }
    });

    println!("{file_name} > {output_path}");
}

/// Entry point.
/// Reconstructs both 3-gram and 4-gram sentences.
fn main() {
    reconstruct("3gram");
    reconstruct("4gram");
}
