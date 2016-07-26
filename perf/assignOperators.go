package perf

import (
	"FrontEnd_WebTools/model"
	"errors"
	"math/rand"
	"fmt"
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


//below function for new assigning of single operator case
//Each sector of each basestation connects to 10 strongest



func AssignSingleOperator(sc *model.Scenario, enFlags []bool) (map[string]interface{}, error) {
 //Single Operator Assigning is different

 // Handling argument nil exception
	if sc == nil || enFlags == nil {
		return nil, errors.New(ARG_NIL)
	}

    var valOps uint
	for i := 0; i < len(enFlags); i++ {
		if enFlags[i] == true {
			valOps=uint(i)
			break
		}
	}

 fmt.Println("VAL OPS ",valOps)

 newOps := make([]uint, len(sc.Users()))

for i := 0; i < len(sc.Users()); i++ {
 newOps[i]=4;
 sc.Users()[i].CurrOp = model.NewOperator(uint(10)) //default operator is 10
}	

var flag uint
var assigned = []uint{}
for i := 0; i < len(sc.BaseStations()); i++ {
	var losses = []float64{} //to have losses for that basestation
    for j:=0;j<len(sc.Users());j++{
      losses=append(losses,sc.Loss(uint(j),uint(i)))
    }	
    losses,ind := sort(losses)
   fmt.Println(ind[0:10])
   //Assigning top 10 connected to each basestation
   var count uint
   count=0
   for j:=0;j<len(sc.Users());j++{
      flag=0
      //checking if already assigned
      for k:=0;k<len(assigned);k++{
          if assigned[k]==ind[j]{
          	   flag=1
             	break
              } 
          } 
          if flag==0{
                    count+=1
			     	sc.Users()[ind[j]].CurrOp = sc.GetStationByID(uint(i)).OwnerOp()
				    assigned=append(assigned,ind[j])
				    newOps[ind[j]] = sc.Users()[ind[j]].CurrOp.ID()
				    //to check if 10 have been assigned				   
            }
        if count==10{
        	break
        }
    }
//the ue's not assigned to any bs have to be assigned a null operator
}


    returnData := map[string]interface{}{}
    returnData["opconn"] = newOps
	return returnData, nil
}