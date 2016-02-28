package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

/// E.g for sending from browser
// Obj={id:"32",Column1:30,column2:-45,column3:-34,column4:-40}
// $.post("http://localhost:8080/update",JSON.stringify(Obj),"json")
type rowdata struct {
	Id                        string
	Column1, Column2, Column3 float64
	Column4                   float64
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

		udpateData := EvaluateMore(rxdata)

		txbytes, _ := json.Marshal(udpateData)

		nbytes, werr := w.Write(txbytes)
		_ = nbytes
		if werr != nil {
			log.Println("I got some error while writing back", werr)
		} else {
			// log.Println("Sent this  ", string(txbytes))
			// log.Printf("Successfully returned %d bytes", nbytes)
		}

	}

}

func EvaluateMore(r rowdata) rowdata {
	// log.Printf("I received this to process %#v", r)
	// Actual algo goes here ....

	r.Column3 = r.Column3 + 25
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
