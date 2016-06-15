package perf

import "FrontEnd_WebTools/model"

func intrStations(mode string, sc *model.Scenario, userID uint, params map[string]interface{}) []uint {
	var bsIds []uint

	switch mode {
	case "FR1":
		bsIds = make([]uint, len(sc.BaseStations()))
		for i := 0; i < len(sc.BaseStations()); i++ {
			bsIds[i] = uint(i)
		}
	case "FR3":
	case "FFR":
	case "AFFR":
	default:
		return nil
	}

	return bsIds
}
