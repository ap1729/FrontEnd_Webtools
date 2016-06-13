package perf

import (
  "FrontEnd_WebTools/model" 
  "fmt"
  )

func Level1(sc *model.Scenario) map[string]interface{} {
	returnData := map[string]interface{}{}

	fmt.Println("Level 1 reached")

    var newoper =[]uint{}
    var bs model.BaseStation
   id:=0
   max:=-1.0
	for i:=0;i<len(sc.LossTable);i++{ //for each ue
		 id=0
		  max=sc.LossTable[i][0]
		for j:=0;j<len(sc.LossTable[0]);j++{ //for all bs
		    
		    for k:=1;k<len(sc.LossTable[0]);k++{
              if max<sc.LossTable[i][k]{
              	max=sc.LossTable[i][k]
              	id=k
                 }
		    }
        }
		  bs=sc.GetStationByID(uint(id))
        newoper=append(newoper,bs.OwnerOp().ID())

       
	}


	returnData["changeColor"]=newoper
	return returnData
}