//go:build !(js || sdl || azul3d || ebiten)

package nes

type DummyAudio struct {
	input      chan int16
	frequency  int
	sampleSize int
}

func NewAudio(frequency int, sampleSize int) (audio *DummyAudio, err error) {
	audio = &DummyAudio{
		input:      make(chan int16),
		frequency:  frequency,
		sampleSize: sampleSize,
	}
	return
}

func (audio *DummyAudio) Input() chan int16 {
	return audio.input
}

func (audio *DummyAudio) Run() {}

func (audio *DummyAudio) TogglePaused() {}

func (audio *DummyAudio) SetSpeed(speed float32) {}

func (audio *DummyAudio) Close() {}
