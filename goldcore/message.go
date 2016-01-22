package goldcore

import (
	"fmt"
	"sync"
)

/*This package defines a messaging System for game.
Any System Can define its own messages
*/

//GMessage : Wrapper around index
type GMessage uint32

//GameMessage : Input and Output for Game Nodes
type GameMessage struct {
	Message GMessage    //User defined string. Allows for switches.
	Payload interface{} //Type can be retrieved from global register
}

func (msg *GameMessage) String() string {
	return fmt.Sprintf("Message: %s\nPayload: %v", gameMessageStringRegistrar[msg.Message], msg.Payload)
}

var gameMessageMutex = &sync.Mutex{}
var gameMessageStringRegistrar = []string{}

//RegisterGameMessage : Use this function to create your own GameMessage
func RegisterGameMessage(msg string) GMessage {
	gameMessageMutex.Lock()
	gameMessageStringRegistrar = append(gameMessageStringRegistrar, msg)
	index := len(gameMessageStringRegistrar) - 1
	gameMessageMutex.Unlock()
	return GMessage(index)
}

var gameMessageBuffer = []GameMessage{}
var gameMessageBufferClear = false
var gameMessageBufferIndex int

//NewGameMessage : Use this fuction to create Game Messages
func NewGameMessage(msg GMessage, payload interface{}) *GameMessage {
	if gameMessageBufferClear {
		gameMessageBufferIndex = 0
		gameMessageBufferClear = false
	}
	if gameMessageBufferIndex >= len(gameMessageBuffer) {
		gameMessageBuffer = append(gameMessageBuffer, GameMessage{Message: msg, Payload: payload})
	} else {
		gameMessageBuffer[gameMessageBufferIndex].Message = msg
		gameMessageBuffer[gameMessageBufferIndex].Payload = payload
	}
	gameMessageBufferIndex++
	return &gameMessageBuffer[gameMessageBufferIndex-1]
}

/*clearGameMessageBuffer : Invalidates every message in the GameMessageBuffer.
It is essential that this only be called after the time it takes for the longerst
GameMessage "path" to be taken" Remember that GameMessages often invoke other
game messages Technically this doesn't have to be called at all depending on the
games runtime*/
func clearGameMessageBuffer() {
	gameMessageBufferClear = true
}

//dumpGameMessageBuffer : Prints all valid game messages
func dumpGameMessageBuffer() string {
	return fmt.Sprintln(gameMessageBuffer[:gameMessageBufferIndex])
}
