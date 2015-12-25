package goldcore

import (
	"runtime"
	"testing"

	"github.com/trustmaster/goflow"
)

type GameWindowListener struct {
	flow.Component
	InputGameMessage <-chan GameMessage
	t                *testing.T
}

func (gW *GameWindowListener) OnInputGameMessage(gM GameMessage) {
	gW.t.Log(gM.Message)
	GameWindowMessage <- gM
}

var GameWindowMessage = make(chan GameMessage)

type GameWindowTester struct {
	flow.Graph
}

func SetupGameWindowTest(t *testing.T) (*GameWindowTester, *GameWindow) {
	n := &GameWindowTester{}
	n.InitGraphState()

	gW := NewGameWindow(800, 800, "Test Window")
	n.Add(gW, "gamewindow")
	n.Add(&GameWindowListener{t: t}, "gamewindowlistener")

	n.Connect("gamewindow", "OutputGameMessage", "gamewindowlistener", "InputGameMessage")

	n.MapInPort("In", "gamewindow", "InputGameMessage")
	return n, gW
}

func TestGameWindow(t *testing.T) {
	gT, _ := SetupGameWindowTest(t)

	in := make(chan GameMessage)
	gT.SetInPort("In", in)
	flow.RunNet(gT)
	in <- GameMessage{Message: "hello"}

	t.Log("Messages Propagated", <-GameWindowMessage)

	in <- GameMessage{Message: WindowRunning}
	for (<-GameWindowMessage).Message != WindowRunning {
		runtime.Gosched()
	}
	t.Log("Messages Propagated")
	in <- GameMessage{Message: WindowClosed}
	for msg := <-GameWindowMessage; msg.Message != WindowClosed; msg = <-GameWindowMessage {
		t.Log("Recieved", msg)
	}
	t.Log("Messages Propagated")
	close(in)
	<-gT.Wait()
}
