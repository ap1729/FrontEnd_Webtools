package perf

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/service"
)

func SinrHeatMap(sc *model.Scenario, hexMap *service.HexMap, p *Params) map[string]interface{} {
	preSinrVals := make([]float64, len(sc.Users()))
	postSinrVals := make([]float64, len(sc.Users()))
	for i := 0; i < len(sc.Users()); i++ {
		if hexMap.FindContainingHex(sc.Users()[i].X(), sc.Users()[i].Y()) == nil {
			preSinrVals[i] = -1000
			postSinrVals[i] = -1000
			continue
		}
		temp := SinrProfile(sc, hexMap, sc.Users()[i].ID(), 0, p)
		preSinrVals[i] = temp["pre"].(float64)
		postSinrVals[i] = temp["post"].(float64)
	}
	return map[string]interface{}{"pre": preSinrVals, "post": postSinrVals}
}
