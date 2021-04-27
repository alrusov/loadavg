package loadavg

import (
	"testing"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------------//

func Test1(t *testing.T) {
	la := Init(3)
	time.Sleep(1100 * time.Millisecond)
	la.Add(1.)

	v := la.Value()
	if v != 0. {
		t.Errorf(`4: Value(): expected 0., got %f`, v)
	}

	time.Sleep(1000 * time.Millisecond)
	v = la.Value()
	if v == 0. { // кривовато
		t.Errorf(`5: Value(): expected 1., got %f`, v)
	}

	for i := 1; i < 5; i++ {
		la.Add(1. * float64(i))
		time.Sleep(1000 * time.Millisecond)
	}

	v = la.Value()
	if v == 0. { // кривовато
		t.Errorf(`6: Value(): expected 2., got %f`, v)
	}
}

//----------------------------------------------------------------------------------------------------------------------------//
