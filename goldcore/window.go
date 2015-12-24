package goldcore

import (
	"fmt"
	"runtime"

	sf "github.com/manyminds/gosfml"
	"github.com/trustmaster/goflow"
)

//Define Window Messages
var (
	WindowCreated                = []byte("game window created")
	WindowClosed                 = []byte("game window closed")
	WindowResized                = []byte("game window resized")
	WindowLostFocus              = []byte("game window lost focus")
	WindowGainedFocus            = []byte("game window gained focus")
	WindowTextEntered            = []byte("game window text entered")
	WindowKeyPressed             = []byte("game window key pressed")
	WindowKeyReleased            = []byte("game window key released")
	WindowMousegWeelMoved        = []byte("game window mouse gWeel moved")
	WindowMouseButtonPressed     = []byte("game window mouse button pressed")
	WindowMouseButtonReleased    = []byte("game window mouse button released")
	WindowMouseMoved             = []byte("game window mouse moved")
	WindowMouseEntered           = []byte("game window mouse entered")
	WindowMouseLeft              = []byte("game window mouse left")
	WindowJoystickButtonPressed  = []byte("game window joystick button pressed")
	WindowJoystickButtonReleased = []byte("game window joystick button released")
	WindowJoystickMoved          = []byte("game window joystick button Moved")
	WindowJoystickConnected      = []byte("game window joystick connected")
	WindowJoystickDisconnected   = []byte("game window joystick disconnected")
)

//WindowObserver : Implementations of this interface get an Window event and the event
//Called on every window event. Decide for yourself if you handle it.
//TODO Create a custom type instead of passing sf.Event
type WindowObserver interface {
	OnNotify(GameMessage)
}

//GameWindow : Wrapper around SFML Window. Additional Functionality for Game
type GameWindow struct {
	renderWindow *sf.RenderWindow
	observers    []WindowObserver
	InputSystem  InputSystem
	game         *Game
	flow.Component
	InputGameMessage  <-chan GameMessage
	OutputGameMessage chan<- GameMessage
}

//OnInputGameMessage : What to do when a message happens
func (gW *GameWindow) OnInputGameMessage(gM GameMessage) {
	fmt.Println(string(gM.Message))

}

//AddObserver : Adds Window Observer to Observer list
func (gW *GameWindow) AddObserver(wO WindowObserver) {
	gW.observers = append(gW.observers, wO)
}

//RemoveObserver : Removes WindowObserver from observer list
func (gW *GameWindow) RemoveObserver(wO WindowObserver) {
	for i, o := range gW.observers {
		if o == wO {
			gW.observers = append(gW.observers[:i], gW.observers[i+1:]...)
			return
		}
	}
}

func (gW *GameWindow) notify(gM GameMessage) {
	for _, o := range gW.observers {
		o.OnNotify(gM)
	}
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
		InputSystem:  NewInputSystem(),
	}
	activeGameWindow = gW
	return gW
}

//PollEvent : Calls handlers for every event since this was last called
func (gW *GameWindow) PollEvent() {
	for event := gW.renderWindow.PollEvent(); event != nil; event = gW.renderWindow.PollEvent() {
		switch ev := event.(type) {
		case sf.EventClosed:
			gW.notify(GameMessage{Message: WindowClosed})
			gW.CloseWindow()
		case sf.EventLostFocus:
			gW.notify(GameMessage{Message: WindowLostFocus})
		case sf.EventGainedFocus:
			gW.notify(GameMessage{Message: WindowLostFocus})
		case sf.EventResized:
			gW.notify(GameMessage{Message: WindowResized, Payload: Vector2u{X: event.(sf.EventResized).Width, Y: event.(sf.EventResized).Height}})
		case sf.EventJoystickButtonPressed:
			gW.notify(GameMessage{Message: WindowJoystickButtonPressed})
		case sf.EventJoystickButtonReleased:
			gW.notify(GameMessage{Message: WindowJoystickButtonReleased})
		case sf.EventJoystickConnected:
			gW.notify(GameMessage{Message: WindowJoystickConnected})
		case sf.EventJoystickDisconnected:
			gW.notify(GameMessage{Message: WindowJoystickDisconnected})
		case sf.EventJoystickMoved:
			gW.notify(GameMessage{Message: WindowJoystickMoved})
		case sf.EventKeyPressed:
			gW.notify(GameMessage{Message: WindowKeyPressed, Payload: SFEventKeyPressedToEventKey(event.(sf.EventKeyPressed))})
			gW.InputSystem.SetKeyPressed(ev)
		case sf.EventKeyReleased:
			gW.notify(GameMessage{Message: WindowKeyReleased, Payload: SFEventKeyReleasedToEventKey(event.(sf.EventKeyReleased))})
			gW.InputSystem.SetKeyReleased(ev)
		case sf.EventTextEntered:
			//TODO Implement
			gW.notify(GameMessage{Message: WindowTextEntered})
		case sf.EventMouseButtonPressed:
			gW.notify(GameMessage{Message: WindowMouseButtonPressed, Payload: SFMouseButtonPressedToEventMouseButton(event.(sf.EventMouseButtonPressed))})
			gW.InputSystem.SetMouseButtonPressed(ev)
		case sf.EventMouseButtonReleased:
			gW.notify(GameMessage{Message: WindowMouseButtonReleased, Payload: SFMouseButtonReleasedToEventMouseButton(event.(sf.EventMouseButtonReleased))})
			gW.InputSystem.SetMouseButtonReleased(ev)
		case sf.EventMouseMoved:
			gW.notify(GameMessage{Message: WindowMouseMoved, Payload: SFEventMouseMovedToEventMouseMoved(event.(sf.EventMouseMoved))})
			gW.InputSystem.SetMouseMove(ev)
		case sf.EventMouseWheelMoved:
			gW.notify(GameMessage{Message: WindowMousegWeelMoved, Payload: SFEventMouseWheelMovedToEventMouseMoved(event.(sf.EventMouseWheelMoved))})
			gW.InputSystem.SetMouseWheelMove(ev)
		case sf.EventMouseEntered:
			gW.notify(GameMessage{Message: WindowMouseEntered})
		case sf.EventMouseLeft:
			gW.notify(GameMessage{Message: WindowMouseLeft})
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
