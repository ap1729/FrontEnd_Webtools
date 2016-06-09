package main

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/service"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"
)

// Global variables that encapsulate all required data
// (Try them out by invoking the suggestion tool by typing the "dot")!
var scenario *model.Scenario
var hexMap *service.HexMap

// E.g for sending from browser
// Obj={id:"32",Column1:30,Column2:-45,Column3:-34,Column4:-40}
// $.post("http://localhost:8080/update",JSON.stringify(Obj),"json")

func initialize() bool {

	// Time stamp 1
	lap1 := time.Now()

	// Read all nodes (BS and UE)
	sc, err := service.ReadNodes("data/OmniLocations.csv")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return false
	}

	// Time stamp 2
	lap2 := time.Now()

	// Import loss values into Scenario object
	err = service.ImportLossTable("data/OmniLosses.csv", sc)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return false
	}

	// Time stamp 3
	lap3 := time.Now()

	// Generate hexagonal cell map of ISD 1000 and upto 3 tiers
	hm := service.NewHexMap(500*2/math.Sqrt(3), 3)

	// Time stamp 4
	lap4 := time.Now()

	// Display execution times
	fmt.Printf("\nPreliminary initialization time estimate:\n")
	fmt.Printf("Location read time: %v\nLosses read time: %v\nCell map init time: %v\n", lap2.Sub(lap1), lap3.Sub(lap2), lap4.Sub(lap3))

	scenario = sc
	hexMap = hm
	return true
}

// Temporary quick fix, changes pending
func handlerroute(w http.ResponseWriter, r *http.Request) {

	// w.Header().Add("Access-Control-Allow-Origin", "*")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only POST method is supported."))
		return
	}
	log.Println("Request Method is:", r.Method)
	// log.Println("Request is originated from  ", r.RemoteAddr)
	// log.Println("Request is originated URL  ", r.RequestURI)
	// log.Println("Request Headers", r.Header)

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Body is:", string(content))

		// Data transfer format (JSON) is equivalent to a generic dictionary: string -> object
		// For more help, see section "Generic JSON with interface{}"
		// at https://blog.golang.org/json-and-go
		var rxData map[string]interface{}
		jerr := json.Unmarshal(content, &rxData)
		if jerr != nil {
			log.Println("Unmarshalling error:", jerr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println("Unmarshalled Data is:", rxData)

		// Each of these types must call calculations defined under package /perf.
		// Each function in perf returns data is a generic string dictionary.
		// Define and pass data to the /perf functions as required - keep them general!
		if rxData["type"] == "A" {
			// TODO: SIR (or whatever)

		} else if rxData["type"] == "B" {
			// TODO: level 1 (or whatever)

		} else if rxData["type"] == "C" {
			// TODO: CDF plot (or whatever)

		} else if rxData["type"] == "D" {
			// TODO: FR3 (or whatever)

		} else {
			fmt.Println("Unknown command")
			return
		}

		// Console feedback
		fmt.Printf("\nUser requested to perform calculations of type \"%v\".\n", rxData["type"])

	}
}

func main() {

	initSuccess := initialize()
	if initSuccess == false {
		fmt.Println("Fatal error! Failed to load data.")
	} else {
		fmt.Println("\nSuccessfully loaded data.")
		fmt.Printf("There are %d BS's and %d UE's.\n", len(scenario.BaseStations), len(scenario.Users))
		fmt.Printf("The loss table is a %d x %d matrix.\n\n", len(scenario.LossTable), len(scenario.LossTable[0]))
	}

	log.Println("\nStarted Server at :8080")
	http.HandleFunc("/update", handlerroute)
	http.Handle("/", http.FileServer(http.Dir(".")))
	err2 := http.ListenAndServe(":8080", nil)
	if err2 != nil {
		log.Println("Error ", err2)
	}
}
