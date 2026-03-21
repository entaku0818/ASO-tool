package imgproc

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestCaptionAreaH(t *testing.T) {
	if h := captionAreaH("bottom"); h != 220 {
		t.Errorf("bottom: got %d, want 220", h)
	}
	if h := captionAreaH("center"); h != 120 {
		t.Errorf("center: got %d, want 120", h)
	}
	if h := captionAreaH(""); h != 120 {
		t.Errorf("default: got %d, want 120", h)
	}
}

func TestGenerate_BottomAlign_Height(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 390, 844))
	for y := 0; y < 844; y++ {
		for x := 0; x < 390; x++ {
			img.Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
		}
	}
	for _, tc := range []struct {
		align string
		wantH int
	}{
		{"bottom", 1152},
		{"center", 1052},
		{"",       1052},
	} {
		r, err := Generate(&GenerateRequest{
			Image: img, Device: "iphone67",
			BGColor: "#4F46E5", TextColor: "#FFFFFF",
			ImageAlign: tc.align,
			Captions:   map[string]string{"ja": "テスト"},
		})
		if err != nil {
			t.Fatalf("align=%q: %v", tc.align, err)
		}
		decoded, err := png.Decode(bytes.NewReader(r.Images["ja"]))
		if err != nil {
			t.Fatalf("decode: %v", err)
		}
		gotH := decoded.Bounds().Dy()
		if gotH != tc.wantH {
			t.Errorf("align=%q: height got %d, want %d", tc.align, gotH, tc.wantH)
		} else {
			t.Logf("align=%q: ✅ height=%d", tc.align, gotH)
		}
	}
}

func TestGenerate_AllPresets(t *testing.T) {
	img := makeImg(390, 844)
	presets := []struct{ name, from, to string }{
		{"ocean",    "#667eea", "#764ba2"},
		{"midnight", "#0f0c29", "#302b63"},
		{"sunset",   "#f093fb", "#f5576c"},
		{"forest",   "#134e5e", "#71b280"},
		{"gold",     "#f7971e", "#ffd200"},
		{"rose",     "#ee0979", "#ff6a00"},
		{"sky",      "#4facfe", "#00f2fe"},
		{"dark",     "#1a1a2e", "#16213e"},
	}
	for _, p := range presets {
		r, err := Generate(&GenerateRequest{
			Image: img, Device: "iphone67",
			BGGradientFrom: p.from, BGGradientTo: p.to, BGGradientDir: "tb",
			TextColor: "#FFFFFF",
			Captions:  map[string]string{"en": p.name},
		})
		if err != nil {
			t.Errorf("preset %s: %v", p.name, err)
			continue
		}
		if len(r.Images["en"]) == 0 {
			t.Errorf("preset %s: empty image", p.name)
			continue
		}
		t.Logf("preset %-10s ✅ %d bytes", p.name, len(r.Images["en"]))
	}
}
