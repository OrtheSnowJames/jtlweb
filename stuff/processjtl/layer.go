package processjtl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

func TranslateStyle(style string, element interface{}) {
	styleParts := strings.Split(style, ";")
	for _, part := range styleParts {
		kv := strings.Split(part, ":")
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "width":
			width, _ := strconv.Atoi(value)
			switch e := element.(type) {
			case *Button:
				e.Width = int32(width)
			case *Text:
				e.Width = int32(width)
			case *TextField:
				e.Width = int32(width)
			case *BaseElement:
				e.Width = int32(width)
			}

		case "height":
			height, _ := strconv.Atoi(value)
			switch e := element.(type) {
			case *Button:
				e.Height = int32(height)
			case *Text:
				e.Height = int32(height)
			case *TextField:
				e.Height = int32(height)
			case *BaseElement:
				e.Height = int32(height)
			}

		case "color":
			colorParts := strings.Split(value, ",")
			if len(colorParts) == 4 {
				r, _ := strconv.Atoi(strings.TrimSpace(colorParts[0]))
				g, _ := strconv.Atoi(strings.TrimSpace(colorParts[1]))
				b, _ := strconv.Atoi(strings.TrimSpace(colorParts[2]))
				a, _ := strconv.Atoi(strings.TrimSpace(colorParts[3]))
				color := sdl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}

				switch e := element.(type) {
				case *Button:
					e.Color = color
				case *TextField:
					e.Color = color
				case *BaseElement:
					e.Color = color
				}
			}

		case "border-color":
			colorParts := strings.Split(value, ",")
			if len(colorParts) == 4 {
				r, _ := strconv.Atoi(strings.TrimSpace(colorParts[0]))
				g, _ := strconv.Atoi(strings.TrimSpace(colorParts[1]))
				b, _ := strconv.Atoi(strings.TrimSpace(colorParts[2]))
				a, _ := strconv.Atoi(strings.TrimSpace(colorParts[3]))
				color := sdl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}

				switch e := element.(type) {
				case *Button:
					e.BorderColor = color
				case *TextField:
					e.BorderColor = color
				case *BaseElement:
					e.BorderColor = color
				}
			}

		case "font-family":
			switch e := element.(type) {
			case *Text:
				e.SetFontFamily(value)
			case *TextField:
				e.FontFamily = value
			case *BaseElement:
				e.FontFamily = value
			}

		case "margin", "padding":
			margin, _ := strconv.Atoi(value)
			if button, ok := element.(*Button); ok {
				button.Margin = int32(margin)
			}

		case "margin-left", "margin-right", "margin-up", "margin-down":
			margin, _ := strconv.Atoi(value)
			if text, ok := element.(*Text); ok {
				switch key {
				case "margin-left":
					text.X += int32(margin)
				case "margin-right":
					text.X -= int32(margin)
				case "margin-up":
					text.Y += int32(margin)
				case "margin-down":
					text.Y -= int32(margin)
				}
			}

		case "center":
			if text, ok := element.(*Text); ok {
				text.Center = value == "true"
			}
		}
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
