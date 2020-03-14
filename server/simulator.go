package server

import "math/rand"

func simulateMonsterFight(E1, E2 int) (survivor int, victory bool) {
	P := getProba(E1, E2, false)
	// FIGHT
	if rand.Float64() < P {
		// Victory
		for i := 0; i < E1; i++ {
			if rand.Float64() < P {
				survivor++
			}
		}
		return survivor, true
	} else {
		// Loss
		for i := 0; i < E2; i++ {
			if rand.Float64() < (1 - P) {
				survivor++
			}
		}
		return survivor, false
	}
}

func simulateHumanFight(E1, E2 int) (survivor int, victory bool) {
	P := getProba(E1, E2, true)
	// FIGHT
	if rand.Float64() < P {
		// Victory
		for i := 0; i < E1+E2; i++ {
			if rand.Float64() < P {
				survivor++
			}
		}
		return survivor, true
	} else {
		// Loss
		for i := 0; i < E2; i++ {
			if rand.Float64() < (1 - P) {
				survivor++
			}
		}
		return survivor, false
	}
}

// getProba reimplements the getProba logic from Board.cs in the C# implementation
func getProba(E1, E2 int, involveHumans bool) float64 {
	// True by property
	if (involveHumans && E1 >= E2) || (!involveHumans && float64(E1) >= 1.5*float64(E2)) {
		return 1
	}

	if E1 == E2 {
		return 0.5
	}

	if E1 < E2 {
		return float64(E1) / (2 * float64(E2))
	} else {
		return (float64(E1) / float64(E2)) - 0.5
	}
}
