package main

import "testing"

func TestSum(t *testing.T) {
	// total := Sum(5, 5)
	// if total != 10 {
	// 	t.Errorf("sum was incorrect, got %d expected %d", total, 10)
	// }
	tables := []struct {
		a int
		b int
		n int
	}{
		{1, 2, 3},
		{2, 2, 4},
		{25, 26, 51},
	}
	for _, item := range tables {
		total := Sum(item.a, item.b)
		if total != item.n {
			t.Errorf("sum was incorrect, got %d expected %d", total, item.n)
		}
	}
}

func TestGetMax(t *testing.T) {
	tables := []struct {
		a int
		b int
		n int
	}{
		{4, 2, 4},
		{3, 2, 3},
		{2, 5, 5},
	}
	for _, item := range tables {
		max := GetMax(item.a, item.b)
		if max != item.n {
			t.Errorf("getMax was incorrect, got %d expected %d", max, item.n)
		}
	}
}
func TestFib(t *testing.T) {
	tables := []struct {
		a int
		n int
	}{
		{1, 1},
		{8, 21},
		{50, 12586269025},
	}
	for _, item := range tables {
		fib := Fibonacci(item.a)
		if fib != item.n {
			t.Errorf("Fibonacci was incorrect, got %d expected %d", fib, item.n)
		}
	}
}

// go test para correr los tests
// go test -cover para ver porcentaje de covertura
// go test -coverprofile=coverage.out para exportarlo a un archivo
// go test -cpuprofile=cpu.out para ver profiling de cpu en un archivo
// go  tool pprof cpu.out para analizarlo
// comando top para ver el mayor uso
// comando pdf para exportarlo a pdf, o web para visualizarlo en la web
