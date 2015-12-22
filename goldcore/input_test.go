package goldcore

import "testing"

var exampleKH = KeyboardHandler{}

func TestKeyboardHandler(t *testing.T) {
	closure := false
	testClosure := func() {
		closure = true
	}
	eK := EventKey{
		Code:  KeyA,
		Shift: true,
	}
	exampleKH.AddEventKey(eK, testClosure)
	exampleKH.check(eK)
	if !closure {
		t.Errorf("KeyHandler failed to call closure. Got closure = %t, want closure = %t", closure, true)
	}
	closure = false
	exampleKH.RemoveEventKey(eK)
	exampleKH.check(eK)
	if closure {
		t.Errorf("KeyHandler failed to remove closure. Removed %s, but KH is %s", eK, exampleKH)
	}

}
