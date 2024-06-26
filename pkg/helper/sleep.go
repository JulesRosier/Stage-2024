package helper

import (
	"math"
	"math/rand"
	"time"
)

// warpSpeed is a boolean that determines whether the sleep durations are in real time or in warp speed
var warpSpeed = true

func GenerateNormalDuration(mean, stdDev float64) float64 {
	// Generate two random numbers from a uniform distribution
	u1 := rand.Float64()
	u2 := rand.Float64()

	// Apply Box-Muller transform to get values from a standard normal distribution
	z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)

	// Scale and shift the values to get values from a normal distribution with given mean and standard deviation
	duration := mean + z0*stdDev

	return duration
}

func RandMinutes(mean, stdDev int) time.Duration {
	t := time.Minute*time.Duration(GenerateNormalDuration(float64(mean), float64(stdDev))) + time.Minute*30
	if warpSpeed {
		t = t / 60
	}
	return t
}
