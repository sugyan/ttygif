package ttygif

import (
	"image"
	"image/color/palette"
	"image/gif"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// GenerateGif creates gif animation
func GenerateGif(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var (
		prevTv TimeVal
		first  = true
		delays []int
		images []*image.Paletted
	)
	tempDir, err := ioutil.TempDir("", "ttygif")
	if err != nil {
		log.Fatal(err)
	}

	reader := NewTtyReader(file)
	for {
		data, err := reader.ReadData()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		// play
		if first {
			first = false
		} else {
			diff := data.TimeVal.Subtract(prevTv)
			delays = append(delays, int((diff.Sec*1000000+diff.Usec)/10000))
		}
		prevTv = data.TimeVal
		print(string(*data.Buffer))

		imageFile, err := CaptureImage(tempDir, data)
		if err != nil {
			log.Fatal(err)
		}

		// read file
		file, err := os.Open(imageFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		img, err := png.Decode(file)
		// add image
		// TODO: concurrently?
		p := image.NewPaletted(img.Bounds(), palette.WebSafe)
		for x := p.Rect.Min.X; x < p.Rect.Max.X; x++ {
			for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
				p.Set(x, y, img.At(x, y))
			}
		}
		images = append(images, p)
	}
	// remove temp images
	if err := os.RemoveAll(tempDir); err != nil {
		log.Fatal(err)
	}

	// create animated GIF
	if len(images) > len(delays) {
		delays = append(delays, 200)
	}

	outFile, err := os.Create("out.gif")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	if err = gif.EncodeAll(outFile, &gif.GIF{
		Image:     images,
		Delay:     delays,
		LoopCount: -1,
	}); err != nil {
		log.Fatal(err)
	}
}
