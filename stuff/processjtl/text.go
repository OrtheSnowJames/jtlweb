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
	Rotation    float64 // Add Rotation field
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
	if t.Rotation != 0 {
		font := GetFontWithSize(t.FontFamily, int(t.FontSize))
		var contents string
		if t.Content == "" {
			contents = "Blank String..."
		} else {
			contents = t.Content
		}

		// Create surface for measuring text dimensions
		surface, err := font.RenderUTF8Blended(contents, t.Color)
		if err != nil {
			return
		}
		textWidth := int32(surface.W)
		textHeight := int32(surface.H)
		surface.Free()

		// Create texture for rotation
		texture, err := Renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
			sdl.TEXTUREACCESS_TARGET, textWidth, textHeight)
		if err != nil {
			return
		}
		defer texture.Destroy()

		// Set render target to texture
		prevTarget := Renderer.GetRenderTarget()
		Renderer.SetRenderTarget(texture)
		Renderer.SetDrawColor(240, 240, 240, 0) // Transparent background
		Renderer.Clear()

		// Render text to texture
		surface, err = font.RenderUTF8Blended(contents, t.Color)
		if err == nil {
			textTexture, err := Renderer.CreateTextureFromSurface(surface)
			if err == nil {
				Renderer.Copy(textTexture, nil, &sdl.Rect{X: 0, Y: 0, W: textWidth, H: textHeight})
				textTexture.Destroy()
			}
			surface.Free()
		}

		// Reset target and draw rotated texture
		Renderer.SetRenderTarget(prevTarget)

		// Calculate position
		x := t.X + int32(shared.OffX) + t.Margin
		y := t.Y + int32(shared.OffY) + t.Margin
		if t.Center {
			windowWidth, windowHeight := Window.GetSize()
			x = (windowWidth-textWidth)/2 + int32(shared.OffX)
			y = (windowHeight-textHeight)/2 + int32(shared.OffY)
		}

		dstRect := &sdl.Rect{
			X: x,
			Y: y,
			W: textWidth,
			H: textHeight,
		}

		texture.SetBlendMode(sdl.BLENDMODE_BLEND)
		Renderer.CopyEx(texture, nil, dstRect, t.Rotation, nil, sdl.FLIP_NONE)

		// Store the actual dimensions
		t.Width = textWidth
		t.Height = textHeight
	} else {
		// Original non-rotated drawing code
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
