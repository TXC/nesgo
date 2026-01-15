//go:build ebiten && !(js || sdl || azul3d)

package nes

import (
	"image/color"
	"log"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	nativeScreenWidth  = 256
	nativeScreenHeight = 240
	scaleFactor        = 2
)

var EbitenPalette = []color.RGBA{
	{0x66, 0x66, 0x66, 0xff},
	{0x00, 0x2A, 0x88, 0xff},
	{0x14, 0x12, 0xA7, 0xff},
	{0x3B, 0x00, 0xA4, 0xff},
	{0x5C, 0x00, 0x7E, 0xff},
	{0x6E, 0x00, 0x40, 0xff},
	{0x6C, 0x06, 0x00, 0xff},
	{0x56, 0x1D, 0x00, 0xff},
	{0x33, 0x35, 0x00, 0xff},
	{0x0B, 0x48, 0x00, 0xff},
	{0x00, 0x52, 0x00, 0xff},
	{0x00, 0x4F, 0x08, 0xff},
	{0x00, 0x40, 0x4D, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0xAD, 0xAD, 0xAD, 0xff},
	{0x15, 0x5F, 0xD9, 0xff},
	{0x42, 0x40, 0xFF, 0xff},
	{0x75, 0x27, 0xFE, 0xff},
	{0xA0, 0x1A, 0xCC, 0xff},
	{0xB7, 0x1E, 0x7B, 0xff},
	{0xB5, 0x31, 0x20, 0xff},
	{0x99, 0x4E, 0x00, 0xff},
	{0x6B, 0x6D, 0x00, 0xff},
	{0x38, 0x87, 0x00, 0xff},
	{0x0C, 0x93, 0x00, 0xff},
	{0x00, 0x8F, 0x32, 0xff},
	{0x00, 0x7C, 0x8D, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0xFF, 0xFE, 0xFF, 0xff},
	{0x64, 0xB0, 0xFF, 0xff},
	{0x92, 0x90, 0xFF, 0xff},
	{0xC6, 0x76, 0xFF, 0xff},
	{0xF3, 0x6A, 0xFF, 0xff},
	{0xFE, 0x6E, 0xCC, 0xff},
	{0xFE, 0x81, 0x70, 0xff},
	{0xEA, 0x9E, 0x22, 0xff},
	{0xBC, 0xBE, 0x00, 0xff},
	{0x88, 0xD8, 0x00, 0xff},
	{0x5C, 0xE4, 0x30, 0xff},
	{0x45, 0xE0, 0x82, 0xff},
	{0x48, 0xCD, 0xDE, 0xff},
	{0x4F, 0x4F, 0x4F, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0xFF, 0xFE, 0xFF, 0xff},
	{0xC0, 0xDF, 0xFF, 0xff},
	{0xD3, 0xD2, 0xFF, 0xff},
	{0xE8, 0xC8, 0xFF, 0xff},
	{0xFB, 0xC2, 0xFF, 0xff},
	{0xFE, 0xC4, 0xEA, 0xff},
	{0xFE, 0xCC, 0xC5, 0xff},
	{0xF7, 0xD8, 0xA5, 0xff},
	{0xE4, 0xE5, 0x94, 0xff},
	{0xCF, 0xEF, 0x96, 0xff},
	{0xBD, 0xF4, 0xAB, 0xff},
	{0xB3, 0xF3, 0xCC, 0xff},
	{0xB5, 0xEB, 0xF2, 0xff},
	{0xB8, 0xB8, 0xB8, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
}

type EbitenEvent interface {
	Event
}

type EbitenKeyEvent struct {
	EbitenEvent
	Key     ebiten.Key
	Pressed bool
}

func (e EbitenKeyEvent) Flag() uint       { return EvGlobal }
func (e EbitenKeyEvent) String() string   { return "EbitenKeyEvent" }
func (e EbitenKeyEvent) Process(nes *NES) {}

type EbitenQuitEvent struct {
	EbitenEvent
}

func (e EbitenQuitEvent) Flag() uint       { return EvGlobal }
func (e EbitenQuitEvent) String() string   { return "EbitenQuitEvent" }
func (e EbitenQuitEvent) Process(nes *NES) { nes.state = Quitting }

type EbitenKeyCallback map[ebiten.Key]func(any) Event

type EbitenVideo struct {
	input         chan []uint8
	width, height int
	scale         int
	x_offset      int
	y_offset      int
	palette       []color.RGBA
	events        chan Event
	framePool     *sync.Pool
	overscan      bool
	caption       string
	fps           float64
	game          *EbitenGame
	running       bool
	screenPool    *sync.Pool
}

func NewVideo(caption string, events chan Event, framePool *sync.Pool, fps float64) (video *EbitenVideo, err error) {

	screenPool := &sync.Pool{New: func() any { return make([]ebiten.Image, 0) }}

	video = &EbitenVideo{
		input:      make(chan []uint8, 128),
		palette:    EbitenPalette,
		events:     events,
		framePool:  framePool,
		overscan:   true,
		caption:    caption,
		fps:        fps,
		width:      nativeScreenWidth * scaleFactor,
		height:     nativeScreenHeight * scaleFactor,
		x_offset:   0,
		y_offset:   0,
		scale:      scaleFactor,
		running:    false,
		screenPool: screenPool,
	}

	ebiten.SetWindowSize(video.width, video.height)
	ebiten.SetWindowClosingHandled(true)
	video.SetCaption(caption)
	game := NewGame(video)
	video.game = game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	return
}

func (video *EbitenVideo) Input() chan []uint8 {
	return video.input
}
func (video *EbitenVideo) Events() chan Event {
	return video.events
}
func (video *EbitenVideo) SetScaleEvent(scale int) {
	width := nativeScreenWidth * scale
	height := nativeScreenHeight * scale
	video.ResizeEvent(width, height)
}
func (video *EbitenVideo) ResizeEvent(width, height int) {
	ebiten.SetWindowSize(width, height)
	video.Reshape(width, height)
}
func (video *EbitenVideo) Reshape(width int, height int) {
	x_offset := 0
	y_offset := 0

	r := ((float64)(height)) / ((float64)(width))

	if r > 0.9375 { // Height taller than ratio
		h := (int)(math.Floor((float64)(0.9375 * (float64)(width))))
		y_offset = (height - h) / 2
		height = h
	} else if r < 0.9375 { // Width wider
		w := (int)(math.Floor((float64)((256.0 / 240.0) * (float64)(height))))
		x_offset = (width - w) / 2
		width = w
	}

	video.width = width
	video.height = height
	video.x_offset = x_offset
	video.y_offset = y_offset
}
func (video *EbitenVideo) frameWidth() int {
	if video.overscan {
		return nativeScreenWidth - 16
	}
	return nativeScreenWidth
}
func (video *EbitenVideo) frameHeight() int {
	if video.overscan {
		return nativeScreenHeight - 16
	}

	return nativeScreenHeight
}

func (video *EbitenVideo) Run() {
	video.running = true
	frame := make([]uint32, 0xf000)

	for video.running {
		select {
		case ev := <-video.game.events:
			var event Event
			switch e := ev.(type) {
			case EbitenQuitEvent:
				video.running = false
				event = &QuitEvent{}
			case EbitenKeyEvent:
				switch e.Key {
				case ebiten.KeyEscape:
					video.running = false
					event = &QuitEvent{}
				case ebiten.KeyGraveAccent:
					video.overscan = !video.overscan
				case ebiten.KeyDigit1:
					video.ResizeEvent(256, 240) // Native 256x240
				case ebiten.KeyDigit2:
					video.ResizeEvent(512, 480) // 2x Scale 512x480
				case ebiten.KeyDigit3:
					video.ResizeEvent(768, 720) // 3x Scale 768x720
				case ebiten.KeyDigit4:
					video.ResizeEvent(1024, 960) // 4x Scale 1024x960
				case ebiten.KeyDigit5:
					video.ResizeEvent(1280, 1200) // 5x Scale 1280x1200
				case ebiten.KeyDigit6:
					video.ResizeEvent(1536, 1440) // 6x Scale 1536x1440
				case ebiten.KeyP:
					event = &PauseEvent{}
				//case ebiten.KeyN:
				//	event = &FrameStepEvent{}
				case ebiten.KeyQ:
					video.running = false
					event = &QuitEvent{}
				case ebiten.KeyL:
					event = &SavePatternTablesEvent{}
				case ebiten.KeyR:
					event = &ResetEvent{}
				case ebiten.KeyS:
					event = &RecordEvent{}
				case ebiten.KeyNumpadAdd:
					event = &AudioRecordEvent{}
				case ebiten.KeyNumpadSubtract:
					event = &AudioStopEvent{}
				case ebiten.KeyO:
					event = &CPUDecodeEvent{}
				case ebiten.KeyI:
					event = &PPUDecodeEvent{}
				case ebiten.Key9:
					event = &ShowBackgroundEvent{}
				case ebiten.Key0:
					event = &ShowSpritesEvent{}
				case ebiten.KeyF1:
					event = &SaveStateEvent{}
				case ebiten.KeyF5:
					event = &LoadStateEvent{}
				case ebiten.KeyF8:
					event = &FPSEvent{2.}
				case ebiten.KeyF9:
					event = &FPSEvent{1.}
				case ebiten.KeyF10:
					event = &FPSEvent{.75}
				case ebiten.KeyF11:
					event = &FPSEvent{.50}
				case ebiten.KeyF12:
					event = &FPSEvent{.25}
				case ebiten.KeyNumpad0:
					event = &MuteEvent{}
				case ebiten.KeyNumpad1:
					event = &MutePulse1Event{}
				case ebiten.KeyNumpad2:
					event = &MutePulse2Event{}
				case ebiten.KeyNumpad3:
					event = &MuteTriangleEvent{}
				case ebiten.KeyNumpad4:
					event = &MuteNoiseEvent{}
				case ebiten.KeyNumpad5:
					event = &MuteDMCEvent{}
				}

				if event == nil && video.running {
					event = &ControllerEvent{
						Button: ebitenKeyToNesButton(ev),
						Down:   ev.(EbitenKeyEvent).Pressed,
					}
				}

				if event != nil {
					go func() { video.events <- event }()
				}
			}
		case fb := <-video.input:
			index := 0
			x, y := 0, 0

			for _, c := range fb {
				if pixelInFrame(x, y, video.overscan) {
					p := video.palette[c]
					frame[index] = (uint32(p.R) << 24) | (uint32(p.G) << 16) | (uint32(p.B) << 8) | uint32(p.A)
					index++
				}

				switch x {
				case 255:
					x = 0
					y++
				default:
					x++
				}
			}
			video.framePool.Put(fb)

			image := ebiten.NewImage(video.frameWidth(), video.frameHeight())
			image.WritePixels(fb)

			video.screenPool.Put(image)
		}
	}
}

func ebitenKeyToNesButton(ev any) Button {
	if k, ok := ev.(EbitenKeyEvent); ok {
		// Controller buttons
		switch k.Key {
		case ebiten.KeyZ: // A
			return A
		case ebiten.KeyX: // B
			return B
		case ebiten.KeyShiftRight: // Select
			return Select
		case ebiten.KeyEnter: // Start
			return Start
		case ebiten.KeyArrowUp: // Up
			return Up
		case ebiten.KeyArrowDown: // Down
			return Down
		case ebiten.KeyArrowLeft: // Left
			return Left
		case ebiten.KeyArrowRight: // Right
			return Right
		}
	}

	return One
}

func (video *EbitenVideo) SetCaption(caption string) {
	video.caption = caption
	ebiten.SetWindowTitle("nesgo - " + video.caption)
}

type EbitenGame struct {
	video                *EbitenVideo
	keys                 []ebiten.Key
	events               chan EbitenEvent
	windowClosingHandled bool
}

func NewGame(video *EbitenVideo) *EbitenGame {
	return &EbitenGame{
		video:  video,
		events: make(chan EbitenEvent),
		keys:   []ebiten.Key{},
	}
}

func (g *EbitenGame) repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

// Update updates a game by one tick. The given argument represents a screen
// image.
//
// Update updates only the game logic and Draw draws the screen.
//
// You can assume that Update is always called TPS-times per second (60 by
// default), and you can assume that the time delta between two Updates is
// always 1 / TPS [s] (1/60[s] by default). As Ebitengine already adjusts the
// number of Update calls, you don't have to measure time deltas in Update by
// e.g. OS timers.
//
// An actual TPS is available by ActualTPS(), and the result might slightly
// differ from your expected TPS, but still, your game logic should stick to
// the fixed time delta and should not rely on ActualTPS() value.
// This API is for just measurement and/or debugging. In the long run, the
// number of Update calls should be adjusted based on the set TPS on average.
//
// An actual time delta between two Updates might be bigger than expected. In
// this case, your game's [Update] or [Draw] takes longer than they should.
// In this case, there is nothing other than optimizing your game
// implementation.
//
// In the first frame, it is ensured that [Update] is called at least once
// before [Draw]. You can use [Update] to initialize the game state.
//
// After the first frame, [Update] might not be called or might be called once
// or more for one frame. The frequency is determined by the current TPS
// (tick-per-second).
//
// If the error returned is nil, game execution proceeds normally.
// If the error returned is Termination, game execution halts, but does not
// return an error from RunGame.
// If the error returned is any other non-nil value, game execution halts and
// the error is returned from RunGame.
func (g *EbitenGame) Update() error {
	if ebiten.IsWindowBeingClosed() {
		g.windowClosingHandled = true
	}

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	for _, key := range g.keys {
		switch key {
		case ebiten.KeyEscape, ebiten.KeyGraveAccent:
			if g.windowClosingHandled && !g.video.running {
				return ebiten.Termination
			}
		default:
			g.video.events <- &EbitenKeyEvent{
				Key:     key,
				Pressed: g.repeatingKeyPressed(key),
			}
		}
	}

	return nil
}

// Draw draws the game screen by one frame.
//
// The provided argument represents a screen image. The updated content is
// adopted as the game screen.
//
// The frequency of [Draw] calls depends on the user's environment, especially
// the monitor's refresh rate.
// For portability, you should not put your game logic in Draw in general.
func (g *EbitenGame) Draw(screen *ebiten.Image) {
	// Drawing is handled in the EbitenVideo's Run method.

	geo := ebiten.GeoM{}
	geo.Scale(float64(g.video.width)/float64(nativeScreenWidth), float64(g.video.height)/float64(nativeScreenHeight))
	op := &ebiten.DrawImageOptions{
		GeoM: geo,
	}
	screen.DrawImage(g.video.screenPool.Get().(*ebiten.Image), op)
}

// Layout accepts a native outside size in device-independent pixels and
// returns the game's logical screen size in pixels. The logical size is used
// for 1) the screen size given at Draw and 2) calculation of the scale from
// the screen to the final screen size.
//
// On desktops, the outside is a window or a monitor (fullscreen mode). On
// browsers, the outside is a body element. On mobiles, the outside is the
// view's size.
//
// Even though the outside size and the screen size differ, the rendering
// scale is automatically adjusted to fit with the outside dimensions.
//
// Layout is called almost every frame.
//
// It is ensured that Layout is invoked before Update is called in the first
// frame.
//
// If Layout returns non-positive numbers, the caller may panic.
//
// You can return a fixed screen size if desired, or you can also return a
// calculated screen size adjusted with the given outside size.
//
// If the game implements the interface LayoutFer, Layout is never called and
// LayoutF is called instead.
func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.video.width, g.video.height
}
