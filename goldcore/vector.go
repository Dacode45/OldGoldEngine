package goldcore

import (
	sf "github.com/manyminds/gosfml"
)

/////////////////////////////////////
///		STRUCTS
/////////////////////////////////////

//Vector2i : 2D integer vector
type Vector2i struct {
	X, Y int
}

//Vector2u : 2D unsigned vector
type Vector2u struct {
	X, Y uint
}

//Vector2f : 2D float vector
type Vector2f struct {
	X, Y float32
}

//Vector3f : 3D float vector
type Vector3f struct {
	X, Y, Z float32
}

/////////////////////////////////////
///		FUNCS
/////////////////////////////////////

/////////////////////////////////////
// Vector2i

//Plus : Returns the sum of two vectors.
func (vec Vector2i) Plus(other Vector2i) Vector2i {
	return Vector2i{X: vec.X + other.X, Y: vec.Y + other.Y}
}

//Minus : Returns the difference of two vectors.
func (vec Vector2i) Minus(other Vector2i) Vector2i {
	return Vector2i{X: vec.X - other.X, Y: vec.Y - other.Y}
}

//Equals : True if x and y of vector are equal
func (vec Vector2i) Equals(other Vector2i) bool {
	return vec.X == other.X && vec.Y == other.Y
}

//SFVector2uToGEVector2i : SFML Vector2u to GoldEngine Vector2U
//TODO Refactor name of these kinds of function to be SF_ToGE_
func SFVector2uToGEVector2i(other sf.Vector2i) Vector2i {
	return Vector2i{X: other.X, Y: other.Y}
}

/////////////////////////////////////
// Vector2u

//Plus : Returns the sum of two vectors.
func (vec Vector2u) Plus(other Vector2u) Vector2u {
	return Vector2u{X: vec.X + other.X, Y: vec.Y + other.Y}
}

//Minus : Returns the difference of two vectors.
func (vec Vector2u) Minus(other Vector2u) Vector2u {
	return Vector2u{X: vec.X - other.X, Y: vec.Y - other.Y}
}

//Equals : True if x and y of vector are equal
func (vec Vector2u) Equals(other Vector2u) bool {
	return vec.X == other.X && vec.Y == other.Y
}

//SFVector2uToGEVector2u : SFML Vector2u to GoldEngine Vector2U
//TODO Refactor name of these kinds of function to be SF_ToGE_
func SFVector2uToGEVector2u(other sf.Vector2u) Vector2u {
	return Vector2u{X: other.X, Y: other.Y}
}

/////////////////////////////////////
// Vector2f

//Plus : Returns the sum of two vectors.
func (vec Vector2f) Plus(other Vector2f) Vector2f {
	return Vector2f{X: vec.X + other.X, Y: vec.Y + other.Y}
}

//Minus : Returns the difference of two vectors.
func (vec Vector2f) Minus(other Vector2f) Vector2f {
	return Vector2f{X: vec.X - other.X, Y: vec.Y - other.Y}
}

//Equals : True if x and y of vector are equal
func (vec Vector2f) Equals(other Vector2f) bool {
	return vec.X == other.X && vec.Y == other.Y
}

//ToSFML : Allows for SFML compatability
func (vec Vector2i) ToSFML() sf.Vector2i {
	return sf.Vector2i{X: vec.X, Y: vec.Y}
}

//ToSFML : Allows for SFML compatability
func (vec Vector2u) ToSFML() sf.Vector2u {
	return sf.Vector2u{X: vec.X, Y: vec.Y}
}

//ToSFML : Allows for SFML compatability
func (vec Vector2f) ToSFML() sf.Vector2f {
	return sf.Vector2f{X: vec.X, Y: vec.Y}
}

//ToSFML : Allows for SFML compatability
func (vec Vector3f) ToSFML() sf.Vector3f {
	return sf.Vector3f{X: vec.X, Y: vec.Y}
}
