package ints

import (
	"encoding/binary"
	"fmt"
	"math/bits"
)

// UInt128 128 bit unsigned integer
type UInt128 [2]uint64

// Less returns true when i < val
func (i *UInt128) Less(val *UInt128) bool {
	if i[0] < val[0] {
		return true
	} else if i[0] > val[0] {
		return false
	}
	return i[1] < val[1]
}

// Equal returns true when i == val
func (i *UInt128) Equal(val *UInt128) bool {
	return *i == *val
}

func (i *UInt128) AsBytes() (result [16]byte) {
	binary.BigEndian.PutUint64(result[0:8], i[0])
	binary.BigEndian.PutUint64(result[8:16], i[1])
	return result
}

// BytesLen returs count of requared bytes (16..0) if 0 then i is 0
func (i *UInt128) BytesLen() int {
	bytes := i.AsBytes()

	for i := 0; i < 16; i++ {
		if bytes[i] > 0 {
			return 16 - i
		}
	}

	return 0
}

func (i *UInt128) FromBytes(buf [16]byte) *UInt128 {
	i[0] = binary.BigEndian.Uint64(buf[0:8])
	i[1] = binary.BigEndian.Uint64(buf[8:16])

	return i
}

func (i UInt128) Link() *UInt128 {
	return &i
}

func (i *UInt128) Copy() UInt128 {
	return UInt128{i[0], i[1]}
}

func (i *UInt128) IsEmpty() bool {
	return i[0] == 0 && i[1] == 0
}
func (i *UInt128) IsUint64() bool {
	return i[0] == 0
}

func UInt128FromBytes(buf [16]byte) (i UInt128) {
	i.FromBytes(buf)
	return i
}

func (i *UInt128) Set(val *UInt128) *UInt128 {
	i[0] = val[0]
	i[1] = val[1]

	return i
}
func (i *UInt128) SetUint64(val *uint64) *UInt128 {
	i[0] = 0
	i[1] = *val

	return i
}

func (i *UInt128) SetBytes(buf []byte) *UInt128 {
	l := len(buf)

	if l == 0 {
		i[0] = 0
		i[1] = 0
	}

	if l >= 8 {
		i[1] = binary.BigEndian.Uint64(buf[l-8 : l])
	} else {
		i[0] = 0
		i[1] = binary.BigEndian.Uint64(BytesAppendForward(buf, 8-l))
		return i
	}

	if l >= 16 {
		i[0] = binary.BigEndian.Uint64(buf[l-8-8 : l-8])
	} else {
		i[0] = binary.BigEndian.Uint64(BytesAppendForward(buf[:l-8], 16-l))
	}

	return i
}

func (i *UInt128) UInt64() uint64 {
	return i[1]
}
func (i *UInt128) Int() int {
	return int(i[1])
}

// Add creates result to the sum i+y
func (i *UInt128) Add(y *UInt128) (result UInt128) {
	var carry uint64
	result[1], carry = bits.Add64(i[1], y[1], carry)
	result[0], _ = bits.Add64(i[0], y[0], carry)
	return result
}

// AddUInt64 creates result to the sum i+y
func (i *UInt128) AddUInt64(y uint64) (result UInt128) {
	var carry uint64
	result[1], carry = bits.Add64(i[1], y, carry)
	result[0], _ = bits.Add64(i[0], 0, carry)
	return result
}

// Add creates result to the sum i+y, carry is 0
func (i *UInt128) AddOverflow(y *UInt128, carry uint64) (result UInt128, carryOut uint64) {
	result[1], carry = bits.Add64(i[1], y[1], carry)
	result[0], carry = bits.Add64(i[0], y[0], carry)
	return result, carry
}

// Sub creates result to the difference i-y
func (i *UInt128) Sub(y *UInt128) (result UInt128) {
	var borrow uint64
	result[1], borrow = bits.Sub64(i[1], y[1], borrow)
	result[0], _ = bits.Sub64(i[0], y[0], borrow)
	return result
}

// Sub creates result to the difference i-y, borrow is 0
func (i *UInt128) SubOverflow(y *UInt128, borrow uint64) (result UInt128, borrowOut uint64) {
	result[1], borrow = bits.Sub64(i[1], y[1], borrow)
	result[0], borrow = bits.Sub64(i[0], y[0], borrow)
	return result, borrow
}

// Mul creates result to the mult i * y
func (i *UInt128) Mul(y *UInt128) (result UInt128) {
	hiIt, loIt := bits.Mul64(i[1], y[1])
	result[0] = hiIt
	result[1] = loIt

	_, loIt = bits.Mul64(i[0], y[1])
	result[0], _ = bits.Add64(result[0], loIt, 0)

	_, loIt = bits.Mul64(i[1], y[0])
	result[0], _ = bits.Add64(result[0], loIt, 0)

	return result
}

// MulUInt64 creates result to the mult i * y
func (i *UInt128) MulUInt64(y uint64) (result UInt128) {
	hiIt, loIt := bits.Mul64(i[1], y)
	result[0] = hiIt
	result[1] = loIt

	_, loIt = bits.Mul64(i[0], y)
	result[0], _ = bits.Add64(result[0], loIt, 0)

	return result
}

// Mul creates result to the mult i * y
func (i *UInt128) MulOverflow(y *UInt128) (hi, lo UInt128) {

	hiIt, loIt := bits.Mul64(i[1], y[1])
	lo[0] = hiIt
	lo[1] = loIt

	hiIt, loIt = bits.Mul64(i[0], y[1])
	var carry uint64
	lo[0], carry = bits.Add64(lo[0], loIt, 0)
	hi[1], carry = bits.Add64(hi[1], hiIt, carry)
	hi[0], _ = bits.Add64(hi[0], 0, carry)

	hiIt, loIt = bits.Mul64(i[1], y[0])
	lo[0], carry = bits.Add64(lo[0], loIt, 0)
	hi[1], carry = bits.Add64(hi[1], hiIt, carry)
	hi[0], _ = bits.Add64(hi[0], 0, carry)

	hiIt, loIt = bits.Mul64(i[0], y[0])
	hi[1], carry = bits.Add64(hi[1], loIt, 0)
	hi[0], _ = bits.Add64(hi[0], hiIt, carry)

	return hi, lo
}

// MoveBitUp creates result << (move)
func (i *UInt128) MoveBitUp(move int) (result UInt128) {
	byMove := move / 64
	biMove := move - byMove*64

	if byMove > 0 {
		result[0] = i[1] << biMove
	} else {
		result[0] = i[0]<<biMove | i[1]>>(64-biMove)
		result[1] = i[1] << biMove
	}

	return result
}

// MoveBitDown creates result >> (move)
func (i *UInt128) MoveBitDown(move int) (result UInt128) {
	byMove := move / 64
	biMove := move - byMove*64

	if byMove > 0 {
		result[1] = i[0] >> biMove
	} else {
		result[1] = i[1]>>biMove | i[0]<<(64-biMove)
		result[0] = i[0] >> biMove
	}

	return result
}

// MoveBitDown creates result >> (1)
func (i *UInt128) MoveBitDown1Internal() {
	i[1] = i[1]>>1 | i[0]<<63
	i[0] = i[0] >> 1

}

// DivUint64 creates result to the quo = i / y; rem = i - div*y
func (i *UInt128) DivUint64(y uint64) (quo UInt128, rem UInt128) {
	var r uint64
	quo[0], r = bits.Div64(0, i[0], y)
	quo[1], rem[1] = bits.Div64(r, i[1], y)
	return quo, rem
}

// Div creates result to the quo = i / y; rem = i - div*y
func (i UInt128) Div(y *UInt128) (quo UInt128, rem UInt128) {
	if y.IsUint64() {
		return i.DivUint64(y.UInt64())
	}

	if i.Equal(y) {
		quo[1] = 1
		return quo, rem
	}

	if y.IsUint64() && y[1] == 1 {
		return i.Copy(), rem
	}

	if i.Less(y) {
		return quo, i.Copy()
	}

	rem = i

	k := bits.Len64(i[0] / y[0])
	// for !rem.Less(y.MoveBitUp(k).Link()) {
	// 	k++
	// }

	diffY := y.MoveBitUp(k)
	diffRes := (&UInt128{0, 1}).MoveBitUp(k)

	for !rem.Less(y) {
		// may be faster
		for rem.Less(&diffY) && k > 0 {
			//diffRes = diffRes.MoveBitDown(1)
			//diffY = diffY.MoveBitDown(1)
			diffRes.MoveBitDown1Internal()
			diffY.MoveBitDown1Internal()
			k--
		}

		rem = rem.Sub(&diffY)
		quo = quo.Add(&diffRes)
	}

	return quo, rem

}

func (i *UInt128) Text(base int) string {
	if base < 2 {
		panic("Base should be not less than 2")
	}
	if base > 61 {
		panic("Base should be not more than 61")
	}

	/*
		todo
			if base == 2 {

			}
			if base == 4 {

			}
			if base == 8 {

			}
			if base == 16 {

			}
	*/

	res := ""
	for !i.IsEmpty() {
		q, r := i.DivUint64(uint64(base))
		res = digits[int(r[1])] + res
		i = &q
	}

	if res == "" {
		res = digits[0]
	}

	return res
}
func (i *UInt128) String() string {
	return i.Text(10)
}

func (i *UInt128) FromText(value string, base int, ignoreFail bool) error {
	return i.FromTextByte([]byte(value), base, ignoreFail)
}
func (i *UInt128) FromTextByte(value []byte, base int, ignoreFail bool) error {
	res := &UInt128{}
	dv := uint64(base)

	for _, b := range value {
		d, ok := searchDigits[b]
		if base <= searchDigitsLowerSize && (!ok || d >= dv) {
			d, ok = searchDigitsLower[b]
		}

		if !ok || d >= dv {
			if !ignoreFail {
				return fmt.Errorf("char %v (%v) does not exists", string([]byte{b}), b)
			}
			continue
		}
		res = res.MulUInt64(dv).Link()
		res = res.AddUInt64(d).Link()
	}

	i.Set(res)

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *UInt128) UnmarshalJSON(input []byte) error {
	if input == nil {
		return nil
	}
	if len(input) == 0 {
		return fmt.Errorf("empty value")
	}
	if input[0] != '"' {
		return i.FromTextByte(input, 10, false)
	} else if len(input) < 3 {
		return fmt.Errorf("empty value (\")")
	}
	return i.FromTextByte(input[1:len(input)-1], 10, false)
}

//MarshalJSON() implements json.Marshaler.
func (i UInt128) MarshalJSON() ([]byte, error) {
	txt := i.Text(10)
	return []byte(txt), nil
}

// Power returns p = i^e
func (i *UInt128) Power(e *UInt128) (p *UInt128) {
	p = &UInt128{0, 1}

	for !e.IsEmpty() {
		if (e[1] & 1) != 0 {
			p = p.Mul(i).Link()
		}
		i = i.Mul(i).Link()
		e.MoveBitDown1Internal()
	}
	return p
}

// GetBit gets bit from pos (0..127)
func (i *UInt128) GetBit(pos int) bool {
	val := pos / 64

	pos = pos - val*64

	mask := uint64(1)
	mask = mask << pos

	return (i[1-val] & mask) > 0
}

// SetBit set bit as val in pos (0..127)
func (i *UInt128) SetBit(pos int, val bool) {
	if val {
		i.SetBit1(pos)
		return
	}
	i.SetBit0(pos)
}

// SetBit1 set bit as 1 in pos (0..127)
func (i *UInt128) SetBit1(pos int) {
	val := pos / 64

	pos = pos - val*64

	mask := uint64(1)
	mask = mask << pos

	if (i[1-val] & mask) > 0 {
		return
	}

	i[1-val] = i[1-val] + mask
}

// SetBit0 set bit as 0 in pos (0..127)
func (i *UInt128) SetBit0(pos int) {
	val := pos / 64

	pos = pos - val*64

	mask := uint64(1)
	mask = mask << pos

	if (i[1-val] & mask) == 0 {
		return
	}

	i[1-val] = i[1-val] - mask
}

// Root r => r^b=i when no resault then nearest r^b<=i
func (i *UInt128) Root(base uint64) (r UInt128) {
	iB := &UInt128{0, base}

	if base == 0 {
		panic("root base should be > 0")
	}
	if base == 1 {
		return i.Copy()
	}
	if i.IsEmpty() {
		return i.Copy()
	}
	if i.Equal(&UInt128{0, 1}) {
		return i.Copy()
	}

	l := i.BytesLen() * 8

	l = l/int(base) + 1

	for l >= 0 {
		r.SetBit1(l)
		pb := r.Power(iB)
		if pb.Equal(i) {
			return r
		}
		for !pb.Less(i) {
			r.SetBit0(l)
			l = l - 1
			if l < 0 {
				break
			}
			r.SetBit1(l)
			pb = r.Power(iB)
			if pb.Equal(i) {
				return r
			}
		}
		l = l - 1
	}

	return r
}

// Sqrt r => r^2=i when no resault then nearest r^2<=i
func (i *UInt128) Sqrt() (r UInt128) {
	return i.Root(2)
}
