package ints

import (
	"encoding/json"
	"testing"
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
		DefaultUuidGenerator.Next()
	}
}
