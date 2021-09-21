package main

import "fmt"

func main() {
	c := make(chan int, 3)

	c <- 1
	c <- 2
	c <- 3
	//c <- 4
	x := 5 / 0
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
	x++
	//fmt.Println(<-c)
}
