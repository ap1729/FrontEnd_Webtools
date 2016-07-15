package perf

import (
	"FrontEnd_WebTools/model"
	"errors"
	"math/rand"
)

// Change the registered operator of each user, based on the enabled operators as specified
// by the flags.
//
// The function always assigns to a user its original operator (as per data) if that operator
// is enabled. If not, it randomly assigns it to one of the enabled operators.
func AssignOperators(sc *model.Scenario, enFlags []bool) (map[string]interface{}, error) {

	// Handling argument nil exception
	if sc == nil || enFlags == nil {
		return nil, errors.New(ARG_NIL)
	}

	valOps := []uint{}
	for i := 0; i < len(enFlags); i++ {
		if enFlags[i] == true {
			valOps = append(valOps, uint(i))
		}
	}
	valN := len(valOps)

	// Handling all-disabled exception
	if valN == 0 {
		return nil, errors.New("No operators were enabled in the flags.")
	}

	rand.Seed(19)
	newOps := make([]uint, len(sc.Users()))
	for i := 0; i < len(sc.Users()); i++ {
		if valN == 1 {
			sc.Users()[i].CurrOp = sc.GetOperatorByID(valOps[0])
		} else {
			sc.Users()[i].CurrOp = sc.Users()[i].DefaultOp()
			if enFlags[sc.Users()[i].CurrOp.ID()] == false {
				sc.Users()[i].CurrOp = sc.GetOperatorByID(uint(valOps[rand.Intn(valN)]))
			}
		}
		newOps[i] = sc.Users()[i].CurrOp.ID()
	}

	returnData := map[string]interface{}{}
	returnData["opconn"] = newOps
	return returnData, nil
}
