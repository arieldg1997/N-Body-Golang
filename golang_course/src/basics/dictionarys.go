package main

import "fmt"

func main() {
	m := make(map[string]int)

	m["Jose"] = 14
	m["Ariel"] = 23
	fmt.Println(m)

	for i, v := range m {
		fmt.Println(i, v)
	}

	value, ok := m["Unknown"]
	fmt.Println(value, ok)

}
