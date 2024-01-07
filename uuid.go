package ints

import (
	"fmt"
)

type Uuid struct {
	UInt128
}

// Less returns true when i < val
func (i Uuid) Less(val Uuid) bool {
	return i.UInt128.Less(val.UInt128)
}

// Equal returns true when i == val
func (i Uuid) Equal(val Uuid) bool {
	return i == val
}

func (i Uuid) Link() *Uuid {
	return &i
}

func (i Uuid) AsUUID() string {
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

// MarshalJSON() implements json.Marshaler.
func (i Uuid) MarshalJSON() ([]byte, error) {
	txt := "\"" + i.AsUUID() + "\""
	return []byte(txt), nil
}

func (i Uuid) Text(base int) string {
	return i.UInt128.Text(base)
}

func UuidFromText(value string, base int, ignoreFail bool) (i Uuid, err error) {
	err = i.FromText(value, base, ignoreFail)

	return i, err
}

func UuidFromTextByte(value []byte, base int, ignoreFail bool) (i Uuid, err error) {
	err = i.FromTextByte(value, base, ignoreFail)

	return i, err
}

func UuidFromTextMust(value string, base int, ignoreFail bool) (i Uuid) {
	err := i.FromText(value, base, ignoreFail)

	if err != nil {
		panic(err)
	}

	return i
}

func UuidFromTextByteMust(value []byte, base int, ignoreFail bool) (i Uuid) {
	err := i.FromTextByte(value, base, ignoreFail)

	if err != nil {
		panic(err)
	}

	return i
}
