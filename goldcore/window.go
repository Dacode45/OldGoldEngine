package goldcore

import (
	"fmt"
	"runtime"

	sf "github.com/manyminds/gosfml"
)

//WindowObserver : Implementations of this interface get an Window event and the event
type WindowObserver interface {
	OnNotify(eventType sf.EventType, event sf.Event)
}

//WindowHandler : Notifies Window Observers of events
type WindowHandler struct {
	observers []WindowObserver
}

func (wH *WindowHandler) AddObserver(wO WindowObserver) {
	wH.observers = append(observers, wO)
	return &wO
}

func (wH *WindowHandler) RemoveObserver(wO *WindowObserver) {
	for i, o := range wH.observers {
		if &o == wO {
			wH.observers = append(wH.observers[:i], wH.observers[i+1:])
			return
		}
	}
}

func (wH *WindowHandler) notify(eventType sf.EventType, event sf.Event) {
	for _, o := range wH.observers {
		o.OnNotify(eventType, event)
	}
}

//GameWindow : Wrapper around SFML Window. Additional Functionality for Game
type GameWindow struct {
	renderWindow  *sf.RenderWindow
	WindowHandler WindowHandler
	game          *Game
}

var activeGameWindow *GameWindow

//NewGameWindow : Creates a new game window. Inactivates any GameWindow, and
//activates the newly created one
func NewGameWindow(width, height uint, name string) *GameWindow {
	if activeGameWindow != nil {
		activeGameWindow.SetActive(false)
	}
	runtime.LockOSThread()
	gW := &GameWindow{
		renderWindow: sf.NewRenderWindow(sf.VideoMode{Width: windowWidth, Height: windowHeight, BitsPerPixel: 32}, name, sf.StyleDefault, sf.DefaultContextSettings()),
	}
	activeGameWindow = gW
	return gW
}

//PollEvent : Calls handlers for every event since this was last called
func (gW *GameWindow) PollEvent() {
	for event := gW.renderWindow.PollEvent(); event != nil; event = gW.renderWindow.PollEvent() {
		switch ev := event.(type) {
		case sf.EventTypeClosed:
			gW.WindowHandler.notify(sf.EventTypeClosed, event)
			gW.CloseWindow()
		case sf.EventTypeLostFocus:
			gW.WindowHandler.notify(sf.EventTypeLostFocus, event)
		case sf.EventTypeGainedFocus:
			gW.WindowHandler.notify(sf.EventTypeGainedFocus, event)
		case sf.EventTypeResized:
			gW.WindowHandler.notify(sf.EventTypeResized, event)
		case sf.EventTypeJoystickButtonPressed:
			//TODO Implement
			fmt.Println()
		case sf.EventTypeJoystickButtonReleased:
			//TODO Implement
			fmt.Println()
		case sf.EventTypeJoystickConnected:
			//TODO Implement
			fmt.Println()
		case sf.EventTypeJoystickDisconnected:
			//TODO Implement
			fmt.Println()
		case sf.EventTypeJoystickMoved:
			//TODO Implement
			fmt.Println()
		case sf.EventTypeKeyPressed:
			//TODO Implement
			fmt.Println()
		case sf.EventTypeKeyReleased:
			//TODO Implement
			fmt.Println()
		case sf.EventTypeTextEntered:
			//TODO Implement
			fmt.Println()
		case sf.EventTypeMouseButtonPressed:
			//TODO Implement
			fmt.Println()
		case sf.EventTypeMouseButtonReleased:
			//TODO Implement
			fmt.Println()
		case sf.EventTypeMouseEntered:
			//TODO Implement
			fmt.Println()
		case sf.EventTypeMouseLeft:
			//TODO Implement
			fmt.Println()
		case sf.EventTypeMouseMoved:
			//TODO Implement
			fmt.Println()

		}
	}
}

//SetActive : Sets the Window as active inactivating the currently active window
func (gW *GameWindow) SetActive(active bool) {
	if activeGameWindow != nil {
		activeGameWindow.SetActive(false)
	}
	if active {
		runtime.LockOSThread()
		activeGameWindow = gW
	} else {
		activeGameWindow = nil
	}
}

//ResizeWindow resizes window by width and height
func (gW *GameWindow) ResizeWindow(width, height uint) {
	if gW.renderWindow != nil {
		renderWindow.SetSize(sf.Vector2u{X: width, Y: height})
	}
}

//SetWindowName sets the name of the current window
func (gW *GameWindow) SetWindowName(newName string) {
	if gW.renderWindow != nil {
		gW.windowName = newName
	}
}

//CloseWindow : Inactivates window, sets active window as nil
func (gW *GameWindow) CloseWindow() {
	gW.SetActive(false)
	gW.renderWindow.Close()
	if activeGameWindow == gw {
		activeGameWindow = nil
	}
}
