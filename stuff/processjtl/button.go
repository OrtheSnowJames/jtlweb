package processjtl

import (
	"fmt"
	"jtlweb/stuff/shared"

	"github.com/veandco/go-sdl2/sdl"
)

type Button struct {
	BaseElement // Embed BaseElement to inherit its methods
	Text        string
	Margin      int32
	OnClick     func()
	wasPressed  bool    // Track if button was previously pressed
	Rotation    float64 // Add Rotation field to Button struct
}

func NewButton(text string, x, y, width, height, margin int32, color, borderColor sdl.Color, onClick func()) *Button {
	return &Button{
		BaseElement: BaseElement{
			X:           x,
			Y:           y,
			Width:       width,
			Height:      height,
			Color:       color,
			BorderColor: borderColor,
			FontFamily:  "DejaVuSans",
		},
		Text:    text,
		Margin:  margin,
		OnClick: onClick,
	}
}

func (b *Button) Draw() {
	// Create a render target for rotation
	if b.Rotation != 0 {
		texture, err := Renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
			sdl.TEXTUREACCESS_TARGET, b.Width, b.Height)
		if err != nil {
			return
		}
		defer texture.Destroy()

		// Set render target to our texture
		prevTarget := Renderer.GetRenderTarget()
		Renderer.SetRenderTarget(texture)
		Renderer.SetDrawColor(240, 240, 240, 255)
		Renderer.Clear()

		// Draw button to texture
		rect := &sdl.Rect{X: 0, Y: 0, W: b.Width, H: b.Height}
		Renderer.SetDrawColor(b.Color.R, b.Color.G, b.Color.B, b.Color.A)
		Renderer.FillRect(rect)

		// Draw the rest of your button content here...

		// Reset render target and draw the rotated texture
		Renderer.SetRenderTarget(prevTarget)
		dstRect := &sdl.Rect{
			X: b.X + int32(shared.OffX),
			Y: b.Y + int32(shared.OffY),
			W: b.Width,
			H: b.Height,
		}
		texture.SetBlendMode(sdl.BLENDMODE_BLEND)
		Renderer.CopyEx(texture, nil, dstRect, b.Rotation, nil, sdl.FLIP_NONE)
	} else {
		// Normal drawing code when no rotation
		rect := &sdl.Rect{X: b.X + int32(shared.OffX), Y: b.Y + int32(shared.OffY), W: b.Width, H: b.Height}

		// Draw button background
		x, y, state := sdl.GetMouseState()
		isHovered := x >= b.X+int32(shared.OffX) && x < b.X+b.Width+int32(shared.OffX) && y >= b.Y+int32(shared.OffY) && y < b.Y+b.Height+int32(shared.OffY)

		if state&sdl.ButtonLMask() != 0 && isHovered {
			// Darken color when clicked
			Renderer.SetDrawColor(
				uint8(float64(b.Color.R)*0.8),
				uint8(float64(b.Color.G)*0.8),
				uint8(float64(b.Color.B)*0.8),
				b.Color.A,
			)
		} else if isHovered {
			// Lighten color when hovered
			Renderer.SetDrawColor(
				uint8(min(255, float64(b.Color.R)*1.2)),
				uint8(min(255, float64(b.Color.G)*1.2)),
				uint8(min(255, float64(b.Color.B)*1.2)),
				b.Color.A,
			)
		} else {
			// Normal color
			Renderer.SetDrawColor(b.Color.R, b.Color.G, b.Color.B, b.Color.A)
		}
		Renderer.FillRect(rect)

		// Draw border
		Renderer.SetDrawColor(b.BorderColor.R, b.BorderColor.G, b.BorderColor.B, b.BorderColor.A)
		Renderer.DrawRect(rect)

		// Render text with fixed font size
		font := GetFont(b.FontFamily)

		surface, err := font.RenderUTF8Blended(b.Text, sdl.Color{R: 0, G: 0, B: 0, A: 255})
		if err == nil {
			texture, err := Renderer.CreateTextureFromSurface(surface)
			if err == nil {
				textRect := &sdl.Rect{
					X: b.X + (b.Width-int32(surface.W))/2 + int32(shared.OffX),
					Y: b.Y + (b.Height-int32(surface.H))/2 + int32(shared.OffY),
					W: int32(surface.W),
					H: int32(surface.H),
				}
				Renderer.Copy(texture, nil, textRect)
				texture.Destroy()
			}
			surface.Free()
		}
	}
}

func (b *Button) GetBaseElement() *BaseElement {
	return &b.BaseElement
}

func (b *Button) CheckClick() {
	x, y, state := sdl.GetMouseState()
	if x >= b.X+int32(shared.OffX) && x < b.X+b.Width+int32(shared.OffX) && y >= b.Y+int32(shared.OffY) && y < b.Y+b.Height+int32(shared.OffY) {
		if state&sdl.ButtonLMask() != 0 {
			// Handle repeating click event
			executeEventHandler(&b.BaseElement, "clickrepeat")

			// Handle one-time click event only when button wasn't pressed before
			if !b.wasPressed {
				executeEventHandler(&b.BaseElement, "click")
				b.wasPressed = true
			}
		} else {
			// Reset the pressed state when mouse button is released
			b.wasPressed = false
		}
	} else {
		// Reset the pressed state when mouse leaves button area
		b.wasPressed = false
	}
}

// Implement the String method for Button
func (b *Button) String() string {
	return fmt.Sprintf("Button{Text: %s, X: %d, Y: %d, Width: %d, Height: %d}", b.Text, b.X, b.Y, b.Width, b.Height)
}
