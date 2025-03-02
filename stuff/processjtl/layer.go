package processjtl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

func TranslateStyle(style string, element interface{}) {
	// Split style into key-value pairs
	styleParts := strings.Split(style, ";")
	for _, part := range styleParts {
		kv := strings.Split(part, ":")
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch v := element.(type) {
		case *Button:
			applyButtonStyle(key, value, v)
		case *Text:
			applyTextStyle(key, value, v)
		case *TextField:
			applyTextFieldStyle(key, value, v)
		case *BaseElement:
			applyBaseElementStyle(key, value, v)
		}
	}
}

func applyButtonStyle(key, value string, button *Button) {
	switch key {
	case "width":
		width, _ := strconv.Atoi(value)
		button.Width = int32(width)
	case "height":
		height, _ := strconv.Atoi(value)
		button.Height = int32(height)
	case "color":
		colorParts := strings.Split(value, ",")
		if len(colorParts) == 4 {
			r, _ := strconv.Atoi(strings.TrimSpace(colorParts[0]))
			g, _ := strconv.Atoi(strings.TrimSpace(colorParts[1]))
			b, _ := strconv.Atoi(strings.TrimSpace(colorParts[2]))
			a, _ := strconv.Atoi(strings.TrimSpace(colorParts[3]))
			button.Color = sdl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
		}
	case "border-color":
		colorParts := strings.Split(value, ",")
		if len(colorParts) == 4 {
			r, _ := strconv.Atoi(strings.TrimSpace(colorParts[0]))
			g, _ := strconv.Atoi(strings.TrimSpace(colorParts[1]))
			b, _ := strconv.Atoi(strings.TrimSpace(colorParts[2]))
			a, _ := strconv.Atoi(strings.TrimSpace(colorParts[3]))
			button.BorderColor = sdl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
		}
	case "margin":
		margin, _ := strconv.Atoi(value)
		button.Margin = int32(margin)
	case "padding":
		padding, _ := strconv.Atoi(value)
		button.Margin = int32(padding)
	}
}

func applyTextFieldStyle(key, value string, tf *TextField) {
	switch key {
	case "width":
		width, _ := strconv.Atoi(value)
		tf.Width = int32(width)
	case "height":
		height, _ := strconv.Atoi(value)
		tf.Height = int32(height)
	case "color":
		colorParts := strings.Split(value, ",")
		if len(colorParts) == 4 {
			r, _ := strconv.Atoi(strings.TrimSpace(colorParts[0]))
			g, _ := strconv.Atoi(strings.TrimSpace(colorParts[1]))
			b, _ := strconv.Atoi(strings.TrimSpace(colorParts[2]))
			a, _ := strconv.Atoi(strings.TrimSpace(colorParts[3]))
			tf.Color = sdl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
		}
	case "font-family":
		tf.FontFamily = value
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
		text.X += int32(margin)
		text.Y += int32(margin)
	case "margin-left":
		margin, _ := strconv.Atoi(value)
		text.X += int32(margin)
	case "margin-right":
		margin, _ := strconv.Atoi(value)
		text.X -= int32(margin)
	case "margin-up":
		margin, _ := strconv.Atoi(value)
		text.Y += int32(margin)
	case "margin-down":
		margin, _ := strconv.Atoi(value)
		text.Y -= int32(margin)
	case "center":
		text.Center = value == "true"
	}
}

func applyBaseElementStyle(key, value string, base *BaseElement) {
	switch key {
	case "width":
		width, _ := strconv.Atoi(value)
		base.Width = int32(width)
	case "height":
		height, _ := strconv.Atoi(value)
		base.Height = int32(height)
	case "color":
		colorParts := strings.Split(value, ",")
		if len(colorParts) == 4 {
			r, _ := strconv.Atoi(strings.TrimSpace(colorParts[0]))
			g, _ := strconv.Atoi(strings.TrimSpace(colorParts[1]))
			b, _ := strconv.Atoi(strings.TrimSpace(colorParts[2]))
			a, _ := strconv.Atoi(strings.TrimSpace(colorParts[3]))
			base.Color = sdl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
		}
	case "border-color":
		colorParts := strings.Split(value, ",")
		if len(colorParts) == 4 {
			r, _ := strconv.Atoi(strings.TrimSpace(colorParts[0]))
			g, _ := strconv.Atoi(strings.TrimSpace(colorParts[1]))
			b, _ := strconv.Atoi(strings.TrimSpace(colorParts[2]))
			a, _ := strconv.Atoi(strings.TrimSpace(colorParts[3]))
			base.BorderColor = sdl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
		}
	case "font-family":
		base.FontFamily = value
	}
}

type CanvasObject interface {
	Draw()
	CheckClick()
	String() string
}

var objects []CanvasObject

// Says to raylib, but really i was too lazy to rename it.
func ToRaylib(jtlcomps []interface{}) []CanvasObject {
	result := make([]CanvasObject, 0)
	yOffset := int32(20)
	margin := int32(20) // Fixed margin instead of screen-relative

	for _, elem := range jtlcomps {
		comp, ok := elem.(map[string]interface{})
		if !ok {
			continue
		}

		key, keyExists := comp["KEY"].(string)
		if !keyExists {
			continue
		}

		// Extract class and id directly from the component
		class, _ := comp["class"].(string)
		id, _ := comp["id"].(string)

		content, _ := comp["Contents"].(string)
		styles, _ := comp["style"].(string)
		parsedStyles := ParseCSS(styles)

		// Add class and id to parsed styles
		if class != "" {
			parsedStyles["class"] = class
		}
		if id != "" {
			parsedStyles["id"] = id
		}

		// Use fixed dimensions instead of screen-relative
		width := int32(200) // Fixed width
		height := int32(40) // Fixed height

		if element := CreateElement(key, content,
			margin, yOffset,
			width, height, parsedStyles, 14); element != nil { // Use fixed font size

			// Debug print
			if baseEl, ok := element.(interface{ GetBaseElement() *BaseElement }); ok {
				fmt.Printf("Created element with class: %s\n", baseEl.GetBaseElement().Class)
			}

			result = append(result, element.(CanvasObject))
			yOffset += height + 20 // Fixed spacing
		}
	}

	return result
}
