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

// TODO: Move these consts to a config file.
// Data generation modes:
// Read CSV - "import"
// Generate manually - "manual"
const dataGenOpt = "import"
const locFilePath = "data/HataMossNodelocations.csv"
const lossFilePath = "data/HataMossLosses.csv"

// Package scope variables that encapsulate all required data
// (Try them out by invoking the suggestion tool by typing the "dot")!
var scenario *model.Scenario
var hexMap *service.HexMap
var opEnable []bool

// A package level object to store return data
// *See notes regarding this in "updatehandler()"
var response *Response

// E.g for sending from browser
// Obj={id:"32",Column1:30,Column2:-45,Column3:-34,Column4:-40}
// $.post("http://localhost:8080/update",JSON.stringify(Obj),"json")

func initialize() bool {

	// Time stamp variables
	var lap1, lap2, lap3 time.Time

	// Time stamp 1
	lap1 = time.Now()

	sb := model.NewScenarioBuilder()
	if dataGenOpt == "import" {
		// Read all nodes (BS and UE)
		suc := service.ReadNodes(sb, locFilePath)
		if suc == false {
			// fmt.Printf("Error: %v", err)
			return false
		}
		// Time stamp 2
		lap2 = time.Now()
		// Import loss values into Scenario object
		suc = service.ReadLossTable(sb, lossFilePath)
		if suc == false {
			// fmt.Printf("Error: %v", err)
			return false
		}

	} else if dataGenOpt == "manual" {
		suc := service.GenerateMap(sb)
		if suc == false {
			return false
		}
		// Time stamp 2
		lap2 = time.Now()
		// Generate losses using path loss model manually
		suc = sb.Seal("calc", nil)
		if suc == false {
			return false
		}

	} else {
		return false
	}

	scenario = sb.Finalize()
	sb = nil
	opEnable = make([]bool, len(scenario.Operators()))
	for i := 0; i < len(scenario.Operators()); i++ {
		opEnable[i] = true
	}

	// Time stamp 3
	lap3 = time.Now()

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
	response = nil
	// TODO: Try making a pointer to a response object and capture it in the defer command.
	// Perhaps then, updating the values in this function will reflect on the deferred
	// execution. (Instead of making a global variable which is absolutely not thread safe IMO.)
	defer sendResponse(&w)
	// Additionally, sending back any response with HTTP codes still invokes the deferred
	// response function, due to which multiple responses may be issued.

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

		var r = NewResponse()

		// This is safe; if the key does not exist, the variable is assigned its default zero value.
		frMode, _ := rxData["frmode"].(string)
		ueID, _ := rxData["node"].(float64)
		curLevel, _ := rxData["level"].(float64)
		intCancelCount, _ := rxData["intcnc"].(float64)
		topN, _ := rxData["topbsno"].(float64)

		if ueID < 0 || topN < 0 || intCancelCount < 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid parameters in request."))
			return
		}

		params := &perf.Params{FrMode: frMode, Level: uint(curLevel), IntCancellers: uint(intCancelCount), OpEnableFlags: opEnable}

		switch rxData["perf"] {
		case "scmeta":
			r.data = service.PackageScenario(scenario)
		case "enop":
			vals := rxData["opflags"].([]interface{})
			for i := 0; i < len(scenario.Operators()); i++ {
				opEnable[i] = vals[i].(float64) == 1
			}
			r.data, r.err = perf.AssignOperators(scenario, opEnable)
		case "lvlchng":
			targetLvl := uint(rxData["params"].(float64))
			r.data, r.err = perf.ChangeLevel(scenario, targetLvl, opEnable)
		case "emer":
			r.data, r.err = perf.EmDownlink(scenario, hexMap, opEnable)
		case "sir":
			fmt.Printf("topN is: %v, and somehow uint(topN) is %v", topN, uint(topN))
			r.data, r.err = perf.SinrProfile(scenario, hexMap, uint(ueID), uint(topN), params)
		case "heatmap":
			r.data = perf.SinrHeatMap(scenario, hexMap, params)
		case "cdf":
			r.data = perf.CDF(scenario, hexMap, params)
		default:
			fmt.Println("\nFATAL: Unknown command")
			return
		}

		// Console feedback
		fmt.Printf("\nUser requested to perform calculations of type \"%v\".\n", rxData["perf"])
		fmt.Printf("\nAn error occured in perf:\n%v", err)
		response = r

	}
}

func sendResponse(w *http.ResponseWriter) {
	rStat := recover()

	var status int
	var msg string
	var data map[string]interface{} = nil

	// If panic takes place, recover status in non-nil.
	if rStat != nil {
		status = 1
		msg = "An unknown error occured at our end."
		// Print stack trace for debugging assitance
		debug.PrintStack()
		fmt.Printf("\n\nRecovered :)\nError encountered: %v\n\n", rStat)
	} else if response == nil {
		status = 1
		msg = "An unknown error occured at our end."
		fmt.Println("\nSomething funny and unexplained happened :/")
	} else {
		if response.err == nil {
			status = 0
			msg = "Successful execution."
		} else {
			status = 2
			msg = "We encountered an error, which has been logged and will be fixed soon."
		}
	}

	if status == 0 {
		data = response.data
	}

	var respMap = map[string]interface{}{}
	respMap["status"] = status
	respMap["data"] = data
	respMap["msg"] = msg
	serializedData, _ := json.Marshal(respMap)
	txbytes, err := (*w).Write(serializedData)
	if err != nil {
		log.Println("Failed to send response to user.\nError: ", err)
	} else {
		log.Println("Response sent:\n", string(txbytes))
	}

}

type Response struct {
	status int
	data   map[string]interface{}
	err    error
}

func NewResponse() *Response {
	return &Response{status: -1, data: nil, err: nil}
}
