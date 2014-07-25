package main

import (
	"bytes"
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"flag"
	"fmt"
	"github.com/sugyan/ttygif"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io"
	"io/ioutil"
	"j4k.co/terminal"
	"os"
	"path/filepath"
)

// colors
var (
	TerminalColorRed          = color.RGBA{0xff, 0x00, 0x00, 0xff}
	TerminalColorGreen        = color.RGBA{0x00, 0xff, 0x00, 0xff}
	TerminalColorBlue         = color.RGBA{0x00, 0x00, 0xff, 0xff}
	TerminalColorYellow       = color.RGBA{0xff, 0xff, 0x00, 0xff}
	TerminalColorMagenta      = color.RGBA{0xff, 0x00, 0xff, 0xff}
	TerminalColorCyan         = color.RGBA{0x00, 0xff, 0xff, 0xff}
	TerminalColorLightGrey    = color.RGBA{0xcc, 0xcc, 0xcc, 0xff}
	TerminalColorDarkGrey     = color.RGBA{0x99, 0x99, 0x99, 0xff}
	TerminalColorLightRed     = color.RGBA{0xff, 0x33, 0x33, 0xff}
	TerminalColorLightGreen   = color.RGBA{0x33, 0xff, 0x33, 0xff}
	TerminalColorLightYellow  = color.RGBA{0xff, 0xff, 0x33, 0xff}
	TerminalColorLightBlue    = color.RGBA{0x33, 0x33, 0xff, 0xff}
	TerminalColorLightMagenta = color.RGBA{0xff, 0x33, 0xff, 0xff}
	TerminalColorLightCyan    = color.RGBA{0x33, 0xff, 0xff, 0xff}
	TerminalColorWhite        = color.RGBA{0xff, 0xff, 0xff, 0xff}
	TerminalColorGray         = color.RGBA{0x7f, 0x7f, 0x7f, 0x7f}
)

var font *truetype.Font

const fontSize = 18

func init() {
	fontData, err := Asset("font/Anonymous Pro Minus.ttf")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	font, err = freetype.ParseFont(fontData)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	input := flag.String("in", "ttyrecord", "input ttyrec file")
	output := flag.String("out", "tty.gif", "output gif file")
	speed := flag.Float64("s", 1.0, "play speed")
	flag.Parse()

	err := generateGIF(*input, *output, *speed)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	absPath, err := filepath.Abs(*output)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%s created!\n", absPath)
}

func generateGIF(input string, output string, speed float64) (err error) {
	// input
	inFile, err := os.Open(input)
	if err != nil {
		return
	}
	defer inFile.Close()

	// virtual terminal
	var state = terminal.State{}
	vt, err := terminal.Create(&state, ioutil.NopCloser(bytes.NewBuffer([]byte{})))
	if err != nil {
		return
	}
	defer vt.Close()

	// read ttyrecord
	reader := ttygif.NewTtyReader(inFile)
	var (
		first  = true
		prevTv ttygif.TimeVal
		images []*image.Paletted
		delays []int
	)
	for {
		var data *ttygif.TtyData
		data, err = reader.ReadData()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return
			}
		}
		var diff ttygif.TimeVal
		if first {
			first = false
		} else {
			diff = data.TimeVal.Subtract(prevTv)
		}
		prevTv = data.TimeVal

		// calc delay and capture
		delay := int(float64(diff.Sec*1000000+diff.Usec)/speed) / 10000
		if delay > 0 {
			var img *image.Paletted
			img, err = capture(&state)
			if err != nil {
				return
			}
			images = append(images, img)
			delays = append(delays, delay)
		}
		// write to vt
		_, err = vt.Write(*data.Buffer)
		if err != nil {
			return
		}
	}

	outFile, err := os.Create(output)
	err = gif.EncodeAll(outFile, &gif.GIF{
		Image: images,
		Delay: delays,
	})
	if err != nil {
		return
	}
	return nil
}

func capture(state *terminal.State) (paletted *image.Paletted, err error) {
	fb := font.Bounds(fontSize)

	paletted = image.NewPaletted(image.Rect(0, 0, 80*int(fb.XMax-fb.XMin)+10, 24*int(fb.YMax-fb.YMin)), palette.WebSafe)
	draw.Draw(paletted, paletted.Bounds(), image.Black, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetDst(paletted)
	c.SetClip(paletted.Bounds())
	for row := 0; row < 24; row++ {
		for col := 0; col < 80; col++ {
			ch, fg, bg := state.Cell(col, row)
			if bg != terminal.DefaultBG {
				var uniform *image.Uniform
				switch bg {
				case terminal.Red:
					uniform = image.NewUniform(TerminalColorRed)
				case terminal.Green:
					uniform = image.NewUniform(TerminalColorGreen)
				case terminal.Blue:
					uniform = image.NewUniform(TerminalColorBlue)
				case terminal.Yellow:
					uniform = image.NewUniform(TerminalColorYellow)
				case terminal.Magenta:
					uniform = image.NewUniform(TerminalColorMagenta)
				case terminal.Cyan:
					uniform = image.NewUniform(TerminalColorCyan)
				case terminal.LightGrey:
					uniform = image.NewUniform(TerminalColorLightGrey)
				case terminal.DarkGrey:
					uniform = image.NewUniform(TerminalColorDarkGrey)
				case terminal.LightRed:
					uniform = image.NewUniform(TerminalColorLightRed)
				case terminal.LightGreen:
					uniform = image.NewUniform(TerminalColorLightGreen)
				case terminal.LightYellow:
					uniform = image.NewUniform(TerminalColorLightYellow)
				case terminal.LightBlue:
					uniform = image.NewUniform(TerminalColorLightBlue)
				case terminal.LightMagenta:
					uniform = image.NewUniform(TerminalColorLightMagenta)
				case terminal.LightCyan:
					uniform = image.NewUniform(TerminalColorLightCyan)
				case terminal.White:
					uniform = image.NewUniform(TerminalColorWhite)
				case terminal.DefaultFG:
					uniform = image.White
				default:
					uniform = image.Transparent
				}
				draw.Draw(paletted, image.Rect(5+col*int(fb.XMax-fb.XMin), row*int(fb.YMax-fb.YMin)-int(fb.YMin), 5+(col+1)*int(fb.XMax-fb.XMin), (row+1)*int(fb.YMax-fb.YMin)-int(fb.YMin)), uniform, image.ZP, draw.Src)
			}
			switch fg {
			case terminal.Red:
				c.SetSrc(image.NewUniform(TerminalColorRed))
			case terminal.Green:
				c.SetSrc(image.NewUniform(TerminalColorGreen))
			case terminal.Blue:
				c.SetSrc(image.NewUniform(TerminalColorBlue))
			case terminal.Yellow:
				c.SetSrc(image.NewUniform(TerminalColorYellow))
			case terminal.Magenta:
				c.SetSrc(image.NewUniform(TerminalColorMagenta))
			case terminal.Cyan:
				c.SetSrc(image.NewUniform(TerminalColorCyan))
			case terminal.LightGrey:
				c.SetSrc(image.NewUniform(TerminalColorLightGrey))
			case terminal.DarkGrey:
				c.SetSrc(image.NewUniform(TerminalColorDarkGrey))
			case terminal.LightRed:
				c.SetSrc(image.NewUniform(TerminalColorLightRed))
			case terminal.LightGreen:
				c.SetSrc(image.NewUniform(TerminalColorLightGreen))
			case terminal.LightYellow:
				c.SetSrc(image.NewUniform(TerminalColorLightYellow))
			case terminal.LightBlue:
				c.SetSrc(image.NewUniform(TerminalColorLightBlue))
			case terminal.LightMagenta:
				c.SetSrc(image.NewUniform(TerminalColorLightMagenta))
			case terminal.LightCyan:
				c.SetSrc(image.NewUniform(TerminalColorLightCyan))
			case terminal.White:
				c.SetSrc(image.NewUniform(TerminalColorWhite))
			case terminal.DefaultFG:
				c.SetSrc(image.White)
			case terminal.DefaultBG:
				c.SetSrc(image.Black)
			default:
				// TODO??
				c.SetSrc(image.NewUniform(TerminalColorGray))
			}
			str := string(ch)
			_, err = c.DrawString(str, freetype.Pt(5+col*int(fb.XMax-fb.XMin), (row+1)*int(fb.YMax-fb.YMin)))
			if err != nil {
				return
			}
		}
	}

	return paletted, nil
}
