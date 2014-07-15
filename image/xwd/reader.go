package xwd

import (
	"encoding/binary"
	"image"
	"image/color"
	"image/color/palette"
	"io"
)

// XWDFileHeader type
type XWDFileHeader struct {
	HeaderSize        uint32
	FileVersion       uint32
	PixmapFormat      uint32
	PixmapDepth       uint32
	PixmapWidth       uint32
	PixmapHeight      uint32
	XOffset           uint32
	ByteOrder         uint32
	BitmapUnit        uint32
	BitmapBitOrder    uint32
	BitmapPad         uint32
	BitsPerPixel      uint32
	BytesPerLine      uint32
	VisualClass       uint32
	RedMask           uint32
	GreenMask         uint32
	BlueMask          uint32
	BitsPerRgb        uint32
	NumberOfColors    uint32
	ColorMapEntries   uint32
	WindowWidth       uint32
	WindowHeight      uint32
	WindowX           uint32
	WindowY           uint32
	WindowBorderWidth uint32
}

// XWDColorMap type
type XWDColorMap struct {
	EntryNumber uint32
	Red         uint16
	Green       uint16
	Blue        uint16
	Flags       uint8
	Padding     uint8
}

// Decode reads a XWD image from r and returns it as an image.Image.
func Decode(r io.Reader) (img image.Image, err error) {
	buf := make([]byte, 100)
	_, err = r.Read(buf)
	if err != nil {
		return
	}
	header := XWDFileHeader{
		HeaderSize:        binary.BigEndian.Uint32(buf[0:4]),
		FileVersion:       binary.BigEndian.Uint32(buf[4:8]),
		PixmapFormat:      binary.BigEndian.Uint32(buf[8:12]),
		PixmapDepth:       binary.BigEndian.Uint32(buf[12:16]),
		PixmapWidth:       binary.BigEndian.Uint32(buf[16:20]),
		PixmapHeight:      binary.BigEndian.Uint32(buf[20:24]),
		XOffset:           binary.BigEndian.Uint32(buf[24:28]),
		ByteOrder:         binary.BigEndian.Uint32(buf[28:32]),
		BitmapUnit:        binary.BigEndian.Uint32(buf[32:36]),
		BitmapBitOrder:    binary.BigEndian.Uint32(buf[36:40]),
		BitmapPad:         binary.BigEndian.Uint32(buf[40:44]),
		BitsPerPixel:      binary.BigEndian.Uint32(buf[44:48]),
		BytesPerLine:      binary.BigEndian.Uint32(buf[48:52]),
		VisualClass:       binary.BigEndian.Uint32(buf[52:56]),
		RedMask:           binary.BigEndian.Uint32(buf[56:60]),
		GreenMask:         binary.BigEndian.Uint32(buf[60:64]),
		BlueMask:          binary.BigEndian.Uint32(buf[64:68]),
		BitsPerRgb:        binary.BigEndian.Uint32(buf[68:72]),
		NumberOfColors:    binary.BigEndian.Uint32(buf[72:76]),
		ColorMapEntries:   binary.BigEndian.Uint32(buf[76:80]),
		WindowWidth:       binary.BigEndian.Uint32(buf[80:84]),
		WindowHeight:      binary.BigEndian.Uint32(buf[84:88]),
		WindowX:           binary.BigEndian.Uint32(buf[88:92]),
		WindowY:           binary.BigEndian.Uint32(buf[92:96]),
		WindowBorderWidth: binary.BigEndian.Uint32(buf[96:100]),
	}
	// window name
	windowName := make([]byte, header.HeaderSize-100)
	_, err = r.Read(windowName)
	if err != nil {
		return
	}
	// not used?
	colorMaps := make([]XWDColorMap, header.ColorMapEntries)
	for i := 0; i < int(header.ColorMapEntries); i++ {
		buf := make([]byte, 12)
		_, err = r.Read(buf)
		if err != nil {
			return
		}
		colorMaps[i] = XWDColorMap{
			EntryNumber: binary.BigEndian.Uint32(buf[0:4]),
			Red:         binary.BigEndian.Uint16(buf[4:6]),
			Green:       binary.BigEndian.Uint16(buf[6:8]),
			Blue:        binary.BigEndian.Uint16(buf[8:10]),
			Flags:       uint8(buf[10]),
			Padding:     uint8(buf[11]),
		}
	}
	// create PalettedImage
	rect := image.Rect(0, 0, int(header.PixmapWidth), int(header.PixmapHeight))
	paletted := image.NewPaletted(rect, palette.WebSafe)
	for x := 0; x < int(header.PixmapHeight); x++ {
		for y := 0; y < int(header.PixmapWidth); y++ {
			buf := make([]byte, 4)
			_, err = r.Read(buf)
			if err != nil {
				return
			}
			paletted.Set(y, x, color.RGBA{
				R: uint8(buf[2]),
				G: uint8(buf[1]),
				B: uint8(buf[0]),
			})
		}
	}
	return paletted, nil
}
