package perf

import (
	"FrontEnd_WebTools/model"
	"fmt"
)

// Retreives the signal loss profile for a user and its interfering stations
// for a given scenario and level of cooperation.
//
// The function arranges the signal values from the strongest to the weakest,
// additionally incorporating any order constraints specified by the level of
// cooperation. The second return value specifies the BaseStation indices
// corresponding to the source of the loss value.
func signalLossProfile(userID uint, sc *model.Scenario, level uint, intrStationIds []uint) ([]float64, []uint) {
	losses := filter(sc.LossProfile(userID), intrStationIds)

	fmt.Printf("Loss values:\n%v\n", losses)

	losses, ind := sort(losses)
	for i := 0; i < len(ind); i++ {
		ind[i] = intrStationIds[ind[i]]
	}

	fmt.Printf("After sort: Interfering station ID's:\n%v\n", ind)
	fmt.Printf("After sort :Loss values:\n%v\n", losses)

	if level == 0 {
		actualOper := sc.GetUserByID(userID).DefaultOp().ID()
		for i := 0; i < len(ind); i++ {
			if sc.GetStationByID(ind[i]).OwnerOp().ID() == actualOper {
				tempInd := ind[i]
				tempLoss := losses[i]
				for k := i; k > 0; k-- {
					ind[k] = ind[k-1]
					losses[k] = losses[k-1]
				}
				ind[0] = tempInd
				losses[0] = tempLoss
				break
			}
		}
		return losses, ind
	}

	if level == 1 {
		return losses, ind
	}

	return nil, nil
}
