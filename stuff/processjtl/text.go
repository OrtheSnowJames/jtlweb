package processjtl

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Text struct {
	BaseElement // Embed BaseElement
	Content     string
	FontSize    int32
}

func NewText(content string, x, y, fontSize int32, color sdl.Color) *Text {
	return &Text{
		BaseElement: BaseElement{
			X:          x,
			Y:          y,
			Width:      0, // Will be set when drawing
			Height:     fontSize,
			Color:      color,
			FontFamily: "DejaVuSans", // default font
		},
		Content:  content,
		FontSize: fontSize,
	}
}

func (t *Text) Draw() {
	font := GetFont(t.FontFamily)
	surface, err := font.RenderUTF8Blended(t.Content, t.Color)
	if err == nil {
		texture, err := Renderer.CreateTextureFromSurface(surface)
		if err == nil {
			textRect := &sdl.Rect{
				X: t.X,
				Y: t.Y,
				W: int32(surface.W),
				H: int32(surface.H),
			}
			Renderer.Copy(texture, nil, textRect)
			texture.Destroy()
		}
		surface.Free()
	}
}

func (t *Text) CheckClick() {
	// Text elements do not handle clicks
}

// SetFontFamily changes the font of the text element
func (t *Text) SetFontFamily(fontFamily string) {
	t.FontFamily = fontFamily
}
