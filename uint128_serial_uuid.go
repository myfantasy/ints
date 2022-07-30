package ints

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

const Mval uint64 = ((1 << 2) << 4) << 8
const Nval uint64 = 1 << 63
const maxStep uint64 = 1 << 14

var DefaultUuidGenerator = &UuidSerialGenerator{}

type Uuid struct {
	UInt128
}

// Less returns true when i < val
func (i Uuid) Less(val Uuid) bool {
	return i.UInt128.Less(&val.UInt128)
}

// Equal returns true when i == val
func (i Uuid) Equal(val Uuid) bool {
	return i == val
}

func (i Uuid) Link() *Uuid {
	return &i
}

type UuidSerialGenerator struct {
	RandomTail UInt128
	Mx         sync.Mutex

	step     uint64
	lastTime uint64
}

func (usg *UuidSerialGenerator) Init() error {
	b := make([]byte, 8)
	_, err := rand.Read(b)

	if err != nil {
		return err
	}

	usg.RandomTail.SetBytes(b)

	usg.RandomTail[1] = (usg.RandomTail[1] << 4) >> 4

	return nil
}

// Next Serial id
// xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx
// M 0100 (version 4)
// N 10 (variant 1)
// tttttttt-tttt-Msss-Nrrr-rrrrrrrrrrrr
// M - 0100 (version 4)
// N - 10ss (variant 1)
// t - time (unix milli)
// s - step
// r - random value (Generated on Init)
func (usg *UuidSerialGenerator) Next() (res Uuid) {
	var overflow bool
	for {
		usg.Mx.Lock()
		res, overflow = usg.NextRaw()
		usg.Mx.Unlock()
		if !overflow {
			return res
		}
		time.Sleep(time.Microsecond * 1)
	}
}

// NextOrDefault gets Next from usg where usg != nil else uses DefaultUuidGenerator.Next()
func (usg *UuidSerialGenerator) NextOrDefault() (res Uuid) {
	if usg != nil {
		return usg.Next()
	}

	return DefaultUuidGenerator.Next()
}

// NextRaw - Next Serial id (requare to lock befor call)
// xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx
// M 0100 (version 4)
// N 10 (variant 1)
// tttttttt-tttt-Msss-Nrrr-rrrrrrrrrrrr
// M - 0100 (version 4)
// N - 10ss (variant 1)
// t - time (unix milli)
// s - step
// r - random value (Generated on Init)
// overflow - may be not correct or unique result
func (usg *UuidSerialGenerator) NextRaw() (res Uuid, overflow bool) {
	a := uint64(time.Now().UnixMilli())

	if usg.lastTime < a {
		usg.lastTime = a
		usg.step = 0
	} else if usg.step >= maxStep {
		overflow = true
	}

	// tttttttt-tttt-
	a = a << 16

	// -Msss-
	m := usg.step>>2 + Mval

	// -N (10ss)
	n := (usg.step&3)<<60 + Nval

	// tttttttt-tttt-Msss-
	res.UInt128[0] = a + m
	// -Nrrr-rrrrrrrrrrrr (-N (10ss))
	res.UInt128[1] = n + usg.RandomTail[1]

	// Next step
	usg.step++

	return res, overflow
}

func (i *Uuid) AsUUID() string {
	bts := i.AsBytes()
	res := make([]byte, 36)

	app := 0
	for i, v := range bts {
		if i == 4 || i == 6 || i == 8 || i == 10 {
			res[i*2+app] = hyphen
			app++
		}
		res[i*2+app] = digitsByte[int(v>>4)]
		res[i*2+1+app] = digitsByte[int((v<<4)>>4)]
	}

	return string(res)
}

func (i Uuid) String() string {
	return i.AsUUID()
}

// NextUUID - Next Serial id from DefaultUuidGenerator
// xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx
// M 0100 (version 4)
// N 10 (variant 1)
// tttttttt-tttt-Msss-Nrrr-rrrrrrrrrrrr
// M - 0100 (version 4)
// N - 10ss (variant 1)
// t - time (unix milli)
// s - step
// r - random value (Generated on Init)
func NextUUID() (res Uuid) {
	return DefaultUuidGenerator.Next()
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Uuid) UnmarshalJSON(input []byte) error {
	if input == nil {
		return nil
	}
	if len(input) == 0 {
		return fmt.Errorf("empty value")
	}
	if input[0] != '"' {
		return i.FromTextByte(input, 16, true)
	} else if len(input) < 3 {
		return fmt.Errorf("empty value (\")")
	}
	return i.FromTextByte(input[1:len(input)-1], 16, true)
}

//MarshalJSON() implements json.Marshaler.
func (i Uuid) MarshalJSON() ([]byte, error) {
	txt := "\"" + i.AsUUID() + "\""
	return []byte(txt), nil
}

func (i *Uuid) TimePart() time.Time {
	t := int64(i.UInt128[0] >> 16)
	return time.Unix(t/1000, t%1000000)
}

func (i *Uuid) StepPart() uint64 {
	t := ((i.UInt128[0] << (64 - 12)) >> (64 - 14)) +
		((i.UInt128[1] << 2) >> (64 - 2))
	return t
}

func (i *Uuid) UniquePart() uint64 {
	t := ((i.UInt128[0] << 4) >> 4)
	return t
}
