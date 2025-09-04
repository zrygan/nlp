package fourier

import (
	"math"

	"github.com/go-audio/audio"
	"gonum.org/v1/gonum/dsp/fourier"
)

func ConvertBufferToFloat64(buffer *audio.IntBuffer) []float64 {
	ret := make([]float64, len(buffer.Data))

	for i, v := range buffer.Data {
		ret[i] = float64(v) / float64(1<<15)
	}

	return ret
}

func Fourier(buffer *audio.IntBuffer) [][]float64 {
	samples := ConvertBufferToFloat64(buffer)
	size := 1024
	hopSize := size / 4

	fft := fourier.NewFFT(size)
	nf := (len(samples)-size)/hopSize + 1
	output := make([][]float64, nf)

	hanning := make([]float64, size)
	for i := range hanning {
		hanning[i] = 0.5 - 0.5*math.Cos(2*math.Pi*float64(i)/float64(size-1))
	}

	for i := range nf {
		frame := make([]float64, size)
		for j := range size {
			frame[j] = samples[i*hopSize+j] * hanning[j]
		}

		result := fft.Coefficients(nil, frame)
		mags := make([]float64, size/2)

		for j := range size / 2 {
			// split the complex number
			R := real(result[j])
			I := imag(result[j])

			// absolute value of C
			mags[j] = math.Sqrt(R*R + I*I)
		}

		output[i] = mags
	}

	return output
}
