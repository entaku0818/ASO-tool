package model

import "time"

type TemplateStyle struct {
	BGColor        string `json:"bg_color,omitempty"`
	BGGradientFrom string `json:"bg_gradient_from,omitempty"`
	BGGradientTo   string `json:"bg_gradient_to,omitempty"`
	BGGradientDir  string `json:"bg_gradient_dir,omitempty"`
	TextColor      string `json:"text_color,omitempty"`
	ImageAlign     string `json:"image_align,omitempty"`
}

type Template struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Category    string        `json:"category"`
	Description string        `json:"description"`
	Device      string        `json:"device"`
	Style       TemplateStyle `json:"style"`
	IsPro       bool          `json:"is_pro"`
	SortOrder   int           `json:"sort_order"`
	CreatedAt   time.Time     `json:"created_at"`
}
