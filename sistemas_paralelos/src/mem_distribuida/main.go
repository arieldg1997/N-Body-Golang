package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

// Communicational struct
type channels struct {
	xpos, ypos, zpos, mass chan []float64
}

func main() {

	// Var declarations
	var (
		N, iterations                                         int
		xdistance, ydistance, zdistance, distance, fmagnitude float64
		xvers, yvers, zvers                                   float64
		index_i, index_j                                      int
		ax, ay, az, dvx, dvy, dvz, dpx, dpy, dpz, fx, fy, fz  float64
		xv, yv, zv                                            []float64
		sqrt_n                                                int
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

	sqrt_n = int(math.Sqrt(float64(N)))
	if (sqrt_n * sqrt_n) != N {
		sqrt_n++
	}

	// fmt.Printf("%d, %d\n", N, iterations)

	// Slices allocations
	xpos := make([]float64, N)
	ypos := make([]float64, N)
	zpos := make([]float64, N)
	mass := make([]float64, N)
	forcesx := make([]float64, N*N)
	forcesy := make([]float64, N*N)
	forcesz := make([]float64, N*N)
	xv = make([]float64, N)
	yv = make([]float64, N)
	zv = make([]float64, N)

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

	// Set senders channels
	var snd_chans [ROUTINES]channels
	for i := range snd_chans {
		ch_pos_x := make(chan []float64, 1)
		ch_pos_y := make(chan []float64, 1)
		ch_pos_z := make(chan []float64, 1)
		ch_masses := make(chan []float64, 1)
		chans_literal := channels{xpos: ch_pos_x, ypos: ch_pos_y, zpos: ch_pos_z, mass: ch_masses}
		snd_chans[i] = chans_literal
	}

	// Set receptions channels
	var rcv_chans [ROUTINES]channels
	for i := range rcv_chans {
		ch_pos_x := make(chan []float64, 1)
		ch_pos_y := make(chan []float64, 1)
		ch_pos_z := make(chan []float64, 1)
		ch_masses := make(chan []float64, 1)
		chans_literal := channels{xpos: ch_pos_x, ypos: ch_pos_y, zpos: ch_pos_z, mass: ch_masses}
		rcv_chans[i] = chans_literal
	}

	// Send masses to go routines
	for i := 0; i < ROUTINES; i++ {
		snd_chans[i].mass <- mass[:]
	}

	// Invoke go routines
	for r := 1; r < ROUTINES; r++ {
		go go_routine(r, snd_chans[r], rcv_chans[r], N, ROUTINES, iterations)
	}

	for it := 0; it < iterations; it++ {
		// for i := 0; i < N; i++ {
		// 	fmt.Printf("%f, %f, %f\n", xpos[i], ypos[i], zpos[i])
		// }
		// Send positions
		for i := 1; i < ROUTINES; i++ {
			snd_chans[i].xpos <- xpos[:]
			snd_chans[i].ypos <- ypos[:]
			snd_chans[i].zpos <- zpos[:]
		}

		// Main Routine job
		for i := 0; i < N/ROUTINES; i++ {
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
					// F = G * ((m1 * m2) / D²)
					fmagnitude = G * ((mass[i] * mass[j]) / (distance * distance))
					// Components of the unity vector
					xvers = xdistance / distance
					yvers = ydistance / distance
					zvers = zdistance / distance
					index_i = i*N + j
					index_j = j*N + i
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

		// Recive and replace forces matrix
		for r := 1; r < ROUTINES; r++ {
			fx := <-rcv_chans[r].xpos
			fy := <-rcv_chans[r].ypos
			fz := <-rcv_chans[r].zpos
			for i := r * N / ROUTINES; i < (r+1)*N/ROUTINES; i++ {
				for j := i + 1; j < N; j++ {
					index_i := i*N + j
					forcesx[index_i] = fx[index_i]
					forcesy[index_i] = fy[index_i]
					forcesz[index_i] = fz[index_i]
				}
			}
			for i := r * N / ROUTINES; i < (r+1)*N/ROUTINES; i++ {
				for j := i + 1; j < N; j++ {
					index_j := j*N + i
					forcesx[index_j] = fx[index_j]
					forcesy[index_j] = fy[index_j]
					forcesz[index_j] = fz[index_j]
				}
			}
		}
		// Send forces matrix
		for i := 1; i < ROUTINES; i++ {
			snd_chans[i].xpos <- forcesx[:]
			snd_chans[i].ypos <- forcesy[:]
			snd_chans[i].zpos <- forcesz[:]
		}

		// Main Routine job
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

			// Speed ​​difference delta v
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

		// Recive and replace new positions
		for i := 1; i < ROUTINES; i++ {
			copy(forcesx[i*N/ROUTINES:(i+1)*N/ROUTINES], (<-rcv_chans[i].xpos))
			copy(forcesy[i*N/ROUTINES:(i+1)*N/ROUTINES], (<-rcv_chans[i].ypos))
			copy(forcesz[i*N/ROUTINES:(i+1)*N/ROUTINES], (<-rcv_chans[i].zpos))
		}
	}

	// Time stop and show results
	elapsed := time.Since(start)
	fmt.Printf("Tiempo en segundos %v\n", elapsed)
}

// BACKGROUND GO ROUTINE
func go_routine(id int, rcv_chans, snd_chans channels, N, ROUTINES, iterations int) {
	// Constants
	const (
		// Universal gravitational constant
		G float64 = 6.674e-11
		// 1 second time differential
		DT float64 = 1
		// Epsilon
		EPS float64 = 3e4
	)

	// Vars declarations
	var index_i, index_j int
	var xdistance, ydistance, zdistance, distance, fmagnitude float64
	var xvers, yvers, zvers float64
	var ax, ay, az, dvx, dvy, dvz, dpx, dpy, dpz, fx, fy, fz float64
	var forcesx, forcesy, forcesz []float64
	var xpos, ypos, zpos, mass []float64
	var xv, yv, zv []float64
	var cota_inf int = id * N / ROUTINES
	var cota_sup int = (id + 1) * N / ROUTINES

	// Slices allocation
	forcesx = make([]float64, N*N)
	forcesy = make([]float64, N*N)
	forcesz = make([]float64, N*N)
	xv = make([]float64, N)
	yv = make([]float64, N)
	zv = make([]float64, N)

	// Recieve mass from main routine
	mass = <-rcv_chans.mass

	for it := 0; it < iterations; it++ {
		// Recieve positions from main routine
		xpos = <-rcv_chans.xpos
		ypos = <-rcv_chans.ypos
		zpos = <-rcv_chans.zpos
		for i := cota_inf; i < cota_sup; i++ {
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
					// F = G * ((m1 * m2) / D²)
					fmagnitude = G * ((mass[i] * mass[j]) / (distance * distance))
					// Components of the unity vector.
					xvers = xdistance / distance
					yvers = ydistance / distance
					zvers = zdistance / distance
					index_i = i*N + j
					index_j = j*N + i
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

		// Send partial forces to main routine
		snd_chans.xpos <- forcesx[:]
		snd_chans.ypos <- forcesy[:]
		snd_chans.zpos <- forcesz[:]

		// Recive all forces from main routine
		forcesx = <-rcv_chans.xpos
		forcesy = <-rcv_chans.ypos
		forcesz = <-rcv_chans.zpos
		for i := cota_inf; i < cota_sup; i++ {
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

			// Speed ​​difference delta v
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

		// Send positions to main routine
		snd_chans.xpos <- xpos[cota_inf:cota_sup]
		snd_chans.ypos <- ypos[cota_inf:cota_sup]
		snd_chans.zpos <- zpos[cota_inf:cota_sup]
	}
}
