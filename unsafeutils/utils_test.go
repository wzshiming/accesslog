package unsafeutils

import (
	"reflect"
	"testing"
)

func TestCache(t *testing.T) {
	if reflect.TypeFor[S]() != reflect.TypeFor[S]() {
		t.FailNow()
	}
}

type S struct {
	_ func()
	A int
	B string
}

func TestFields(t *testing.T) {
	got := Fields[S]()
	want := []string{
		"A",
		"B",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGet(t *testing.T) {
	s := S{
		A: 10,
		B: "hello",
	}
	a, _ := Get[int](&s, "A")
	if a != 10 {
		t.Errorf("got %v, want %v", a, 10)
	}
	b, _ := Get[string](&s, "B")
	if b != "hello" {
		t.Errorf("got %v, want %v", b, "hello")
	}
}

func TestSet(t *testing.T) {
	s := S{}
	_ = Set[int](&s, "A", 10)
	if s.A != 10 {
		t.Errorf("got %v, want %v", s.A, 10)
	}
	_ = Set[string](&s, "B", "hello")
	if s.B != "hello" {
		t.Errorf("got %v, want %v", s.B, "hello")
	}
}
