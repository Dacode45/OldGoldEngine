package goldcore

/////////////////////////////////////
///		CONSTS
/////////////////////////////////////

const (
	KeyA         = iota ///< The A key
	KeyB                ///< The B key
	KeyC                ///< The C key
	KeyD                ///< The D key
	KeyE                ///< The E key
	KeyF                ///< The F key
	KeyG                ///< The G key
	KeyH                ///< The H key
	KeyI                ///< The I key
	KeyJ                ///< The J key
	KeyK                ///< The K key
	KeyL                ///< The L key
	KeyM                ///< The M key
	KeyN                ///< The N key
	KeyO                ///< The O key
	KeyP                ///< The P key
	KeyQ                ///< The Q key
	KeyR                ///< The R key
	KeyS                ///< The S key
	KeyT                ///< The T key
	KeyU                ///< The U key
	KeyV                ///< The V key
	KeyW                ///< The W key
	KeyX                ///< The X key
	KeyY                ///< The Y key
	KeyZ                ///< The Z key
	KeyNum0             ///< The 0 key
	KeyNum1             ///< The 1 key
	KeyNum2             ///< The 2 key
	KeyNum3             ///< The 3 key
	KeyNum4             ///< The 4 key
	KeyNum5             ///< The 5 key
	KeyNum6             ///< The 6 key
	KeyNum7             ///< The 7 key
	KeyNum8             ///< The 8 key
	KeyNum9             ///< The 9 key
	KeyEscape           ///< The Escape key
	KeyLControl         ///< The left Control key
	KeyLShift           ///< The left Shift key
	KeyLAlt             ///< The left Alt key
	KeyLSystem          ///< The left OS specific key: window (Windows and Linux), apple (MacOS X), ...
	KeyRControl         ///< The right Control key
	KeyRShift           ///< The right Shift key
	KeyRAlt             ///< The right Alt key
	KeyRSystem          ///< The right OS specific key: window (Windows and Linux), apple (MacOS X), ...
	KeyMenu             ///< The Menu key
	KeyLBracket         ///< The [ key
	KeyRBracket         ///< The ] key
	KeySemiColon        ///< The ; key
	KeyComma            ///< The , key
	KeyPeriod           ///< The . key
	KeyQuote            ///< The ' key
	KeySlash            ///< The / key
	KeyBackSlash        ///< The \ key
	KeyTilde            ///< The ~ key
	KeyEqual            ///< The = key
	KeyDash             ///< The - key
	KeySpace            ///< The Space key
	KeyReturn           ///< The Return key
	KeyBack             ///< The Backspace key
	KeyTab              ///< The Tabulation key
	KeyPageUp           ///< The Page up key
	KeyPageDown         ///< The Page down key
	KeyEnd              ///< The End key
	KeyHome             ///< The Home key
	KeyInsert           ///< The Insert key
	KeyDelete           ///< The Delete key
	KeyAdd              ///< +
	KeySubtract         ///< -
	KeyMultiply         ///< *
	KeyDivide           ///< /
	KeyLeft             ///< Left arrow
	KeyRight            ///< Right arrow
	KeyUp               ///< Up arrow
	KeyDown             ///< Down arrow
	KeyNumpad0          ///< The numpad 0 key
	KeyNumpad1          ///< The numpad 1 key
	KeyNumpad2          ///< The numpad 2 key
	KeyNumpad3          ///< The numpad 3 key
	KeyNumpad4          ///< The numpad 4 key
	KeyNumpad5          ///< The numpad 5 key
	KeyNumpad6          ///< The numpad 6 key
	KeyNumpad7          ///< The numpad 7 key
	KeyNumpad8          ///< The numpad 8 key
	KeyNumpad9          ///< The numpad 9 key
	KeyF1               ///< The F1 key
	KeyF2               ///< The F2 key
	KeyF3               ///< The F3 key
	KeyF4               ///< The F4 key
	KeyF5               ///< The F5 key
	KeyF6               ///< The F6 key
	KeyF7               ///< The F7 key
	KeyF8               ///< The F8 key
	KeyF9               ///< The F9 key
	KeyF10              ///< The F10 key
	KeyF11              ///< The F11 key
	KeyF12              ///< The F12 key
	KeyF13              ///< The F13 key
	KeyF14              ///< The F14 key
	KeyF15              ///< The F15 key
	KeyPause            ///< The Pause key

	KeyCount ///< Keep last -- the total number of keyboard keys
)

//KeyCode : defines the numerical value of a given key
type KeyCode int

var globalKeyboard = map[KeyCode]bool{
	KeyA:         false, ///< The A key
	KeyB:         false, ///< The B key
	KeyC:         false, ///< The C key
	KeyD:         false, ///< The D key
	KeyE:         false, ///< The E key
	KeyF:         false, ///< The F key
	KeyG:         false, ///< The G key
	KeyH:         false, ///< The H key
	KeyI:         false, ///< The I key
	KeyJ:         false, ///< The J key
	KeyK:         false, ///< The K key
	KeyL:         false, ///< The L key
	KeyM:         false, ///< The M key
	KeyN:         false, ///< The N key
	KeyO:         false, ///< The O key
	KeyP:         false, ///< The P key
	KeyQ:         false, ///< The Q key
	KeyR:         false, ///< The R key
	KeyS:         false, ///< The S key
	KeyT:         false, ///< The T key
	KeyU:         false, ///< The U key
	KeyV:         false, ///< The V key
	KeyW:         false, ///< The W key
	KeyX:         false, ///< The X key
	KeyY:         false, ///< The Y key
	KeyZ:         false, ///< The Z key
	KeyNum0:      false, ///< The 0 key
	KeyNum1:      false, ///< The 1 key
	KeyNum2:      false, ///< The 2 key
	KeyNum3:      false, ///< The 3 key
	KeyNum4:      false, ///< The 4 key
	KeyNum5:      false, ///< The 5 key
	KeyNum6:      false, ///< The 6 key
	KeyNum7:      false, ///< The 7 key
	KeyNum8:      false, ///< The 8 key
	KeyNum9:      false, ///< The 9 key
	KeyEscape:    false, ///< The Escape key
	KeyLControl:  false, ///< The left Control key
	KeyLShift:    false, ///< The left Shift key
	KeyLAlt:      false, ///< The left Alt key
	KeyLSystem:   false, ///< The left OS specific key: window (Windows and Linux), apple (MacOS X), ...
	KeyRControl:  false, ///< The right Control key
	KeyRShift:    false, ///< The right Shift key
	KeyRAlt:      false, ///< The right Alt key
	KeyRSystem:   false, ///< The right OS specific key: window (Windows and Linux), apple (MacOS X), ...
	KeyMenu:      false, ///< The Menu key
	KeyLBracket:  false, ///< The [ key
	KeyRBracket:  false, ///< The ] key
	KeySemiColon: false, ///< The ; key
	KeyComma:     false, ///< The , key
	KeyPeriod:    false, ///< The . key
	KeyQuote:     false, ///< The ' key
	KeySlash:     false, ///< The / key
	KeyBackSlash: false, ///< The \ key
	KeyTilde:     false, ///< The ~ key
	KeyEqual:     false, ///< The = key
	KeyDash:      false, ///< The - key
	KeySpace:     false, ///< The Space key
	KeyReturn:    false, ///< The Return key
	KeyBack:      false, ///< The Backspace key
	KeyTab:       false, ///< The Tabulation key
	KeyPageUp:    false, ///< The Page up key
	KeyPageDown:  false, ///< The Page down key
	KeyEnd:       false, ///< The End key
	KeyHome:      false, ///< The Home key
	KeyInsert:    false, ///< The Insert key
	KeyDelete:    false, ///< The Delete key
	KeyAdd:       false, ///< +
	KeySubtract:  false, ///< -
	KeyMultiply:  false, ///< *
	KeyDivide:    false, ///< /
	KeyLeft:      false, ///< Left arrow
	KeyRight:     false, ///< Right arrow
	KeyUp:        false, ///< Up arrow
	KeyDown:      false, ///< Down arrow
	KeyNumpad0:   false, ///< The numpad 0 key
	KeyNumpad1:   false, ///< The numpad 1 key
	KeyNumpad2:   false, ///< The numpad 2 key
	KeyNumpad3:   false, ///< The numpad 3 key
	KeyNumpad4:   false, ///< The numpad 4 key
	KeyNumpad5:   false, ///< The numpad 5 key
	KeyNumpad6:   false, ///< The numpad 6 key
	KeyNumpad7:   false, ///< The numpad 7 key
	KeyNumpad8:   false, ///< The numpad 8 key
	KeyNumpad9:   false, ///< The numpad 9 key
	KeyF1:        false, ///< The F1 key
	KeyF2:        false, ///< The F2 key
	KeyF3:        false, ///< The F3 key
	KeyF4:        false, ///< The F4 key
	KeyF5:        false, ///< The F5 key
	KeyF6:        false, ///< The F6 key
	KeyF7:        false, ///< The F7 key
	KeyF8:        false, ///< The F8 key
	KeyF9:        false, ///< The F9 key
	KeyF10:       false, ///< The F10 key
	KeyF11:       false, ///< The F11 key
	KeyF12:       false, ///< The F12 key
	KeyF13:       false, ///< The F13 key
	KeyF14:       false, ///< The F14 key
	KeyF15:       false, ///< The F15 key
	KeyPause:     false, ///< The Pause key

}

//TODO Implent ability to block a keycode from calling commands.
var globalKeyboardBlockedSet = []KeyCode{}
