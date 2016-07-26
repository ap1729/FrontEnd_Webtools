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
    
    //If single operator case,then interfering stations are only 57
    //interfering stations are only those of CurrOp of user
    if NoUsers(p.OpEnableFlags)==1{
    	fmt.Println("SINGLE OPERATOR CASE REACHED\n")
    	for i:=0;i<4;i++{
    		p.OpEnableFlags[i]=false
    	}
    	if sc.Users()[userID].CurrOp.ID() != 10{
    	//10 is default operator	
    	p.OpEnableFlags[sc.Users()[userID].CurrOp.ID()]=true
       }else{
       	fmt.Println("This User is Currently not assigned to any operator")
       }

    }
   
   if sc.Users()[userID].CurrOp.ID() != 10{
   //if the user is assigned to an operator

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
	rxPows := make([]float64, len(losses))
	for i := 0; i < len(losses); i++ {
		// Add transmit power to each loss value
		rxPows[i] = losses[i] + 46
	}
	sinrVals, err := sinr(rxPows, p.IntCancellers)
	if err != nil {
		return nil, fmt.Errorf("Could not calculate SINR:\n%v", err.Error())
	}

	// Assigning key-value pairs to return object
	returnData["pre"] = sinrVals[0]
	returnData["post"] = sinrVals[1]
	returnData["roi"] = sinrVals[2]
	// To avoid slicing the array beyond its bounds, which may occur unobviously.
	if profileTopN > uint(len(rxPows)) {
		profileTopN = uint(len(rxPows))
	}
	returnData["bsid"] = bsId[0:profileTopN]
	returnData["sir"] = rxPows[0:profileTopN]

	return returnData, nil
    }else{
    	returnData["pre"] = -1000.00
	    returnData["post"]= -1000.00
	    returnData["roi"] = -1000.00
	    returnData["bsid"]=nil
	    returnData["sir"]=nil
	return returnData,nil
    }
}
