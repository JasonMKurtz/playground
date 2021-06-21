package app

import "fmt"

func Greet(name string) (string, error) {
	return fmt.Sprintf("Hello %s", name), nil
}
