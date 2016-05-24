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
)
var Pathloss = [][]float64{} //global variable to store path loss as a 2d array

// E.g for sending from browser
// Obj={id:"32",Column1:30,Column2:-45,Column3:-34,Column4:-40}
// $.post("http://localhost:8080/update",JSON.stringify(Obj),"json")
type rowdata struct {
	Node  int
        Level int
	
}//data which is coming from user

type returndata struct{
        SIR []float64
        PrS float64
        PoS float64
        ROI float64
}// structure for returning data to front end

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
		udpateData := EvaluateMore(rxdata)
		txbytes, _ := json.Marshal(udpateData)

		nbytes, werr := w.Write(txbytes)
		_ = nbytes
		if werr != nil {
			log.Println("I got some error while writing back", werr)
		} else {
			 log.Println("Sent this  ", string(txbytes))
			// log.Printf("Successfully returned %d bytes", nbytes)
		}

	}

}

func EvaluateMore(r rowdata) returndata {
	// log.Printf("I received this to process %#v", r)
	// Actual algo goes here ....
       row:=Pathloss[r.Node]
      //row variable has the path loss for that node 
      //sort it in descending order to get SIR for level 1
      //check if level is 0 ,then do a cyclic shift 
      //return array of top 10
     //Also Calculate SINR's and send back

     var returnobj returndata
     returnobj.SIR=row 
     returnobj.PrS=13.13
     returnobj.PoS=57.57



     fmt.Println(returnobj)
       return returnobj
}

func main() {


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
             if count!=0{//to only print top 2 rows
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
//       fmt.Println(Pathloss)
	log.Println("Started Server at :8080")
	http.HandleFunc("/update", handlerroute)
	http.Handle("/", http.FileServer(http.Dir(".")))
	err1 := http.ListenAndServe(":8080", nil)
	if err1 != nil {
		log.Println("Error ", err1)
	}
}
