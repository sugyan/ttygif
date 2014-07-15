package ttygif

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/gif"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// TtyPlayCaptureProcessor type, which is the implementation of TtyPlayProcessor
type TtyPlayCaptureProcessor struct {
	timestamp TimeVal
	captured  []*CapturedImage
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
	imgPath := filepath.Join(t.tempDir, fmt.Sprintf("%03d_%06d", t.timestamp.Sec, t.timestamp.Usec))
	captured, err := CaptureImage(imgPath)
	if err != nil {
		return
	}
	t.captured = append(t.captured, captured)
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

// GetPalettedImages returns paletted images from saved capture images
func (t *TtyPlayCaptureProcessor) GetPalettedImages() (images []*image.Paletted, err error) {
	w := sync.WaitGroup{}
	images = make([]*image.Paletted, len(t.captured))
	for ii, cc := range t.captured {
		w.Add(1)
		go func(i int, captured *CapturedImage) {
			defer w.Done()
			// open image
			var file *os.File
			file, err = os.Open(captured.path)
			if err != nil {
				return
			}
			defer file.Close()
			// decode
			var img image.Image
			img, err = captured.decoder(file)
			if err != nil {
				return
			}
			images[i] = image.NewPaletted(img.Bounds(), palette.WebSafe)
			for x := images[i].Rect.Min.X; x < images[i].Rect.Max.X; x++ {
				for y := images[i].Rect.Min.Y; y < images[i].Rect.Max.Y; y++ {
					images[i].Set(x, y, img.At(x, y))
				}
			}
		}(ii, cc)
	}
	w.Wait()
	return images, nil
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
	images, err := capturer.GetPalettedImages()
	if err != nil {
		return
	}

	file, err := os.Create(outFile)
	if err != nil {
		return
	}
	defer file.Close()

	err = gif.EncodeAll(file, &gif.GIF{
		Image:     images,
		Delay:     capturer.delays,
		LoopCount: -1,
	})
	if err != nil {
		return
	}
	return nil
}
