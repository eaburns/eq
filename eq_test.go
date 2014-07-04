package eq

import (
	"math/big"
	"testing"
)

type s struct {
	x, Y, Z int
}

type i struct {
	X interface{}
}

type bigInt struct{ *big.Int }

func (a bigInt) Eq(c interface{}) bool {
	b, ok := c.(bigInt)
	return ok && a.Cmp(b.Int) == 0
}

func TestDeep(t *testing.T) {
	tests := []struct {
		u, v interface{}
		eq   bool
	}{
		// Nil an non-nill
		{nil, nil, true},
		{1, nil, false},
		{nil, 1, false},

		// Bool
		{true, true, true},
		{true, false, false},

		// Int
		{0, 0, true},
		{0, 1, false},

		// Unsigned int
		{uint(0), uint(0), true},
		{uint(0), uint(1), false},

		// Float
		{0.0, 0.0, true},
		{0.0, 1.0, false},

		// Complex
		{0.0 + 3i, 0.0 + 3i, true},
		{0.0 + 3i, 1.0 + 3i, false},
		{0.0 + 3i, 0.0 + 4, false},

		// Equal-typed arrays
		{[...]int{}, [...]int{}, true},
		{[...]int{4, 5, 6}, [...]int{4, 5, 6}, true},
		{[...]int{4, 5, 6}, [...]int{5, 6, 7}, false},

		// Slices
		{[]int{}, []int{}, true},
		{[]int{4, 5, 6}, []int{4, 5, 6}, true},
		{[]int{4, 5, 6}, []int{5, 6, 7}, false},
		{[]int{}, []int{4, 5, 6, 7}, false},
		{[]int{4, 5, 6, 7}, []int{}, false},
		{[]int{4, 5, 6}, []int{4, 5, 6, 7}, false},

		// Pointers
		{&[]int{}, &[]int{}, true},
		{&[]int{4, 5, 6}, &[]int{4, 5, 6}, true},
		{&[]int{4, 5, 6}, &[]int{5, 6, 7}, false},
		{&[]int{}, &[]int{4, 5, 6, 7}, false},
		{&[]int{4, 5, 6, 7}, &[]int{}, false},
		{&[]int{4, 5, 6}, &[]int{4, 5, 6, 7}, false},

		// Structs
		{s{x: 0, Y: 1, Z: 2}, s{x: 0, Y: 1, Z: 2}, true},
		{s{x: 0, Y: 1, Z: 2}, s{x: 1, Y: 1, Z: 2}, true},
		{s{x: 0, Y: 1, Z: 2}, s{x: 0, Y: 0, Z: 2}, false},
		{s{x: 0, Y: 1, Z: 2}, s{x: 0, Y: 1, Z: 3}, false},

		// Interfaces—wrap in i{} to make values a nested interface{} type
		{i{0}, i{0}, true},
		{i{i{0}}, i{i{0}}, true},
		{i{0}, i{1}, false},
		{i{i{0}}, i{i{1}}, false},

		{bigInt{big.NewInt(5)}, bigInt{big.NewInt(5)}, true},
		{bigInt{big.NewInt(5)}, bigInt{big.NewInt(6)}, false},
	}
	for _, test := range tests {
		eq := Deep(test.u, test.v)
		if eq == test.eq {
			continue
		}
		t.Errorf("expected ExportedDeepEqual(\n%v,\n%v\n) == %t, got %t",
			test.u, test.v, test.eq, eq)
	}
}
