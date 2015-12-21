package goldframework

import (
	"fmt"
	"runtime"

	sf "github.com/manyminds/gosfml"
	"github.com/sadlil/gologger"
)

func init() {
	runtime.LockOSThread()
}

var logger = gologger.GetLogger(gologger.CONSOLE, gologger.ColoredLog)

var (
	windowWidth  uint = 800
	windowHeight uint = 600
	windowName        = "Gold Engine"
	renderWindow *sf.RenderWindow
)

const (
	//NoWindow: Either Create Window was never called or the window was destroyed.
	NoWindow = "No Window Available"
)

//CreateWindow creates a Game Window, all Inputs are taken from it. Calling it
//destroys previous window
func CreateWindow(width, height uint, name string) {
	if renderWindow != nil {
		renderWindow.Close()
		logger.Warn(fmt.Sprintf("Window %s has been closed", windowName))
		renderWindow = nil
	}
	windowWidth = width
	windowHeight = height
	windowName = name
	renderWindow = sf.NewRenderWindow(sf.VideoMode{Width: windowWidth, Height: windowHeight, BitsPerPixel: 32}, windowName, sf.StyleDefault, sf.DefaultContextSettings())
}

//ResizeWindow resizes window by width and height
func ResizeWindow(width, height uint) {
	if renderWindow != nil {
		windowWidth = width
		windowHeight = height
		renderWindow.SetSize(sf.Vector2u{X: width, Y: height})
	} else {
		logger.Warn(NoWindow)
	}
}

//SetWindowName sets the name of the current window
func SetWindowName(newName string) {
	if renderWindow != nil {
		windowName = newName
	} else {
		logger.Warn(NoWindow)
	}
}
