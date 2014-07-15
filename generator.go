package ttygif

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/gif"
	"io/ioutil"
	"os"
)

// TtyPlayCaptureProcessor type, which is the implementation of TtyPlayProcessor
type TtyPlayCaptureProcessor struct {
	timestamp TimeVal
	images    []*image.Paletted
	delays    []int
	tempDir   string
	speed     float32
}

// Process captures and append images
func (t *TtyPlayCaptureProcessor) Process(diff TimeVal) (err error) {
	delay := int(float32(diff.Sec*1000000+diff.Usec)/t.speed) / 10000
	if delay == 0 {
		return nil
	}
	t.delays = append(t.delays, delay)
	t.timestamp = t.timestamp.Add(diff)

	// capture and append to images
	tmpFileName := fmt.Sprintf("%03d_%06d", t.timestamp.Sec, t.timestamp.Usec)
	img, err := CaptureImage(t.tempDir, tmpFileName)
	if err != nil {
		return
	}
	paletted := image.NewPaletted(img.Bounds(), palette.WebSafe)
	for x := paletted.Rect.Min.X; x < paletted.Rect.Max.X; x++ {
		for y := paletted.Rect.Min.Y; y < paletted.Rect.Max.Y; y++ {
			paletted.Set(x, y, img.At(x, y))
		}
	}
	if err != nil {
		return
	}
	t.images = append(t.images, paletted)
	return nil
}

// NewTtyPlayCaptureProcessor returns TtyPlayCaptureProcessor instance
func NewTtyPlayCaptureProcessor(speed float32) (t *TtyPlayCaptureProcessor, err error) {
	tempDir, err := ioutil.TempDir("", "ttygif")
	if err != nil {
		return
	}
	return &TtyPlayCaptureProcessor{
		tempDir: tempDir,
		speed:   speed,
	}, nil
}

// RemoveTempDirectory remove tempDir
func (t *TtyPlayCaptureProcessor) RemoveTempDirectory() (err error) {
	if err = os.RemoveAll(t.tempDir); err != nil {
		return
	}
	return nil
}

// GifGenerator type
type GifGenerator struct {
	speed float32
}

// NewGifGenerator returns GifGenerator instance
func NewGifGenerator() *GifGenerator {
	return &GifGenerator{speed: 1.0}
}

// Speed sets the speed
func (g *GifGenerator) Speed(speed float32) {
	g.speed = speed
}

// Generate creates gif animation
func (g *GifGenerator) Generate(inFile string, outFile string) (err error) {
	capturer, err := NewTtyPlayCaptureProcessor(g.speed)
	if err != nil {
		return err
	}
	defer capturer.RemoveTempDirectory()

	player := NewTtyPlayer()
	player.Processor(capturer)
	err = player.Play(inFile)
	if err != nil {
		return
	}

	file, err := os.Create(outFile)
	if err != nil {
		return
	}
	defer file.Close()

	err = gif.EncodeAll(file, &gif.GIF{
		Image:     capturer.images,
		Delay:     capturer.delays,
		LoopCount: -1,
	})
	if err != nil {
		return
	}
	return nil
}
