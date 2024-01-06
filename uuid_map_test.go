package ints

import (
	"encoding/json"
	"testing"
)

func TestUuidMap(t *testing.T) {
	m := UuidMap[int]{
		NextUUID(): 1,
		NextUUID(): 2,
	}

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Error(err)
	}

	var m2 UuidMap[int]

	err = json.Unmarshal(b, &m2)
	if err != nil {
		t.Error(err)
	}

	for k, v := range m {
		v2, ok := m2[k]

		if !ok {
			t.Fatalf("not found key %v after unmarshal ```%v```", k, string(b))
		}

		if v != v2 {
			t.Fatalf("values are different %v, %v for key: `%v` after unmarshal ```%v```", v, v2, k, string(b))
		}
	}
}
