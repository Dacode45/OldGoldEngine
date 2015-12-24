package goldcore

import (
	"fmt"
	"runtime"

	sf "github.com/manyminds/gosfml"
	"github.com/trustmaster/goflow"
)

//Define Window Messages
var (
	//Messages for Output
	WindowCreated = []byte("game window created")
	WindowClosed  = []byte("game window closed")
	//Payload Vector2u
	//In: Resizes the Window
	//Out: New Size of Window
	WindowResized     = []byte("game window resized")
	WindowLostFocus   = []byte("game window lost focus")
	WindowGainedFocus = []byte("game window gained focus")
	//Payload EventTextEntered
	//In: Notifies Text Entered Observers
	//Out: EventTextEntered object
	WindowTextEntered = []byte("game window text entered")
	//Payload EventKey
	//In: Calls KeyPressed Command
	//Out: EventKey object
	WindowKeyPressed = []byte("game window key pressed")
	//Payload EventKey
	//In: Calls KeyReleased Command
	//Out: EventKey object
	WindowKeyReleased = []byte("game window key released")
	//Payload EventMouseWheelMoved
	//In: Notifies MouseWheelMoved Observers
	//Out: EventMouseWheelMoved object
	WindowMouseWheelMoved = []byte("game window mouse gWeel moved")
	//Payload EventMouseButton
	//In: Calls MouseButtonPressed Command
	//Out: EventMouseButton
	WindowMouseButtonPressed = []byte("game window mouse button pressed")
	//Payload EventMouseButton
	//In: Calls MouseButtonReleased Command
	//Out: EventMouseButton
	WindowMouseButtonReleased = []byte("game window mouse button released")
	//Payload EventMouseMoved
	//In: Notifies Mouse Moved Observers
	//Out: EventMouseMoved object
	WindowMouseMoved             = []byte("game window mouse moved")
	WindowMouseEntered           = []byte("game window mouse entered")
	WindowMouseLeft              = []byte("game window mouse left")
	WindowJoystickButtonPressed  = []byte("game window joystick button pressed")
	WindowJoystickButtonReleased = []byte("game window joystick button released")
	WindowJoystickMoved          = []byte("game window joystick button Moved")
	WindowJoystickConnected      = []byte("game window joystick connected")
	WindowJoystickDisconnected   = []byte("game window joystick disconnected")

	//Rendering
	WindowStopped  = []byte("game window stopped")
	WindowPaused   = []byte("game window paused")
	WindowRunning  = []byte("game window stopped")
	WindowSpinning = []byte("game window spinning WARNING: Game window is allowed to render, but has not been given the render signal. You should Deactivate or Stop the Game window if you don't want to render")
	WindowRendered = []byte("game window rendered")
)

//WindowObserver : Implementations of this interface get an Window event and the event
//Called on every window event. Decide for yourself if you handle it.
type WindowObserver interface {
	OnNotify(GameMessage)
}

const (
	//RenderStopped :  Stopped State
	RenderStopped = 0
	//RenderPaused : Paused State
	RenderPaused = 1
	//RenderRunning : Running State
	RenderRunning = 2
)

//GameWindow : Wrapper around SFML Window. Additional Functionality for Game
type GameWindow struct {
	renderWindow *sf.RenderWindow
	renderState  chan int
	currentState int
	stopped      bool
	wait         chan int
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

	gW := &GameWindow{
		renderWindow: sf.NewRenderWindow(sf.VideoMode{Width: width, Height: height, BitsPerPixel: 32}, name, sf.StyleDefault, sf.DefaultContextSettings()),
		InputSystem:  NewInputSystem(),
	}
	go gW.render()
	return gW
}

//PollEvent : Calls handlers for every event since this was last called
//TODO Make sure notify only called once
func (gW *GameWindow) PollEvent() {
	for event := gW.renderWindow.PollEvent(); event != nil; event = gW.renderWindow.PollEvent() {
		switch ev := event.(type) {
		case sf.EventClosed:
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
			gW.InputSystem.SetTextEntered(ev)
			gW.notify(GameMessage{Message: WindowTextEntered, Payload: SFEventTextEnteredToEventTextEntered(event.(sf.EventTextEntered))})
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
			gW.notify(GameMessage{Message: WindowMouseWheelMoved, Payload: SFEventMouseWheelMovedToEventMouseMoved(event.(sf.EventMouseWheelMoved))})
			gW.InputSystem.SetMouseWheelMove(ev)
		case sf.EventMouseEntered:
			gW.notify(GameMessage{Message: WindowMouseEntered})
		case sf.EventMouseLeft:
			gW.notify(GameMessage{Message: WindowMouseLeft})
		}
	}
}

//TODO Figure out how wait event should work
//Or decide whether it is even needed with the observer pattern

//GetCurrentRenderState : Returns Current RenderSTate
func (gW *GameWindow) GetCurrentRenderState() int {
	return gW.currentState
}

//IsStopped : Checks whether the game is stopped.
//Render State should be RenderStopped if true
func (gW *GameWindow) IsStopped() bool {
	return gW.stopped
}

//Stop : Sets the Window as active inactivating the currently active window
func (gW *GameWindow) Stop() {
	gW.Deactivate()
	gW.renderState <- RenderStopped
}

//Deactivate : Pauses but does not stop the window.
func (gW *GameWindow) Deactivate() {
	if !gW.stopped {
		if activeGameWindow == gW {
			activeGameWindow = nil
		}
		//Can't go from stoped state to pause state this way

		gW.renderState <- RenderPaused
	}
}

//Activate : Activates window it will start rendering now.
//If Window was previously stopped, must call start
func (gW *GameWindow) Activate() {
	if !gW.stopped {
		if activeGameWindow != nil {
			activeGameWindow.Deactivate()
		}

		gW.renderState <- RenderRunning
	}
}

//Start : Starts the window after it being stoped
func (gW *GameWindow) Start() {
	if gW.stopped && gW.IsOpen() {
		go gW.render()
	}
}

//NextFrame : Allows rendering of next frame. Don't know what to do with
//data yet but its there for good measure
func (gW *GameWindow) NextFrame(data int) {
	gW.wait <- data
}

//render : Called on Creation of game window.
//All rendering is done in this thread wheile
//you can pollEvents at anytime. Ensures only
//one opengl context at a time
func (gW *GameWindow) render() {
	gW.stopped = false
	gW.currentState = RenderPaused //Begin in paused state
	gW.renderWindow.SetActive(false)
	for gW.renderWindow.IsOpen() {
		select {
		case gW.currentState = <-gW.renderState:
			switch gW.currentState {
			case RenderStopped:
				gW.stopped = true
				gW.renderWindow.SetActive(false)
				gW.notify(GameMessage{Message: WindowStopped})
				return
			case RenderRunning:
				runtime.LockOSThread()
				gW.renderWindow.SetActive(true)
				gW.notify(GameMessage{Message: WindowRunning})
			case RenderPaused:
				gW.renderWindow.SetActive(false)
				runtime.UnlockOSThread()
				gW.notify(GameMessage{Message: WindowPaused})
			}
		default:
			//Don't Starve :) the schedular
			runtime.Gosched()
			if gW.currentState == RenderPaused {
				break
			}
			select {
			case <-gW.wait:
				gW.renderWindow.Display()
			default:
				//GameWindow Spinning is bad. Address
				gW.notify(GameMessage{Message: WindowSpinning})
			}
			//Actual work
		}
	}
}

//CloseWindow : Inactivates window, sets active window as nil
func (gW *GameWindow) CloseWindow() {
	gW.Stop()
	gW.renderWindow.Close()
	gW.notify(GameMessage{Message: WindowClosed})
}

//SetSize resizes window by width and height
func (gW *GameWindow) SetSize(size Vector2u) {
	gW.renderWindow.SetSize(size.ToSFML())
}

//GetSize : Returns size of GameWindow
func (gW *GameWindow) GetSize() (size Vector2u) {
	return SFVector2uToGEVector2u(gW.renderWindow.GetSize())
}

//SetTitle sets the name of the current window
func (gW *GameWindow) SetTitle(newName string) {
	gW.renderWindow.SetTitle(newName)
}

//GetPosition : Returns Position of Game Window
func (gW *GameWindow) GetPosition() (pos Vector2i) {
	return SFVector2uToGEVector2i(gW.renderWindow.GetPosition())
}

//SetPosition : SetPosition of Game Window
func (gW *GameWindow) SetPosition(pos Vector2i) {
	gW.renderWindow.SetPosition(pos.ToSFML())
}

//IsOpen : Checks if game window is open
func (gW *GameWindow) IsOpen() bool {
	return gW.renderWindow.IsOpen()
}

//TODO Implement the more advanced features fo Render window
