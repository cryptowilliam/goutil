package gdesktop

import (
	"bytes"
	"github.com/kbinani/screenshot"
	"image/png"
)

type CompLvl int

const (
	DefaultCompression = CompLvl(png.DefaultCompression)
	NoCompression      = CompLvl(png.NoCompression)
	BestSpeed          = CompLvl(png.BestSpeed)
	BestCompression    = CompLvl(png.BestCompression)
)

// get display count
func GetActiveScreenCount() int {
	return screenshot.NumActiveDisplays()
}

// get screenshot PNG image
// screenIdx is a 0-base index
func ScreenShot(screenIndex int, lvl CompLvl) ([]byte, error) {
	bounds := screenshot.GetDisplayBounds(screenIndex)
	rgbBuf, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}
	pngBuf := new(bytes.Buffer)

	var e png.Encoder
	e.CompressionLevel = png.CompressionLevel(lvl)
	err = e.Encode(pngBuf, rgbBuf)
	if err != nil {
		return nil, err
	}
	return pngBuf.Bytes(), nil
}
