package perf

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/service"
)

func SinrHeatMap(sc *model.Scenario, hexMap *service.HexMap, p *Params) map[string]interface{} {
	preSinrVals := make([]float64, len(sc.Users()))
	postSinrVals := make([]float64, len(sc.Users()))
	for i := 0; i < len(sc.Users()); i++ {
		vals, err := SinrProfile(sc, hexMap, sc.Users()[i].ID(), 0, p)
		if err != nil {
			// TODO: As of now, heat map ignores error and assigns them a fallback value.
			preSinrVals[i] = -1000
			postSinrVals[i] = -1000
		} else {
			preSinrVals[i] = vals["pre"].(float64)
			postSinrVals[i] = vals["post"].(float64)
		}

	}
	return map[string]interface{}{"pre": preSinrVals, "post": postSinrVals}
}
