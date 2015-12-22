package goldcore

import (
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
	Alt     bool    //< Is the Alt key pressed?
	Control bool    //< Is the Control key pressed?
	Shift   bool    //< Is the Shift key pressed?
	System  bool    //< Is the System key pressed?
	Pressed bool    //< Is the button being pressed (true) or released (false)
}

//ModifierKeyCode : Used to check if for modifiers only. Example inorder to check if only the shift
//key was pressed make an EventKey with Code set to ModifierCode and Shift set to true
const ModifierKeyCode = -1

//KeyboardHandler : Will IMMEDIATELY call a function once key is pressed as long
//as the handler or the set is active. It is up to the user to ensure syncronization
//KeyboardHandlers only work as a part of a key set
type KeyboardHandler struct {
	keysPressed map[EventKey]InputCommand
	active      bool
}

//AddEventKey : Binds Command to Event Key object.
func (kH *KeyboardHandler) AddEventKey(eK EventKey, kC KeyCommand) {
	kH.keysPressed[eK] = kC
}

//RemoveEventKey : Removes Event Key and its effects
func (kH *KeyboardHandler) RemoveEventKey(eK EventKey) {
	delete(kH.keysPressed, eK)
}

//IsActive : Checks if keyboardHandler is active.
func (kH *KeyboardHandler) IsActive() bool {
	return kH.active
}

//check : Internally called. Will call Key Command if available
func (kH *KeyboardHandler) check(eK *EventKey) {
	if cmd, ok := kH.keysPressed[&eK]; ok {
		cmd.Execute()
	}
}

//KeyboardSet :Set of actual KeyboardHandlers called on input. You can enable
//and disable keyboardhandlers rather than constantly adding and removing them
type KeyboardSet struct {
	keyboardHandlers []KeyboardHandler
	numActive        int
	active           bool
}

//IsActive : Check if Keyboard Set is active
func (kS *KeyboardSet) IsActive() bool {
	return kS.active
}

//AddHandler : Give a Handler and get a pointer back to the handler in memory
//Automatically makes it active so set active will need to be called afterwards
func (kS *KeyboardSet) AddHandler(kH KeyboardHandler) *KeyboardHandler {
	kS.keyboardHandlers = append(ks.keyboardHandlers[:kS.numActive], kH, kS.keyboardHandlers[ks.numActive:])
	kS.numActive++
	return &kH
}

//RemoveHandler : Removes KeyboardHandler
func (kS *KeyboardSet) RemoveHandler(kH *KeyboardHandler) {
	//This is much faster if this array is actually in memory
	for i, k := range kS.keyboardHandlers {
		if &k == kH {
			kS.keyboardHandlers = append(kS.keyboardHandlers[:i], kS.keyboardHandlers[i+1:])
			kS.numActive--
			return
		}
	}
}

//SetActive : Makes the given keyboardHandler active
func (kS *KeyboardSet) SetActive(kH *KeyboardHandler, active bool) {
	for i, k := range kS.keyboardHandlers {
		if &k == kH {
			//only do something if different
			if kH.active != active {
				if k.active { //if active put it in the active part of the array
					temp := kS.keyboardHandlers[kS.numActive]
					kS.keyboardHandlers[kS.numActive] = k
					kS.keyboardHandlers[i] = temp
					kS.numActive++
				} else { //inactive move to end of array
					ks.keyboardHandlers = append(kS.keyboardHandlers[:i], kS.keyboardHandlers[i+1], k)
					kS.numActive--
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

//KeySets : array of keyboardSets.
type KeySets []KeyboardSet

var numActiveKeySets = 0
var keyboardSets = make(KeySets, 0)

//AddKeyboardSet Registers this set globally. You will no longer have to add
//and re-add sets. Just enable and disable them
func AddKeyboardSet(kS KeyboardSet) {
	keyboardSets = append(keyboardSets, kS)
	numActiveKeySets++
	return &kS
}

//RemoveKeyboardSet : Unregister KeyboardSet globally
func RemoveKeyboardSet(kS *KeyboardSet) {
	for i, k := range keyboardSets {
		if &k == kS {
			keyboardSets = append(keyboardSets[:i], keyboardSets[i+1:])
			numActiveKeySets--
			return
		}
	}
}

//SetKeyboardSetActive : Enable or disable KeybaordSet
func SetKeyboardSetActive(kH *KeyboardSet, active bool) {
	for i, k := range keyboardSets {
		if &k == kS {
			//only do something if different
			if kS.active != active {
				if k.active { //if active put it in the active part of the array
					temp := keyboardSets[kS.numActiveKeySets]
					keyboardSets[kS.numActive] = k
					keyboardSets[i] = temp
					numActiveKeySets++
				} else { //inactive move to end of array
					keyboardSets = append(keyboardSets[:i], keyboardSets[i+1], k)
					numActiveKeySets--
				}
			}
			return
		}
	}
}

//CheckKeyboardSet : Calls the keybaord command for all associated event keys
//Passes pointer to event to the actual handlers.
func CheckKeyboardSet(eK *EventKey) {
	for i := 0; i < numActiveKeySets; i++ {
		keyboardSets[i].check(eK)
	}
}

//SetKeyPressed : Normally called from keyboard but you are allowed to fake it.
//Sets flags for pressing
func SetKeyPressed(event sf.EventKeyPressed) {
	eK := EventKey{
		Code:    event.Code,
		Alt:     event.Alt != 0,
		Control: event.Control != 0,
		Shift:   event.Shift != 0,
		System:  event.System != 0,
		Pressed: true}
	eKModifier := eK
	eKModifier.Code = ModifierKeyCode
	CheckKeyboardSet(&eK)
	//TODO Feels wrong to call this function twice. Should profilie latter to see
	//if there are fewer cache misses the second time than the first.
	CheckKeyboardSet(&eKModifier)
	registerEventKey(event, true)
}

//SetKeyReleased : Normally called from keyboard but you are allowed to fake it.
//Sets flags for pressing
func SetKeyReleased(event sf.EventKeyReleased) {
	eK := EventKey{
		Code:    event.Code,
		Alt:     event.Alt != 0,
		Control: event.Control != 0,
		Shift:   event.Shift != 0,
		System:  event.System != 0,
		Pressed: false}
	eKModifier := eK
	eKModifier.Code = ModifierKeyCode
	CheckKeyboardSet(eK)
	//TODO Feels wrong to call this function twice. Should profilie latter to see
	//if there are fewer cache misses the second time than the first.
	CheckKeyboardSet(eKModifier)
	registerEventKey(eK, false)
}

//registerEventKey : internally sets key pressed
//TODO, I should double check wheter all keys sync up. Esspecialy R&LShifts, Control, Alt
func registerEventKey(eK EventKey) {
	globalKeyboard[eK.Code] = eK.pressed
	if eK.Alt != 0 {
		globalKeyboard[eK.Alt] = eK.pressed
	}
	if eK.Control != 0 {
		globalKeyboard[eK.Control] = eK.pressed
	}
	if eK.Shift != 0 {
		globalKeyboard[eK.Shift] = eK.pressed

	}
	if eK.System != 0 {
		globalKeyboard[eK.System] = eK.pressed

	}
}

//GetKeyPressed : Tells if key is currently held down. Since This does not check
//wheter this is the firt frame it is down, users mus set that flag themselves
func GetKeyPressed(KeyCode) {
	pressed, _ := globalKeyboard[KeyCode]
	return pressed
}

//MouseHandler : TODO implement
type MouseHandler struct {
}
