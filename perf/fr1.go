package perf

import (
	"FrontEnd_WebTools/model"
	"fmt"
)

func FR1(sc *model.Scenario, userID uint, level uint, intrCancelCount uint, profileTopN uint) map[string]interface{} {
	returnData := map[string]interface{}{}
	///sc.LossTable[userID] gives row of pathloss
	fmt.Println("FR1 () got data")
	row := sc.LossProfile(userID)
	values := make([]float64, len(row))
	keys := make([]int, len(row))
	for i := 0; i < len(row); i++ {
		keys[i] = i
		values[i] = row[i] + 46.0
	}

	

	temp1 := 0
	temp2 := 0.0
	for i := 0; i < len(keys); i++ {
		for j := 1; j < len(keys); j++ {
		
			if values[j-1] < values[j] {
				temp1 = keys[j-1]
				keys[j-1] = keys[j]
				keys[j] = temp1

				temp2 = values[j-1]
				values[j-1] = values[j]
				values[j] = temp2
			}
		}
	}
	//sorting of values done

	if level == 0 {
		//if level is zero do cyclic shift
		actualOper := sc.GetUserByID(userID).DefaultOp().ID()

		for i := 0; i < len(keys); i++ {
			if actualOper == sc.GetStationByID(uint(keys[i])).OwnerOp().ID() { //operator of bs
				//cyclic shift
				temp1 := keys[i]
				temp2 := values[i]
				for k := i; k > 0; k-- {
					keys[k] = keys[k-1]
					values[k] = values[k-1]
				}
				keys[0] = temp1
				values[0] = temp2
				break
			}

			//do cyclic shift

		}
	} //level0 cyclic shift over

	//operator of bs
	var opid = []uint{}
	for i := 0; i < int(profileTopN); i++ {
		opid = append(opid, sc.GetStationByID(uint(keys[i])).OwnerOp().ID())
	}

	//now to calculate SINR and ROI
	returnData["operno"] = opid
	returnData["SINR"] = SINR_ROI(values, intrCancelCount)
	returnData["BSid"] = keys[0:profileTopN]
	returnData["SIR"] = values[0:profileTopN]
	//operator id list for bargraph color should be done here

	return returnData
}
