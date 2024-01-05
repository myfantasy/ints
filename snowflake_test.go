package ints

import (
	"testing"
	"time"
)

func BenchmarkSnowflakeGenerator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NextID()
	}
}

func TestLimitID(t *testing.T) {
	a1 := NextID()
	time.Sleep(1 * time.Millisecond)
	ls := LimitID(time.Now())
	time.Sleep(1 * time.Millisecond)
	a2 := NextID()

	if ls < a1 {
		t.Fatalf("ls should be more then a1 but not: ls: %v a1: %v", ls, a1)
	}

	if a2 < ls {
		t.Fatalf("ls should be less then a2 but not: ls: %v a2: %v", ls, a2)
	}
}

func TestGetTimeFromID(t *testing.T) {
	t0 := time.Now()
	time.Sleep(2 * time.Millisecond)
	a1 := NextID()
	tm := GetTimeFromID(a1)
	time.Sleep(2 * time.Millisecond)
	t1 := time.Now()

	if tm.Before(t0) {
		t.Fatalf("tm should be after then t0 but not: tm: %v t0: %v", tm, t0)
	}

	if tm.After(t1) {
		t.Fatalf("tm should be befor then t1 but not: tm: %v t0: %v", tm, t1)
	}
}
