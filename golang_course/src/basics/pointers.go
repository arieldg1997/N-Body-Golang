package main

import "fmt"

type pc struct {
	ram   int
	disk  int
	brand string
}

func (myPc pc) ping() {
	fmt.Println(myPc.brand, "Pong")
}

func (myPc *pc) duplicateRam() {
	myPc.ram = myPc.ram * 2
}

func (myPc pc) String() string {
	return fmt.Sprintf("Tengo %d GB de Ram, %d GB de disco y es una %s", myPc.ram, myPc.disk, myPc.brand)
}

func main() {
	a := 50
	b := &a
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(*b)
	*b = 100
	fmt.Println(a)
	// a y b son alias

	myPc := pc{ram: 16, disk: 200, brand: "hp"}
	myPc.ping()
	fmt.Println(myPc)
	myPc.duplicateRam()
	fmt.Println(myPc)

}
