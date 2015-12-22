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

type WindowEvent sf.EventType

const (
  WindowClosed = sf.EventTypeClosed
  WindowResized = sf.EventTypeResized
  WindowLostFocus = sf.EventTypeLostFocus
  WindowGainedFocus = sf.EventTypeGainedFocus
)

type WindowObserver interface{
  OnNotify(event WindowEvent)
}

var windowObservers = make([]*WindowObserver,10)

func AddWindowObserver(observer *WindowObserver){
  windowObservers = append(windowObservers, observer)
}

func RemoveWindowObserver(observer *WindowObserver) bool{
  for i, o := range windowObservers{
    if o == observer{
      windowObservers = append(windowObservers[:i], windowObservers[i+1:])
      return true
    }
  }
  return false
}

func RemoveAllWindowObservers(){
  windowObservers := make([]*WindowObserver, 10)
}

func notifyWindowObservers(event WindowEvent){
  for _, o := range windowObservers {
    o.OnNotify(event);
  }
}


//CreateWindow creates a Game Window, all Inputs are taken from it. Calling it
//destroys previous window
func CreateWindow(width, height uint, name string) {
	DestroyWindow()
	windowWidth = width
	windowHeight = height
	windowName = name
	renderWindow = sf.NewRenderWindow(sf.VideoMode{Width: windowWidth, Height: windowHeight, BitsPerPixel: 32}, windowName, sf.StyleDefault, sf.DefaultContextSettings())
}

//DestroyWindow: Closes current window if available and removes all window observers
//TODO Decide if you want to keep window observers since a new window may be created
func DestroyWindow(){
  if renderWindow != nil {
    notifyWindowObservers(WindowClosed)
    RemoveAllWindowObservers()
		renderWindow.Close()
		logger.Warn(fmt.Sprintf("Window %s has been closed", windowName))
		renderWindow = nil
	}
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

//PollEvent: Checks all events that have been queued since the last call, and passes
//them to the approprite handler. For Key Mouse and joystick check input.go.
//Handles all window events natively.
func PollEvent(){
  for event := renderWindow.PollEvent(); event != nil; event = renderWindow.PollEvent(){
    switch ev = event.(type){
    case: sf.EventTypeClosed:
      DestroyWindow()
    case: sf.EventTypeGainedFocus:
    case: sf.EventTypeJoystickButtonPressed:
    case: sf.EventTypeJoystickButtonReleased:
    case: sf.EventTypeJoystickConnected:
    case: sf.EventTypeJoystickMoved:
    case: sf.EventTypeKeyPressed:
			SetKeyPressed()
    case: sf.EventTypeKeyReleased:
    case: sf.EventTypeLostFocus:
    case: sf.EventTypeMouseButtonPressed:
    case: sf.EventTypeMouseButtonReleased:
    case: sf.EventTypeMouseEntered:
    case: sf.EventTypeMouseLeft:
    case: sf.EventTypeMouseMoved:
    case: sf.EventTypeMouseWheelMoved:
    case: sf.EventTypeResized:
      ResizeWindow(ev.Width, ev.Height)
    case: sf.EventTypeTextEntered:
    }
  }
}
