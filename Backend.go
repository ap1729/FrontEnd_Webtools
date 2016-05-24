package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
        //"math"
        //"encoding/csv"
         //"fmt"
         //"os"
	 //"strconv"
)

type rowdata struct {
        node string 
	ueno int 
        operno float64
//ueno is no of ue 
}//data which is coming from user







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
		log.Println("body is ", string(data))//body is proper
               
		var rxdata rowdata
		jerr := json.Unmarshal(data, &rxdata)
		if jerr != nil {
			log.Println("Unmarshalling error ", jerr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

                log.Println("UnMarshalled recieved Data",rxdata)//some problem unmarshalling	





		udpateData:= EvaluateMore(rxdata)
                log.Println("Updated Data",udpateData)
		txbytes, _ := json.Marshal(udpateData)
                log.Println("Marshalled updated Data",string(txbytes))
		nbytes, werr := w.Write(txbytes)
		_ = nbytes
		if werr != nil {
			log.Println("I got some error while writing back", werr)
		} else {
			 log.Println("Sent this  ",string(txbytes))
			// log.Printf("Successfully returned %d bytes", nbytes)
		}

	}

}


func EvaluateMore(r rowdata) rowdata{
 log.Printf("I received this to process %#v", r)
//calculation done here
r.ueno=2
r.operno=5.0
 log.Printf("After process %#v", r)
return r
}



func main() {


log.Println("Started Server at :8080")
	http.HandleFunc("/update", handlerroute)
	http.Handle("/", http.FileServer(http.Dir(".")))
	err1 := http.ListenAndServe(":8080", nil)
	if err1 != nil {
		log.Println("Error ", err1)
         }

}








