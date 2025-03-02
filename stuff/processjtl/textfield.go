package processjtl

import (
	"fmt"
	"jtlweb/stuff/shared"

	"github.com/veandco/go-sdl2/sdl"
)

type TextField struct {
	BaseElement
	Text     string
	Active   bool
	Focused  bool
	OnSubmit func(string)
}

func NewTextField(x, y, width, height int32, color, borderColor sdl.Color) *TextField {
	// Get initial window size for centering
	w, h := Window.GetSize()

	return &TextField{
		BaseElement: BaseElement{
			X:           int32(w)/2 - width/2,  // Center horizontally
			Y:           int32(h)/2 - height/2, // Center vertically
			Width:       width,
			Height:      height,
			Color:       color,
			BorderColor: borderColor,
			FontFamily:  "DejaVuSans", // Set default font
		},
		Active:   true,
		Focused:  false,
		OnSubmit: nil,
	}
}

func (t *TextField) Draw() {
	rect := &sdl.Rect{X: t.X + int32(shared.OffX), Y: t.Y + int32(shared.OffY), W: t.Width, H: t.Height}

	// Draw background with proper color handling for focus state
	if t.Focused {
		// Lighten color when focused
		Renderer.SetDrawColor(
			uint8(min(255, float64(t.Color.R)*1.1)),
			uint8(min(255, float64(t.Color.G)*1.1)),
			uint8(min(255, float64(t.Color.B)*1.1)),
			t.Color.A,
		)
	} else {
		Renderer.SetDrawColor(t.Color.R, t.Color.G, t.Color.B, t.Color.A)
	}
	Renderer.FillRect(rect)

	// Draw border
	if t.Focused {
		// redify border when active
		Renderer.SetDrawColor(255, 0, 0, 255)
	} else {
		Renderer.SetDrawColor(t.BorderColor.R, t.BorderColor.G, t.BorderColor.B, t.BorderColor.A)
	}
	Renderer.DrawRect(rect)

	// Render text
	if t.Text != "" {
		surface, err := GetFont(t.FontFamily).RenderUTF8Blended(t.Text, sdl.Color{R: 0, G: 0, B: 0, A: 255})
		if err == nil {
			texture, err := Renderer.CreateTextureFromSurface(surface)
			if err == nil {
				textRect := &sdl.Rect{
					X: t.X + 5 + int32(shared.OffX),
					Y: t.Y + (t.Height-int32(surface.H))/2 + int32(shared.OffY),
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

// Helper function to prevent color overflow (same as in button.go)
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func (t *TextField) SetFocus(focused bool) {
	t.Active = focused
	t.Focused = focused
}

func (t *TextField) CheckClick() {
	// Get current mouse state
	x, y, state := sdl.GetMouseState()

	// Only handle mouse press (not release)
	if state&sdl.ButtonLMask() != 0 {
		if x >= t.X+int32(shared.OffX) && x < t.X+t.Width+int32(shared.OffX) && y >= t.Y+int32(shared.OffY) && y < t.Y+t.Height+int32(shared.OffY) {
			t.SetFocus(true)
			t.Active = true
		} else {
			t.SetFocus(false)
			t.Active = false
		}
	}
}

func (t *TextField) HandleInput(event *sdl.KeyboardEvent) bool {
	if !t.Focused || event.Type != 768 {
		return false
	}

	if event.Keysym.Sym == sdl.K_ESCAPE {
		t.SetFocus(false)
		return false
	}

	switch event.Keysym.Sym {
	case sdl.K_BACKSPACE:
		if len(t.Text) > 0 {
			t.Text = t.Text[:len(t.Text)-1]
		}
	case sdl.K_RETURN:
		if t.OnSubmit != nil {
			t.OnSubmit(t.Text)
		}
		return true
	default:
		if event.Keysym.Sym >= 32 && event.Keysym.Sym <= 126 {
			t.Text += string(event.Keysym.Sym)
		}
	}
	return false
}

// Implement the String method for TextField
func (t *TextField) String() string {
	return fmt.Sprintf("TextField{Text: %s, X: %d, Y: %d, Width: %d, Height: %d}", t.Text, t.X, t.Y, t.Width, t.Height)
}
