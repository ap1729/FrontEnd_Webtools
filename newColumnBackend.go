package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
)

/// E.g for sending from browser
// Obj={id:"32",Column1:30,Column2:-45,Column3:-34,Column4:-40}
// $.post("http://localhost:8080/update",JSON.stringify(Obj),"json")
type rowdata struct {
	Id                        string
	Column1, Column2, Column3 float64
	Column4, Column5          float64
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
			log.Println("Sent this  ", string(txbytes))
			// log.Printf("Successfully returned %d bytes", nbytes)
		}

	}

}

func EvaluateMore(r rowdata) rowdata {
	// log.Printf("I received this to process %#v", r)
	// Actual algo goes here ....
	var a, b, dist1, dist2 float64
	if (r.Column5 == 1) || (r.Column5 == 2) {
		a = (r.Column1) / 30
		b = (r.Column2 + 5) / 30

		dist1 = ((a-0.25)*(a-0.25) + (b-0.5)*(b-0.5))
		dist2 = ((a-0.75)*(a-0.75) + (b-0.5)*(b-0.5))

		r.Column4 = (0.2) * ((1 / dist1) + (1 / dist2))
		r.Column3 = 500 + 80*math.Log(math.Abs(dist1*dist2))
		if r.Column3 < 0 {
			r.Column3 = -r.Column3
		}
		if math.IsInf(r.Column3, 0) || math.IsInf(r.Column4, 0) {
			r.Column3 = -1
			r.Column4 = -1
		}

		log.Printf("Returned %f %f %f %f",r.Column1,r.Column2,r.Column3,r.Column4)

	} else {

		log.Printf("BASESTATION")
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