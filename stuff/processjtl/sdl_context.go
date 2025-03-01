package processjtl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Fonts    map[string]*ttf.Font
)

const (
	pathtoassets    = "/assets/"
	defaultFontSize = 14
)

func InitSDL() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("SDL Init Error: %v", err)
	}

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("TTF Init Error: %v", err)
	}

	var err error
	Window, err = sdl.CreateWindow("JTL Webview",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600,
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return fmt.Errorf("window Creation Error: %v", err)
	}

	Renderer, err = sdl.CreateRenderer(Window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return fmt.Errorf("renderer Creation Error: %v", err)
	}

	// Initialize fonts map
	Fonts = make(map[string]*ttf.Font)

	exeDir, err := GetExeDir()
	if err != nil {
		return fmt.Errorf("error getting exe directory: %v", err)
	}

	// Load DejaVuSans font
	Fonts["DejaVuSans"], err = ttf.OpenFont(fmt.Sprintf("%s%sDejaVuSans.ttf", exeDir, pathtoassets), defaultFontSize)
	if err != nil {
		return fmt.Errorf("DejaVuSans font loading error: %v", err)
	}

	// Load JetBrains Mono font
	Fonts["JetBrainsMono"], err = ttf.OpenFont(fmt.Sprintf("%s%sJetBrainsMono-Regular.ttf", exeDir, pathtoassets), defaultFontSize)
	if err != nil {
		return fmt.Errorf("JetBrainsMono font loading error: %v", err)
	}

	return nil
}

func CleanupSDL() {
	for _, font := range Fonts {
		if font != nil {
			font.Close()
		}
	}
	if Renderer != nil {
		Renderer.Destroy()
	}
	if Window != nil {
		Window.Destroy()
	}
	ttf.Quit()
	sdl.Quit()
}

func GetExeDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}

// GetFont returns the requested font or falls back to DejaVuSans
func GetFont(fontFamily string) *ttf.Font {
	if font, ok := Fonts[fontFamily]; ok {
		return font
	}
	return Fonts["DejaVuSans"] // fallback to default
}
