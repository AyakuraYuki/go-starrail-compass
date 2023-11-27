package ng

import "fmt"

func Example_parseRingGroup() {
	rg, err := parseRingGroup("o")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s, %s", rg.Name(), rg.ShortName())
	// Output:
	// Outer, o
}

func Example_parseRingGroups() {
	rgs, err := parseRingGroups("mo,io,i")
	if err != nil {
		panic(err)
	}
	for _, rg := range rgs {
		fmt.Printf("%s, %s\n", rg.Name(), rg.ShortName())
	}
	// Output:
	// OuterMiddle, om
	// OuterInner, oi
	// Inner, i
}

func Example_parseRing() {
	tests := []string{
		"3+2",
		"0-1",
	}
	for _, test := range tests {
		ring, err := parseRing(test)
		if err != nil {
			panic(err)
		}
		fmt.Printf("location: %d, speed: %+d\n", ring.Location, ring.Speed)
	}
	// Output:
	// location: 3, speed: +2
	// location: 0, speed: -1
}
