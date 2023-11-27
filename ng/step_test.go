package ng

import "fmt"

func ExampleSteps_String() {
	steps := Steps{
		{RingGroup: Inner, Count: 1},
		{RingGroup: OuterMiddle, Count: 3},
		{RingGroup: MiddleInner, Count: 2},
		{RingGroup: Outer, Count: 1},
		{RingGroup: Inner, Count: 1},
		{RingGroup: Middle, Count: 0},
	}
	fmt.Println(steps.String())
	// Output:
	// i2,mi2,o1,om3
}
