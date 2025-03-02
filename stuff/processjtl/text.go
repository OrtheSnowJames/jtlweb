package processjtl

import (
	"fmt"
	"jtlweb/stuff/shared"

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

func (t *Text) Draw() {
	font := GetFontWithSize(t.FontFamily, int(t.FontSize))
	var contents string
	if t.Content == "" {
		contents = "Blank String..."
	} else {
		contents = t.Content
	}
	surface, err := font.RenderUTF8Blended(contents, t.Color)
	if err == nil {
		texture, err := Renderer.CreateTextureFromSurface(surface)
		if err == nil {
			// Update width and height based on rendered text
			textWidth := int32(surface.W)
			textHeight := int32(surface.H)

			// Center the text if the center style is applied
			x := t.X + int32(shared.OffX) + t.Margin
			y := t.Y + int32(shared.OffY) + t.Margin
			if t.Center {
				windowWidth, windowHeight := Window.GetSize()
				x = (windowWidth-textWidth)/2 + int32(shared.OffX)
				y = (windowHeight-textHeight)/2 + int32(shared.OffY)
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

// Implement the String method for Text
func (t *Text) String() string {
	return fmt.Sprintf("Text{Content: %s, X: %d, Y: %d, Width: %d, Height: %d}", t.Content, t.X, t.Y, t.Width, t.Height)
}
