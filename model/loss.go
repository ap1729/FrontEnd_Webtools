package model

import . "math"

func HataLoss(x1, y1, x2, y2 float64) float64 {
	var dist float64     // Link Distance (km)
	var hB float64 = 20  // Height of BS (m)
	var hM float64 = 5   // Height of MS (m)
	var fr float64 = 500 // Frequency (MHz)
	var cH float64       // Antenna correction factor
	const citySize string = "large"

	dist = euclideanDist(x1, y1, x2, y2) / 1000
	if citySize == "small" || citySize == "medium" {
		cH = 0.8 + (1.1*Log10(fr)-0.7)*hM - 1.56*Log10(fr)
	} else if citySize == "large" {
		if fr >= 150 && fr <= 200 {
			cH = 8.29*Pow(Log10(1.54*hM), 2) - 1.1
		} else if fr > 200 && fr <= 1500 {
			cH = 3.2*Pow(Log10(11.75*hM), 2) - 4.97
		} else {
			cH = NaN()
		}
	} else {
		cH = NaN()
	}

	return 69.55 + 26.16*Log10(fr) - 13.82*Log10(hB) - cH +
		(44.9-6.55*Log10(hB))*Log10(dist)

}

func euclideanDist(x1, y1, x2, y2 float64) float64 {
	return Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
}
