package gql

import "testing"

func TestUnmarshalUint64(t *testing.T) {
	type test struct {
		In  interface{}
		Out uint64
		Err bool
	}
	tt := []test{
		{"", 0, false},
		{"0", 0, false},
		{"1", 1, false},
		{`""`, 0, false},
	}

	for i, v := range tt {
		r, err := UnmarshalUint64(v.In)

		if err != nil && !v.Err {
			t.Fatalf("Test %d: Failed %v", i, err)
		}

		if r != v.Out {
			t.Fatalf("Test %d: Failed %d != %d", i, r, v.Out)
		}
	}
}

func TestUnmarshalUint32(t *testing.T) {
	type test struct {
		In  interface{}
		Out uint32
	}
	tt := []test{
		{"", 0},
		{"0", 0},
		{"1", 1},
		{`""`, 0},
	}

	for i, v := range tt {
		r, err := UnmarshalUint32(v.In)
		if err != nil {
			t.Fatalf("Test %d: Failed %v", i, err)
		}
		if r != v.Out {
			t.Fatalf("Test %d: Failed %d != %d", i, r, v.Out)
		}
	}
}
