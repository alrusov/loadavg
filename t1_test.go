package loadavg

import (
	"testing"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------------//

func Test1(t *testing.T) {
	la, err := Init(0)
	if err == nil {
		t.Errorf(`1: Init(0): expect error, got no error`)
	}

	la, err = Init(1)
	if err != nil {
		t.Errorf(`2: Init(1): %s`, err.Error())
	}

	la, err = Init(3)
	if err != nil {
		t.Errorf(`3: Init(3): %s`, err.Error())
	} else {
		time.Sleep(1100 * time.Millisecond)
		la.Add(1.)

		v := la.Value()
		if v != 0. {
			t.Errorf(`4: Value(): expect 0., got %f`, v)
		}

		time.Sleep(1000 * time.Millisecond)
		v = la.Value()
		if v == 0. { // кривовато
			t.Errorf(`5: Value(): expect 1., got %f`, v)
		}

		for i := 1; i < 5; i++ {
			la.Add(1. * float64(i))
			time.Sleep(1000 * time.Millisecond)
		}

		v = la.Value()
		if v == 0. { // кривовато
			t.Errorf(`6: Value(): expect 2., got %f`, v)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------------//
