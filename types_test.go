package poly_test

import (
	"fmt"
	"testing"

	"github.com/ykalchevskiy/poly"
)

type T1 struct{}

func (T1) TypeName() string { return "T1" }

type T2 struct{}

func (T2) TypeName() string { return "T2" }

type T3 struct{}

func (T3) TypeName() string { return "T3" }

type T4 struct{}

func (T4) TypeName() string { return "T4" }

func ExampleTypeList() {
	var tl poly.TypeList[T1, poly.TypeList[T2, poly.TypeList[T3, poly.TypeListLast]]]

	for _, t := range tl.Types() {
		fmt.Println(t.Name)
	}

	// Output:
	// T1
	// T2
	// T3
}

func ExampleTypes1() {
	var tl poly.Types1[T1]

	for _, t := range tl.Types() {
		fmt.Println(t.Name)
	}

	// Output:
	// T1
}

func ExampleTypes2() {
	var tl poly.Types2[T1, T2]

	for _, t := range tl.Types() {
		fmt.Println(t.Name)
	}

	// Output:
	// T1
	// T2
}

func ExampleTypes3() {
	var tl poly.Types3[T1, T2, T3]

	for _, t := range tl.Types() {
		fmt.Println(t.Name)
	}

	// Output:
	// T1
	// T2
	// T3
}

func TestTypes4(t *testing.T) {
	var tl poly.Types4[T1, T2, T3, T4]

	for i, typ := range tl.Types() {
		if name := fmt.Sprintf("T%d", i+1); name != typ.Name {
			t.Errorf("expected %s, got %s", name, typ.Name)
		}
	}
}
