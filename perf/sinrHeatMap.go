package perf

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/service"
	"math"
	//"fmt"
)

func SinrHeatMap(sc *model.Scenario, hexMap *service.HexMap, p *Params) map[string]interface{} {
	preSinrVals := make([]float64, len(sc.Users()))
	postSinrVals := make([]float64, len(sc.Users()))
	var preSumRate float64 = 0.0
	var postSumRate float64 = 0.0
	var CenterPostRate float64=0.0
	var CenterPreRate float64=0.0
	for i := 0; i < len(sc.Users()); i++ {
		if sc.Users()[i].CurrOp.ID() == 10{
            //if currop not assigned
            preSinrVals[i] = -1000
			postSinrVals[i] = -1000
		}else{
		vals, err := SinrProfile(sc, hexMap, sc.Users()[i].ID(), 0, p)
		

		if err != nil {
			// TODO: As of now, heat map ignores error and assigns them a fallback value.
			preSinrVals[i] = -1000
			postSinrVals[i] = -1000
		} else{	
			preSinrVals[i] = vals["pre"].(float64)
			preSumRate += math.Log10(1 + math.Pow(10, preSinrVals[i]/10))
			postSinrVals[i] = vals["post"].(float64)
			postSumRate += math.Log10(1 + math.Pow(10, postSinrVals[i]/10))
		}
       }
	}//for loop over

//to find rate only for center cell
	centerUsers := hexMap.FindContainedUsers(9)
	for i:=0;i<len(centerUsers);i++ {
		//loops over all users in center cell
		if centerUsers[i].CurrOp.ID() != 10{
          vals, err := SinrProfile(sc, hexMap, centerUsers[i].ID(), 0, p)
         if err == nil {
			CenterPreRate += math.Log10(1 + math.Pow(10, vals["pre"].(float64)/10))
			CenterPostRate += math.Log10(1 + math.Pow(10, vals["post"].(float64)/10))
		     }


		}	
	}

	return map[string]interface{}{"pre": preSinrVals, "post": postSinrVals, "preSumRate": preSumRate, "postSumRate": postSumRate,"centerPostRate":CenterPostRate,"centerPreRate":CenterPreRate}
}
