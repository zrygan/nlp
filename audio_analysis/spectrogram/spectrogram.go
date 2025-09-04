package spectrogram

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
)

func color(val float64) (r, g, b float64) {
	// val expected in [0, 1]
	if val < 0 {
		val = 0
	}
	if val > 1 {
		val = 1
	}

	r = math.Min(math.Max(1.5-math.Abs(4*val-3), 0), 1)
	g = math.Min(math.Max(1.5-math.Abs(4*val-2), 0), 1)
	b = math.Min(math.Max(1.5-math.Abs(4*val-1), 0), 1)
	return
}

func MakeSpectrogram(fftOutput [][]float64, fileName string, emotionCode string) error {
	width := len(fftOutput)
	if width == 0 {
		return fmt.Errorf("fftOutput has zero width")
	}
	height := len(fftOutput[0])
	if height == 0 {
		return fmt.Errorf("fftOutput has zero height")
	}

	dirPath := filepath.Join("data", emotionCode)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	window := gg.NewContext(width, height)

	// Compute min/max dB
	minDB, maxDB := math.MaxFloat64, -math.MaxFloat64
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			v := 20 * math.Log10(fftOutput[i][j]+1e-6)
			if v < minDB {
				minDB = v
			}
			if v > maxDB {
				maxDB = v
			}
		}
	}

	for x := range width {
		for y := range height {
			valDB := 20 * math.Log10(fftOutput[x][y]+1e-6)
			val := (valDB - minDB) / (maxDB - minDB)
			val = math.Pow(val, 0.3)

			r, g, b := color(val)
			window.SetRGB(r, g, b)
			window.SetPixel(x, height-y-1)
		}
	}
	outFile := filepath.Join(dirPath, strings.TrimSuffix(filepath.Base(fileName), ".wav")+".png")
	if err := window.SavePNG(outFile); err != nil {
		return fmt.Errorf("failed to save PNG: %v", err)
	}

	return nil
}
