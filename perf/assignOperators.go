package perf

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/service"
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
//Take each cell's users(approx 200)
//divide into 4 parts of 50 each
//take 12 sectors of 4 operator's stations
//to each group of 50 assign one operator and connect top10 to each sector
///leave 20 in each group as vacant



func AssignSingleOperator(sc *model.Scenario,hexMap *service.HexMap, enFlags []bool) (map[string]interface{}, error) {
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
var total uint
total=0
var ASSIGNED =[]float64{}
var usersPerBS uint
for i:=0;i<19;i++{
	//for each hexagon
   ue:=hexMap.FindContainedUsers(uint(i))
   bs:=hexMap.FindContainedStations(uint(i))
   //ue's and bs in that cell are got
   //bs should be 12 in number
   
   usersPerBS =uint(len(ue)/4)
 //number of users associated to each basestation .approx 50

   for j:=0;j<4;j++{
   	//going by operator ,each sector has at max 10
     	//var assigned = []uint{}
   	     for t:=0;t<3;t++{
   		//for each sector
   		     var losses = []float64{} //to have losses for that basestation
   		
   		     fmt.Println("        BASESTaTION:",bs[3*j+t].ID())
   
   	       for k:=j*int(usersPerBS);k<(j+1)*int(usersPerBS);k++{
                 losses=append(losses,sc.Loss(uint(ue[k].ID()),uint(bs[3*j+t].ID())))
   	                 }


   	   //losses is now got for all 50 users assigned to that operator 

        losses,ind := sort(losses)
        ind1 := make([]uint, len(ind))
        
        for l:=0;l<len(ind);l++{
        	//ind[l]=ue[l+int(j*int(usersPerBS))].ID()
        	ind1[l]=ue[ind[l]].ID()

        	//fmt.Println(ind[l],ind1[l],losses[l])
        }  
         //ind1 has ue id's 

     var count uint  
     count=0
     
     for b:=0;b<len(ind1);b++{
       flag=0
      //checking if already assigned
		      for k:=0;k<len(ASSIGNED);k++{
		          if ASSIGNED[k]==float64(ind1[b]){
		          	   flag=1
		             	break
		              } 
		          } 
		       if flag==0{
             //assigning now
		      // 	fmt.Println(bs[3*j+t].OwnerOp())
		       	    count+=1
		       	    total+=1
		       	    ASSIGNED=append(ASSIGNED,float64(ind1[b]))	
			     	sc.Users()[ind1[b]].CurrOp = bs[3*j+t].OwnerOp()
				   // assigned=append(assigned,ind1[b]) 
				    newOps[ind1[b]] = sc.Users()[ind1[b]].CurrOp.ID()
		       }
		       if count==10{
		       	break
		       }
       }


     
   	 }
   }

   
}//each hexagon

ASSIGNED,e :=sort(ASSIGNED)
e[0]=0.0
fmt.Println(ASSIGNED)
fmt.Println("TOTAL",total)


/*
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
*/

    returnData := map[string]interface{}{}
    returnData["opconn"] = newOps
	return returnData, nil
}