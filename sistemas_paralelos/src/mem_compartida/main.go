package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

// Consts
const (
	// Universal gravitational constant
	G float64 = 6.674e-11
	// 1 second time differential
	DT float64 = 1
	// Epsilon
	EPS float64 = 3e4
	// Number of Go Routines
	ROUTINES = 2
)

// Var declarations
var (
	N, iterations             int
	xpos, ypos, zpos          []float64
	xv, yv, zv                []float64
	forcesx, forcesy, forcesz []float64
	mass                      []float64
)

func main() {

	//Parameters handling
	args := os.Args[1:]
	if len(args) == 2 {
		var err1, err2 error
		N, err1 = strconv.Atoi(args[0])
		iterations, err2 = strconv.Atoi(args[1])
		if err1 != nil || err2 != nil {
			err := fmt.Errorf("program %s arguments must be integers", os.Args[0])
			log.Fatal(err)

		}
	} else {
		err := fmt.Errorf("program %s arguments must be N and Iterations", os.Args[0])
		log.Fatal(err)
	}
	// fmt.Printf("%d, %d\n", N, iterations)

	// Slices allocations

	xpos = make([]float64, N)
	ypos = make([]float64, N)
	zpos = make([]float64, N)
	xv = make([]float64, N)
	yv = make([]float64, N)
	zv = make([]float64, N)
	mass = make([]float64, N)
	forcesx = make([]float64, N*N)
	forcesy = make([]float64, N*N)
	forcesz = make([]float64, N*N)
	var xdistance, ydistance, zdistance float64
	var distance float64
	var fmagnitude float64
	var xvers, yvers, zvers float64
	var ax, ay, az, dvx, dvy, dvz, dpx, dpy, dpz float64
	var fx, fy, fz float64
	var sqrt_n int = int(math.Sqrt(float64(N)))
	if (sqrt_n * sqrt_n) != N {
		sqrt_n++
	}

	// Distance between bodies
	var dist int = 1000000

	// Bodies data
	for i := 0; i < N; i++ {
		xpos[i] = float64((i % sqrt_n) * dist)
		ypos[i] = float64(dist * i)
		zpos[i] = 5000
		mass[i] = 5.97e25
	}

	// Timer start
	start := time.Now()

	// Initilize wait groups
	var wg sync.WaitGroup
	var wg_intern sync.WaitGroup

	for it := 0; it < iterations; it++ {
		wg.Add(ROUTINES - 1)
		wg_intern.Add(ROUTINES)

		// Invoke go routines
		for r := 1; r < ROUTINES; r++ {
			go go_routine(r, &wg, &wg_intern)
		}

		// for i := 0; i < N; i++ {
		// 	fmt.Printf("%f, %f, %f\n", xpos[i], ypos[i], zpos[i])
		// }

		// Main Routine Job
		for i := 0; i < N/ROUTINES; i++ {
			for j := i + 1; j < N; j++ {
				// Distance vector with direction Body_i ---> Body_j
				xdistance = xpos[j] - xpos[i]
				ydistance = ypos[j] - ypos[i]
				zdistance = zpos[j] - zpos[i]
				// Distance between bodies magnitude
				distance = math.Sqrt((xdistance * xdistance) + (ydistance * ydistance) + (zdistance * zdistance))
				// Ignore collisions
				if distance < EPS*6 {
					// Components of the force vector exerted by j on i.
					forcesx[i*N+j] = 0.00
					forcesy[i*N+j] = 0.00
					forcesz[i*N+j] = 0.00
					// Opposite forces (symmetric), exerted by i on j.
					forcesx[j*N+i] = 0.00
					forcesy[j*N+i] = 0.00
					forcesz[j*N+i] = 0.00
				} else {
					// Magnitude of the gravitational attractive force
					// F = G * ((m1 * m2) / D??)
					fmagnitude = G * ((mass[i] * mass[j]) / (distance * distance))
					// Components of the unity vector
					xvers = xdistance / distance
					yvers = ydistance / distance
					zvers = zdistance / distance
					index_i := i*N + j
					index_j := j*N + i
					// Components of the force vector exerted by j on i.
					forcesx[index_i] = (xvers) * fmagnitude
					forcesy[index_i] = (yvers) * fmagnitude
					forcesz[index_i] = (zvers) * fmagnitude
					// Opposite forces (symmetric), exerted by i on j.
					forcesx[index_j] = ((-1) * xvers) * fmagnitude
					forcesy[index_j] = ((-1) * yvers) * fmagnitude
					forcesz[index_j] = ((-1) * zvers) * fmagnitude
				}
			}
		}

		// Intern barrier, waiting for the calculation of all forces
		wg_intern.Done()
		wg_intern.Wait()

		for i := 0; i < N/ROUTINES; i++ {
			fx = 0.0
			fy = 0.0
			fz = 0.0
			// Aceleration vector
			// |M = F * A| ---> |A = F / M|
			left_index := i * N
			for b := 0; b < N; b++ {
				if b != i {
					index := left_index + b
					fx += forcesx[index]
					fy += forcesy[index]
					fz += forcesz[index]
				}
			}

			// Acceleration difference
			ax = fx / mass[i] * DT
			ay = fy / mass[i] * DT
			az = fz / mass[i] * DT

			// Speed ??????difference delta v
			dvx = ax * DT
			dvy = ay * DT
			dvz = az * DT

			// Position difference
			// Leapfrog scheme (initial speed + (1/2 * new speed)) * DT
			dpx = (xv[i] + (dvx / 2)) * DT
			dpy = (yv[i] + (dvy / 2)) * DT
			dpz = (zv[i] + (dvz / 2)) * DT

			// Position update
			xpos[i] += dpx
			ypos[i] += dpy
			zpos[i] += dpz

			// Velocity update
			xv[i] = (xv[i] + dvx) / 2
			yv[i] = (yv[i] + dvy) / 2
			zv[i] = (zv[i] + dvz) / 2
		}

		// Wait for all the bodies to move
		wg.Wait()
	}

	// Time stop and show results
	elapsed := time.Since(start)
	fmt.Printf("Tiempo en segundos %v\n", elapsed)
}

// BACKGROUND GO ROUTINE
func go_routine(id int, wg, wg_intern *sync.WaitGroup) {
	defer wg.Done()

	// Vars declarations
	var xdistance, ydistance, zdistance float64
	var distance float64
	var fmagnitude float64
	var xvers, yvers, zvers float64
	var ax, ay, az, dvx, dvy, dvz, dpx, dpy, dpz float64
	var fx, fy, fz float64

	for i := id * N / ROUTINES; i < (id+1)*N/ROUTINES; i++ {
		// fmt.Printf("%f, %f, %f\n", xpos[i], ypos[i], zpos[i])
		for j := i + 1; j < N; j++ {
			// Distance vector with direction Body_i ---> Body_j
			xdistance = xpos[j] - xpos[i]
			ydistance = ypos[j] - ypos[i]
			zdistance = zpos[j] - zpos[i]
			// Distance between bodies magnitude
			distance = math.Sqrt((xdistance * xdistance) + (ydistance * ydistance) + (zdistance * zdistance))
			// Ignore collisions
			if distance < EPS*6 {
				// Components of the force vector exerted by j on i.
				forcesx[i*N+j] = 0.00
				forcesy[i*N+j] = 0.00
				forcesz[i*N+j] = 0.00
				// Opposite forces (symmetric), exerted by i on j.
				forcesx[j*N+i] = 0.00
				forcesy[j*N+i] = 0.00
				forcesz[j*N+i] = 0.00
			} else {
				// Magnitude of the gravitational attractive force
				// F = G * ((m1 * m2) / D??)
				fmagnitude = G * ((mass[i] * mass[j]) / (distance * distance))
				// Components of the unity vector.
				xvers = xdistance / distance
				yvers = ydistance / distance
				zvers = zdistance / distance
				index_i := i*N + j
				index_j := j*N + i
				// Componentes del vector fuerza ejercida por j sobre i.
				forcesx[index_i] = (xvers) * fmagnitude
				forcesy[index_i] = (yvers) * fmagnitude
				forcesz[index_i] = (zvers) * fmagnitude
				// Opposite forces (symmetric), exerted by i on j.
				forcesx[index_j] = ((-1) * xvers) * fmagnitude
				forcesy[index_j] = ((-1) * yvers) * fmagnitude
				forcesz[index_j] = ((-1) * zvers) * fmagnitude
			}
		}
	}

	// Intern barrier, waiting for the calculation of all forces
	wg_intern.Done()
	wg_intern.Wait()

	for i := id * N / ROUTINES; i < (id+1)*N/ROUTINES; i++ {
		fx = 0.0
		fy = 0.0
		fz = 0.0
		// Aceleration vector
		// |M = F * A| ---> |A = F / M|
		left_index := i * N
		for b := 0; b < N; b++ {
			if b != i {
				index := left_index + b
				fx += forcesx[index]
				fy += forcesy[index]
				fz += forcesz[index]
			}
		}

		// Acceleration difference
		ax = fx / mass[i] * DT
		ay = fy / mass[i] * DT
		az = fz / mass[i] * DT

		// Speed ??????difference delta v
		dvx = ax * DT
		dvy = ay * DT
		dvz = az * DT

		// Position difference
		// Leapfrog scheme (initial speed + (1/2 * new speed)) * DT
		dpx = (xv[i] + (dvx / 2)) * DT
		dpy = (yv[i] + (dvy / 2)) * DT
		dpz = (zv[i] + (dvz / 2)) * DT

		// Position update
		xpos[i] += dpx
		ypos[i] += dpy
		zpos[i] += dpz

		// Velocity update
		xv[i] = (xv[i] + dvx) / 2
		yv[i] = (yv[i] + dvy) / 2
		zv[i] = (zv[i] + dvz) / 2
	}
}
