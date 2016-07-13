package perf

import "FrontEnd_WebTools/model"

func SinrHeatMap(sc *model.Scenario, frMode string, level uint, intrCancelCount uint, opEnable []bool, params map[string]interface{}) map[string]interface{} {
	preSinrVals := make([]float64, len(sc.Users()))
	postSinrVals := make([]float64, len(sc.Users()))
	for i := 0; i < len(sc.Users()); i++ {
		temp := SinrProfile(sc, frMode, sc.Users()[i].ID(), level, intrCancelCount, 0, opEnable, params)
		preSinrVals[i] = temp["pre"].(float64)
		postSinrVals[i] = temp["post"].(float64)
	}
	return map[string]interface{}{"pre": preSinrVals, "post": postSinrVals}
}
