package imgproc

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"strings"

	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed fonts/NotoSansJP-Regular.ttf
var fontData []byte

// DevicePreset defines the canvas dimensions for a device type.
type DevicePreset struct {
	Width        int
	Height       int
	FrameColor   string
	CornerRadius float64
}

var DevicePresets = map[string]DevicePreset{
	"iphone67": {Width: 393, Height: 852, FrameColor: "#1a1a1a", CornerRadius: 55},
	"iphone65": {Width: 414, Height: 896, FrameColor: "#1a1a1a", CornerRadius: 55},
	"ipad":     {Width: 1024, Height: 1366, FrameColor: "#1a1a1a", CornerRadius: 20},
}

// GenerateRequest holds all inputs for screenshot generation.
type GenerateRequest struct {
	Image          image.Image
	Device         string
	BGColor        string
	BGGradientFrom string // if set, use gradient instead of solid (from color)
	BGGradientTo   string // gradient to color
	BGGradientDir  string // "tb" (top→bottom), "lr" (left→right), "tlbr" (diagonal)
	TextColor      string
	Captions       map[string]string // langCode → caption text
	ImageAlign     string            // "center" (default) or "bottom"
}

// GenerateResult contains the generated PNG bytes per language.
type GenerateResult struct {
	Images map[string][]byte // langCode → PNG bytes
}

// parseHexColor parses a CSS hex color string (#RRGGBB) into 0-1 float64 components.
func parseHexColor(s string) (r, g, b float64) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return 0, 0, 0
	}
	var ri, gi, bi uint8
	_, _ = fmt.Sscanf(s, "%02x%02x%02x", &ri, &gi, &bi)
	return float64(ri) / 255, float64(gi) / 255, float64(bi) / 255
}

// newFontFace parses the embedded NotoSansJP font and returns a face at the given size.
func newFontFace(size float64) (font.Face, error) {
	f, err := opentype.Parse(fontData)
	if err != nil {
		return nil, fmt.Errorf("parse font: %w", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size: size,
		DPI:  72,
	})
	if err != nil {
		return nil, fmt.Errorf("new face: %w", err)
	}
	return face, nil
}

// scaleImageCover scales src to cover the target w×h area (cover fit),
// returning a cropped image.RGBA of exactly w×h.
func scaleImageCover(src image.Image, w, h int) image.Image {
	srcBounds := src.Bounds()
	srcW := float64(srcBounds.Dx())
	srcH := float64(srcBounds.Dy())

	scale := math.Max(float64(w)/srcW, float64(h)/srcH)
	dstW := int(math.Round(srcW * scale))
	dstH := int(math.Round(srcH * scale))

	// Scale into a larger image
	scaled := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	draw.BiLinear.Scale(scaled, scaled.Bounds(), src, srcBounds, draw.Over, nil)

	// Crop to center
	x0 := (dstW - w) / 2
	y0 := (dstH - h) / 2
	cropped := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(cropped, cropped.Bounds(), scaled, image.Pt(x0, y0), draw.Over)
	return cropped
}

const (
	padding     = 40
	screenInset = 8
)

// captionAreaH returns the caption area height based on alignment.
// bottom-align uses a taller caption area so the device clearly "sits at the bottom".
func captionAreaH(imageAlign string) int {
	if imageAlign == "bottom" {
		return 220
	}
	return 120
}

// Generate creates a framed screenshot for each language that has a non-empty caption.
// It mirrors the logic from the frontend ScreenshotGenerator.tsx drawFrame() function.
func Generate(req *GenerateRequest) (*GenerateResult, error) {
	preset, ok := DevicePresets[req.Device]
	if !ok {
		preset = DevicePresets["iphone67"]
	}

	capH := captionAreaH(req.ImageAlign)
	W := preset.Width + padding*2
	H := preset.Height + padding*2 + capH

	result := &GenerateResult{Images: make(map[string][]byte)}

	for lang, caption := range req.Captions {
		if caption == "" {
			continue
		}

		dc := gg.NewContext(W, H)

		// 1. Background fill (solid or gradient)
		if req.BGGradientFrom != "" && req.BGGradientTo != "" {
			var x0, y0, x1, y1 float64
			switch req.BGGradientDir {
			case "lr":
				x0, y0, x1, y1 = 0, 0, float64(W), 0
			case "tlbr":
				x0, y0, x1, y1 = 0, 0, float64(W), float64(H)
			default: // "tb"
				x0, y0, x1, y1 = 0, 0, 0, float64(H)
			}
			grad := gg.NewLinearGradient(x0, y0, x1, y1)
			r1, g1, b1 := parseHexColor(req.BGGradientFrom)
			r2, g2, b2 := parseHexColor(req.BGGradientTo)
			grad.AddColorStop(0, color.RGBA{R: uint8(r1 * 255), G: uint8(g1 * 255), B: uint8(b1 * 255), A: 255})
			grad.AddColorStop(1, color.RGBA{R: uint8(r2 * 255), G: uint8(g2 * 255), B: uint8(b2 * 255), A: 255})
			dc.SetFillStyle(grad)
			dc.DrawRectangle(0, 0, float64(W), float64(H))
			dc.Fill()
			dc.SetFillStyle(gg.NewSolidPattern(color.Black))
		} else {
			br, bg, bb := parseHexColor(req.BGColor)
			dc.SetRGB(br, bg, bb)
			dc.Clear()
		}

		// 2. Caption text (upper capH area, centered)
		// W/13 gives ~36px at iphone67 width — noticeably bolder than W/18.
		// NotoSansJP only has the Regular weight embedded, so we draw twice
		// with a 0.5px x-offset to simulate a heavier stroke (fake bold).
		fontSize := math.Max(float64(W)/13, 14)
		face, err := newFontFace(fontSize)
		if err != nil {
			return nil, err
		}
		dc.SetFontFace(face)
		tr, tg, tb := parseHexColor(req.TextColor)
		dc.SetRGB(tr, tg, tb)
		lineSpacing := 1.4
		maxTextW := float64(W - padding*2)
		captionY := float64(capH) / 2
		// First pass
		dc.DrawStringWrapped(caption, float64(W)/2, captionY,
			0.5, 0.5, maxTextW, lineSpacing, gg.AlignCenter)
		// Second pass offset by 0.5px — fake bold for Regular-only font
		dc.DrawStringWrapped(caption, float64(W)/2+0.5, captionY,
			0.5, 0.5, maxTextW, lineSpacing, gg.AlignCenter)

		// 3. Device frame (outer rounded rectangle)
		fx := float64(padding)
		fy := float64(capH)
		if req.ImageAlign == "bottom" {
			// bottom-align: device sits near the bottom; capH=220 gives breathing room above
			fy = float64(H - padding - preset.Height)
		}
		fw := float64(preset.Width)
		fh := float64(preset.Height)
		cr := preset.CornerRadius

		fr, fg2, fb := parseHexColor(preset.FrameColor)
		dc.SetRGB(fr, fg2, fb)
		dc.DrawRoundedRectangle(fx, fy, fw, fh, cr)
		dc.Fill()

		// 4. Screen area with rounded clip (inset from frame)
		ir := math.Max(cr-screenInset, 4)
		sx := fx + screenInset
		sy := fy + screenInset
		sw := fw - screenInset*2
		sh := fh - screenInset*2

		dc.DrawRoundedRectangle(sx, sy, sw, sh, ir)
		dc.Clip()

		// 5. Draw screenshot image (cover fit inside screen area)
		if req.Image != nil {
			scaled := scaleImageCover(req.Image, int(sw), int(sh))
			dc.DrawImage(scaled, int(sx), int(sy))
		} else {
			dc.SetRGB(0.9, 0.9, 0.9)
			dc.DrawRectangle(sx, sy, sw, sh)
			dc.Fill()
		}

		dc.ResetClip()

		// 6. Encode as PNG
		var buf bytes.Buffer
		if err := png.Encode(&buf, dc.Image()); err != nil {
			return nil, fmt.Errorf("png encode [%s]: %w", lang, err)
		}
		result.Images[lang] = buf.Bytes()
	}

	return result, nil
}
