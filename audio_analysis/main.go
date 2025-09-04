package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zrygan.nlp/audio_analysis/decoder"
	"github.com/zrygan.nlp/audio_analysis/fourier"
	"github.com/zrygan.nlp/audio_analysis/spectrogram"
)

func audioAnalysis(f *os.File, fn string, emotion string) {
	spectrogram.MakeSpectrogram(
		fourier.Fourier(
			decoder.DecodeWAV(f),
		),
		fn,
		emotion,
	)
}

func main() {
	dir, err := os.ReadDir("data")
	if err != nil {
		panic(err)
	}

	for _, entry := range dir {
		if entry.IsDir() {
			fmt.Println("📁 ", entry.Name())

			emotionDir, err := os.ReadDir(filepath.Join("data", entry.Name()))
			if err != nil {
				panic(err)
			}

			for _, audioFile := range emotionDir {
				fn := audioFile.Name()
				fmt.Printf("\t🔉 %s\n", fn)

				af, err := os.Open(filepath.Join("data", entry.Name(), fn))
				if err != nil {
					panic(err)
				}

				audioAnalysis(af, fn, entry.Name())

				af.Close()
			}
		}

		fmt.Println()
	}
}
