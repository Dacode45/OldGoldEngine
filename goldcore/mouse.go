package goldcore

import (
	sf "github.com/manyminds/gosfml"
)

const (
	MouseLeft     = iota ///< The left mouse button
	MouseRight           ///< The right mouse button
	MouseMiddle          ///< The middle (wheel) mouse button
	MouseXButton1        ///< The first extra mouse button
	MouseXButton2        ///< The second extra mouse button

	MouseButtonCount ///< Keep last -- the total number of mouse buttons
)

//MouseButton : 5 possible values
type MouseButton int

/////////////////////////////////////
///		FUNCTIONS
/////////////////////////////////////

//IsMouseButtonPressed : Check if a mouse button is pressed
//
// 	button: Button to check
func IsMouseButtonPressed(button MouseButton) bool {
	return sf.IsMouseButtonPressed(sf.MouseButton(button))
}

//MouseSetPosition : Set the current position of the mouse
//
// This function sets the current position of the mouse
// cursor relative to the given window, or desktop if nil is passed.
//
// 	position:   New position of the mouse
// 	relativeTo: Reference window
func MouseSetPosition(position Vector2i, relativeTo GameWindow) {

	sf.MouseSetPosition(position.ToSFML(), relativeTo.renderWindow)
}

//MouseGetPosition : Get the current position of the mouse
//
// This function returns the current position of the mouse
// cursor relative to the given window, or desktop if nil is passed.
//
// 	relativeTo: Reference window
func MouseGetPosition(relativeTo GameWindow) Vector2i {
	pos := sf.MouseGetPosition(relativeTo.renderWindow)
	return Vector2i{X: pos.X, Y: pos.Y}
}
