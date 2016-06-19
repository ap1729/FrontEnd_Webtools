package perf

import "FrontEnd_WebTools/model"

// TODO: Truncating the profile values as per end-user request is wrongly handled
// here. The job of this function is to solely compute results independent of
// front-end display options.

// Computes the SINR values and the received power profile for a user.
//
// Parameter description:
//
// frMode - frequency-reuse mode
// intrCancelCount - the number of interferers to cancel
// profileTopN - for how many top stations the power profile must be returned
func SinrProfile(sc *model.Scenario, frMode string, userID uint, level uint, intrCancelCount uint, profileTopN uint, params map[string]interface{}) map[string]interface{} {
	returnData := map[string]interface{}{}

	intStatIds := intrStations(frMode, sc, userID, params)
	losses, sortInd := signalLossProfile(userID, sc, level, intStatIds)

	op := make([]uint, len(sortInd))
	bsId := make([]uint, len(sortInd))
	for i := 0; i < len(sortInd); i++ {
		losses[i] += 46
		bsId[i] = intStatIds[sortInd[i]]
		op[i] = sc.GetStationByID(bsId[i]).OwnerOp().ID()
	}

	// Calculate SINR and ROI
	returnData["operno"] = op[0:profileTopN]
	returnData["SINR"] = sinr(losses, intrCancelCount)
	returnData["BSid"] = bsId[0:profileTopN]
	returnData["SIR"] = losses[0:profileTopN]

	return returnData
}
