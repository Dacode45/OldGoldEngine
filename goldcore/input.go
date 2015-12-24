package goldcore

import (
	"fmt"

	sf "github.com/manyminds/gosfml"
)

//Sets Flags For Player Input. Defines interfaces for interpreting input
//Input can either be goten directly, or (the better way) use input commands
//tied directly to when a button is pressed

//KeyCommand :Create a struct with an Execute function. Allows Command
//to have state but not really recommended or necessary
type KeyCommand func()

//EventKey : Gives You a keycode + all modifier keys.
//In order to check modifiers only set Code to ModifierCode
type EventKey struct {
	Code    KeyCode //< Code of the key that has been pressed
	Alt     KeyCode //< Is the Alt key pressed?
	Control KeyCode //< Is the Control key pressed?
	Shift   KeyCode //< Is the Shift key pressed?
	System  KeyCode //< Is the System key pressed?
	Pressed bool    //< Is the button being pressed (true) or released (false)
}

func (eK *EventKey) String() string {
	return fmt.Sprintf("Code: %d, Alt: %d, Control: %d, Shift: %d, System: %d, Pressed: %t",
		eK.Code, eK.Alt, eK.Control, eK.Shift, eK.System, eK.Pressed)
}

//ModifierKeyCode : Used to check if for modifiers only. Example inorder to check if only the shift
//key was pressed make an EventKey with Code set to ModifierCode and Shift set to true
const ModifierKeyCode = -1

//KeyboardHandler : Will IMMEDIATELY call a function once key is pressed as long
//as the handler or the set is active. It is up to the user to ensure syncronization
//KeyboardHandlers only work as a part of a key set
type KeyboardHandler struct {
	keysPressed map[EventKey]KeyCommand
	key         uint
}

//AddEventKey : Binds Command to Event Key object.
func (kH *KeyboardHandler) AddEventKey(eK EventKey, kC KeyCommand) {
	kH.keysPressed[eK] = kC
}

//RemoveEventKey : Removes Event Key and its effects
func (kH *KeyboardHandler) RemoveEventKey(eK EventKey) {
	delete(kH.keysPressed, eK)
}

//check : Internally called. Will call Key Command if available
func (kH *KeyboardHandler) check(eK *EventKey) {
	if cmd, ok := kH.keysPressed[*eK]; ok {
		//fmt.Println("hello")

		go cmd()
	}
}

//NewKeyboardHandler : Creates a new Keyboard Handler
func NewKeyboardHandler() KeyboardHandler {
	kH := KeyboardHandler{keysPressed: make(map[EventKey]KeyCommand)}
	return kH
}

//KeyboardSet :Set of actual KeyboardHandlers called on input. You can enable
//and disable keyboardhandlers rather than constantly adding and removing them
type KeyboardSet struct {
	keyboardHandlers []KeyboardHandler
	numActive        int
	key              uint
	nextIndex        uint
}

//NewKeyboardSet : Creates a new KeyboardSet
func NewKeyboardSet() KeyboardSet {
	kS := KeyboardSet{keyboardHandlers: make([]KeyboardHandler, 0)}
	return kS
}

//MouseButtonIsActive : Check if the keyboardhandler with this key is active
func (kS *KeyboardSet) MouseButtonIsActive(key uint) bool {
	for i, k := range kS.keyboardHandlers {
		if k.key == key {
			return i < kS.numActive
		}
	}
	return false
}

//AddHandler : Give a Handler and get a pointer back to the handler in memory
//Automatically makes it active so set active will need to be called afterwards
func (kS *KeyboardSet) AddHandler(kH KeyboardHandler) (key uint) {
	kS.keyboardHandlers = append(kS.keyboardHandlers, KeyboardHandler{})
	copy(kS.keyboardHandlers[kS.numActive+1:], kS.keyboardHandlers[kS.numActive:])
	kH.key = kS.nextIndex
	kS.nextIndex++
	kS.keyboardHandlers[kS.numActive] = kH
	kS.numActive++
	return kH.key

}

//RemoveHandler : Removes KeyboardHandler
func (kS *KeyboardSet) RemoveHandler(key uint) {
	//This is much faster if this array is actually in memory
	for i, k := range kS.keyboardHandlers {
		if k.key == key {
			kS.keyboardHandlers = append(kS.keyboardHandlers[:i], kS.keyboardHandlers[i+1:]...)
			if i < kS.numActive {

				kS.numActive--
			}
			return
		}
	}
}

//SetActive : Makes the given keyboardHandler active
func (kS *KeyboardSet) SetActive(key uint, active bool) {
	for i, k := range kS.keyboardHandlers {
		if k.key == key {
			//only do something if different
			if (i < kS.numActive) != active {
				if !(i < kS.numActive) { //if active put it in the active part of the array
					//fmt.Println("hello")
					temp := kS.keyboardHandlers[kS.numActive]
					kS.keyboardHandlers[kS.numActive] = k
					kS.keyboardHandlers[i] = temp
					kS.numActive++
				} else { //inactive move to end of array
					kS.keyboardHandlers = append(kS.keyboardHandlers[:i], kS.keyboardHandlers[i+1:]...)
					kS.keyboardHandlers = append(kS.keyboardHandlers, k) //TODO there has to be a better way
					kS.numActive--
					//fmt.Println(kS.keyboardHandlers)
				}
			}
			return
		}
	}
}

//check : internally called. Calles KeybaordHandler check for every active one
func (kS *KeyboardSet) check(ek *EventKey) {
	for i := 0; i < kS.numActive; i++ {
		kS.keyboardHandlers[i].check(ek)
	}
}

//KeyboardDispatcher : Handles keyboardSets for an InputSystem
type KeyboardDispatcher struct {
	numActiveKeySets int
	nextKeySetIndex  uint
	keyboardSets     []KeyboardSet
}

//NewKeyboardDispatcher : Creates a new KeyboardDispatcher
func NewKeyboardDispatcher() KeyboardDispatcher {
	return KeyboardDispatcher{keyboardSets: make([]KeyboardSet, 0)}
}

//AddKeyboardSet Registers this set globally. You will no longer have to add
//and re-add sets. Just enable and disable them
func (kD *KeyboardDispatcher) AddKeyboardSet(kS KeyboardSet) (key uint) {

	kD.keyboardSets = append(kD.keyboardSets, KeyboardSet{})
	copy(kD.keyboardSets[kD.numActiveKeySets+1:], kD.keyboardSets[kD.numActiveKeySets:])
	kS.key = kD.nextKeySetIndex
	kD.nextKeySetIndex++
	kD.keyboardSets[kD.numActiveKeySets] = kS
	kD.numActiveKeySets++

	return kS.key
}

//RemoveKeyboardSet : Unregister KeyboardSet globally
func (kD *KeyboardDispatcher) RemoveKeyboardSet(key uint) {
	for i, k := range kD.keyboardSets {
		if k.key == key {
			kD.keyboardSets = append(kD.keyboardSets[:i], kD.keyboardSets[i+1:]...)
			if i < kD.numActiveKeySets {

				kD.numActiveKeySets--
			}
			return
		}
	}
}

//SetKeyboardSetActive : Enable or disable KeybaordSet
func (kD *KeyboardDispatcher) SetKeyboardSetActive(key uint, active bool) {
	for i, k := range kD.keyboardSets {
		if k.key == key {
			//only do something if different
			if (i < kD.numActiveKeySets) != active {
				if !(i < kD.numActiveKeySets) { //if inactive put it in the active part of the array
					temp := kD.keyboardSets[kD.numActiveKeySets]
					kD.keyboardSets[kD.numActiveKeySets] = k
					kD.keyboardSets[i] = temp
					kD.numActiveKeySets++
				} else { //inactive move to end of array
					kD.keyboardSets = append(kD.keyboardSets[:i], kD.keyboardSets[i+1:]...)
					kD.keyboardSets = append(kD.keyboardSets, k)
					kD.numActiveKeySets--
				}
			}
			return
		}
	}
}

//CheckKeyboardSet : Calls the keybaord command for all associated event keys
//Passes pointer to event to the actual handlers.
func (kD *KeyboardDispatcher) CheckKeyboardSet(eK *EventKey) {

	for i := 0; i < kD.numActiveKeySets; i++ {

		kD.keyboardSets[i].check(eK)
	}
}

//SFEventKeyPressedToEventKey : Converts Sf Key to Event Key
func SFEventKeyPressedToEventKey(event sf.EventKeyPressed) EventKey {
	return EventKey{
		Code:    KeyCode(event.Code),
		Alt:     KeyCode(event.Alt),
		Control: KeyCode(event.Control),
		Shift:   KeyCode(event.Shift),
		System:  KeyCode(event.System),
		Pressed: true}
}

//SFEventKeyReleasedToEventKey : Convert Sf Key to Event Key
func SFEventKeyReleasedToEventKey(event sf.EventKeyReleased) EventKey {
	return EventKey{
		Code:    KeyCode(event.Code),
		Alt:     KeyCode(event.Alt),
		Control: KeyCode(event.Control),
		Shift:   KeyCode(event.Shift),
		System:  KeyCode(event.System),
		Pressed: false}
}

//EventKeyToSFEventKeyPressed : Convert Event key to sf.EventKey
//Don't recommend using this function outside of situations where you want to
//simulate keyboard input
func EventKeyToSFEventKeyPressed(event EventKey) sf.EventKeyPressed {
	return sf.EventKeyPressed{
		Code:    sf.KeyCode(event.Code),
		Alt:     int(event.Alt),
		Control: int(event.Control),
		Shift:   int(event.Shift),
		System:  int(event.System),
	}
}

//EventKeyToSFEventKeyReleased : Convert Event key to sf.EventKey
//Don't recommend using this function outside of situations where you want to
//simulate keyboard input
func EventKeyToSFEventKeyReleased(event EventKey) sf.EventKeyReleased {
	return sf.EventKeyReleased{
		Code:    sf.KeyCode(event.Code),
		Alt:     int(event.Alt),
		Control: int(event.Control),
		Shift:   int(event.Shift),
		System:  int(event.System),
	}
}

//////////////////////////////////////////////////////
///Mouse Button

//MouseButtonCommand : What is called on mouse clicked
type MouseButtonCommand func(button MouseButton, x int, y int)

//EventMouseButton : Called once when mouse is clicked, called again when mouse stops clicking
type EventMouseButton struct {
	Button  MouseButton //< Code of the button that has been pressed
	Clicked bool        //< Whether Button is Pressed or Released
}

//SFMouseButtonPressedToEventMouseButton sfml to Mouse buton
func SFMouseButtonPressedToEventMouseButton(eM sf.EventMouseButtonPressed) EventMouseButton {
	return EventMouseButton{Button: MouseButton(eM.Button), Clicked: true}
}

//SFMouseButtonReleasedToEventMouseButton sfml to mouse button
func SFMouseButtonReleasedToEventMouseButton(eM sf.EventMouseButtonReleased) EventMouseButton {
	return EventMouseButton{Button: MouseButton(eM.Button), Clicked: false}
}

//EventMouseButtonToSFMouseButtonPressed mousebutton to sfml
func EventMouseButtonToSFMouseButtonPressed(eM EventMouseButton, x, y int) sf.EventMouseButtonPressed {
	return sf.EventMouseButtonPressed{Button: sf.MouseButton(eM.Button), X: x, Y: y}
}

//EventMouseButtonToSFMouseButtonReleased mousebutton to sfml
func EventMouseButtonToSFMouseButtonReleased(eM EventMouseButton, x, y int) sf.EventMouseButtonReleased {
	return sf.EventMouseButtonReleased{Button: sf.MouseButton(eM.Button), X: x, Y: y}
}

//MouseButtonHandler : TODO implement
type MouseButtonHandler struct {
	buttonsClicked map[EventMouseButton]MouseButtonCommand
	key            uint
}

//AddMouseButton : Adds command
func (mH *MouseButtonHandler) AddMouseButton(eM EventMouseButton, eC MouseButtonCommand) {
	mH.buttonsClicked[eM] = eC
}

//RemoveMouseButton : removes command
func (mH *MouseButtonHandler) RemoveMouseButton(eM EventMouseButton) {
	delete(mH.buttonsClicked, eM)
}

//check : calls command
func (mH *MouseButtonHandler) check(eM *EventMouseButton, x, y int) {
	if cmd, ok := mH.buttonsClicked[*eM]; ok {
		go cmd(eM.Button, x, y)
	}
}

//NewMouseButtonHandler : Makes new Button Handler
func NewMouseButtonHandler() MouseButtonHandler {
	return MouseButtonHandler{buttonsClicked: make(map[EventMouseButton]MouseButtonCommand)}
}

//MouseButtonSet : Same as a KeyboardSet
type MouseButtonSet struct {
	buttonHandlers []MouseButtonHandler
	numActive      int
	key            uint
	nextIndex      uint
}

//NewMouseButtonSet : Creates a new mouse button set
func NewMouseButtonSet() MouseButtonSet {
	return MouseButtonSet{buttonHandlers: make([]MouseButtonHandler, 0)}
}

//MouseButtonIsActive : Checks if key is active
func (mS *MouseButtonSet) MouseButtonIsActive(key uint) bool {
	for i := 0; i < mS.numActive; i++ {
		if mS.buttonHandlers[i].key == key {
			return true
		}
	}
	return false
}

//AddHandler : Give a Handler and get a pointer back to the handler in memory
//Automatically makes it active so set active will need to be called afterwards
func (mS *MouseButtonSet) AddHandler(mH MouseButtonHandler) (key uint) {
	mS.buttonHandlers = append(mS.buttonHandlers, MouseButtonHandler{})
	copy(mS.buttonHandlers[mS.numActive+1:], mS.buttonHandlers[mS.numActive:])
	mH.key = mS.nextIndex
	mS.nextIndex++
	mS.buttonHandlers[mS.numActive] = mH
	mS.numActive++
	return mH.key

}

//RemoveHandler : Removes MouseButtonHandler
func (mS *MouseButtonSet) RemoveHandler(key uint) {
	//This is much faster if this array is actually in memory
	for i, k := range mS.buttonHandlers {
		if k.key == key {
			mS.buttonHandlers = append(mS.buttonHandlers[:i], mS.buttonHandlers[i+1:]...)
			if i < mS.numActive {

				mS.numActive--
			}
			return
		}
	}
}

//SetActive : Makes the given MouseButtonHandler active
func (mS *MouseButtonSet) SetActive(key uint, active bool) {
	for i, k := range mS.buttonHandlers {
		if k.key == key {
			//only do something if different
			if (i < mS.numActive) != active {
				if !(i < mS.numActive) { //if active put it in the active part of the array
					//fmt.Println("hello")
					temp := mS.buttonHandlers[mS.numActive]
					mS.buttonHandlers[mS.numActive] = k
					mS.buttonHandlers[i] = temp
					mS.numActive++
				} else { //inactive move to end of array
					mS.buttonHandlers = append(mS.buttonHandlers[:i], mS.buttonHandlers[i+1:]...)
					mS.buttonHandlers = append(mS.buttonHandlers, k) //TODO there has to be a better way
					mS.numActive--
					//fmt.Println(mS.buttonHandlers)
				}
			}
			return
		}
	}
}

//check : internally called. Calles KeybaordHandler check for every active one
func (mS *MouseButtonSet) check(ek *EventMouseButton, x, y int) {
	for i := 0; i < mS.numActive; i++ {
		mS.buttonHandlers[i].check(ek, x, y)
	}
}

//MouseButtonDispatcher : Handles MouseButtons for an Input System
type MouseButtonDispatcher struct {
	numActiveMouseButtonSets int
	nextMouseButtonSetIndex  uint
	mouseButtonSets          []MouseButtonSet
}

//NewMouseButtonDispatcer : Create a new MouseButtonDispatcher
func NewMouseButtonDispatcer() MouseButtonDispatcher {
	//DONOT CHANGE INITIAL SIZE FROM ZERO
	return MouseButtonDispatcher{mouseButtonSets: make([]MouseButtonSet, 0)}
}

//AddMouseButtonSet Registers this set globally. You will no longer have to add
//and re-add sets. Just enable and disable them
func (mD *MouseButtonDispatcher) AddMouseButtonSet(mS MouseButtonSet) (key uint) {

	mD.mouseButtonSets = append(mD.mouseButtonSets, MouseButtonSet{})
	copy(mD.mouseButtonSets[mD.numActiveMouseButtonSets+1:], mD.mouseButtonSets[mD.numActiveMouseButtonSets:])
	mS.key = mD.nextMouseButtonSetIndex
	mD.nextMouseButtonSetIndex++
	mD.mouseButtonSets[mD.numActiveMouseButtonSets] = mS
	mD.numActiveMouseButtonSets++
	return mS.key
}

//RemoveMouseButtonSet : Unregister MouseButtonSet globally
func (mD *MouseButtonDispatcher) RemoveMouseButtonSet(key uint) {
	for i, k := range mD.mouseButtonSets {
		if k.key == key {
			mD.mouseButtonSets = append(mD.mouseButtonSets[:i], mD.mouseButtonSets[i+1:]...)
			if i < mD.numActiveMouseButtonSets {

				mD.numActiveMouseButtonSets--
			}
			return
		}
	}
}

//SetMouseButtonSetActive : Enable or disable KeybaordSet
func (mD *MouseButtonDispatcher) SetMouseButtonSetActive(key uint, active bool) {
	for i, k := range mD.mouseButtonSets {
		if k.key == key {
			//only do something if different
			if (i < mD.numActiveMouseButtonSets) != active {
				if !(i < mD.numActiveMouseButtonSets) { //if inactive put it in the active part of the array
					temp := mD.mouseButtonSets[mD.numActiveMouseButtonSets]
					mD.mouseButtonSets[mD.numActiveMouseButtonSets] = k
					mD.mouseButtonSets[i] = temp
					mD.numActiveMouseButtonSets++
				} else { //inactive move to end of array
					mD.mouseButtonSets = append(mD.mouseButtonSets[:i], mD.mouseButtonSets[i+1:]...)
					mD.mouseButtonSets = append(mD.mouseButtonSets, k)
					mD.numActiveMouseButtonSets--
				}
			}
			return
		}
	}
}

//CheckMouseButtonSet : Calls the keybaord command for all associated event keys
//Passes pointer to event to the actual handlers.
func (mD *MouseButtonDispatcher) CheckMouseButtonSet(eM *EventMouseButton, x, y int) {
	for i := 0; i < mD.numActiveMouseButtonSets; i++ {
		mD.mouseButtonSets[i].check(eM, x, y)
	}
}

//////////////////////////////////////////////////////////
//MouseMove
//Uses Observer rather than the Command Pattern

//EventMouseMoved : Called once when mouse is clicked, called again when mouse stops clicking
type EventMouseMoved struct {
	X int //< X position of the mouse pointer, relative to the left of the owner window
	Y int //< Y position of the mouse pointer, relative to the top of the owner window
}

//EventMouseMovedToSFML : Converts EventMouseMoved to the sfml version
func (eM *EventMouseMoved) EventMouseMovedToSFML() sf.EventMouseMoved {
	return sf.EventMouseMoved{X: eM.X, Y: eM.Y}
}

//SFEventMouseMovedToEventMouseMoved sfml to Mouse buton
func SFEventMouseMovedToEventMouseMoved(eM sf.EventMouseMoved) EventMouseMoved {
	return EventMouseMoved{X: eM.X, Y: eM.Y}
}

//MouseMoveObserver : what is called on mouse move
type MouseMoveObserver interface {
	OnMouseMove(event EventMouseMoved)
}

//MouseMovedHandler : Observer pattern
type MouseMovedHandler struct {
	observers []MouseMoveObserver
}

//NewMouseMovedHandler : Returns MouseMoveHandler
func NewMouseMovedHandler() MouseMovedHandler {
	return MouseMovedHandler{observers: make([]MouseMoveObserver, 0)}
}

//AddMouseMoveObserver : Adds command
func (mH *MouseMovedHandler) AddMouseMoveObserver(mO MouseMoveObserver) {
	mH.observers = append(mH.observers, mO)
}

//RemoveMouseMoveObserver : removes command
func (mH *MouseMovedHandler) RemoveMouseMoveObserver(mO MouseMoveObserver) {
	for i, o := range mH.observers {
		if o == mO {
			mH.observers = append(mH.observers[:i], mH.observers[i+1:]...)
			return
		}
	}
}

func (mH *MouseMovedHandler) notify(eM EventMouseMoved) {
	for _, o := range mH.observers {
		go o.OnMouseMove(eM)
	}
}

///////////////////////////////////////////////
//Mouse Wheel Scroll

//Uses Observer rather than the Command Pattern

//EventMouseWheelMoved : Called once when mouse is clicked, called again when mouse stops clicking
type EventMouseWheelMoved struct {
	Delta int
	X     int //< X position of the mouse pointer, relative to the left of the owner window
	Y     int //< Y position of the mouse pointer, relative to the top of the owner window
}

//EventMouseWheelMovedToSFML : Converts EventMouseMoved to the sfml version
func (eM *EventMouseWheelMoved) EventMouseWheelMovedToSFML() sf.EventMouseWheelMoved {
	return sf.EventMouseWheelMoved{X: eM.X, Y: eM.Y, Delta: eM.Delta}
}

//SFEventMouseWheelMovedToEventMouseMoved sfml to Mouse buton
func SFEventMouseWheelMovedToEventMouseMoved(eM sf.EventMouseWheelMoved) EventMouseWheelMoved {
	return EventMouseWheelMoved{X: eM.X, Y: eM.Y, Delta: eM.Delta}
}

//MouseWheelMoveObserver : what is called on mouse move
type MouseWheelMoveObserver interface {
	OnMouseMove(EventMouseWheelMoved)
}

//MouseWheelMovedHandler : TODO implement
type MouseWheelMovedHandler struct {
	observers []MouseWheelMoveObserver
}

//NewMouseWheelMovedHandler : New MouseWheelMovedHandler
func NewMouseWheelMovedHandler() MouseWheelMovedHandler {
	return MouseWheelMovedHandler{observers: make([]MouseWheelMoveObserver, 0)}
}

//AddMouseWheelMoveObserver : Adds command
func (mH *MouseWheelMovedHandler) AddMouseWheelMoveObserver(mO MouseWheelMoveObserver) {
	mH.observers = append(mH.observers, mO)
}

//RemoveMouseWheelMoveObserver : removes command
func (mH *MouseWheelMovedHandler) RemoveMouseWheelMoveObserver(mO MouseWheelMoveObserver) {
	for i, o := range mH.observers {
		if o == mO {
			mH.observers = append(mH.observers[:i], mH.observers[i+1:]...)
			return
		}
	}
}

func (mH *MouseWheelMovedHandler) notify(eM EventMouseWheelMoved) {
	for _, o := range mH.observers {
		go o.OnMouseMove(eM)
	}
}

///////////////////////////////////////////////
//TextEntered

//Uses Observer rather than the Command Pattern

//EventTextEntered : Called once text entered
type EventTextEntered struct {
	Char rune //value of the rune
}

//ToSFML : Converts EventTextEntered to the sfml version
//TODO : Refactor all code to fit this naming scheme
func (eM *EventTextEntered) ToSFML() sf.EventTextEntered {
	return sf.EventTextEntered{Char: eM.Char}
}

//SFEventTextEnteredToEventTextEntered sfml to text entered
func SFEventTextEnteredToEventTextEntered(eM sf.EventTextEntered) EventTextEntered {
	return EventTextEntered{Char: eM.Char}
}

//TextEnteredObserver : what is called on text entered
type TextEnteredObserver interface {
	OnTextEntered(EventTextEntered)
}

//TextEnteredHandler : TODO implement
type TextEnteredHandler struct {
	observers []TextEnteredObserver
}

//NewTextEnteredHandler : New NewTextEnteredHandler
func NewTextEnteredHandler() TextEnteredHandler {
	return TextEnteredHandler{observers: make([]TextEnteredObserver, 0)}
}

//AddTextEnteredObserver : Adds command
func (mH *TextEnteredHandler) AddTextEnteredObserver(mO TextEnteredObserver) {
	mH.observers = append(mH.observers, mO)
}

//RemoveTextEnteredObserver : removes command
func (mH *TextEnteredHandler) RemoveTextEnteredObserver(mO TextEnteredObserver) {
	for i, o := range mH.observers {
		if o == mO {
			mH.observers = append(mH.observers[:i], mH.observers[i+1:]...)
			return
		}
	}
}

//notify : Tells all observers everything
func (mH *TextEnteredHandler) notify(eM EventTextEntered) {
	for _, o := range mH.observers {
		go o.OnTextEntered(eM)
	}
}

//InputSystem Has dispatchers for mouse and keyboard as well as notifies observers
//for mouse scroll and others
type InputSystem struct {
	keyboardDispatcher     KeyboardDispatcher
	mouseButtonDispatcher  MouseButtonDispatcher
	mouseWheelMovedHandler MouseWheelMovedHandler
	mouseMovedHandler      MouseMovedHandler
	textEnteredHandler     TextEnteredHandler
}

//NewInputSystem : Creates a New Input System
func NewInputSystem() InputSystem {
	return InputSystem{
		keyboardDispatcher:     NewKeyboardDispatcher(),
		mouseButtonDispatcher:  NewMouseButtonDispatcer(),
		mouseWheelMovedHandler: NewMouseWheelMovedHandler(),
		mouseMovedHandler:      NewMouseMovedHandler(),
		textEnteredHandler:     NewTextEnteredHandler(),
	}
}

//SetKeyPressed : Normally called from keyboard but you are allowed to fake it.
//Sets flags for pressing
func (iS *InputSystem) SetKeyPressed(event sf.EventKeyPressed) {

	eK := SFEventKeyPressedToEventKey(event)
	eKModifier := eK
	eKModifier.Code = ModifierKeyCode
	//go func() {
	iS.keyboardDispatcher.CheckKeyboardSet(&eK)
	//TODO Feels wrong to call this function twice. Should profilie latter to see
	//if there are fewer cache misses the second time than the first.
	iS.keyboardDispatcher.CheckKeyboardSet(&eKModifier)
	//}()
}

//SetKeyReleased : Normally called from keyboard but you are allowed to fake it.
//Sets flags for pressing
func (iS *InputSystem) SetKeyReleased(event sf.EventKeyReleased) {
	eK := SFEventKeyReleasedToEventKey(event)
	eKModifier := eK
	eKModifier.Code = ModifierKeyCode
	//TODO I've pretty much set a limit that the number of inputs processed is
	//O(numActiveKeySets*E[active handlers per set])
	//is this worth putting in a go function. Bet
	//go func() {
	iS.keyboardDispatcher.CheckKeyboardSet(&eK)
	//TODO Feels wrong to call this function twice. Should profilie latter to see
	//if there are fewer cache misses the second time than the first.
	iS.keyboardDispatcher.CheckKeyboardSet(&eKModifier)
	//}()
}

//SetMouseButtonPressed : Normally called from mouseButton but you are allowed to fake it.
//Sets flags for pressing
func (iS *InputSystem) SetMouseButtonPressed(event sf.EventMouseButtonPressed) {
	eM := SFMouseButtonPressedToEventMouseButton(event)
	iS.mouseButtonDispatcher.CheckMouseButtonSet(&eM, event.X, event.Y)

}

//SetMouseButtonReleased : Normally called from mouseButton but you are allowed to fake it.
//Sets flags for pressing
func (iS *InputSystem) SetMouseButtonReleased(event sf.EventMouseButtonReleased) {
	eM := SFMouseButtonReleasedToEventMouseButton(event)
	iS.mouseButtonDispatcher.CheckMouseButtonSet(&eM, event.X, event.Y)

}

//SetMouseMove : Sets the Mouse Move
func (iS *InputSystem) SetMouseMove(eM sf.EventMouseMoved) {
	iS.mouseMovedHandler.notify(SFEventMouseMovedToEventMouseMoved(eM))
}

//SetMouseWheelMove : Sets the Mouse Move
func (iS *InputSystem) SetMouseWheelMove(eM sf.EventMouseWheelMoved) {
	iS.mouseWheelMovedHandler.notify(SFEventMouseWheelMovedToEventMouseMoved(eM))
}

//SetTextEntered : Sets text entered
func (iS *InputSystem) SetTextEntered(eT sf.EventTextEntered) {
	iS.textEnteredHandler.notify(SFEventTextEnteredToEventTextEntered(eT))
}
