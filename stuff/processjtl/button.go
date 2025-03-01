package processjtl

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Button struct {
	BaseElement // Embed BaseElement to inherit its methods
	Text        string
	Margin      int32
	OnClick     func()
	wasPressed  bool // Track if button was previously pressed
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
	rect := &sdl.Rect{X: b.X, Y: b.Y, W: b.Width, H: b.Height}

	// Draw button background
	x, y, state := sdl.GetMouseState()
	isHovered := x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height

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

	// Render text
	surface, err := GetFont(b.FontFamily).RenderUTF8Blended(b.Text, sdl.Color{R: 0, G: 0, B: 0, A: 255})
	if err == nil {
		texture, err := Renderer.CreateTextureFromSurface(surface)
		if err == nil {
			textRect := &sdl.Rect{
				X: b.X + (b.Width-int32(surface.W))/2,
				Y: b.Y + (b.Height-int32(surface.H))/2,
				W: int32(surface.W),
				H: int32(surface.H),
			}
			Renderer.Copy(texture, nil, textRect)
			texture.Destroy()
		}
		surface.Free()
	}
}

func (b *Button) GetBaseElement() *BaseElement {
	return &b.BaseElement
}

func (b *Button) CheckClick() {
	x, y, state := sdl.GetMouseState()
	if x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height {
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
