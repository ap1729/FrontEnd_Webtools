package perf

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/service"
	"errors"
	"fmt"
)

// Retreives the signal loss profile for a user and its interfering stations
// for a given scenario and level of cooperation.
//
// The function arranges the signal values from the strongest to the weakest,
// additionally incorporating any order constraints specified by the level of
// cooperation. The second return value specifies the BaseStation indices
// corresponding to the source of the loss value.
func lossProfile(sc *model.Scenario, hexMap *service.HexMap, userID uint, intrStationIds []uint, p *Params) ([]float64, []uint, error) {

	// Handling argument nil exception
	if sc == nil || hexMap == nil || p == nil || intrStationIds == nil {
		return nil, nil, errors.New(ARG_NIL)
	}

	losses := filter(sc.LossProfile(userID), intrStationIds)
	losses, ind := sort(losses)
	// The sort indices were created for the loss array, and are sequential from
	// 0 to len(losses). Hence, map the index to the actual BaseStation ID:
	for i := 0; i < len(ind); i++ {
		ind[i] = intrStationIds[ind[i]]
	}



//changes start
if p.Level == 0 {

     	actualBS := sc.GetUserByID(userID).BS0().ID()
     			for i :=0;i<len(ind);i++{
			
			if sc.GetStationByID(ind[i]).ID()== actualBS{
				fmt.Println("AAAAA",actualBS)
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
		return losses, ind, nil
	}

	if p.Level == 1 {
	
	actualBS := sc.GetUserByID(userID).BS1()
	var actualBSid uint
	if actualBS !=nil{
		actualBSid=actualBS.ID()
	}else{
		actualBSid=sc.GetUserByID(userID).BS0().ID()
	}
     			for i:=0;i<len(ind);i++{
		//	fmt.Println("AAAAA",actualBS)
			if sc.GetStationByID(ind[i]).ID()== actualBSid{
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
		return losses, ind, nil

	}

//changes end

/*
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
		return losses, ind, nil
	}

	if p.Level == 1 {
		return losses, ind, nil
	}
*/
	// As of now, we only level 0 and 1 implementations exist.
	return nil, nil, fmt.Errorf(LVL_INV_FMT, p.Level)
}
