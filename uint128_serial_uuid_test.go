package ints

import (
	"encoding/json"
	"testing"
	"time"
)

type testStructUuid struct {
	A Uuid
	B *Uuid
	C *Uuid
}

func TestUInt128ToUUID(t *testing.T) {

	s := testStructUuid{
		A: DefaultUuidGenerator.Next(),
		B: DefaultUuidGenerator.Next().Link(),
	}

	b, err := json.Marshal(s)
	if err != nil {
		t.Errorf("fail Marshal %v", err)
	}

	res := testStructUuid{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		t.Errorf("fail Unmarshal %v", err)
	}

	if s.A != res.A {
		t.Errorf("fail A is not equal %v != %v (%v)", s.A, res.A, string(b))
	}

	if *s.B != *res.B {
		t.Errorf("fail B is not equal %v != %v (%v)", s.B, res.B, string(b))
	}

	if res.C != nil {
		t.Errorf("fail C should be nil %v (%v)", res.C, string(b))
	}
}

func BenchmarkUInt128Next(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NextUUID()
	}
}

func TestLimitSerialUUID(t *testing.T) {
	a1 := NextUUID()
	time.Sleep(1 * time.Millisecond)
	ls := LimitSerialUUID(time.Now())
	time.Sleep(1 * time.Millisecond)
	a2 := NextUUID()

	if ls.Less(a1) {
		t.Fatalf("ls should be more then a1 but not: ls: %v a1: %v", ls, a1)
	}

	if a2.Less(ls) {
		t.Fatalf("ls should be less then a2 but not: ls: %v a2: %v", ls, a2)
	}
}

func TestGetTimeFormSerialUUID(t *testing.T) {
	t0 := time.Now()
	time.Sleep(2 * time.Millisecond)
	a1 := NextUUID()
	tm := GetTimeFormSerialUUID(a1)
	time.Sleep(2 * time.Millisecond)
	t1 := time.Now()

	if tm.Before(t0) {
		t.Fatalf("tm should be after then t0 but not: tm: %v t0: %v", tm, t0)
	}

	if tm.After(t1) {
		t.Fatalf("tm should be befor then t1 but not: tm: %v t0: %v", tm, t1)
	}
}
