package service

import (
	"FrontEnd_WebTools/model"
	. "math"
	"math/rand"
)

func GenerateMap(sb *model.ScenarioBuilder) bool {
	if sb.IsSealed() {
		return false
	}

	// Read node locations, as we need the BS locations from file
	sbTemp := model.NewScenarioBuilder()
	ReadNodes(sbTemp, "data/SectorLocations.csv")
	dummyLosses := make([][]float64, sbTemp.NumUsersAdded())
	for i := range dummyLosses {
		dummyLosses[i] = make([]float64, sbTemp.NumStationsAdded())
	}
	sbTemp.Seal("import", dummyLosses)
	scTemp := sbTemp.Finalize()

	// Add all BaseStations to the ScenarioBuilder
	for i := 0; i < len(scTemp.BaseStations()); i++ {
		bs := scTemp.BaseStations()[i]
		opId := bs.OwnerOp().ID()
		if !sb.OperatorExists(opId) {
			sb.AddOperator(opId)
		}
		sb.AddNode("BS", bs.X(), bs.Y(), bs.Ht(), opId)
	}

	// Add users at uniform locations throughout the map
	isd := float64(500)
	side := isd * 2 / Sqrt(3)
	step := float64(50)

	// From manual calculations, equation of border line to 3-tier cell map is:
	// x = (|y| - 4*side) / sqrt(3)
	// We can tightly fit the map by varying y as [-4*side, 4*side]

	numOp := len(scTemp.Operators())
	var i, j float64
	for i = -4 * side; i <= 4*side; i += step {
		y := i
		iniX := (Abs(y) - 8*side) / Sqrt(3)
		for j = iniX; j < Abs(iniX); j += step {
			x := j
			sb.AddNode("UE", x, y, 0, scTemp.GetOperatorByID(uint(rand.Intn(numOp))).ID())
		}
	}

	return true
}
