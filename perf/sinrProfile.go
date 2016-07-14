package perf

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/service"
)

// TODO: Truncating the profile values as per end-user request is wrongly handled
// here. The job of this function is to solely compute results independent of
// front-end display options.

// Computes the SINR values and the received power profile for a user.
//
// Parameter description:
//
// p - The parameters of the model to be simulated
// profileTopN - for how many top stations the power profile must be returned
func SinrProfile(sc *model.Scenario, hexMap *service.HexMap, userID uint, profileTopN uint, p *Params) map[string]interface{} {
	returnData := map[string]interface{}{}

	// What stations interfere with the current user, given the system parameters
	intStatIds := intrStations(sc, hexMap, userID, p)
	// The loss profile and corresponding BaseStation source ID's
	losses, bsId := lossProfile(sc, hexMap, userID, intStatIds, p)

	// Calculate SINR and ROI
	sinrVals := sinr(losses, p.IntCancellers)
	returnData["pre"] = sinrVals[0]
	returnData["post"] = sinrVals[1]
	returnData["roi"] = sinrVals[2]

	if profileTopN > uint(len(losses)) {
		profileTopN = uint(len(losses))
	}
	returnData["bsid"] = bsId[0:profileTopN]
	returnData["sir"] = losses[0:profileTopN]

	return returnData
}
