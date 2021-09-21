package main

import (
	pk "curso_go_basico/src/mypackage"
	"fmt"
)

func main() {
	// Constantes
	const pi float64 = 3.14
	const pi2 = 3.14
	fmt.Println(pi, pi2)

	// Variables
	base := 12
	var altura int = 14
	var area int
	area = base * altura
	fmt.Println(area)

	// Zero values
	var a int
	var b float64
	var c string
	var d bool
	fmt.Println(a, b, c, d)

	fmt.Println(holaMundo("Ariel"))

	// Array
	var array [4]int
	array[0] = 1
	array[1] = 2

	fmt.Println(array, len(array), cap(array))

	// Slice

	slice := []int{0, 1, 2, 3, 4, 5, 6}

	fmt.Println(slice, len(slice), cap(slice))

	// Metodos en el slice
	fmt.Println(slice[0])
	fmt.Println(slice[:3])
	fmt.Println(slice[2:4])
	fmt.Println(slice[4:])

	// Append
	slice = append(slice, 7)
	fmt.Println(slice)

	// Append nueva lista

	newSlice := []int{8, 9, 10}
	slice = append(slice, newSlice...)
	fmt.Println(slice)

	var myCar pk.CarPublic
	myCar.Brand = "Ferrari"
	fmt.Println(myCar)

	// var myOtherCar pk.carPrivate
	// myOtherCar.brand = "Ferrari"
	// fmt.Println(myOtherCar)

	pk.PrintMessagePublic("Caca")
	// pk.printMessagePrivate("Caca")
}

func holaMundo(a string) string {
	return fmt.Sprintf("Hola %s", a)
}
