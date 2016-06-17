package perf

import (
	"FrontEnd_WebTools/model"
	"fmt"
)

func ChangeLevel(sc *model.Scenario, targetLvl uint) map[string]interface{} {
	returnData := map[string]interface{}{}

	fmt.Println("Level 1 reached")

	newOper := make([]uint, len(sc.Users()))
	id := 0
	max := -1.0
	for i := 0; i < len(sc.Users()); i++ { //for each ue
		id = 0
		max = sc.Loss(uint(i), 0)
		for j := 0; j < len(sc.BaseStations()); j++ { //for all bs
			if max < sc.Loss(uint(i), uint(j)) {
				max = sc.Loss(uint(i), uint(j))
				id = j
			}
		}
		newOper[i] = sc.GetStationByID(uint(id)).OwnerOp().ID()
	}

	returnData["changeColor"] = newOper
	return returnData
}
