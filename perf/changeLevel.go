package perf

import (
	"FrontEnd_WebTools/model"
	"errors"
	"fmt"
	"math"
)

// Returns the operators that users connect to, if the cooperation level was set
// to targetLvl. This is a cosmetic change, and no internal data is modified.
func ChangeLevel(sc *model.Scenario, targetLvl uint, opEnable []bool) (map[string]interface{}, error) {

	if sc == nil || opEnable == nil {
		return nil, errors.New(ARG_NIL)
	}

	returnData := map[string]interface{}{}
	newOper := make([]int, len(sc.Users()))

	if targetLvl == 0 {
		for i := 0; i < len(newOper); i++ {
			newOper[i] = int(sc.Users()[i].CurrOp.ID())
		}

	} else if targetLvl == 1 {
		for i := 0; i < len(newOper); i++ {
			id := -1
			max := math.Inf(-1)
			// Finding the strongest BaseStation seen by the user:
			for j := 0; j < len(sc.BaseStations()); j++ {
				if opEnable[sc.BaseStations()[j].OwnerOp().ID()] == true && max < sc.Loss(uint(i), uint(j)) {
					max = sc.Loss(uint(i), uint(j))
					id = j
				}
			}
			// PONDER: If there is no BaseStation found, assign it an ID of -1, and not panic (?)
			if id < 0 {
				newOper[i] = -1
			} else {
				newOper[i] = int(sc.GetStationByID(uint(id)).OwnerOp().ID())
			}
		}
	} else {
		return nil, fmt.Errorf(LVL_INV_FMT, targetLvl)
	}

	returnData["opconn"] = newOper
	return returnData, nil
}
