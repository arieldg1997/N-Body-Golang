package main

import "fmt"

func main() {
	slice := []string{"hola", "mundo", "!"}
	for i, elem := range slice {
		fmt.Println(i, elem)
	}
	for _, elem := range slice {
		fmt.Println(elem)
	}
	for i := range slice {
		fmt.Println(i)
	}

}
