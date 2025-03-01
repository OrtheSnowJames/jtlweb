package processjtl

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ElementCreator func(content string, x, y, width, height int32, styles map[string]string) UIElement

var elementCreators = map[string]ElementCreator{
	"button":    createButton,
	"p":         createText,
	"textfield": createTextField,
}

func RegisterElement(elementType string, creator ElementCreator) {
	elementCreators[elementType] = creator
}

func CreateElement(elementType string, content string, x, y, width, height int32, styles map[string]string) UIElement {
	creator, exists := elementCreators[elementType]
	if !exists {
		return nil
	}
	return creator(content, x, y, width, height, styles)
}

func createButton(content string, x, y, width, height int32, styles map[string]string) UIElement {
	button := NewButton(content, x, y, width, height, 20,
		sdl.Color{R: 200, G: 200, B: 200, A: 255},
		sdl.Color{R: 100, G: 100, B: 100, A: 255},
		nil) // Set onClick to nil initially

	// Set class and id first
	if class, ok := styles["class"]; ok {
		button.Class = class
	}
	if id, ok := styles["id"]; ok {
		button.ID = id
	}

	// Then apply other styles
	for key, value := range styles {
		if key != "class" && key != "id" {
			TranslateStyle(key+":"+value, button)
		}
	}
	return button
}

// Add createTextField function
func createTextField(content string, x, y, width, height int32, styles map[string]string) UIElement {
	textField := NewTextField(x, y, width, height,
		sdl.Color{R: 255, G: 255, B: 255, A: 255},
		sdl.Color{R: 100, G: 100, B: 100, A: 255})
	textField.Text = content

	for key, value := range styles {
		TranslateStyle(key+":"+value, textField)
	}
	return textField
}

// Add createText function
func createText(content string, x, y, width, height int32, styles map[string]string) UIElement {
	text := NewText(content, x, y, int32(float32(height)*0.6),
		sdl.Color{R: 0, G: 0, B: 0, A: 255})

	for key, value := range styles {
		TranslateStyle(key+":"+value, text)
	}
	return text
}
