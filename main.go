package main

import (
	"encoding/json"
	"fmt"
	"jtlweb/stuff/processjtl"
	"jtlweb/stuff/shared"
	"os"
	"path/filepath"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

type AppState int

var config map[string]interface{}
var openPath string

const (
	StateInput AppState = iota
	StateRendering
)

func main() {
	if err := processjtl.InitSDL(); err != nil {
		panic(err)
	}
	defer processjtl.CleanupSDL()
	_, err := readConf()
	if err != nil {
		fmt.Println("Error reading config file: ", err)
		os.Exit(1)
	}

	state := StateInput
	var objects []processjtl.CanvasObject
	var winlock *processjtl.Locker

	// Set the callback to update main objects
	processjtl.SetUpdateMainObjectsCallback(func(newObjects []processjtl.CanvasObject) {
		objects = newObjects
	})

	textField := processjtl.NewTextField(100, 250, 600, 40,
		sdl.Color{R: 255, G: 255, B: 255, A: 255},
		sdl.Color{R: 100, G: 100, B: 100, A: 255})

	running := true
	for running {
		if shared.OpenPath != openPath {
			shared.OpenPath = openPath
			fmt.Printf("shared.OpenPath: %v\n", shared.OpenPath)
		}
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_RESIZED {
					// Update text field position on window resize
					w, h := processjtl.Window.GetSize()
					textField.X = int32(w)/2 - textField.Width/2
					textField.Y = int32(h)/2 - textField.Height/2
				}
			case *sdl.MouseButtonEvent:
				if state == StateInput {
					textField.CheckClick() // Add this line to handle mouse clicks
				}
			case *sdl.KeyboardEvent:
				if e.Keysym.Sym == sdl.K_ESCAPE && state == StateRendering {
					state = StateInput
					textField.Text = ""
				} else if state == StateInput {
					if textField.HandleInput(e) {
						// Handle file loading
						content, err := os.ReadFile(textField.Text)
						if err == nil {
							textField.Text, err = getFullPath(textField.Text)
							if err != nil {
								fmt.Println("Error getting full path: ", err)
								continue
							}
							openPath = textField.Text
							shared.OpenPath = openPath

							// Clear error handling and debug output
							winlock, objects = processjtl.MakeWebview(string(content))
							if winlock != nil {
							}
							if objects == nil {
								fmt.Println("No objects created from JTL document")
								continue
							}

							fmt.Printf("Created %d objects\n", len(objects))
							state = StateRendering
						} else {
							fmt.Printf("Error reading file: %v\n", err)
						}
					}
				}
			}
		}

		processjtl.Renderer.SetDrawColor(240, 240, 240, 255)
		processjtl.Renderer.Clear()

		switch state {
		case StateInput:
			drawInputState(textField)
		case StateRendering:
			drawRenderingState(objects)
		}

		processjtl.Renderer.Present()
		sdl.Delay(16) // Cap at ~60 FPS
	}
}

func readConf() (string, error) {
	exedir, err := processjtl.GetExeDir()

	if err != nil {
		return "", err
	}

	file, err := os.ReadFile(exedir + "/conf.json")
	if err != nil {
		return "", fmt.Errorf("error reading config file: %v", err)
	}

	filestring := string(file)

	// unmarshal JSON
	err = json.Unmarshal([]byte(filestring), &config)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return filestring, nil
}

func drawInputState(textField *processjtl.TextField) {
	// Draw input prompt text centered above the text field
	surface, err := processjtl.Fonts[config["defaultUrlTextboxFont"].(string)].RenderUTF8Blended("Enter JTL file path:",
		sdl.Color{R: 0, G: 0, B: 0, A: 255})
	if err == nil {
		texture, err := processjtl.Renderer.CreateTextureFromSurface(surface)
		if err == nil {
			w, _ := processjtl.Window.GetSize()
			textRect := &sdl.Rect{
				X: int32(w)/2 - int32(surface.W)/2,
				Y: textField.Y - 30,
				W: int32(surface.W),
				H: int32(surface.H),
			}
			textField.FontFamily = config["defaultUrlTextboxFont"].(string)
			processjtl.Renderer.Copy(texture, nil, textRect)
			texture.Destroy()
		}
		surface.Free()
	}

	textField.Draw()
}

func drawRenderingState(objects []processjtl.CanvasObject) {
	processjtl.ObjectsMutex.Lock()
	localObjects := make([]processjtl.CanvasObject, len(objects))
	copy(localObjects, objects)
	processjtl.ObjectsMutex.Unlock()

	if len(localObjects) == 0 {
		fmt.Println("No objects to draw")
		return
	}

	for _, obj := range localObjects {
		if obj == nil {
			continue
		}
		obj.Draw()
		obj.CheckClick()
	}
}

func getFullPath(inputPath string) (string, error) {
	// Expand the user's home directory (~) to its full path
	expandedPath := inputPath
	if strings.HasPrefix(inputPath, "~") {
		expandedPath = strings.Replace(inputPath, "~", os.Getenv("HOME"), 1)
	}

	// If the inputPath is not absolute, make it absolute
	if !filepath.IsAbs(expandedPath) {
		// Convert the path to an absolute path
		absPath, err := filepath.Abs(expandedPath)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path: %v", err)
		}
		return absPath, nil
	}

	// If the inputPath is already absolute, return it as is
	return expandedPath, nil
}
