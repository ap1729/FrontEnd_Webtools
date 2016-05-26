package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
   //     "math"
        "encoding/csv"
         "fmt"
         "os"
	 "strconv"
  //       "sort"
)

type TestRecord struct {
         col1 int
         node  string
          x float64
         y  float64
 	m1 float64
         m2  float64
 }//records for data from operator
var bs []TestRecord  //to store bs data in a struct
var ue []TestRecord  //to store ue data in a struct

type oper struct{
  operno int
  ueno int 
  bsno int
}//structure for each operator info ,i.e how many ue's and bs belong for each operator
var op []oper


var Pathloss = [][]float64{} //global variable to store path loss as a 2d array
var Pathlossdata = map[int]float64{}// global dictionary which stores SIRdata map





// E.g for sending from browser
// Obj={id:"32",Column1:30,Column2:-45,Column3:-34,Column4:-40}
// $.post("http://localhost:8080/update",JSON.stringify(Obj),"json")
type rowdata struct {
        Type    string
	Node    int
        Level   int
        TopBsno int
	
}//data which is coming from user

type returndata struct{
        SIR []float64
        PrS float64
        PoS float64
        ROI float64
        Bsid []int
}// structure for returning data to front end for SIR

type returndata1 struct{
    Changecolor []int

}// structure for returning data to front end for level1


func handlerroute(w http.ResponseWriter, r *http.Request) {

	// w.Header().Add("Access-Control-Allow-Origin", "*")
	if r.Method != "POST" {

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only POST method is supported"))
		return
	}
	log.Println("Request Method is ", r.Method)
	// log.Println("Request is originated from  ", r.RemoteAddr)
	// log.Println("Request is originated URL  ", r.RequestURI)
	// log.Println("Request Headers", r.Header)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("body is ", string(data))



		var rxdata rowdata
		jerr := json.Unmarshal(data, &rxdata)
		if jerr != nil {
			log.Println("Unmarshalling error ", jerr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
                log.Println("Marshalled Data",rxdata)	
                var updateData returndata
                var updateData1 returndata1
         if string(data[9])=="A"{
//pass SIR function
  //log.Println("DETECTED")
       updateData =SIR(rxdata)
       txbytes, _ := json.Marshal(updateData)
        nbytes, werr := w.Write(txbytes)
		_ = nbytes
		if werr != nil {
			log.Println("I got some error while writing back", werr)
		} else {
			 log.Println("Sent this  ", string(txbytes))
			// log.Printf("Successfully returned %d bytes", nbytes)
		}
 
         } else if string(data[9])=="B"{
		updateData1=level1(rxdata)
        //     log.Println("UPDATED",updateData1)//check if function works
              txbytes1, _ := json.Marshal(updateData1)
               nbytes1, werr1 := w.Write(txbytes1)
		_ = nbytes1
		if werr1 != nil {
			log.Println("I got some error while writing back", werr1)
		} else {
			 log.Println("Sent this  ", string(txbytes1))
			// log.Printf("Successfully returned %d bytes", nbytes)
		}                  

                } 



		
		

	}

}

/*
func EvaluateMore(r rowdata) returndata{
if r.Type=="A"{
  fmt.Println("BBBBBB")
  return SIR(r)
 }else{//make function for level1 stuff
  return SIR(r)
 }
}

*/






func SIR(r rowdata) returndata {
	// log.Printf("I received this to process %#v", r)
	// Actual algo goes here ....
// below is SIR profile
    count :=0
    for _,i := range Pathloss[r.Node]{
     //fmt.Println("AA",i)
      Pathlossdata[count]=i
    count+=1
    }
 //fmt.Println(Pathlossdata)
 keys:= []int{}
values:= []float64{}
  for key,value := range Pathlossdata{
     keys =append(keys,key)
     values =append(values,value)
  } 
// fmt.Println("KeYS",keys) 
// fmt.Println(values)
temp1:=0
temp2:=0.0
  for i:=0;i< len(keys);i++{
   for j:=1;j<len(keys);j++{
      if values[j-1]<values[j]{
         temp1=keys[j-1]
         keys[j-1]=keys[j]
         keys[j]=temp1
  
         temp2=values[j-1] 
         values[j-1]=values[j]
         values[j]=temp2
       }
   }
 //  fmt.Println("\n",i)
 } 
fmt.Println(keys)
fmt.Println(values)


if r.Level==0{
  lb:=0
  ub:=0
  currop:=-1  
  for i:=0;i <len(op);i++ {
   ub+=op[i].ueno
   if r.Node >= lb && r.Node<ub {
        currop=i
        break
      } else{
     lb+=op[i].ueno
     }
 
  }//for loop over
//currop+1 is curr operator no

//for cyclic shift 
for i:=0;i<r.TopBsno;i++{
   //now to find which operator keys[i] belongs to
   lb:=0
   ub:=0
   currop1:=-1
   for j:=0;j<len(op);j++{
     ub+=op[j].bsno
      if keys[i]>=lb && keys[i]<ub{
       currop1=j  
      
    //   fmt.Println("AAAAAAAA",currop1+1,keys[i])  ,to check
        
       if currop==currop1{
       break}
    }else{
        lb+=op[j].bsno
          }
 
    }// inner for loop is over,finding bs operator
if currop1== currop{
   //currop1 is the operator number
  //cyclic shift is done here
    temp1:=keys[i]
    temp2:=values[i]
    for k:=i;k>0;k--{
    keys[k]=keys[k-1]
    values[k]=values[k-1]
    }
   keys[0]=temp1
   values[0]=temp2
    break
   } 
  }//for loop is over



}// for level 0 ,cyclic shift



//values have descending order of path loss
//keys have bs id 's in that order
       
      //row variable has the path loss for that node 
      //sort it in descending order to get SIR for level 1
      //check if level is 0 ,then do a cyclic shift 
      //return array of top 10 by default
     //Also Calculate SINR's and send back

   //  sort.Float64s(row)//function which sorts data in ascending order

//Ramanan Code here

     var returnobj returndata
     returnobj.SIR=values[0:r.TopBsno]
     returnobj.PrS=13.13
     returnobj.PoS=57.57
     returnobj.Bsid=keys[0:r.TopBsno]
     fmt.Println(returnobj)

       return returnobj
}



func level1(r rowdata) returndata1{
//function to return array of operator numbers to change colour
var returnobj1 returndata1

          
max:=0.0
id:=0
//Using 2D array Pathloss[][]
for i:=0;i<len(Pathloss);i++ {
//loop for each row
  max=Pathloss[i][0]
  id=0
  for j:=0;j<76;j++{
//for all elements in one row
   if max<Pathloss[i][j]{
     id=j       
     max=Pathloss[i][j]
   }

  }

//need to make array to return here
//id has bsno ,to find which operator it belongs to
lb:=0
ub:=0
for i:=0;i<len(op);i++{
 ub+=op[i].bsno
 if lb<=id && id<ub{
    returnobj1.Changecolor=append(returnobj1.Changecolor,i)      //adding to array to return
  break
   }

 lb+=op[i].bsno
 }


}
//fmt.Println(id,max,"\n")
//fmt.Println(op) 


return returnobj1
}


















func main() {

//pathloss csv
csvfile, err := os.Open("Converted.csv")
         if err != nil {
                 fmt.Println(err)
                 os.Exit(1)
         }
         defer csvfile.Close()
         reader := csv.NewReader(csvfile)
         reader.FieldsPerRecord = -1
         rawCSVdata, err := reader.ReadAll()
         if err != nil {
                 fmt.Println(err)
                 os.Exit(1)
         }
         // sanity check, display to standard output
         //for _, each := range rawCSVdata {
          //       fmt.Printf("col1 : %s node : %s,\nx:%f\n\n,y:%f,m1:%f,m2:%f \n\n\n", each[0], each[1],each[2],each[3],each[4],each[5])
        //}
         // now, safe to move raw CSV data to struct
        count:=0
         
         for _, each := range rawCSVdata {
             if count!=0{//to not print first row
           //   fmt.Printf("row",each) 
              temp:=[]float64{} 
                for i := 0; i < 76; i++ { //this part is hardcoded ,later will make it indpt
              
               a, err := strconv.ParseFloat(each[i], 64)            
                 if err==nil{
                     temp= append(temp,a)
                           }
                 }//for loop over
             Pathloss =append(Pathloss,temp)
              
                 } 
           count+=1

         }//for loop of csv parse is over





//Nodelocations csv file

csvfile1, err1 := os.Open("Nodelocations.csv")

         if err1 != nil {
                 fmt.Println(err1)
                 os.Exit(1)
         }
         defer csvfile1.Close()
         reader1 := csv.NewReader(csvfile1)
         reader1.FieldsPerRecord = -1
         rawCSVdata1, err1 := reader1.ReadAll()
         if err1 != nil {
                 fmt.Println(err1)
                 os.Exit(1)
         }
    count =0
var one oper
one.operno=0
one.ueno=0
one.bsno=0
op =append(op,one)

           switch1:=0 //to know when bs are read fully
           bscount:=0
           bscurrop:=1
           uecount:=0
           uecurrop:=1
for _, each1 := range rawCSVdata1 {
    
             if count!=0{//to not print first row
//each is of formatof col1,node,x,y,m1,m2
  var oneRecord TestRecord
   b, err := strconv.Atoi(each1[0])
  if err==nil{
   oneRecord.col1=b
  }
oneRecord.node=each1[1]
     
 c, err := strconv.ParseFloat(each1[2], 64)
  if err==nil{
   oneRecord.y=c
  }
 d, err := strconv.ParseFloat(each1[3], 64)
  if err==nil{
   oneRecord.x=d
  }
 e, err := strconv.ParseFloat(each1[4], 64)
  if err==nil{
   oneRecord.m1=e
  }
 f, err := strconv.ParseFloat(each1[5], 64)
  if err==nil{
   oneRecord.m2=f
  }

//fmt.Println(oneRecord)
//oneRecord created

  if string(each1[1][0])=="B"{
    //to add to bs
bs=append(bs,oneRecord)

a := int(each1[1][2]-48)
    if bscurrop == a{
      bscount+=1
    }  else {
       
        op[bscurrop-1].bsno=bscount
        op[bscurrop-1].operno=bscurrop
        bscount=1
        bscurrop+=1 
       var one oper
	one.operno=0
	one.ueno=0
	one.bsno=0
        op =append(op,one)

    } 


}else{ //to add to ue
 if switch1==0{
        op[bscurrop-1].bsno=bscount
        op[bscurrop-1].operno=bscurrop
  }else{
  switch1=1
 }

 ue=append(ue,oneRecord)       
 a :=int(each1[1][2]-48)
//fmt.rintln("\nBB",a)
 if uecurrop == a{
      uecount+=1
    }  else {
       
        op[uecurrop-1].ueno=uecount
        uecount=1
        uecurrop+=1 
   }

}

    count+=1
 } else{
    //fmt.Printf("\nBB",each1)
    count+=1
   }
}//csv parse over
op[uecurrop-1].ueno=uecount
uecount=1
uecurrop+=1 
//op now has for all operators number of bs and ues 
	log.Println("Started Server at :8080")
	http.HandleFunc("/update", handlerroute)
	http.Handle("/", http.FileServer(http.Dir(".")))
	err2 := http.ListenAndServe(":8080", nil)
	if err2 != nil {
		log.Println("Error ", err2)
	}
}
