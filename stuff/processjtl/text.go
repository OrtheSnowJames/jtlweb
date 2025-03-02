package processjtl

import (
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

type Text struct {
	BaseElement // Embed BaseElement
	Content     string
	FontSize    int32
	Margin      int32
	Center      bool
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

func applyTextStyle(key, value string, text *Text) {
	switch key {
	case "font-family":
		text.SetFontFamily(value)
	case "width":
		width, _ := strconv.Atoi(value)
		text.Width = int32(width)
	case "height":
		height, _ := strconv.Atoi(value)
		text.Height = int32(height)
	case "margin":
		margin, _ := strconv.Atoi(value)
		text.Margin = int32(margin)
	case "center":
		text.Center = value == "true"
	}
}

func (t *Text) Draw() {
	font := GetFontWithSize(t.FontFamily, int(t.FontSize))

	surface, err := font.RenderUTF8Blended(t.Content, t.Color)
	if err == nil {
		texture, err := Renderer.CreateTextureFromSurface(surface)
		if err == nil {
			// Update width and height based on rendered text
			textWidth := int32(surface.W)
			textHeight := int32(surface.H)

			// Center the text if the center style is applied
			x := t.X
			y := t.Y
			if t.Center {
				windowWidth, windowHeight := Window.GetSize()
				x = (windowWidth - textWidth) / 2
				y = (windowHeight - textHeight) / 2
			}

			textRect := &sdl.Rect{
				X: x,
				Y: y,
				W: textWidth,
				H: textHeight,
			}
			Renderer.Copy(texture, nil, textRect)
			texture.Destroy()

			// Store the actual dimensions
			t.Width = textWidth
			t.Height = textHeight
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
