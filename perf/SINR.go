package perf

import (
	//"fmt"
	"math"
)

// Function for pre,post SINR and ROI
func sinr(values []float64, intrCancelCount uint) []float64 {
	// Values is recieved power
	var calc = []float64{}
	//PreProcessing SINR calcuation
	//fmt.Println("got values",intrCancelCount)

	//Pre Processing SINR
	var sum float64 = 0.0
	for i := 1; i < len(values); i++ {
		sum += math.Pow(10, (values[i] / 10))
	}
	calc = append(calc, values[0]-10*math.Log10(sum)) //PreProcessing Term

	//Post Processing SINR
	sum = 0.0
	for i := int(intrCancelCount) + 1; i < len(values); i++ {
		sum += math.Pow(10, (values[i] / 10))
	}
	calc = append(calc, values[0]-10*math.Log10(sum)) //Post SINR
	calc = append(calc, 10*math.Log10(sum))           //ROI
	return calc
}
