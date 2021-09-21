package mypackage

import "fmt"

// CarPublic Car con acceso publico
type CarPublic struct {
	Brand string
	Year  int
}

// carPrivate Car con acceso privado
type carPrivate struct {
	brand string
	year  int
}

// PrintMessagePublic imprime de forma publica
func PrintMessagePublic(text string) {
	fmt.Println(text)
}

// printMessagePrivate imprime de forma privado
func printMessagePrivate(text string) {
	fmt.Println(text)
}
