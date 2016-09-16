package perf

import (
	"FrontEnd_WebTools/model"
	"math/rand"
	"fmt"
)

func DestroyBS(sc *model.Scenario,mode float64)  (map[string]interface{},error) {
	//function to destroy basestations based on a probability value given 
	//currently destroys BS0
	var prob float64
	var total uint
	returnData := map[string]interface{}{}

	total=0
	active:= make([]uint, len(sc.BaseStations()))
	prob=0.4*mode //probability that it will be destroyed
	for i:=0;i<len(sc.BaseStations());i++{
		if rand.Float64()<prob{
			sc.GetStationByID(uint(i)).Destroyed=1;
			active[i]=1;
			total=total+1
		}else{
			active[i]=0;
			sc.GetStationByID(uint(i)).Destroyed=0;//in case it was destroyed previously
		}
	}

	returnData["active"]=active
	fmt.Println("Destroyed ",total,"\n")
	return returnData,nil
 }