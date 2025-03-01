package processjtl

import (
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
		}
	}
}

func applyTextStyle(key, value string, text *Text) {
	switch key {
	case "font-family":
		text.SetFontFamily(value)
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

type CanvasObject interface {
	Draw()
	CheckClick()
}

var objects []CanvasObject

// Says to raylib, but really i was too lazy to rename it.
func ToRaylib(jtlcomps []interface{}) []CanvasObject {
	result := make([]CanvasObject, 0)
	w, h := Window.GetSize()
	screenWidth := float32(w)
	screenHeight := float32(h)
	yOffset := float32(20)
	margin := screenWidth * 0.02

	for _, elem := range jtlcomps {
		comp, ok := elem.(map[string]interface{})
		if !ok {
			continue
		}

		key, keyExists := comp["KEY"].(string)
		if !keyExists {
			continue
		}

		content, _ := comp["Contents"].(string)
		styles, _ := comp["style"].(string)
		parsedStyles := ParseCSS(styles)

		width := int32(screenWidth * 0.25)
		height := int32(screenHeight * 0.08)

		if element := CreateElement(key, content,
			int32(margin), int32(yOffset),
			width, height, parsedStyles); element != nil {
			result = append(result, element.(CanvasObject))
			yOffset += float32(height) + (screenHeight * 0.02)
		}
	}

	return result
}
