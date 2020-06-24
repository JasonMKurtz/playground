package main

import (
	"fmt"
)

func main() {
	j := &JRegex{`(?P<g>.+)`, "foobar"}
	fmt.Printf("Named groups: %v\n", j.GetNamedGroups()["g"])
	fmt.Printf("Numbered groups: %v\n", j.GetGroups()[0])
}
