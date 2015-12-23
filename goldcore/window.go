package goldcore

import (
	"fmt"
	"runtime"

	sf "github.com/manyminds/gosfml"
)

//WindowObserver : Implementations of this interface get an Window event and the event
//Called on every window event. Decide for yourself if you handle it.
//TODO Create a custom type instead of passing sf.Event
type WindowObserver interface {
	OnNotify(eventType sf.EventType, event sf.Event)
}

//WindowHandler : Notifies Window Observers of events
type WindowHandler struct {
	observers []WindowObserver
}

//AddObserver : Adds Window Observer to Observer list
func (wH *WindowHandler) AddObserver(wO WindowObserver) {
	wH.observers = append(wH.observers, wO)
}

//RemoveObserver : Removes WindowObserver from observer list
func (wH *WindowHandler) RemoveObserver(wO WindowObserver) {
	for i, o := range wH.observers {
		if o == wO {
			wH.observers = append(wH.observers[:i], wH.observers[i+1:]...)
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
		renderWindow: sf.NewRenderWindow(sf.VideoMode{Width: width, Height: height, BitsPerPixel: 32}, name, sf.StyleDefault, sf.DefaultContextSettings()),
	}
	activeGameWindow = gW
	return gW
}

//PollEvent : Calls handlers for every event since this was last called
func (gW *GameWindow) PollEvent() {
	for event := gW.renderWindow.PollEvent(); event != nil; event = gW.renderWindow.PollEvent() {
		switch ev := event.(type) {
		case sf.EventClosed:
			gW.WindowHandler.notify(sf.EventTypeClosed, event)
			gW.CloseWindow()
		case sf.EventLostFocus:
			gW.WindowHandler.notify(sf.EventTypeLostFocus, event)
		case sf.EventGainedFocus:
			gW.WindowHandler.notify(sf.EventTypeGainedFocus, event)
		case sf.EventResized:
			gW.WindowHandler.notify(sf.EventTypeResized, event)
		case sf.EventJoystickButtonPressed:
			//TODO Implement
			fmt.Println()
		case sf.EventJoystickButtonReleased:
			//TODO Implement
			fmt.Println()
		case sf.EventJoystickConnected:
			//TODO Implement
			fmt.Println()
		case sf.EventJoystickDisconnected:
			//TODO Implement
			fmt.Println()
		case sf.EventJoystickMoved:
			//TODO Implement
			fmt.Println()
		case sf.EventKeyPressed:
			go SetKeyPressed(ev)
			fmt.Println()
		case sf.EventKeyReleased:
			go SetKeyReleased(ev)
			fmt.Println()
		case sf.EventTextEntered:
			//TODO Implement
			fmt.Println()
		case sf.EventMouseButtonPressed:
			SetMouseButtonPressed(ev)
		case sf.EventMouseButtonReleased:
			SetMouseButtonReleased(ev)
		case sf.EventMouseMoved:
			globalMouseMovedHandler.SetMouseMove(ev)
		case sf.EventMouseWheelMoved:
			//TODO Implement
			fmt.Println()
		case sf.EventMouseEntered:
			gW.WindowHandler.notify(sf.EventTypeMouseEntered, event)
		case sf.EventMouseLeft:
			gW.WindowHandler.notify(sf.EventTypeMouseLeft, event)
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
		gW.renderWindow.SetSize(sf.Vector2u{X: width, Y: height})
	}
}

//SetWindowTitle sets the name of the current window
func (gW *GameWindow) SetWindowTitle(newName string) {
	if gW.renderWindow != nil {
		gW.renderWindow.SetTitle(newName)
	}
}

//CloseWindow : Inactivates window, sets active window as nil
func (gW *GameWindow) CloseWindow() {
	gW.SetActive(false)
	gW.renderWindow.Close()
	if activeGameWindow == gW {
		activeGameWindow = nil
	}
}
