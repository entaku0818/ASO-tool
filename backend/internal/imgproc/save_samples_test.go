package imgproc

import (
	"image/png"
	"bytes"
	"os"
	"testing"
)

func TestSaveSamples(t *testing.T) {
	img := makeImg(390, 844)
	if err := os.MkdirAll("/tmp/screenshot_designer_qa", 0755); err != nil {
		t.Fatal(err)
	}

	// ① bottom vs center
	for _, align := range []string{"bottom", "center"} {
		r, _ := Generate(&GenerateRequest{
			Image: img, Device: "iphone67",
			BGGradientFrom: "#667eea", BGGradientTo: "#764ba2", BGGradientDir: "tb",
			TextColor: "#FFFFFF", ImageAlign: align,
			Captions: map[string]string{"ja": "画像位置: " + align},
		})
		if err := os.WriteFile("/tmp/screenshot_designer_qa/align_"+align+"_new.png", r.Images["ja"], 0644); err != nil {
			t.Logf("write failed: %v", err)
		}
		decoded, _ := png.Decode(bytes.NewReader(r.Images["ja"]))
		t.Logf("align=%-6s: %dx%dpx", align, decoded.Bounds().Dx(), decoded.Bounds().Dy())
	}

	// ③ 全8プリセット
	presets := []struct{ id, from, to string }{
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
		r, _ := Generate(&GenerateRequest{
			Image: img, Device: "iphone67",
			BGGradientFrom: p.from, BGGradientTo: p.to, BGGradientDir: "tb",
			TextColor: "#FFFFFF", ImageAlign: "bottom",
			Captions: map[string]string{"ja": p.id + "プリセット"},
		})
		if err := os.WriteFile("/tmp/screenshot_designer_qa/preset_"+p.id+".png", r.Images["ja"], 0644); err != nil {
			t.Logf("write failed: %v", err)
		}
	}
	t.Log("サンプル保存完了: /tmp/screenshot_designer_qa/")
}
