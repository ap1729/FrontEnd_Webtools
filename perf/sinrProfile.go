package perf

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/service"
	"errors"
	"fmt"
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
func SinrProfile(sc *model.Scenario, hexMap *service.HexMap, userID uint, profileTopN uint, p *Params) (map[string]interface{}, error) {

	// Handling argument nil exception
	if sc == nil || hexMap == nil || p == nil {
		return nil, errors.New(ARG_NIL)
	}

	returnData := map[string]interface{}{}

	// What stations interfere with the current user, given the system parameters
	intStatIds, err := intrStations(sc, hexMap, userID, p)
	if err != nil {
		return nil, fmt.Errorf("Interfering stations could not be determined:\n%v", err.Error())
	}
	// The loss profile and corresponding BaseStation source ID's
	losses, bsId, err := lossProfile(sc, hexMap, userID, intStatIds, p)
	if err != nil {
		return nil, fmt.Errorf("Loss profile could not be evaluated:\n%v", err.Error())
	}
	// Calculate SINR and ROI
	sinrVals, err := sinr(losses, p.IntCancellers)
	if err != nil {
		return nil, fmt.Errorf("Could not calculate SINR:\n%v", err.Error())
	}

	// Assigning key-value pairs to return object
	returnData["pre"] = sinrVals[0]
	returnData["post"] = sinrVals[1]
	returnData["roi"] = sinrVals[2]
	// To avoid slicing the array beyond its bounds, which may occur unobviously.
	if profileTopN > uint(len(losses)) {
		profileTopN = uint(len(losses))
	}
	returnData["bsid"] = bsId[0:profileTopN]
	returnData["sir"] = losses[0:profileTopN]

	return returnData, nil
}
