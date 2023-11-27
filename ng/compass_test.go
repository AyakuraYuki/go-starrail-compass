package ng

import (
	"fmt"
	"testing"
)

func ExampleRing_String() {
	ring := &Ring{Location: 1, Speed: 3}
	fmt.Println(ring.String())
	ring = &Ring{Location: 0, Speed: -2}
	fmt.Println(ring.String())
	// Output:
	// 1+3
	// 0-2
}

func ExampleCompass_String() {
	compass := &Compass{
		InnerRing:  Ring{Location: 0, Speed: 1},
		MiddleRing: Ring{Location: 4, Speed: -4},
		OuterRing:  Ring{Location: 0, Speed: 2},
		RingGroups: []RingGroup{
			OuterInner,
			OuterMiddle,
			MiddleInner,
		},
	}
	fmt.Println(compass.String())
	// Output:
	// 0+2,4-4,0+1/mi,oi,om
}

func TestParseCompass(t *testing.T) {
	compass, err := ParseCompass("0+2,3-3,0+3/mi,om,oi")
	if err != nil {
		t.Fatal(err)
	}

	expectedCompass := Compass{
		OuterRing:  Ring{Location: 0, Speed: 2},
		MiddleRing: Ring{Location: 3, Speed: -3},
		InnerRing:  Ring{Location: 0, Speed: 3},
		RingGroups: []RingGroup{MiddleInner, OuterMiddle, OuterInner},
	}

	if compass.String() != expectedCompass.String() {
		t.Fatalf("unexpected result: %#v (expected: %#v)", compass.String(), expectedCompass.String())
	}
}
