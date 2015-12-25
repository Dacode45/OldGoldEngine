package goldcore

/*This package defines a messaging System for game.
Any System Can define its own messages
*/

//GameMessage : Input and Output for Game Nodes
type GameMessage struct {
	Message string      //User defined string. Allows for switches.
	Payload interface{} //Type can be retrieved from global register
}
