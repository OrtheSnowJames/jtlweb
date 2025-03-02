package processjtl

import (
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

// UIElement defines the base interface for all UI elements
type UIElement interface {
	Draw()
	CheckClick()
	GetPosition() (int32, int32)
	SetPosition(x, y int32)
	GetSize() (int32, int32)
	SetSize(width, height int32)
	String() string // Add String method
	AddStyle(style string)
	RemoveAllStyle()
}

// BaseElement provides common functionality for UI elements
type BaseElement struct {
	X, Y          int32
	Width, Height int32
	Color         sdl.Color
	BorderColor   sdl.Color
	FontFamily    string
	Class         string
	ID            string
	EventHandlers map[string]string
	Styles        map[string]string // Add Styles map
}

func (b *BaseElement) GetPosition() (int32, int32) {
	return b.X, b.Y
}

func (b *BaseElement) SetPosition(x, y int32) {
	b.X = x
	b.Y = y
}

func (b *BaseElement) GetSize() (int32, int32) {
	return b.Width, b.Height
}

func (b *BaseElement) SetSize(width, height int32) {
	b.Width = width
	b.Height = height
}

func (b *BaseElement) SetEventHandler(event, handler string) {
	if b.EventHandlers == nil {
		b.EventHandlers = make(map[string]string)
	}
	b.EventHandlers[event] = handler
}

func (b *BaseElement) GetEventHandler(event string) string {
	if b.EventHandlers == nil {
		return ""
	}
	return b.EventHandlers[event]
}

func (b *BaseElement) AddStyle(style string) {
	if b.Styles == nil {
		b.Styles = make(map[string]string)
	}
	kv := strings.Split(style, ":")
	if len(kv) == 2 {
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		b.Styles[key] = value
		TranslateStyle(key+":"+value, b)
	}
}

func (b *BaseElement) RemoveAllStyle() {
	b.Styles = make(map[string]string)
	// Reset to default styles if necessary
}
