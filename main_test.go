package objDict

import (
	"testing"
	"time"
)

var od = New(10 * time.Minute)

type nested struct {
	f4_1 string
}
type tester struct {
	F1 string
	F2 map[string]interface{}
	F3 interface{}
	F4 nested
}

var t1 = tester{
	F1: "foo",
	F2: map[string]interface{}{"foo1": "bar1"},
	F3: nil,
	F4: nested{"bar1"},
}
var t1_f1 tester
var t1_f2 tester = t1
var t1_f3 = t1
var t1_f4 = t1
var Testee = []*tester{
	&t1_f1,
	&t1_f2,
	&t1_f3,
	&t1_f4,
}

func TestInit(t *testing.T) {
	t1_f1.F1 = "difFoo"
	t1_f2.F2 = map[string]interface{}{"Foo1": "difBar"}
	t1_f3.F3 = "123"
	t1_f4.F4 = nested{"Foo1"}
}

func TestGet(t *testing.T) {
	TestInit(t)
	od.Set(t1, -1)
	for idx, obj := range Testee {
		od.Set(*obj, idx+1)
		if i_t1, ok := od.Get(t1); !ok || i_t1 != -1 {
			t.Errorf("t1 changed to %d, want %d.", i_t1, -1)
		}
	}
	od.Set(t1, -2)
	if i_t1, ok := od.Get(t1); !ok || i_t1 != -2 {
		t.Errorf("t1 not changed")
	}
	for idx, obj := range Testee {
		if i_ti, ok := od.Get(*obj); !ok || i_ti != idx+1 {
			t.Errorf("Rewrite by t1: t%d", idx+1)
		}
	}

	od.Set(t1, -10)
	od.SaveFile(".cache")
}

func TestSaveFile(t *testing.T) {
	od.LoadFile(".cache")
	if ans, ok := od.Get(t1); !ok {
		t.Errorf("savefile")
	} else {
		if ans != -10 {
			t.Errorf("savefile")
		}
	}
}
