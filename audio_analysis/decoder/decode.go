package decoder

import (
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func DecodeWAV(file *os.File) *audio.IntBuffer {
	decoder := wav.NewDecoder(file)
	if !decoder.IsValidFile() {
		panic("Invalid WAV file")
	}

	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		panic(err)
	}

	return buffer
}
