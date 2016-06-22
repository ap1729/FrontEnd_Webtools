package perf

import (
	"FrontEnd_WebTools/model"
	"math/rand"
)

func AssignOperators(sc *model.Scenario, enFlags []bool) map[string]interface{} {

	valOps := []uint{}
	for i := 0; i < len(enFlags); i++ {
		if enFlags[i] == true {
			valOps = append(valOps, uint(i))
		}
	}

	valN := len(valOps)
	rand.Seed(19)
	newOps := make([]uint, len(sc.Users()))

	for i := 0; i < len(sc.Users()); i++ {
		sc.Users()[i].CurrOp = sc.Users()[i].DefaultOp()
		if enFlags[sc.Users()[i].CurrOp.ID()] == false {
			sc.Users()[i].CurrOp = sc.GetOperatorByID(uint(rand.Intn(valN-1)) + 1)
			newOps[i] = sc.Users()[i].CurrOp.ID()
		}
	}

	returnData := map[string]interface{}{}
	returnData["opconn"] = newOps
	return returnData
}
