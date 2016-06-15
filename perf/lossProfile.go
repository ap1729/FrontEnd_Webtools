package perf

import (
	"FrontEnd_WebTools/model"
)

func signalLossProfile(userID uint, sc *model.Scenario, level uint, intrStationIds []uint) ([]float64, []uint) {
	losses := Filter(sc.LossProfile(userID), intrStationIds)
	losses, ind := Sort(losses)

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
