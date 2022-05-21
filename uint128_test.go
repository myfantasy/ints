package ints

import "testing"

func TestUInt128DivAndMult(t *testing.T) {
	a := &UInt128{110, 110}
	b := &UInt128{0, 100}

	q, r := a.Div(b)
	k := &UInt128{1, 1844674407370955162}
	m := b.Mul(k)

	if !m.Equal(a.Sub(&r).Link()) {
		t.Errorf("source %v bytes %v result %v[%v] (%v)\n", a, b, q, &q, r)
		t.Errorf("source %v bytes %v result %v \n", b, k, m)
	}
}
