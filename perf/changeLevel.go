package perf

import (
	"FrontEnd_WebTools/model"
	"fmt"
	"math"
)

// Returns the operators that users connect to, if the cooperation level was set
// to targetLvl. This is a cosmetic change, and no internal data is modified.
func ChangeLevel(sc *model.Scenario, targetLvl uint, opEnable []bool) map[string]interface{} {
	returnData := map[string]interface{}{}
	newOper := make([]uint, len(sc.Users()))

	if targetLvl == 0 {
		for i := 0; i < len(newOper); i++ {
			newOper[i] = sc.Users()[i].CurrOp.ID()
		}

	} else if targetLvl == 1 {
		for i := 0; i < len(newOper); i++ {
			id := -1
			max := math.Inf(-1)
			for j := 0; j < len(sc.BaseStations()); j++ { //for all bs
				if opEnable[sc.BaseStations()[j].OwnerOp().ID()] == true && max < sc.Loss(uint(i), uint(j)) {
					max = sc.Loss(uint(i), uint(j))
					id = j
				}
			}
			newOper[i] = sc.GetStationByID(uint(id)).OwnerOp().ID()
		}

	} else {
		newOper = nil
	}

	fmt.Println("Level 1 reached")

	returnData["opconn"] = newOper
	return returnData
}
