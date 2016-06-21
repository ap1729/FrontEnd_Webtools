package main

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/perf"
	"FrontEnd_WebTools/service"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"runtime/debug"
	"time"
)

// Package scope variables that encapsulate all required data
// (Try them out by invoking the suggestion tool by typing the "dot")!
var scenario *model.Scenario
var hexMap *service.HexMap

// A package level object to store return data
var responseData map[string]interface{}

// E.g for sending from browser
// Obj={id:"32",Column1:30,Column2:-45,Column3:-34,Column4:-40}
// $.post("http://localhost:8080/update",JSON.stringify(Obj),"json")

func initialize() bool {

	// Time stamp 1
	lap1 := time.Now()

	// Read all nodes (BS and UE)
	sb := model.NewScenarioBuilder()
	suc := service.ReadNodes(sb, "data/SectorLocations.csv")
	if suc == false {
		// fmt.Printf("Error: %v", err)
		return false
	}

	// Time stamp 2
	lap2 := time.Now()

	// Import loss values into Scenario object
	suc = service.ReadLossTable(sb, "data/SectorLosses.csv")
	if suc == false {
		// fmt.Printf("Error: %v", err)
		return false
	}
	scenario = sb.Finalize()
	sb = nil

	// Time stamp 3
	lap3 := time.Now()

	// Generate hexagonal cell map of ISD 1000 and upto 3 tiers
	hm := service.NewHexMap(500*2/math.Sqrt(3), 3)
	hexMap = hm
	hexMap.AssociateStations(scenario.BaseStations())
	hexMap.AssociateUsers(scenario.Users())

	// Time stamp 4
	lap4 := time.Now()

	// Display execution times
	fmt.Printf("\nPreliminary initialization time estimate:\n")
	fmt.Printf("Location read time: %v\nLosses read time: %v\nCell map init time: %v\n", lap2.Sub(lap1), lap3.Sub(lap2), lap4.Sub(lap3))

	return true
}

func main() {

	initSuccess := initialize()
	if initSuccess == false {
		fmt.Println("Fatal error! Failed to load data.")
	} else {
		fmt.Println("\nSuccessfully loaded data.")
		fmt.Printf("There are %d BS's and %d UE's.\n", len(scenario.BaseStations()), len(scenario.Users()))
	}

	log.Println("\nStarted Server at :8080")
	http.HandleFunc("/update", updateHandler)
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error: ", err)
	}
}

// Temporary quick fix, changes pending
func updateHandler(w http.ResponseWriter, r *http.Request) {
	// Allow Cross-Origin Requests
	w.Header().Add("Access-Control-Allow-Origin", "*")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only POST method is supported."))
		return
	}
	// Request info display
	log.Println("Request Method is:", r.Method)
	log.Println("Request is originated from  ", r.RemoteAddr)
	log.Println("Request is originated URL  ", r.RequestURI)
	log.Println("Request Headers", r.Header)

	// A safety net for handling panics
	defer sendResponse(&w)

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

		// JSON structure:
		// frmode: Frequency-Reuse mode (Ex. "FR1", "FR3")
		// node: User ID (uint)
		// level: Cooperation Level (uint)
		// intcnc: Interference canceller count (uint)
		// topbsno: The top N stations who's profile is shown (uint)
		// perf: The performance metric to evaluate (Ex. "cdf", "sir", "lvlchng")
		// opflags: The flags that specify which operators are active (array of binary)
		// params: Any additional details that the perf function may need

		var returnData map[string]interface{}

		switch rxData["perf"] {
		case "scmeta":
			returnData = service.PackageScenario(scenario)
		case "lvlchng":
			targetLvl := uint(rxData["params"].(float64))
			returnData = perf.ChangeLevel(scenario, targetLvl)
			fmt.Println("Level Change complete.")
		case "sir":
			frMode := rxData["frmode"].(string)
			ueID := uint(rxData["node"].(float64))
			level := uint(rxData["level"].(float64))
			intCancelCount := uint(rxData["intcnc"].(float64))
			topN := uint(rxData["topbsno"].(float64))

			var params map[string]interface{}
			if frMode == "FR3" || frMode == "FFR" {
				params = map[string]interface{}{}
				params["hexmap"] = hexMap
			}
			if frMode == "FFR" {
				params["intcnc"] = intCancelCount
			}

			returnData = perf.SinrProfile(scenario, frMode, ueID, level, intCancelCount, topN, params)
			fmt.Println("SIR calculation complete.")
		case "cdf":
			frMode := rxData["frmode"].(string)
			intCancelCount := uint(rxData["intcnc"].(float64))
			var params map[string]interface{}
			if frMode == "FR3" || frMode == "FFR" {
				params = map[string]interface{}{}
				params["hexmap"] = hexMap
			}
			if frMode == "FFR" {
				params["intcnc"] = intCancelCount
			}
			returnData = perf.CDF(scenario, frMode, intCancelCount, params)
			fmt.Println("CDF calc done")
		default:
			fmt.Println("Unknown command")
			return
		}

		// Console feedback
		fmt.Printf("\nUser requested to perform calculations of type \"%v\".\n", rxData["perf"])
		responseData = returnData

	}
}

func sendResponse(w *http.ResponseWriter) {
	rStat := recover()
	response := map[string]interface{}{}

	if rStat == nil {
		fmt.Printf("\nSuccesful execution.\n")
		response["status"] = 0
		response["data"] = responseData

	} else {
		debug.PrintStack()
		fmt.Printf("\n\nRecovered :)\nError encountered: %v\n\n", rStat)
		response["status"] = 1
		response["data"] = ""
	}

	serializedData, _ := json.Marshal(response)
	txbytes, err := (*w).Write(serializedData)
	if err != nil {
		log.Println("Failed to send response to user.\nError: ", err)
	} else {
		log.Println("Response sent:\n", string(txbytes))
	}

}
