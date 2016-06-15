package perf

import "FrontEnd_WebTools/model"

func FR1(sc *model.Scenario, userID uint, level uint, intrCancelCount uint, profileTopN uint) map[string]interface{} {
	returnData := map[string]interface{}{}

	intStatIds := intrStations("FR1", sc, userID, nil)
	losses, bsId := signalLossProfile(userID, sc, level, intStatIds)

	op := make([]uint, len(bsId))
	for i := 0; i < len(bsId); i++ {
		losses[i] += 46
		op[i] = sc.GetStationByID(bsId[i]).OwnerOp().ID()
	}

	// Calculate SINR and ROI
	returnData["operno"] = op
	returnData["SINR"] = SINR_ROI(losses, intrCancelCount)
	returnData["BSid"] = bsId[0:profileTopN]
	returnData["SIR"] = losses[0:profileTopN]

	return returnData
}
