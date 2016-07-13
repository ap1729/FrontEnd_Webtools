package perf

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/service"
)

// Retreives the signal loss profile for a user and its interfering stations
// for a given scenario and level of cooperation.
//
// The function arranges the signal values from the strongest to the weakest,
// additionally incorporating any order constraints specified by the level of
// cooperation. The second return value specifies the BaseStation indices
// corresponding to the source of the loss value.
func lossProfile(sc *model.Scenario, hexMap *service.HexMap, userID uint, intrStationIds []uint, p *Params) ([]float64, []uint) {
	losses := filter(sc.LossProfile(userID), intrStationIds)
	losses, ind := sort(losses)
	// The sort indices were created for the loss array, and are sequential from
	// 0 to len(losses). Hence, map the index to the actual BaseStation ID:
	for i := 0; i < len(ind); i++ {
		ind[i] = intrStationIds[ind[i]]
	}

	if p.Level == 0 {
		actualOper := sc.GetUserByID(userID).CurrOp.ID()
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

	if p.Level == 1 {
		return losses, ind
	}

	return nil, nil
}
