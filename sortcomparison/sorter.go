package main

import (
	"fmt"
	"sort"
)

type Person struct {
	name string
	age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s is %d years old", p.name, p.age)
}

func main() {
	people := []Person{
		Person{name: "Jason", age: 28},
		Person{name: "Sabah", age: 27},
	}

	sort.Slice(people, func(i, j int) bool {
		return people[i].age < people[j].age
	})

	fmt.Printf("%v\n", people)
}
