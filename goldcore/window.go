package goldcore

import (
	"runtime"

	sf "github.com/manyminds/gosfml"
	"github.com/trustmaster/goflow"
)

//Define Window Messages
const (
	//Messages for Output
	WindowCreated = "game window created"
	WindowClosed  = "game window closed"
	//Payload Vector2u
	//In: Resizes the Window
	//Out: New Size of Window
	WindowResized     = "game window resized"
	WindowLostFocus   = "game window lost focus"
	WindowGainedFocus = "game window gained focus"
	//Payload EventTextEntered
	//In: Notifies Text Entered Observers
	//Out: EventTextEntered object
	WindowTextEntered = "game window text entered"
	//Payload EventKey
	//In: Calls KeyPressed Command
	//Out: EventKey object
	WindowKeyPressed = "game window key pressed"
	//Payload EventKey
	//In: Calls KeyReleased Command
	//Out: EventKey object
	WindowKeyReleased = "game window key released"
	//Payload EventMouseWheelMoved
	//In: Notifies MouseWheelMoved Observers
	//Out: EventMouseWheelMoved object
	WindowMouseWheelMoved = "game window mouse gWeel moved"
	//Payload EventMouseButtonWrapper
	//In: Calls MouseButtonPressed Command
	//Out: EventMouseButtonWrapper
	WindowMouseButtonPressed = "game window mouse button pressed"
	//Payload EventMouseButtonWrapper
	//In: Calls MouseButtonReleased Command
	//Out: EventMouseButtonWrapper
	WindowMouseButtonReleased = "game window mouse button released"
	//Payload EventMouseMoved
	//In: Notifies Mouse Moved Observers
	//Out: EventMouseMoved object
	WindowMouseMoved             = "game window mouse moved"
	WindowMouseEntered           = "game window mouse entered"
	WindowMouseLeft              = "game window mouse left"
	WindowJoystickButtonPressed  = "game window joystick button pressed"
	WindowJoystickButtonReleased = "game window joystick button released"
	WindowJoystickMoved          = "game window joystick button Moved"
	WindowJoystickConnected      = "game window joystick connected"
	WindowJoystickDisconnected   = "game window joystick disconnected"

	//Rendering
	WindowStarted   = "game window started"
	WindowStopped   = "game window stopped"
	WindowPaused    = "game window paused"
	WindowCantPause = "game window paused WARNING: Can't pause a closed or stopped window"
	WindowRunning   = "game window running"
	WindowSpinning  = "game window spinning WARNING: Game window is allowed to render, but has not been given the render signal. You should Deactivate or Stop the Game window if you don't want to render"
	//Payload : int. TODO figure out what to do with this int
	WindowNextFrame = "game window next frame"
	WindowRendered  = "game window rendered"
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
//TODO change stopped param to running param
type GameWindow struct {
	renderWindow         *sf.RenderWindow
	renderState          chan int
	renderStateProcessed chan int
	currentState         int
	stopped              bool
	wait                 chan int
	observers            []WindowObserver
	InputSystem          InputSystem
	game                 *Game
	flow.Component
	InputGameMessage  <-chan GameMessage
	OutputGameMessage chan<- GameMessage
}

//GameWindowMessageBufferSize : Number of messages to keep in buffer
const GameWindowMessageBufferSize = 5

//OnInputGameMessage : What to do when a message happens
func (gW *GameWindow) OnInputGameMessage(gM GameMessage) {
	//	fmt.Println("Window Received", gM)
	switch gM.Message {
	//Poll Events
	case WindowClosed:
		//Close Window
		gW.CloseWindow()
	case WindowKeyPressed:
		gW.InputSystem.SetKeyPressed(EventKeyToSFEventKeyPressed(gM.Payload.(EventKey)))
	case WindowKeyReleased:
		gW.InputSystem.SetKeyReleased(EventKeyToSFEventKeyReleased(gM.Payload.(EventKey)))
	case WindowMouseButtonPressed:
		gW.InputSystem.SetMouseButtonPressed(gM.Payload.(EventMouseButtonWrapper).ToSFMLPressed())
	case WindowMouseButtonReleased:
		gW.InputSystem.SetMouseButtonReleased(gM.Payload.(EventMouseButtonWrapper).ToSFMLReleased())
	case WindowMouseMoved:
		gW.InputSystem.SetMouseMove(gM.Payload.(EventMouseMoved).EventMouseMovedToSFML())
	case WindowMouseWheelMoved:
		gW.InputSystem.SetMouseWheelMove(gM.Payload.(EventMouseWheelMoved).EventMouseWheelMovedToSFML())
	case WindowTextEntered:
		gW.InputSystem.SetTextEntered(gM.Payload.(EventTextEntered).ToSFML())

		//Rendering
	case WindowStopped:
		gW.Stop()
	case WindowPaused:
		gW.Deactivate()
	case WindowRunning:
		gW.Start()
	case WindowNextFrame:
		gW.NextFrame(gM.Payload.(int))
	}
	gW.notify(gM)
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
	//fmt.Println(gM)
	gW.OutputGameMessage <- gM
	//fmt.Println("From Window", gM)
	//fmt.Println(<-gW.OutputGameMessage)
	for _, o := range gW.observers {
		o.OnNotify(gM)
	}
}

var activeGameWindow *GameWindow

//NewGameWindow : Creates a new game window. Inactivates any GameWindow, and
//activates the newly created one
//TODO figure out how to maek the WindowCreated Message Work
//Note, I can't create the windows.
func NewGameWindow(width, height uint, name string) *GameWindow {

	gW := &GameWindow{
		renderWindow: sf.NewRenderWindow(sf.VideoMode{Width: width, Height: height, BitsPerPixel: 32}, name, sf.StyleDefault, sf.DefaultContextSettings()),
		InputSystem:  NewInputSystem(),
		observers:    make([]WindowObserver, 0),
	}
	gW.renderWindow.SetActive(false)
	gW.stopped = true
	return gW
}

//PollEvent : Calls handlers for every event since this was last called
//TODO Make sure notify only called once
func (gW *GameWindow) PollEvent() {
	for event := gW.renderWindow.PollEvent(); event != nil; event = gW.renderWindow.PollEvent() {
		switch ev := event.(type) {
		case sf.EventClosed:
			gW.CloseWindow()
			gW.notify(GameMessage{Message: WindowClosed})
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
			gW.notify(GameMessage{Message: WindowMouseButtonPressed, Payload: SFMouseButtonPressedToEventMouseButtonWrapper(event.(sf.EventMouseButtonPressed))})
			gW.InputSystem.SetMouseButtonPressed(ev)
		case sf.EventMouseButtonReleased:
			gW.notify(GameMessage{Message: WindowMouseButtonReleased, Payload: SFMouseButtonReleasedToEventMouseButtonWrapper(event.(sf.EventMouseButtonReleased))})
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
	//fmt.Println("stop")
	gW.renderState <- RenderStopped
	<-gW.renderStateProcessed
	//fmt.Println("end stop")
	close(gW.renderState)
	close(gW.renderStateProcessed)
	gW.notify(GameMessage{Message: WindowStopped})
}

//Deactivate : Pauses but does not stop the window.
func (gW *GameWindow) Deactivate() {
	if !gW.stopped {
		if activeGameWindow == gW {
			activeGameWindow = nil
		}
		//Can't go from stoped state to pause state this way

		gW.renderState <- RenderPaused
		<-gW.renderStateProcessed
		gW.notify(GameMessage{Message: WindowPaused})

	}
}

//Activate : Activates window it will start rendering now.
//If Window was previously stopped, must call start
func (gW *GameWindow) Activate() {
	if !gW.stopped {
		if activeGameWindow != nil {
			activeGameWindow.Deactivate()
		}
		activeGameWindow = gW
		gW.renderState <- RenderRunning
		<-gW.renderStateProcessed
		//fmt.Println("In Activate")
		gW.notify(GameMessage{Message: WindowRunning})
	}
}

//Start : Starts the window after it being stoped
func (gW *GameWindow) Start() {
	if gW.stopped {
		gW.renderState = make(chan int)
		gW.renderStateProcessed = make(chan int)
		go gW.render()
		<-gW.renderStateProcessed
		gW.notify(GameMessage{Message: WindowStarted})

	}
	gW.Activate()
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
	//fmt.Println("\n\nIS WINDOW OPEN", gW.renderWindow.IsOpen())
	//gW.renderState = make(chan int)
	gW.renderStateProcessed <- RenderRunning
	for gW.renderWindow.IsOpen() {
		select {
		case gW.currentState = <-gW.renderState:
			//		fmt.Println("Recieved Render State", gW.currentState)
			switch gW.currentState {
			case RenderStopped:
				gW.stopped = true
				gW.renderWindow.SetActive(false)
				runtime.UnlockOSThread()
				gW.renderStateProcessed <- RenderStopped
				//close(gW.renderState)
				return
			case RenderRunning:
				runtime.LockOSThread()
				gW.renderWindow.SetActive(true)
				gW.renderStateProcessed <- RenderRunning
			case RenderPaused:
				gW.renderWindow.SetActive(false)
				runtime.UnlockOSThread()
				gW.renderStateProcessed <- RenderPaused
			}
		default:
			//Don't Starve :) the scheduler
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
