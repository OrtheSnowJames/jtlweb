package processjtl

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ElementCreator func(content string, x, y, width, height int32, styles map[string]string, baseFontSize int32) UIElement

var elementCreators = map[string]ElementCreator{
	"button":    createButton,
	"p":         createText,
	"textfield": createTextField,
	"div":       createDiv, // Add div creator
}

func RegisterElement(elementType string, creator ElementCreator) {
	elementCreators[elementType] = creator
}

func CreateElement(elementType string, content string, x, y, width, height int32, styles map[string]string, baseFontSize int32) UIElement {
	creator, exists := elementCreators[elementType]
	if !exists {
		return nil
	}

	element := creator(content, x, y, width, height, styles, baseFontSize)

	// If children exist in the interface (from JTL parser)
	if baseEl, ok := element.(interface{ GetBaseElement() *BaseElement }); ok {
		// Apply parent's styles to all children (with lower priority)
		for key, value := range styles {
			if key != "class" && key != "id" {
				for _, child := range baseEl.GetBaseElement().Children {
					// Only apply style if child doesn't already have it
					if childBase, ok := child.(interface{ GetBaseElement() *BaseElement }); ok {
						if _, exists := childBase.GetBaseElement().Styles[key]; !exists {
							child.AddStyle(key + ":" + value)
						}
					}
				}
			}
		}
	}

	return element
}

func createButton(content string, x, y, width, height int32, styles map[string]string, baseFontSize int32) UIElement {
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
func createTextField(content string, x, y, width, height int32, styles map[string]string, baseFontSize int32) UIElement {
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
func createText(content string, x, y, width, height int32, styles map[string]string, baseFontSize int32) UIElement {
	text := NewText(content, x, y, baseFontSize,
		sdl.Color{R: 0, G: 0, B: 0, A: 255})

	for key, value := range styles {
		TranslateStyle(key+":"+value, text)
	}
	return text
}

// Add createDiv function
func createDiv(content string, x, y, width, height int32, styles map[string]string, baseFontSize int32) UIElement {
	div := NewDiv(x, y, width, height)

	// Set class and id if present
	if class, ok := styles["class"]; ok {
		div.Class = class
	}
	if id, ok := styles["id"]; ok {
		div.ID = id
	}

	// Apply styles
	for key, value := range styles {
		if key != "class" && key != "id" {
			TranslateStyle(key+":"+value, div)
		}
	}

	return div
}
