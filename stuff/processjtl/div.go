package processjtl

import (
	"fmt"
	"jtlweb/stuff/shared"

	"github.com/veandco/go-sdl2/sdl"
)

type Div struct {
	BaseElement
}

func NewDiv(x, y, width, height int32) *Div {
	return &Div{
		BaseElement: BaseElement{
			X:           x,
			Y:           y,
			Width:       width,
			Height:      height,
			Color:       sdl.Color{R: 255, G: 255, B: 255, A: 255},
			BorderColor: sdl.Color{R: 200, G: 200, B: 200, A: 255},
			FontFamily:  "DejaVuSans",
		},
	}
}

func (d *Div) Draw() {
	// Draw the div background if color is set
	rect := &sdl.Rect{
		X: d.X + int32(shared.OffX),
		Y: d.Y + int32(shared.OffY),
		W: d.Width,
		H: d.Height,
	}

	// Draw background
	Renderer.SetDrawColor(d.Color.R, d.Color.G, d.Color.B, d.Color.A)
	Renderer.FillRect(rect)

	// Draw border if border color is set
	Renderer.SetDrawColor(d.BorderColor.R, d.BorderColor.G, d.BorderColor.B, d.BorderColor.A)
	Renderer.DrawRect(rect)

	// Draw all children
	for _, child := range d.Children {
		child.Draw()
	}
}

func (d *Div) CheckClick() {
	// Check clicks for all children
	for _, child := range d.Children {
		child.CheckClick()
	}
}

func (d *Div) GetBaseElement() *BaseElement {
	return &d.BaseElement
}

func (d *Div) String() string {
	return fmt.Sprintf("Div{X: %d, Y: %d, Width: %d, Height: %d, Children: %d}",
		d.X, d.Y, d.Width, d.Height, len(d.Children))
}
