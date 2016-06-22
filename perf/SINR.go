package perf

import (
	// "fmt"
	"math"
)

// TODO: Cleanup function, fix bug of dBm -> linear conversion
// (Technically, it doesnt matter when we convert it back, e-3 and e+3 cancel)

// Computes the pre- and post- processing SINR values, and additionally returns
// remaining interference power (ROI) for the given array of signal powers (dBm).
//
// The function returns an array of three elements, as follows:
//
// [pre SINR, post SINR, ROI]
//
// Note: This function does not order the signal profile and computes results
// directly on the array in the same order of elements it is passed.
func sinr(values []float64, intrCancelCount uint) []float64 {
	// Values is recieved signal powers
	var calc = []float64{}

	//Pre Processing SINR
	var sum float64 = 0.0
	for i := 1; i < len(values); i++ {
		sum += math.Pow(10, (values[i] / 10))
	}
	calc = append(calc, values[0]-10*math.Log10(sum)) //Pre SINR

	//Post Processing SINR
	sum = 0.0
	for i := int(intrCancelCount) + 1; i < len(values); i++ {
		sum += math.Pow(10, (values[i] / 10))
	}
	calc = append(calc, values[0]-10*math.Log10(sum)) //Post SINR
	calc = append(calc, 10*math.Log10(sum))           //ROI
	return calc
}
