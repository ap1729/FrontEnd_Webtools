package perf

import (
	"FrontEnd_WebTools/model"
	"fmt"
)

func Level1(sc *model.Scenario) map[string]interface{} {
	returnData := map[string]interface{}{}

	fmt.Println("Level 1 reached")

	var newOper = []uint{}
	id := 0
	max := -1.0
	for i := 0; i < len(sc.Users()); i++ { //for each ue
		id = 0
		max = sc.Loss(uint(i), 0)
		for j := 0; j < len(sc.BaseStations()); j++ { //for all bs
			for k := 1; k < len(sc.BaseStations()); k++ {
				if max < sc.Loss(uint(i), uint(k)) {
					max = sc.Loss(uint(i), uint(k))
					id = k
				}
			}
		}
		newOper = append(newOper, sc.GetStationByID(uint(id)).OwnerOp().ID())
	}

	returnData["changeColor"] = newOper
	return returnData
}
