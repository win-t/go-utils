package structural_test

import (
	"fmt"
	"testing"

	"github.com/win-t/go-utils/structural"
)

func ExampleAssign() {
	type aTyp struct {
		A int
		B string
		C string
		d int
	}

	type bTyp struct {
		A int
		B string
		C int
		d int
	}

	a := aTyp{1, "2", "3", 4}
	b := bTyp{10, "11", 13, 14}

	structural.Assign(&a, b)

	fmt.Println(a.A, a.B, a.C, a.d)
	// Output: 10 11 3 4
}

func TestAssign(t *testing.T) {
	type ta struct {
		A string
		B int
	}

	type tb struct {
		A string
		B int
	}

	a := ta{A: "Hello", B: 10}
	b := tb{A: "Hai", B: 12}

	structural.Assign(&a, b)

	if a.A != b.A || a.B != b.B {
		t.FailNow()
	}
}

func TestAssignViaPointer(t *testing.T) {
	type ta struct {
		A string
		B int
	}

	type tb struct {
		A string
		B int
	}

	a := ta{A: "Hello", B: 10}
	b := tb{A: "Hai", B: 12}

	structural.Assign(&a, &b)

	if a.A != b.A || a.B != b.B {
		t.FailNow()
	}
}

func TestAssignMissingLeft(t *testing.T) {
	type ta struct {
		A string
	}

	type tb struct {
		A string
		B int
	}

	a := ta{A: "Hello"}
	b := tb{A: "Hai", B: 12}

	structural.Assign(&a, b)

	if a.A != b.A {
		t.FailNow()
	}
}

func TestAssignMissingRight(t *testing.T) {
	type ta struct {
		A string
		B int
	}

	type tb struct {
		A string
	}

	a := ta{A: "Hello", B: 10}
	b := tb{A: "Hai"}

	structural.Assign(&a, b)

	if a.A != b.A {
		t.FailNow()
	}
}

func TestAssignNonExportedField(t *testing.T) {
	type ta struct {
		A string
		b int
	}
	type tb struct {
		A string
		b int
	}

	a := ta{A: "Hello", b: 10}
	b := tb{A: "Hai", b: 12}

	structural.Assign(&a, b)

	if a.A != b.A || a.b == b.b {
		t.FailNow()
	}
}

func TestAssignDifferentType(t *testing.T) {
	type ta struct {
		A string
		B string
	}

	type tb struct {
		A string
		B int
	}

	a := ta{A: "Hello", B: "10"}
	b := tb{A: "Hai", B: 12}

	structural.Assign(&a, b)

	if a.A != b.A {
		t.FailNow()
	}
}

func TestAssingNonPointer(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.FailNow()
		}
	}()

	type ta struct {
		A string
		B int
	}

	type tb struct {
		A string
		B int
	}

	a := ta{A: "Hello", B: 10}
	b := tb{A: "Hai", B: 12}

	structural.Assign(a, b)
}

func TestAssignInvalidSource(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.FailNow()
		}
	}()

	type ta struct {
		A string
		B int
	}

	a := ta{A: "Hello", B: 10}

	structural.Assign(&a, nil)
}
