package perf

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/service"
	"math"
	"math/cmplx"
	"math/rand"
)

func EmDownlink(sc model.Scenario, hexMap service.HexMap, opEnable []bool) map[string]interface{} {
	nBS := len(sc.BaseStations())
	nUE := len(sc.Users())

	// Find out how many operators are active
	nOp := 0
	for i := 0; i < len(opEnable); i++ {
		if opEnable[i] {
			nOp++
		}
	}

	// Frequency allocations to each cell as per frequency planning
	var fSet [19]uint
	// Anonymous function to map sector and frequency set to actual frequency of eNodeB
	var idToFreq func(uint, uint) uint
	// Assign the frequencies based on single or multi operator
	if nOp == 1 {
		fSet = [19]uint{}
		idToFreq = func(fs, opId uint) uint { return 0 }
	} else {
		fSet = [19]uint{2, 1, 3, 1, 3, 2, 1, 3, 2, 1, 3, 2, 1, 3, 2, 1, 2, 1, 3}
		idToFreq = func(fs, opId uint) uint { return 4*(fs-1) + opId }
	}

	rxPows := make([]float64, nUE)
	for i := 0; i < nUE; i++ {
		freqTaps := make([]complex128, 12)
		for j := 0; j < nBS; j++ {
			ph := 2 * math.Pi * rand.Float64()
			Eb := math.Sqrt(sc.Loss(uint(i), uint(j)))
			rWave := complex(Eb*math.Cos(ph), Eb*math.Sin(ph))
			bs := sc.GetStationByID(uint(j))
			rootHex := hexMap.FindContainingHex(bs.X(), bs.Y())
			// Significant optimization needed, we can directly loop through hexmap
			// instead of searching for containing hex for each BS.
			freqTaps[idToFreq(fSet[rootHex.ID], bs.OwnerOp().ID())] += rWave
		}
		rxPows[i] = 0
		for j := 0; j < 12; j++ {
			rxPows[i] += math.Pow(cmplx.Abs(freqTaps[j]), 2)
		}
	}

	return map[string]interface{}{"rxpow": rxPows}
}
