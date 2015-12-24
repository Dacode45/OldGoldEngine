package goldcore

import (
	"fmt"
	"sync"
	"testing"
)

func TestKeyboardHandler(t *testing.T) {
	var exampleKH = NewKeyboardHandler()
	closure := false
	cChan := make(chan bool)
	testClosure := func() {
		closure = true
		cChan <- true
	}
	eK := EventKey{
		Code:  KeyA,
		Shift: KeyLShift,
	}
	exampleKH.AddEventKey(eK, testClosure)
	exampleKH.check(&eK)
	<-cChan
	if !closure {
		t.Errorf("KeyHandler failed to call closure. Got closure = %t, want closure = %t", closure, true)
	}
	closure = false
	exampleKH.RemoveEventKey(eK)
	exampleKH.check(&eK)
	if closure {
		t.Errorf("KeyHandler failed to remove closure. Removed %s, but KH is %s", eK, exampleKH)
	}

}

func TestKeyboardSet(t *testing.T) {

	closure1 := false
	closure2 := false
	closure3 := false
	cChan := make(chan bool)
	testClosure1 := func() {
		closure1 = true
		cChan <- true
	}
	testClosure2 := func() {
		closure2 = true
		cChan <- true
	}
	testClosure3 := func() {
		closure3 = true
		cChan <- true
	}
	eK1 := EventKey{
		Code:  KeyA,
		Shift: KeyLShift,
	}
	eK2 := EventKey{
		Code: KeyA,
	}
	exampleKH1 := NewKeyboardHandler()
	exampleKH1.AddEventKey(eK1, testClosure1)
	exampleKH2 := NewKeyboardHandler()
	exampleKH2.AddEventKey(eK1, testClosure2)
	exampleKH3 := NewKeyboardHandler()
	exampleKH3.AddEventKey(eK2, testClosure3)

	exampleKS := NewKeyboardSet()
	p1 := exampleKS.AddHandler(exampleKH1)
	if !exampleKS.MouseButtonIsActive(p1) {
		t.Errorf("Failed to set p1 active")
	}
	p2 := exampleKS.AddHandler(exampleKH2)
	if !exampleKS.MouseButtonIsActive(p2) {
		t.Errorf("Failed to set p2 active")
	}
	p3 := exampleKS.AddHandler(exampleKH3)
	if !exampleKS.MouseButtonIsActive(p3) {
		t.Errorf("Failed to set p3 active")
	}

	exampleKS.check(&eK1)
	<-cChan
	<-cChan
	if !(closure1 && closure2 && !closure3) {
		t.Errorf("Events Failed to propagate. Expected closure1 = true, closure2 = true, closure 3 = false, got: %t, %t, %t", closure1, closure2, closure3)
	}
	closure1 = false
	closure2 = false
	closure3 = false

	exampleKS.SetActive(p2, false)
	exampleKS.check(&eK1)
	<-cChan
	if !(closure1 && !closure2 && !closure3) {
		t.Errorf("Failed to inactivate handler. Expected closure1 = true, closure2 = false, closure 3 = false, got: %t, %t, %t", closure1, closure2, closure3)
	}

	closure1 = false
	closure2 = false
	closure3 = false

	exampleKS.SetActive(p2, true)
	exampleKS.check(&eK1)
	<-cChan
	<-cChan
	if !(closure1 && closure2 && !closure3) {
		t.Errorf("Failed to activate handler. Expected closure1 = true, closure2 = true, closure 3 = false, got: %t, %t, %t", closure1, closure2, closure3)
	}

	closure1 = false
	closure2 = false
	closure3 = false

	exampleKS.RemoveHandler(p2)
	exampleKS.check(&eK1)
	<-cChan
	if !(closure1 && !closure2 && !closure3) {
		t.Errorf("Failed to remove handler. Expected closure1 = true, closure2 = false, closure 3 = false, got: %t, %t, %t", closure1, closure2, closure3)
	}
}

type ClosureStruct struct {
	initialState int
	State        int
	NumChanges   int
	Changed      bool
	Wait         chan<- string
}

func (cS *ClosureStruct) Change() {
	defer func() {
		cS.Wait <- fmt.Sprintf("HI from closure: %d\n", cS.initialState)
	}()
	if !cS.Changed {
		cS.initialState = cS.State
		cS.Changed = true
	}
	cS.State++
	cS.NumChanges++
}
func (cS *ClosureStruct) Reset() {
	cS.State = cS.initialState
	cS.NumChanges = 0
	cS.Changed = false
}

func (cS *ClosureStruct) String() string {
	return fmt.Sprintf("State : %d, NumChanges %d", cS.State, cS.NumChanges)
}
func TestGlobalKeyboardSet(t *testing.T) {
	iS := NewInputSystem()
	limit := 2 * int(KeyCount)
	eKs := make([]EventKey, limit)
	//Divide list into pressed and not pressed
	for i := 0; i < limit; i++ {
		eKs[i] = EventKey{Code: KeyCode(i % (limit / 2)), Pressed: i < KeyCount}
	}
	cChan := make(chan string, limit)
	closures := make([]ClosureStruct, limit)
	handlers := make([]KeyboardHandler, limit)
	handlerKeys := make([]uint, limit)
	handlersPerSet := 10
	sets := make([]KeyboardSet, limit/handlersPerSet)
	setKeys := make([]uint, limit/handlersPerSet)
	for i := 0; i < limit/handlersPerSet; i++ {
		sets[i] = NewKeyboardSet()
	}
	for i := 0; i < limit; i++ {
		closures[i] = ClosureStruct{State: i % (limit / 2), Wait: cChan}
		handlers[i] = NewKeyboardHandler()
		handlers[i].AddEventKey(eKs[i], closures[i].Change)
		handlerKeys[i] = sets[i%(limit/handlersPerSet)].AddHandler(handlers[i])
	}
	for i := 0; i < limit/handlersPerSet; i++ {
		setKeys[i] = iS.keyboardDispatcher.AddKeyboardSet(sets[i])
	}

	//test sets
	//First Check All Keys work
	//pressed
	for i := 0; i < limit/2; i++ {
		testEK := EventKeyToSFEventKeyPressed(eKs[i])
		iS.SetKeyPressed(testEK)
		t.Logf("Testing Key, %#v, Pressed\n", eKs[i])
		t.Log(<-cChan)
		//Check if the registered closures were called
		if !(closures[i].State == i+1 && closures[i+limit/2].State == i && closures[i].NumChanges == 1 && closures[i+limit/2].NumChanges == 0) {
			t.Errorf("KeyPressed did not register. Expected State :%d NumChanges :%d\n\t,"+
				"but Pressed Closure = %#v \n\t and "+
				"Unpressed Closure = %#v\n", i+1, 1, closures[i], closures[i+limit/2])
		}
		closures[i].Reset()
	}
	//Unpressed
	for i := limit / 2; i < limit; i++ {
		testEK := EventKeyToSFEventKeyReleased(eKs[i])
		iS.SetKeyReleased(testEK)
		t.Logf("Testing Key, %#v, Released\n", eKs[i])
		t.Log(<-cChan)
		//Check if the registered closures were called
		if !(closures[i].State == (i%(limit/2))+1 && closures[i%(limit/2)].State == (i%(limit/2)) && closures[i].NumChanges == 1 && closures[i%(limit/2)].NumChanges == 0) {
			t.Errorf("KeyReleased did not register. Expected State :%d NumChanges :%d\n\t,"+
				"but Unpressed Closure = %#v \n\t and "+
				"Pressed Closure = %#v\n", i+1, 1, closures[i], closures[i%(limit/2)])
		}
		closures[i].Reset()
	}
	//Inactivate Half of the sets and try again
	for i, s := range setKeys {
		if i%2 == 0 {

			iS.keyboardDispatcher.SetKeyboardSetActive(s, false)
		}
	}

	for i := 0; i < limit/2; i++ {
		testEK := EventKeyToSFEventKeyPressed(eKs[i])
		iS.SetKeyPressed(testEK)
		t.Logf("Testing Key, %#v, Pressed\n", eKs[i])
		if (i%(limit/handlersPerSet))%2 != 0 {

			t.Log(<-cChan)
		}
		//Check if the registered closures were called
		if (i%(limit/handlersPerSet))%2 == 0 {
			//These should be disabled
			if !(closures[i].State == i && closures[i+limit/2].State == i && closures[i].NumChanges == 0 && closures[i+limit/2].NumChanges == 0) {
				t.Errorf("KeyPressed did not register. Expected State :%d NumChanges :%d\n\t,"+
					"but Pressed Closure = %#v \n\t and "+
					"Unpressed Closure = %#v\n", i, 0, closures[i], closures[i+limit/2])
			}
		} else {
			if !(closures[i].State == i+1 && closures[i+limit/2].State == i && closures[i].NumChanges == 1 && closures[i+limit/2].NumChanges == 0) {
				t.Errorf("KeyPressed did not register. Expected State :%d NumChanges :%d\n\t,"+
					"but Pressed Closure = %#v \n\t and "+
					"Unpressed Closure = %#v\n", i+1, 1, closures[i], closures[i+limit/2])
			}
		}
		closures[i].Reset()
	}
	//Reactivate Half of the sets and try again
	for _, s := range setKeys {

		iS.keyboardDispatcher.SetKeyboardSetActive(s, true)

	}
	for i := 0; i < limit/2; i++ {
		testEK := EventKeyToSFEventKeyPressed(eKs[i])
		iS.SetKeyPressed(testEK)
		t.Logf("Testing Key, %#v, Pressed\n", eKs[i])
		t.Log(<-cChan)
		//Check if the registered closures were called
		if !(closures[i].State == i+1 && closures[i+limit/2].State == i && closures[i].NumChanges == 1 && closures[i+limit/2].NumChanges == 0) {
			t.Errorf("KeyPressed did not register. Expected State :%d NumChanges :%d\n\t,"+
				"but Pressed Closure = %#v \n\t and "+
				"Unpressed Closure = %#v\n", i+1, 1, closures[i], closures[i+limit/2])
		}
		closures[i].Reset()
	}
	//Remove half of the sets and try again
	for i, s := range setKeys {
		if i%2 == 0 {

			iS.keyboardDispatcher.RemoveKeyboardSet(s)
		}
	}
	for i := 0; i < limit/2; i++ {
		testEK := EventKeyToSFEventKeyPressed(eKs[i])
		iS.SetKeyPressed(testEK)
		t.Logf("Testing Key, %#v, Pressed\n", eKs[i])
		if (i%(limit/handlersPerSet))%2 != 0 {

			t.Log(<-cChan)
		}
		//Check if the registered closures were called
		if (i%(limit/handlersPerSet))%2 == 0 {
			//These should be disabled
			if !(closures[i].State == i && closures[i+limit/2].State == i && closures[i].NumChanges == 0 && closures[i+limit/2].NumChanges == 0) {
				t.Errorf("KeyPressed did not register. Expected State :%d NumChanges :%d\n\t,"+
					"but Pressed Closure = %#v \n\t and "+
					"Unpressed Closure = %#v\n", i, 0, closures[i], closures[i+limit/2])
			}
		} else {
			if !(closures[i].State == i+1 && closures[i+limit/2].State == i && closures[i].NumChanges == 1 && closures[i+limit/2].NumChanges == 0) {
				t.Errorf("KeyPressed did not register. Expected State :%d NumChanges :%d\n\t,"+
					"but Pressed Closure = %#v \n\t and "+
					"Unpressed Closure = %#v\n", i+1, 1, closures[i], closures[i+limit/2])
			}
		}
		closures[i].Reset()
	}
}

func TestMouseButtonHandler(t *testing.T) {
	var exampleMH = NewMouseButtonHandler()
	closure := Vector2i{}
	cChan := make(chan bool)
	testClosure := func(button MouseButton, x, y int) {
		closure = Vector2i{X: x, Y: y}
		cChan <- true
	}
	eM := EventMouseButton{
		Button:  MouseLeft,
		Clicked: true,
	}
	exampleMH.AddMouseButton(eM, testClosure)
	exampleMH.check(&eM, 5, 5)
	<-cChan
	if !closure.Equals(Vector2i{5, 5}) {
		t.Errorf("MouseHandler failed to call closure. Got closure = %t, want closure = %t", closure, Vector2i{5, 5})
	}
	closure = Vector2i{}
	exampleMH.RemoveMouseButton(eM)
	exampleMH.check(&eM, 5, 5)
	if closure.Equals(Vector2i{5, 5}) {
		t.Errorf("MouseHandler call closure. Got closure = %t, want closure = %t", closure, Vector2i{5, 5})
	}

}

func TestMouseButtonSet(t *testing.T) {
	closure1 := Vector2i{}
	closure2 := Vector2i{}
	closure3 := Vector2i{}
	cChan := make(chan bool)
	testClosure1 := func(button MouseButton, x, y int) {
		closure1 = Vector2i{x, y}
		cChan <- true
	}
	testClosure2 := func(button MouseButton, x, y int) {
		closure2 = Vector2i{x, y}
		cChan <- true
	}
	testClosure3 := func(button MouseButton, x, y int) {
		closure3 = Vector2i{x, y}
		cChan <- true
	}
	eM1 := EventMouseButton{
		Button:  MouseLeft,
		Clicked: true,
	}
	eM2 := EventMouseButton{
		Button:  MouseRight,
		Clicked: true,
	}
	exampleMH1 := NewMouseButtonHandler()
	exampleMH1.AddMouseButton(eM1, testClosure1)
	exampleMH2 := NewMouseButtonHandler()
	exampleMH2.AddMouseButton(eM1, testClosure2)
	exampleMH3 := NewMouseButtonHandler()
	exampleMH3.AddMouseButton(eM2, testClosure3)

	exampleMS := NewMouseButtonSet()
	p1 := exampleMS.AddHandler(exampleMH1)
	if !exampleMS.MouseButtonIsActive(p1) {
		t.Errorf("Failed to set p1 active")
	}
	p2 := exampleMS.AddHandler(exampleMH2)
	if !exampleMS.MouseButtonIsActive(p2) {
		t.Errorf("Failed to set p2 active")
	}
	p3 := exampleMS.AddHandler(exampleMH3)
	if !exampleMS.MouseButtonIsActive(p3) {
		t.Errorf("Failed to set p3 active")
	}

	check := Vector2i{5, 5}
	exampleMS.check(&eM1, 5, 5)
	<-cChan
	<-cChan
	if !(closure1.Equals(check) && closure2.Equals(check) && !closure3.Equals(check)) {
		t.Errorf("Events Failed to propagate. Expected closure1 = true, closure2 = true, closure 3 = Vector2i{}, got: %#v, %#v, %#v", closure1, closure2, closure3)
	}
	closure1 = Vector2i{}
	closure2 = Vector2i{}
	closure3 = Vector2i{}

	exampleMS.SetActive(p2, false)
	exampleMS.check(&eM1, 5, 5)
	<-cChan
	if !(closure1.Equals(check) && !closure2.Equals(check) && !closure3.Equals(check)) {
		t.Errorf("Failed to inactivate handler. Expected closure1 = true, closure2 = Vector2i{}, closure 3 = Vector2i{}, got: %#v, %#v, %#v", closure1, closure2, closure3)
	}

	closure1 = Vector2i{}
	closure2 = Vector2i{}
	closure3 = Vector2i{}

	exampleMS.SetActive(p2, true)
	exampleMS.check(&eM1, 5, 5)
	<-cChan
	<-cChan
	if !(closure1.Equals(check) && closure2.Equals(check) && !closure3.Equals(check)) {
		t.Errorf("Failed to activate handler. Expected closure1 = true, closure2 = true, closure 3 = Vector2i{}, got: %#v, %#v, %#v", closure1, closure2, closure3)
	}

	closure1 = Vector2i{}
	closure2 = Vector2i{}
	closure3 = Vector2i{}

	exampleMS.RemoveHandler(p2)
	exampleMS.check(&eM1, 5, 5)
	<-cChan
	if !(closure1.Equals(check) && !closure2.Equals(check) && !closure3.Equals(check)) {
		t.Errorf("Failed to remove handler. Expected closure1 = true, closure2 = Vector2i{}, closure 3 = Vector2i{}, got: %#v, %#v, %#v", closure1, closure2, closure3)
	}
}

type MouseStruct struct {
	initialState Vector2i
	State        Vector2i
	NumChanges   int
	Changed      bool
	Wait         chan<- string
}

func (cS *MouseStruct) Change(mB MouseButton, x, y int) {
	defer func() {
		cS.Wait <- fmt.Sprintf("HI from closure: %d\n", cS.initialState)
	}()
	if !cS.Changed {
		cS.initialState = cS.State
		cS.Changed = true
	}
	cS.State = cS.State.Plus(Vector2i{1, 1})
	cS.NumChanges++
}
func (cS *MouseStruct) Reset() {
	cS.State = cS.initialState
	cS.NumChanges = 0
	cS.Changed = false
}

func TestGlobalMouseButtonSet(t *testing.T) {
	iS := NewInputSystem()
	limit := 2 * int(MouseButtonCount)
	eMs := make([]EventMouseButton, limit)
	//Divide list into pressed and not pressed
	for i := 0; i < limit; i++ {
		eMs[i] = EventMouseButton{Button: MouseButton(i % (limit / 2)), Clicked: i < MouseButtonCount}
	}
	cChan := make(chan string, limit)
	buttons := make([]MouseStruct, limit)
	handlers := make([]MouseButtonHandler, limit)
	handlerKeys := make([]uint, limit)
	handlersPerSet := 5
	sets := make([]MouseButtonSet, limit/handlersPerSet)
	setKeys := make([]uint, limit/handlersPerSet)
	for i := 0; i < limit/handlersPerSet; i++ {
		sets[i] = NewMouseButtonSet()
	}
	for i := 0; i < limit; i++ {
		buttons[i] = MouseStruct{State: Vector2i{i % (limit / 2), i % (limit / 2)}, Wait: cChan}
		handlers[i] = NewMouseButtonHandler()
		handlers[i].AddMouseButton(eMs[i], buttons[i].Change)
		handlerKeys[i] = sets[i%(limit/handlersPerSet)].AddHandler(handlers[i])
	}
	for i := 0; i < limit/handlersPerSet; i++ {
		setKeys[i] = iS.mouseButtonDispatcher.AddMouseButtonSet(sets[i])
	}

	//test sets
	//First Check All Keys work
	//pressed
	for i := 0; i < limit/2; i++ {
		testES := EventMouseButtonToSFMouseButtonPressed(eMs[i], 1, 1)
		iS.SetMouseButtonPressed(testES)
		t.Logf("Testing Key, %#v, Pressed\n", eMs[i])
		t.Log(<-cChan)
		//Check if the registered buttons were called
		if !(buttons[i].State.Equals(Vector2i{i + 1, i + 1}) && buttons[i+limit/2].State.Equals(Vector2i{i, i}) && buttons[i].NumChanges == 1 && buttons[i+limit/2].NumChanges == 0) {
			t.Errorf("KeyPressed did not register. Expected State :%d NumChanges :%d\n\t,"+
				"but Pressed Closure = %#v \n\t and "+
				"Unpressed Closure = %#v\n", i+1, 1, buttons[i], buttons[i+limit/2])
		}
		buttons[i].Reset()
	}
	//Unpressed
	for i := limit / 2; i < limit; i++ {
		testEK := EventMouseButtonToSFMouseButtonReleased(eMs[i], 1, 1)
		iS.SetMouseButtonReleased(testEK)
		t.Logf("Testing Key, %#v, Released\n", eMs[i])
		t.Log(<-cChan)
		//Check if the registered buttons were called
		if !(buttons[i].State.Equals(Vector2i{(i % (limit / 2)) + 1, (i % (limit / 2)) + 1}) && buttons[i%(limit/2)].State.Equals(Vector2i{(i % (limit / 2)), (i % (limit / 2))}) && buttons[i].NumChanges == 1 && buttons[i%(limit/2)].NumChanges == 0) {
			t.Errorf("KeyReleased did not register. Expected State :%d NumChanges :%d\n\t,"+
				"but Unpressed Closure = %#v \n\t and "+
				"Pressed Closure = %#v\n", i+1, 1, buttons[i], buttons[i%(limit/2)])
		}
		buttons[i].Reset()
	}
	//Inactivate Half of the sets and try again
	for i, s := range setKeys {
		if i%2 == 0 {

			iS.mouseButtonDispatcher.SetMouseButtonSetActive(s, false)
		}
	}

	for i := 0; i < limit/2; i++ {
		testES := EventMouseButtonToSFMouseButtonPressed(eMs[i], 1, 1)
		iS.SetMouseButtonPressed(testES)
		t.Logf("Testing Key, %#v, Pressed\n", eMs[i])
		if (i%(limit/handlersPerSet))%2 != 0 {

			t.Log(<-cChan)
		}
		//Check if the registered buttons were called
		if (i%(limit/handlersPerSet))%2 == 0 {
			//These should be disabled
			if !(buttons[i].State.Equals(Vector2i{i, i}) && buttons[i+limit/2].State.Equals(Vector2i{i, i}) && buttons[i].NumChanges == 0 && buttons[i+limit/2].NumChanges == 0) {
				t.Errorf("KeyPressed did not register. Expected State :%d NumChanges :%d\n\t,"+
					"but Pressed Closure = %#v \n\t and "+
					"Unpressed Closure = %#v\n", i, 0, buttons[i], buttons[i+limit/2])
			}
		} else {
			if !(buttons[i].State.Equals(Vector2i{i + 1, i + 1}) && buttons[i+limit/2].State.Equals(Vector2i{i, i}) && buttons[i].NumChanges == 1 && buttons[i+limit/2].NumChanges == 0) {
				t.Errorf("KeyPressed did not register. Expected State :%d NumChanges :%d\n\t,"+
					"but Pressed Closure = %#v \n\t and "+
					"Unpressed Closure = %#v\n", i+1, 1, buttons[i], buttons[i+limit/2])
			}
		}
		buttons[i].Reset()
	}
	//Reactivate Half of the sets and try again
	for _, s := range setKeys {

		iS.mouseButtonDispatcher.SetMouseButtonSetActive(s, true)

	}
	for i := 0; i < limit/2; i++ {
		testES := EventMouseButtonToSFMouseButtonPressed(eMs[i], 1, 1)
		iS.SetMouseButtonPressed(testES)
		t.Logf("Testing Key, %#v, Pressed\n", eMs[i])
		t.Log(<-cChan)
		//Check if the registered buttons were called
		if !(buttons[i].State.Equals(Vector2i{i + 1, i + 1}) && buttons[i+limit/2].State.Equals(Vector2i{i, i}) && buttons[i].NumChanges == 1 && buttons[i+limit/2].NumChanges == 0) {
			t.Errorf("KeyPressed did not register. Expected State :%d NumChanges :%d\n\t,"+
				"but Pressed Closure = %#v \n\t and "+
				"Unpressed Closure = %#v\n", i+1, 1, buttons[i], buttons[i+limit/2])
		}
		buttons[i].Reset()
	}
	//Remove half of the sets and try again
	for i, s := range setKeys {
		if i%2 == 0 {

			iS.mouseButtonDispatcher.RemoveMouseButtonSet(s)
		}
	}
	for i := 0; i < limit/2; i++ {
		testES := EventMouseButtonToSFMouseButtonPressed(eMs[i], 1, 1)
		iS.SetMouseButtonPressed(testES)
		t.Logf("Testing Key, %#v, Pressed\n", eMs[i])
		if (i%(limit/handlersPerSet))%2 != 0 {

			t.Log(<-cChan)
		}
		//Check if the registered buttons were called
		if (i%(limit/handlersPerSet))%2 == 0 {
			//These should be disabled
			if !(buttons[i].State.Equals(Vector2i{i, i}) && buttons[i+limit/2].State.Equals(Vector2i{i, i}) && buttons[i].NumChanges == 0 && buttons[i+limit/2].NumChanges == 0) {
				t.Errorf("KeyPressed did not register. Expected State :%d NumChanges :%d\n\t,"+
					"but Pressed Closure = %#v \n\t and "+
					"Unpressed Closure = %#v\n", i, 0, buttons[i], buttons[i+limit/2])
			}
		} else {
			if !(buttons[i].State.Equals(Vector2i{i + 1, i + 1}) && buttons[i+limit/2].State.Equals(Vector2i{i, i}) && buttons[i].NumChanges == 1 && buttons[i+limit/2].NumChanges == 0) {
				t.Errorf("KeyPressed did not register. Expected State :%d NumChanges :%d\n\t,"+
					"but Pressed Closure = %#v \n\t and "+
					"Unpressed Closure = %#v\n", i+1, 1, buttons[i], buttons[i+limit/2])
			}
		}
		buttons[i].Reset()
	}
}

type MouseMoveStruct struct {
	State  Vector2i
	toCall func()
}

func (mS *MouseMoveStruct) OnMouseMove(event EventMouseMoved) {
	//fmt.Println("HI")
	mS.State = Vector2i{event.X, event.Y}
	mS.toCall()
}

func TestMouseMovedHandler(t *testing.T) {
	var wg sync.WaitGroup
	closures := make([]MouseMoveStruct, 10)
	mH := NewMouseMovedHandler()
	closure := func() {
		wg.Done()
	}
	for i := range closures {
		closures[i] = MouseMoveStruct{toCall: closure}
		mH.AddMouseMoveObserver(&closures[i])
	}

	wg.Add(10)
	mH.notify(EventMouseMoved{5, 5})
	wg.Wait()
	for _, c := range closures {
		if !c.State.Equals(Vector2i{5, 5}) {
			t.Errorf("MouseMoved Failed. Expected %#v got %#v", Vector2i{5, 5}, c.State)
		}
	}
}

type MouseWheelMoveStruct struct {
	State  Vector2i
	toCall func()
}

func (mS *MouseWheelMoveStruct) OnMouseMove(event EventMouseWheelMoved) {
	//fmt.Println("HI")
	mS.State = Vector2i{event.X, event.Y}
	mS.toCall()
}

func TestMouseWheelMovedHandler(t *testing.T) {
	var wg sync.WaitGroup
	closures := make([]MouseWheelMoveStruct, 10)
	mH := NewMouseWheelMovedHandler()
	closure := func() {
		wg.Done()
	}
	for i := range closures {
		closures[i] = MouseWheelMoveStruct{toCall: closure}
		mH.AddMouseWheelMoveObserver(&closures[i])
	}

	wg.Add(10)
	mH.notify(EventMouseWheelMoved{5, 5, 5})
	wg.Wait()
	for _, c := range closures {
		if !c.State.Equals(Vector2i{5, 5}) {
			t.Errorf("MouseMoved Failed. Expected %#v got %#v", Vector2i{5, 5}, c.State)
		}
	}
}
