package imgproc

import (
	"image"
	"image/color"
	"testing"
)

func makeImg(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
		}
	}
	return img
}

func TestGenerate_NilImage(t *testing.T) {
	r, err := Generate(&GenerateRequest{
		Image: nil, Device: "iphone67", BGColor: "#4F46E5", TextColor: "#FFFFFF",
		Captions: map[string]string{"ja": "テスト"},
	})
	if err != nil {
		t.Fatalf("nil image should not error: %v", err)
	}
	if len(r.Images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(r.Images))
	}
}

func TestGenerate_EmptyCaptions(t *testing.T) {
	r, err := Generate(&GenerateRequest{
		Image: makeImg(390, 844), Device: "iphone67", BGColor: "#4F46E5", TextColor: "#FFFFFF",
		Captions: map[string]string{},
	})
	if err != nil {
		t.Fatalf("empty captions should not error: %v", err)
	}
	if len(r.Images) != 0 {
		t.Fatalf("expected 0 images, got %d", len(r.Images))
	}
}

func TestGenerate_AllEmptyCaptionValues(t *testing.T) {
	r, err := Generate(&GenerateRequest{
		Image: makeImg(390, 844), Device: "iphone67", BGColor: "#4F46E5", TextColor: "#FFFFFF",
		Captions: map[string]string{"ja": "", "en": ""},
	})
	if err != nil {
		t.Fatalf("all-empty caption values should not error: %v", err)
	}
	if len(r.Images) != 0 {
		t.Fatalf("expected 0 images (all captions empty), got %d", len(r.Images))
	}
}

func TestGenerate_InvalidHexColor(t *testing.T) {
	// Invalid hex should silently fallback to 0,0,0 (black) — verify no panic/error
	_, err := Generate(&GenerateRequest{
		Image: makeImg(390, 844), Device: "iphone67", BGColor: "notacolor", TextColor: "bad",
		Captions: map[string]string{"en": "test"},
	})
	if err != nil {
		t.Fatalf("invalid hex color should not error (silent fallback): %v", err)
	}
}

func TestGenerate_UnknownDevice(t *testing.T) {
	// Unknown device should fallback to iphone67
	r, err := Generate(&GenerateRequest{
		Image: makeImg(390, 844), Device: "galaxy_s25", BGColor: "#FFFFFF", TextColor: "#000",
		Captions: map[string]string{"en": "fallback device"},
	})
	if err != nil {
		t.Fatalf("unknown device should fallback to iphone67: %v", err)
	}
	if len(r.Images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(r.Images))
	}
}

func TestGenerate_NilCaptionsMap(t *testing.T) {
	// nil captions map — should not panic
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("panic on nil captions: %v", rec)
		}
	}()
	r, err := Generate(&GenerateRequest{
		Image: makeImg(390, 844), Device: "iphone67", BGColor: "#4F46E5", TextColor: "#FFFFFF",
		Captions: nil,
	})
	if err != nil {
		t.Fatalf("nil captions should not error: %v", err)
	}
	if len(r.Images) != 0 {
		t.Fatalf("expected 0 images for nil captions, got %d", len(r.Images))
	}
}

func TestGenerate_LongCaption(t *testing.T) {
	long := "これは非常に長いキャプションテキストです。テキストの折り返し処理が正しく動作するかどうかを確認するために使用します。さらに長くしてみましょう。"
	_, err := Generate(&GenerateRequest{
		Image: makeImg(390, 844), Device: "iphone67", BGColor: "#000000", TextColor: "#FFFFFF",
		Captions: map[string]string{"ja": long},
	})
	if err != nil {
		t.Fatalf("long caption should not error: %v", err)
	}
}

func TestGenerate_iPad(t *testing.T) {
	r, err := Generate(&GenerateRequest{
		Image: makeImg(2048, 2732), Device: "ipad", BGColor: "#FFFFFF", TextColor: "#333333",
		Captions: map[string]string{"ja": "iPadテスト", "en": "iPad test"},
	})
	if err != nil {
		t.Fatalf("ipad device should work: %v", err)
	}
	if len(r.Images) != 2 {
		t.Fatalf("expected 2 images, got %d", len(r.Images))
	}
}
