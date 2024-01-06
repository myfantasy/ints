package ints

import "encoding/json"

type UuidMap[T any] map[Uuid]T

var _ json.Marshaler = UuidMap[struct{}]{}

func (mT *UuidMap[T]) UnmarshalJSON(input []byte) error {
	var tmp map[string]T

	err := json.Unmarshal(input, &tmp)
	if err != nil {
		return err
	}

	if tmp == nil {
		return nil
	}

	m := make(UuidMap[T])
	*mT = m

	for k, v := range tmp {
		var kT Uuid
		err = kT.FromText(k, 16, true)
		if err != nil {
			return err
		}

		m[kT] = v
	}

	return nil
}

func (m UuidMap[T]) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}

	res := []byte("{")
	e := false
	for k, v := range m {
		if e {
			res = append(res, []byte(",")...)
		}

		b, err := k.MarshalJSON()
		if err != nil {
			return nil, err
		}
		res = append(res, b...)
		res = append(res, []byte(":")...)
		b, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
		res = append(res, b...)

		e = true
	}
	res = append(res, []byte("}")...)
	return res, nil
}
