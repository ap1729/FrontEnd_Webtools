package perf

import (
	"FrontEnd_WebTools/model"
	"fmt"
)

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
	fmt.Printf("Interfering station ID's:\n%v\n", intStatIds)

	losses, bsId := signalLossProfile(userID, sc, level, intStatIds)

	op := make([]uint, len(bsId))
	for i := 0; i < len(bsId); i++ {
		losses[i] += 46
		op[i] = sc.GetStationByID(bsId[i]).OwnerOp().ID()
	}

	// Calculate SINR and ROI
	returnData["operno"] = op[0:profileTopN]
	returnData["SINR"] = sinr(losses, intrCancelCount)
	returnData["BSid"] = bsId[0:profileTopN]
	returnData["SIR"] = losses[0:profileTopN]

	return returnData
}
