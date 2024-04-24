package helper

import "math/rand"

//function that returns a random number int the closed interval [0,1]
func GetRandomNumber() float64 {
	r := rand.Float64()
	if r < 0.01 {
		r = GetRandomNumber()
	}
	return r
}
