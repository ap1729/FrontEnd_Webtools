package perf

import (
	"errors"
	"math"
)

// Computes the pre- and post- processing SINR rxPows, and additionally returns
// remaining interference power (ROI) for the given array of signal powers (dBm).
//
// The function returns an array of three elements, as follows:
//
// [pre SINR, post SINR, ROI]
//
// Note: This function does not order the signal profile and computes results
// directly on the array in the same order of elements it is passed.
func sinr(rxPows []float64, intrCancelCount uint) ([]float64, error) {

	// Handling argument nil exception
	if rxPows == nil {
		return nil, errors.New(ARG_NIL)
	}
	if len(rxPows) == 0 {
		return nil, errors.New("The power array argument is empty.")
	}

	// Values is the recieved signal powers
	var vals = make([]float64, 3)
	var noise float64 = 1e-9 // Watts

	// Pre Processing SINR calculation
	var sum float64 = 0.0
	for i := 1; i < len(rxPows); i++ {
		sum += math.Pow(10, (rxPows[i]/10))
	}
	vals[0] = rxPows[0] - 10*math.Log10((sum+noise)) // Pre SINR

	// Post Processing SINR and ROI calculation
	sum = 0.0
	for i := int(intrCancelCount) + 1; i < len(rxPows); i++ {
		sum += math.Pow(10, (rxPows[i]/10))
	}
	vals[1] = rxPows[0] - 10*math.Log10((sum+noise)) // Post SINR
	vals[2] = 10 * math.Log10(sum)                   // ROI

	return vals, nil
}
