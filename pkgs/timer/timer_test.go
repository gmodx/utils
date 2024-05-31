package timer

import (
	"testing"
	"time"
)

func Test_Tick(t *testing.T) {
	var i = 0

	iPlusPlus := func() {
		i++
	}

	err := Tick(0, 100*time.Millisecond, iPlusPlus, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	<-time.After(50 * time.Millisecond)
	if i != 1 {
		t.Errorf("i, expected is %v, current: %v", 1, i)
		return
	}

	<-time.After(100 * time.Millisecond)
	if i != 2 {
		t.Errorf("i, expected is %v, current: %v", 2, i)
		return
	}

	<-time.After(100 * time.Millisecond)
	if i != 3 {
		t.Errorf("i, expected is %v, current: %v", 3, i)
		return
	}
}
