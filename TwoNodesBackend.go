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
)

 type TestRecord struct {
         col1 string
         node  string
          x float64
         y  float64
 	m1 float64
         m2  float64
 }
var allRecords []TestRecord//global
//object for storing values


/// E.g for sending from browser
// Obj={id:"32",Column1:30,Column2:-45,Column3:-34,Column4:-40}
// $.post("http://localhost:8080/update",JSON.stringify(Obj),"json")
type rowdata struct {
	Node,Col                        int
	Column1, Column2, Column3 ,Column4 float64
        X1, Y1 ,X2 ,Y2   float64 
        
//node is indication of Basestation or node
//Column1 to 4 are x,y,m1 and m2 values originally
//X1,X2,Y1,Y2 are x and y coordinates of basestation locations
}

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
                //log.Println("Marshalled Data",rxdata)	
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

func EvaluateMore(r rowdata) rowdata {
	// log.Printf("I received this to process %#v", r)
	// Actual algo goes here ....
        
      if r.Node!=0{
//as r.X1,r.X2,r.Y1 and r.Y2 are zero ,am manually setting it in next line
     //  r.X1=7.5
     //r.Y1=10
     //r.X2=22.5
    //   r.Y2=10
        var a,b,dist1,dist2,BSX1,BSY1,BSX2,BSY2 float64
        a=(r.Column1)/30
        b=(r.Column2+5)/30
        BSX1=(r.X1)/30
        BSY1=(r.Y1+5)/30
	BSX2=(r.X2)/30
	BSY2=(r.Y2+5)/30
        
	dist1=((a-BSX1)*(a-BSX1)+(b-BSY1)*(b-BSY1)) 
	dist2=((a-BSX2)*(a-BSX2)+(b-BSY2)*(b-BSY2)) 
     
	r.Column4=(0.2)*((1/dist1)+(1/dist2))
	r.Column3=500+80*math.Log(math.Abs(dist1*dist2))
        if r.Column3<0{r.Column3=0} 
    allRecords[r.Col].x=r.X1
allRecords[r.Col].y=r.Y1
allRecords[r.Col].m1=r.Column3
allRecords[r.Col].m2=r.Column4




   //      log.Printf("%f %f",r.X1,r.Y1);
       //   log.Printf("%f,%f,%f,%f,%f,%f	",BSX1,BSY1,r.Column1,r.Column2,a,b)
 //log.Printf("Returned %f %f %f %f",r.Column1,r.Column2,r.Column3,r.Column4)
       } else { 
//Basesation stuff of saving session and BS coords
 //log.Println("BASESTATION!!!")
//Need to return bunch of data
 } 
       return r
       

}

func main() {

csvfile, err := os.Open("testfile5.csv")
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
         var oneRecord TestRecord
         
         for _, each := range rawCSVdata {
              if(each[0]!="col1"){         //to ensure first row of csv does not parse in
                 oneRecord.col1 = each[0]
                 oneRecord.node = each[1]
 a, err := strconv.ParseFloat(each[2], 64)
  if err==nil{
   oneRecord.x=a
//fmt.Printf("\nQQQQ%f\n",oneRecord.x)
  }
 b, err := strconv.ParseFloat(each[3], 64)
  if err==nil{
   oneRecord.y=b
  }
 c, err := strconv.ParseFloat(each[4], 64)
  if err==nil{
   oneRecord.m1=c
  }
 d, err := strconv.ParseFloat(each[5], 64)
  if err==nil{
   oneRecord.m2=d
  }
		
  allRecords = append(allRecords, oneRecord)
         }
         }//for loop of csv over
//fmt.Println(allRecords[0:3])
         // second sanity check
//Above part is csv parsing and storing all original values in allRecords object









	log.Println("Started Server at :8080")
	http.HandleFunc("/update", handlerroute)
	http.Handle("/", http.FileServer(http.Dir(".")))
	err1 := http.ListenAndServe(":8080", nil)
	if err1 != nil {
		log.Println("Error ", err1)
	}
}
