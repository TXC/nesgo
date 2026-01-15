//go:build ebiten && !(js || sdl || azul3d)

package nes

import (
	"encoding/binary"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// EbitenAudioReader implements io.Reader to provide continuous audio streaming
type EbitenAudioReader struct {
	input      chan int16
	buffer     []byte
	bufferSize int
	pos        int
}

func NewEbitenAudioReader(sampleSize int) *EbitenAudioReader {
	return &EbitenAudioReader{
		input:      make(chan int16, sampleSize*2),
		buffer:     make([]byte, 0, sampleSize*4),
		bufferSize: sampleSize * 2, // 2 bytes per int16 sample
	}
}

func (r *EbitenAudioReader) Read(p []byte) (n int, err error) {
	// Fill buffer with samples from channel
	for n < len(p) {
		// If we've consumed all buffered data, get new samples
		if r.pos >= len(r.buffer) {
			// Read samples from channel and convert to bytes
			r.buffer = r.buffer[:0]
			r.pos = 0

			// Fill buffer up to bufferSize
			for len(r.buffer) < r.bufferSize && len(r.buffer) < len(p)*2 {
				select {
				case sample := <-r.input:
					// Convert int16 to bytes (little-endian)
					b := make([]byte, 2)
					binary.LittleEndian.PutUint16(b, uint16(sample))
					r.buffer = append(r.buffer, b...)
				default:
					// No more samples available right now
					if len(r.buffer) == 0 {
						// Return silence if no data
						for i := n; i < len(p); i++ {
							p[i] = 0
						}
						return len(p), nil
					}
					goto readBuffer
				}
			}
		}

	readBuffer:
		// Copy from buffer to output
		copied := copy(p[n:], r.buffer[r.pos:])
		r.pos += copied
		n += copied
	}

	return n, nil
}

type EbitenAudio struct {
	sampleSize int
	input      io.Reader
	ctx        *audio.Context
	player     *audio.Player
	reader     *EbitenAudioReader
}

func NewAudio(frequency int, sampleSize int) (*EbitenAudio, error) {
	ctx := audio.NewContext(frequency)

	// Create the audio reader that will provide continuous streaming
	reader := NewEbitenAudioReader(sampleSize)

	// Create player from the reader
	player, err := ctx.NewPlayer(reader)
	if err != nil {
		return nil, err
	}

	return &EbitenAudio{
		input:      reader,
		sampleSize: sampleSize,
		ctx:        ctx,
		player:     player,
		reader:     reader,
	}, nil
}

func (ea *EbitenAudio) Input() chan int16 {
	return ea.reader.input
}

func (ea *EbitenAudio) Run() {
	// The player continuously reads from the reader,
	// which pulls from the reader.input channel.
	// This method can be empty or handle coordination if needed.
	// The actual audio streaming happens via the reader's Read method.
	if !ea.player.IsPlaying() {
		ea.player.Play()
	}
}

func (ea *EbitenAudio) TogglePaused() {
	if ea.player.IsPlaying() {
		ea.player.Pause()
	} else {
		ea.player.Play()
	}
}

func (ea *EbitenAudio) SetSpeed(speed float32) {
	// Ebiten v2 doesn't have built-in playback speed control
	// Would need to resample audio data to implement this
}

func (ea *EbitenAudio) Close() {
	if ea.player != nil {
		ea.player.Pause()
		_ = ea.player.Close()
	}
}
