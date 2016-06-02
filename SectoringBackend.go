package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
        "math"
        "encoding/csv"
         "fmt"
         "os"
	 "strconv"
         "sort"
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

type hexagons struct{
   x float64  
   y float64//x,y of hex center
   ue []int //list of ue's
   bs []int // list of bs

}

var hex []hexagons
var adj = [][]int{} //adjacency matrix

func hexcenter(){
//assuming center is at 0,0
//we are assuming a 5x5 grid of hexagons
//hence left top most hexagon has center at (-2000,500 +(2000/1.73))
xbase:=-2000.00
ybase:=(3000/math.Sqrt(3)) //coordinates of top left hexagon

for i:=0;i<5;i++{
  
   for j:=0;j<5;j++{
   var onehex hexagons
   onehex.x=float64(xbase+float64(j*1000-(i%2)*500))
   onehex.y=ybase-float64(i)*(1500/math.Sqrt(3))
   //fmt.Println("\nInside hexagon" ,i,j)
//   fmt.Println("Center is ",onehex.x,onehex.y)
   for k:=0;k<len(bs);k++{
//finding bs inside hexagon
      if insidehex(onehex.x-bs[k].x,onehex.y-bs[k].y)==1{ //giving relative coords
   onehex.bs=append(onehex.bs,k)
         }
   }

for l:=0;l<len(ue);l++{
//finding ue inside hexagon
if insidehex(onehex.x-ue[l].x,onehex.y-ue[l].y)==1{
 onehex.ue=append(onehex.ue,l)
 }
}
  hex = append(hex,onehex) //appends to global variable

    }//j loop over
  }// i loop over

for i:=0;i<len(hex);i++{
 temp:=[]int{}
  for j:=0;j<len(hex);j++{
  if j==i{
     temp=append(temp,0)
    }else{
   if math.Abs(math.Sqrt((hex[i].x-hex[j].x)*(hex[i].x-hex[j].x) + (hex[i].y-hex[j].y)*(hex[i].y-hex[j].y)) -1000) <0.001{ //checks if distance is close
     temp=append(temp,1)
     } else{
      temp=append(temp,0)
      }
     }
    }//j loop over
   adj=append(adj,temp)

   }
//fmt.Println(adj)
}//function over


func insidehex(relx float64,rely float64) int{
//gettig relative coordinates
//assuming hex side len is 1000/math.Sqrt(3)
if math.Abs(float64(relx*(math.Cos(0)) - rely*math.Sin(0)))<500.00  && math.Abs(float64(relx*(math.Cos(math.Pi/3)) - rely*math.Sin(math.Pi/3))) <500.00 && math.Abs(float64(relx*(math.Cos(-math.Pi/3)) - rely*math.Sin(-math.Pi/3))) <500.00 {
 return 1
 }else{
   return 0
  }
}
	
func adjhex(curr int) []int {
//first tier
temp:=[]int{}
for i:=0;i<len(adj);i++{
  if adj[curr][i]==1{
    temp=append(temp,i)
   }
 }
return temp
}

func secondtier(curr int) []int{
temp:=[]int{}
f:=0
first:=adjhex(curr)
for i:=0;i<len(adj);i++{
   if adj[curr][i]==1{
      for j:=0;j<len(adj);j++{
        if adj[i][j]==1{
          f=0
        // j is a neighbour of neighbour ,but can still be in first tier or original
            if j!=curr{
              for k:=0;k<len(first);k++{
                  if j==first[k]{
                   f=1
                   break             
                 }
               }
            if f==0{
             //check if already in array
            for t:=0;t<len(temp);t++{ 
                if temp[t]==j{
                  f=1
                  break
                  }
                 }
               }
      if f==0{temp=append(temp,j)}else{f=0}
              
             }
           }
      }
    }
  }

return temp
}
//global variables all above 



// E.g for sending from browser
// Obj={id:"32",Column1:30,Column2:-45,Column3:-34,Column4:-40}
// $.post("http://localhost:8080/update",JSON.stringify(Obj),"json")
type rowdata struct {
        Type    string
	Node    int
        Level   int
        TopBsno int
	Noise float64
        Topx int
}//data which is coming from user

type returndata struct{
        SIR []float64
        PrS float64
        PoS float64
        ROI float64
        Bsid []int
       Operno []int //for bar graph
}// structure for returning data to front end for SIR

type returndata1 struct{
    Changecolor []int

}// structure for returning data to front end for level1

type returndata2 struct{
    X []float64
    Y []float64
}//structure for returning data to front end for CDF

type returndata3 struct{
    PrS float64
    PoS float64
    ROI float64
}//structure for returning FR3 data






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

//return data types
                var updateData returndata
                var updateData1 returndata1
                var updateData2 returndata2
                var updateData3 returndata3

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

                } else if string(data[9])=="C"{
                 updateData2 = CDF()
             //CDF plot
                    txbytes2, _ := json.Marshal(updateData2)
               nbytes2, werr2 := w.Write(txbytes2)
		_ = nbytes2
		if werr2 != nil {
			log.Println("I got some error while writing back", werr2)
		} else {
			 log.Println("Sent this  ", string(txbytes2))
			// log.Printf("Successfully returned %d bytes", nbytes)
		}                  
              } else if string(data[9])=="D"{
 fmt.Println("FR3")
updateData3=FR3(rxdata)


txbytes3, _ := json.Marshal(updateData3)
               nbytes3, werr3 := w.Write(txbytes3)
		_ = nbytes3
		if werr3 != nil {
			log.Println("I got some error while writing back", werr3)
		} else {
			 log.Println("Sent this  ", string(txbytes3))
			// log.Printf("Successfully returned %d bytes", nbytes)
		}  


             }else{
            fmt.Println("Unknown command")
            }

	}

}

func operatorbybs(bsid int) int{
//which operator by passing bsid
lb:=0
ub:=0
i:=0
for i=0;i<len(op);i++{
 ub+=op[i].bsno //sectoring
 if bsid>=lb && bsid<ub{
   break
  }
 lb+=op[i].bsno//sectoring
 }

return i
}




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
 //fmt.Println("KeYS",len(keys)) 
 //fmt.Println(values)
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
//fmt.Println(keys)
//fmt.Println(values)


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
var received_power_db_arr = []float64{} //array variable to store received power from the base station (hardcoded no. of base station)
var other_received_power_lin_arr = []float64{} //array variable to store received power linear scale from the other base station (hardcoded /no. of base station)
var pre_processing_sinr_db float64
var post_processing_sinr_db float64
var sum_interferers_lin float64 =0
//fmt.Println(values)
for i := 0; i < len(bs); i++ {
	received_power_db_arr =append(received_power_db_arr,values[i]+46.0) //pathloss+hardcoded transmit power (for one UE)
	values[i] = values[i]+46.0
	}
max_received_power_dB := received_power_db_arr[0]
other_received_power_lin_arr=append(other_received_power_lin_arr,0.0)
for i := 1; i < len(bs); i++ {
	other_received_power_lin_arr =append(other_received_power_lin_arr, math.Pow(10,received_power_db_arr[i]/10)) //converting to linear scale	
	sum_interferers_lin += math.Pow(10,received_power_db_arr[i]/10)
	}
//fmt.Println("\n",received_power_db_arr)
//fmt.Println("\n\n\n\n\n\n",other_received_power_lin_arr)
//fmt.Println("\n\n\n\n\n\n",sum_interferers_lin)	
pre_processing_sinr_db = max_received_power_dB-10*math.Log10(sum_interferers_lin + math.Pow(10,r.Noise/10)) 
//fmt.Println("\n\n\n\n\n\n pre-processing SINR",pre_processing_sinr_db)
//post processign SINR calculation
var sum_interferers_cancel_lin float64 =0
var num_interferers_cancel int = r.Topx 
for i := num_interferers_cancel+1; i < len(bs); i++ {
	other_received_power_lin_arr[i] = math.Pow(10,received_power_db_arr[i]/10) //converting to linear scale	
	sum_interferers_cancel_lin +=  other_received_power_lin_arr[i]
	}
post_processing_sinr_db = max_received_power_dB-10*math.Log10(sum_interferers_cancel_lin + math.Pow(10,r.Noise/10)) //noise level -90dBm hardcoded
//fmt.Println("\n\n\n\n\n\n post-processing SINR",post_processing_sinr_db)

     var returnobj returndata


for i:=0;i<r.TopBsno;i++{
returnobj.Operno=append(returnobj.Operno,operatorbybs(keys[i]))
}
     returnobj.SIR=values[0:r.TopBsno]
     returnobj.PrS=pre_processing_sinr_db
     returnobj.PoS=post_processing_sinr_db
     returnobj.Bsid=keys[0:r.TopBsno]
     returnobj.ROI =10*math.Log10(sum_interferers_cancel_lin)//value here 
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
  for j:=0;j<len(bs);j++{ 
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
 fmt.Println("\n",i," ",operatorbybs(id))
  break
   }

 lb+=op[i].bsno
 }
fmt.Println(i,"  ",id, "   ",i	)

}
//fmt.Println(id,max,"\n")
//fmt.Println(op) 


return returnobj1
}


func CDF() returndata2{

var returnobj2 returndata2
//Get pre and post SINR
var temp rowdata
var getdata returndata

        temp.Type="A"
	temp.Node=0    //which ue  
        temp.Level=0
        temp.TopBsno=10 
	temp.Noise=-90
        temp.Topx=3
//Ramanan changes begin
total_num_ue_CDF := 100                            // total number of UE for CDF calculation is hardcoded
var pre_sinr_arr_dB = []float64{} //array variable to store pre processing SINR for number of UEs considered
var post_sinr_arr_dB = []float64{} //array variable to store post processing SINR for number of UEs considered
//getdata has PrS ,PoS
//fmt.Println("\nBBB",getdata.PrS,"  ",getdata.PoS,"\n")
for i:=0;i<total_num_ue_CDF;i++{ // number of UEs used for caculate CDF hardcoded
getdata = SIR(temp)
pre_sinr_arr_dB = append(pre_sinr_arr_dB,getdata.PrS)
post_sinr_arr_dB=append(post_sinr_arr_dB,getdata.PoS)
temp.Node+=1
}
//fmt.Println("\n pre_sinr_arr_dB \n",pre_sinr_arr_dB,"\n")
sort.Float64s(pre_sinr_arr_dB)//function which sorts data in ascending order
//fmt.Println("\n pre_sinr_arr_dB_after_sorting \n",pre_sinr_arr_dB,"\n")
var min_pre_sinr_dB int
var max_pre_sinr_dB int
min_pre_sinr_dB = int(math.Floor(pre_sinr_arr_dB[0]))
max_pre_sinr_dB = int(math.Ceil(pre_sinr_arr_dB[total_num_ue_CDF-1]))
//fmt.Println("\n post_sinr_arr_dB \n",post_sinr_arr_dB,"\n")
//fmt.Println("\n check \n",max_pre_sinr_dB-min_pre_sinr_dB+1)
//fmt.Println("\n min_sinr_dB \n",min_pre_sinr_dB,"\n")
//fmt.Println("\n max_sinr_dB \n",max_pre_sinr_dB,"\n")
var pre_sinr_cdf float64
cdf_threshold:=min_pre_sinr_dB   // %%%%%%%%%%%%%%threshold is kept as integer as of now%%%%%%%%%%%%%%%%
//for i:=0;i<(max_pre_sinr_dB-min_pre_sinr_dB+1);i++{ // x axis range
for cdf_threshold <= max_pre_sinr_dB{ // x axis range  
returnobj2.X = append(returnobj2.X,float64(cdf_threshold))
pre_sinr_count_ue := 0
for j:=0;j<total_num_ue_CDF;j++{ //  calculateCDF
if pre_sinr_arr_dB[j] <= float64(cdf_threshold){ 	  
pre_sinr_count_ue = pre_sinr_count_ue+1
}
}
pre_sinr_cdf = float64(pre_sinr_count_ue)/float64(total_num_ue_CDF)
//fmt.Println("\n check values \n",pre_sinr_cdf)
returnobj2.Y=append(returnobj2.Y,pre_sinr_cdf)
cdf_threshold = cdf_threshold+1
}
//fmt.Println("\n pre processing sinr cdf level0\n ",returnobj2.X,"\n",returnobj2.Y)
return returnobj2 

//Ramanan changes end 
}

func FR3(r rowdata) returndata3 {
var returnobj3 returndata3
fmt.Println(r.Node,"  ",r.Topx)
//bs,ue,hex,adj
//r.Node is ue no
//r.Topx
returnobj3.PrS=13.13
returnobj3.PoS=1729.22
returnobj3.ROI=12.12
return returnobj3
}













func main() {

//Nodelocations csv file

csvfile1, err1 := os.Open("Sectorlocations.csv")

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
    count := 0
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
    //fmt.Println("AAA",each1)
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
   oneRecord.x=c
  }
 d, err := strconv.ParseFloat(each1[3], 64)
  if err==nil{
   oneRecord.y=d
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


//pathloss csv
csvfile, err := os.Open("SIR.csv")
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
        count=0
         for _, each := range rawCSVdata {

             if count!=0{//to not print first row
           //   fmt.Printf("row",each) 
              temp:=[]float64{} 
                for i := 0; i < len(bs); i++ { 
               a, err := strconv.ParseFloat(each[i], 64)            
                 if err==nil{
                     temp= append(temp,a)
                           }
                 }//for loop over
             Pathloss =append(Pathloss,temp)
              
                 } 
           count+=1

         }//for loop of csv parse is over

fmt.Println(len(Pathloss[0]))
fmt.Println("operator info",op)
fmt.Println("BS no",len(bs)) // no of basestations
fmt.Println("UE no",len(ue)) // no of ues
hexcenter() //creating hex center for all info including center,and ue and bs id's inside each hexagon and adjacency matrix
//fmt.Println(adjhex(12))//testing
//fmt.Println(secondtier(12))//testing
/*//to see hex info
for q:=0;q<len(hex);q++{
 fmt.Println(q,hex[q].x,hex[q].y,hex[q].ue,hex[q].bs)
}
*/

	log.Println("Started Server at :8080")
	http.HandleFunc("/update", handlerroute)
	http.Handle("/", http.FileServer(http.Dir(".")))
	err2 := http.ListenAndServe(":8080", nil)
	if err2 != nil {
		log.Println("Error ", err2)
	}
}