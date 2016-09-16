package perf

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/service"
	"errors"
	"math"
	"math/cmplx"
	"math/rand"
	"fmt"
)

func EmDownlink(sc *model.Scenario, hexMap *service.HexMap, opEnable []bool) (map[string]interface{}, error) {

	if sc == nil || hexMap == nil || opEnable == nil {
		return nil, errors.New(ARG_NIL)
	}

	nBS := len(sc.BaseStations())
	nUE := len(sc.Users())

	// Find out how many operators are active
	nOp := 0
	for i := 0; i < len(opEnable); i++ {
		if opEnable[i] {
			nOp++
		}
	}

	// Frequency allocations to each cell as per frequency planning
	var fSet [19]uint
	// Anonymous function to map sector and frequency set to actual frequency of eNodeB
	var idToFreq func(uint, uint, uint) uint
	// Assign the frequencies based on single or multi operator
	if nOp == 1 {
		fSet = [19]uint{2, 1, 2, 4, 3, 4, 3, 2, 1, 2, 1, 2, 3, 4, 3, 4, 2, 1, 2}
		idToFreq = func(fs, opId, sect uint) uint { return 3*(fs-1) + sect }
	} else if nOp == 4 {
		fSet = [19]uint{2, 1, 3, 1, 3, 2, 1, 3, 2, 1, 3, 2, 1, 3, 2, 1, 2, 1, 3}
		idToFreq = func(fs, opId, sect uint) uint { return 4*(fs-1) + opId }
	} else {
		return nil, errors.New("Invalid enable flags; Only single or four operators are supported.")
	}

	rxPows := make([]float64, nUE)
	for i := 0; i < nUE; i++ {
		freqTaps := make([]complex128, 12)
		for j := 0; j < nBS; j++ {
			if !opEnable[sc.GetStationByID(uint(j)).OwnerOp().ID()] {
				continue
			}
			ph := 2 * math.Pi * rand.Float64()
			Eb := math.Sqrt(math.Pow(10, (sc.Loss(uint(i), uint(j))+46)/10-3))
			rWave := complex(Eb*math.Cos(ph), Eb*math.Sin(ph))

			// fmt.Printf("Loss: %v, Phase: %v\nReceived fields: %v\n", sc.Loss(uint(i), uint(j)), ph, rWave)

			bs := sc.GetStationByID(uint(j))
			rootHex := hexMap.FindContainingHex(bs.X(), bs.Y())
			// Significant optimization needed, we can directly loop through hexmap
			// instead of searching for containing hex for each BS.

			// fmt.Printf("Detected frequency: %v, for fset: %v, opID: %v and sector: %v", idToFreq(fSet[rootHex.ID],
			//	bs.OwnerOp().ID(), bs.ID()%3), fSet[rootHex.ID], bs.OwnerOp().ID(), bs.ID()%3)

			freqTaps[idToFreq(fSet[rootHex.ID], bs.OwnerOp().ID(), bs.ID()%3)] += rWave
		}
		rxPows[i] = 0
		for j := 0; j < 12; j++ {
			rxPows[i] += math.Pow(cmplx.Abs(freqTaps[j]), 2)
		}
		rxPows[i] = 10 * math.Log10(rxPows[i]*1000)
	}

	return map[string]interface{}{"rxpow": rxPows}, nil
}


func EmDownlink1(sc *model.Scenario, hexMap *service.HexMap,p *Params, optype string) (map[string]interface{}, error) {
  //also think of randomly destroying basestations
	//also have a height profile for some users
  if sc == nil || hexMap == nil  {
		return nil, errors.New(ARG_NIL)
	}


//	nBS := len(sc.BaseStations())
	nUE := len(sc.Users())
    rxPows := make([]float64, nUE)
    enableFlags := []bool{true,true,true,true}

	if optype=="single" {
		///single operator case, needs to take only 57 loss values for each ue
		//take top 3 as numerator
		//then multiply all else  with random phase and sum with noise to get denominator
     for i := 0; i < nUE; i++ {
     	if sc.GetUserByID(uint(i)).CurrOp.ID()!=10{
		enableFlags[0]=false
		enableFlags[1]=false
    	enableFlags[2]=false
    	enableFlags[3]=false

	    enableFlags[sc.GetUserByID(uint(i)).CurrOp.ID()]=true
        intStatIds, err := intrStations(sc, hexMap, uint(i), p,enableFlags,optype)
        if err!=nil{
       return nil, fmt.Errorf("Interfering stations could not be determined:\n%v", err.Error())
       }
      //gets interfering station id's , now do everything for loss
		       if(len(intStatIds)!=57){
		      fmt.Println("\n",len(intStatIds))}
	

		for j:=0;j<len(intStatIds);j++{
           if sc.GetStationByID(uint(j)).Destroyed==1{
          // 	fmt.Println("Basestaion ",j," Destroyed\n")
           	 continue
           }//below do the calculation
          
          /* 
           losses, bsId, err := lossProfile(sc, hexMap,uint(i), intStatIds, p)
	       if err != nil {
		      return nil, fmt.Errorf("Loss profile could not be evaluated:\n%v", err.Error())
	       }*/

           /*
           	ph := 2 * math.Pi * rand.Float64()  //random phase
			Eb := math.Sqrt(math.Pow(10, (sc.Loss(uint(i), uint(j))+46)/10-3))
			rWave := complex(Eb*math.Cos(ph), Eb*math.Sin(ph))
			*/ // to add random phase to signal
		 }


	   }//check if operator exist
     }//loop over all users
	} else if optype=="multi" {
		///single operator case, needs to take all 228 loss values for each ue
		//take top 12 as numerator
		//then multiply all else  with random phase and sum with noise to get denominator
		enableFlags[0]=true
		enableFlags[1]=true
    	enableFlags[2]=true
    	enableFlags[3]=true
		    for i := 0; i < nUE; i++ {

		        intStatIds, err := intrStations(sc, hexMap,uint(i), p,enableFlags,optype)
		          //gets interfering station id's , now do everything for loss
		           
		        if err !=nil {
		          return nil, fmt.Errorf("Interfering stations could not be determined:\n%v", err.Error())
		          }
                   if(len(intStatIds)!=228){
		          fmt.Println("%d\n",len(intStatIds))}

		          for j:=0;j<len(intStatIds);j++{
           if sc.GetStationByID(uint(j)).Destroyed==1{
           	fmt.Println("Basestaion ",j," Destroyed\n")
           }//in else case do the calculation
		   }


	  }
	} else {
		return nil, errors.New("Invalid enable flags; Only single or four operators are supported.")
	}






	return map[string]interface{}{"rxpow": rxPows}, nil
//  return map[string]interface{}{"rxpow": rxPows}, nil //returns array of powers
}


