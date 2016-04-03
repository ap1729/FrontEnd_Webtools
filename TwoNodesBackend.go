package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
        "math"
)

/// E.g for sending from browser
// Obj={id:"32",Column1:30,Column2:-45,Column3:-34,Column4:-40}
// $.post("http://localhost:8080/update",JSON.stringify(Obj),"json")
type rowdata struct {
	Id                        string
	Column1, Column2, Column3 float64
	Column4, x1, y1 ,x2 ,y2   float64 
        
//Id is indication of Basestation or node
//Column1 to 4 are x,y,m1 and m2 values originally
//x1,x2,y1,y2 are x and y coordinates of basestation locations
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

func EvaluateMore(r rowdata) rowdata {
	// log.Printf("I received this to process %#v", r)
	// Actual algo goes here ....
        
      if r.Id!="0"{
log.Printf("AAAAQ  	%f %f ",r.x1,r.x2);
//as r.x1,r.x2,r.y1 and r.y2 are zero ,am manually setting it in next line
       r.x1=10
       r.y1=8
       r.x2=10
       r.y2=12
        var a,b,dist1,dist2,BSx1,BSy1,BSx2,BSy2 float64
        a=(r.Column1)/30
        b=(r.Column2+5)/30
        BSx1=(r.x1)/30
        BSy1=(r.y1+5)/30
	BSx2=(r.x2)/30
	BSy2=(r.y2+5)/30
        
	dist1=((a-BSx1)*(a-BSx1)+(b-BSy1)*(b-BSy1)) 
	dist2=((a-BSx2)*(a-BSx2)+(b-BSy2)*(b-BSy2)) 
     
	r.Column4=(0.2)*((1/dist1)+(1/dist2))
	r.Column3=500+80*math.Log(math.Abs(dist1*dist2))
        if r.Column3<0{r.Column3=-r.Column3} 
         log.Printf("%f %f",r.x1,r.y1);
       //   log.Printf("%f,%f,%f,%f,%f,%f	",BSx1,BSy1,r.Column1,r.Column2,a,b)
 //log.Printf("Returned %f %f %f %f",r.Column1,r.Column2,r.Column3,r.Column4)
       } else { 
 log.Println("BASESTATION!!!",r.Id)
//Need to return bunch of data
 } 
       return r
       

}

func main() {
	log.Println("Started Server at :8080")
	http.HandleFunc("/update", handlerroute)
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error ", err)
	}
}
