package ints

import (
	"strings"
	"testing"
)

func TestUInt128DivAndMult(t *testing.T) {
	a := UInt128{110, 110}
	b := UInt128{0, 100}

	q, r := a.Div(b)
	k := UInt128{1, 1844674407370955162}
	m := b.Mul(k)

	if !m.Equal(a.Sub(r)) {
		t.Errorf("source %v bytes %v result %v[%v] (%v)\n", a, b, q, &q, r)
		t.Errorf("source %v bytes %v result %v \n", b, k, m)
	}
}

func TestUInt128Text(t *testing.T) {
	a := UInt128{110, 110}

	s := a.Text(62)

	if s != "CZF1bTWWf8Oa" {
		t.Errorf("incorrect convert base 62 should %v but %v", "CZF1bTWWf8Oa", s)
	}

	b, err := UInt128FromText(s, 62, false)
	if err != nil {
		t.Error(err)
	}

	if !b.Equal(a) {
		t.Errorf("values should be equal %v %v", a, b)
	}

	s = a.Text(10)

	if s != "2029141848108050677870" {
		t.Errorf("incorrect convert base 10 should %v but %v", "2029141848108050677870", s)
	}

	b, err = UInt128FromText(s, 10, false)
	if err != nil {
		t.Error(err)
	}

	if !b.Equal(a) {
		t.Errorf("values should be equal %v %v", a, b)
	}

	s = a.Text(16)

	if s != "6e000000000000006e" {
		t.Errorf("incorrect convert base 16 should %v but %v", "6e000000000000006e", s)
	}

	b, err = UInt128FromText(s, 16, false)
	if err != nil {
		t.Error(err)
	}

	if !b.Equal(a) {
		t.Errorf("values should be equal %v %v", a, b)
	}

	s = a.Text(2)

	if s != "11011100000000000000000000000000000000000000000000000000000000001101110" {
		t.Errorf("incorrect convert base 2 should %v but %v", "11011100000000000000000000000000000000000000000000000000000000001101110", s)
	}

	b, err = UInt128FromText(s, 2, false)
	if err != nil {
		t.Error(err)
	}

	if !b.Equal(a) {
		t.Errorf("values should be equal %v %v", a, b)
	}

	s = a.Text(36)

	if s != "bw8gv58mqmzazy" {
		t.Errorf("incorrect convert base 36 should %v but %v", "bw8gv58mqmzazy", s)
	}

	b, err = UInt128FromText(s, 36, false)
	if err != nil {
		t.Error(err)
	}

	if !b.Equal(a) {
		t.Errorf("values should be equal %v %v", a, b)
	}

	b, err = UInt128FromText(strings.ToUpper(s), 36, false)
	if err != nil {
		t.Error(err)
	}

	if !b.Equal(a) {
		t.Errorf("values should be equal %v %v", a, b)
	}
}
